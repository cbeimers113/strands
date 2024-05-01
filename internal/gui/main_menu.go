package gui

import (
	"github.com/g3n/engine/gui"

	"cbeimers113/strands/internal/gui/component"
	"cbeimers113/strands/internal/io/file"
	"cbeimers113/strands/internal/state"
)

// Register the main menu
func (g *Gui) registerMainMenu() {
	views[MainMenu] = viewControls{
		open: func(closeOthers bool) {
			if closeOthers {
				g.closeViews()
			}

			var (
				width, height int = g.App.GetSize()
				w, h          float32
				nextY         float32
			)

			g.startButton = gui.NewButton("Enter Simulation")
			w, h = g.startButton.Width(), g.startButton.Height()
			g.startButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
			g.startButton.SetUserData(MainMenu)
			g.startButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.State.SetInMenu(false)
				views[SimulationView].open(true)
				g.Keyboard.Clear()
				g.EnableKeyboard(false)
			})
			g.Scene.Add(g.startButton)
			nextY = g.startButton.Position().Y + g.startButton.Height() + 5

			g.saveButton = gui.NewButton("Save Simulation")
			w = g.saveButton.Width()
			g.saveButton.SetPosition((float32(width)-w)/2, nextY)
			g.saveButton.SetUserData(MainMenu)
			g.saveButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.openSaveDialog()
			})
			g.Scene.Add(g.saveButton)
			nextY = g.saveButton.Position().Y + g.saveButton.Height() + 5

			g.settingsButton = gui.NewButton("Settings")
			w = g.settingsButton.Width()
			g.settingsButton.SetPosition((float32(width)-w)/2, nextY)
			g.settingsButton.SetUserData(MainMenu)
			g.settingsButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.Keyboard.Clear()
				g.EnableKeyboard(false)
				Open(ConfigMenu, true)
			})
			g.Scene.Add(g.settingsButton)
			nextY = g.settingsButton.Position().Y + g.settingsButton.Height() + 5

			g.exitButton = gui.NewButton("Exit")
			w = g.exitButton.Width()
			g.exitButton.SetPosition((float32(width)-w)/2, nextY)
			g.exitButton.SetUserData(MainMenu)
			g.exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.Keyboard.Clear()
				g.EnableKeyboard(false)
				g.App.Exit()
			})
			g.Scene.Add(g.exitButton)

			if g.dialogOpen {
				g.openSaveDialog()
			}

			g.State.SetInMenu(true)
		},

		close: func() {
			g.Keyboard.Clear()
			g.EnableKeyboard(false)
			g.Scene.Remove(g.startButton)
			g.Scene.Remove(g.saveButton)
			g.Scene.Remove(g.saveNameField)
			g.Scene.Remove(g.submitButton)
			g.Scene.Remove(g.settingsButton)
			g.Scene.Remove(g.exitButton)
			g.dialogOpen = false
		},

		refresh: func() {
			// Update the input field text if the dialog is opened
			if g.dialogOpen && g.saveNameField.Button != nil && g.saveNameField.Button.Label != nil {
				var (
					x0 = g.saveNameField.Position().X
					y0 = g.saveNameField.Position().Y
					x1 = x0 + g.saveNameField.Width()
					y1 = y0 + g.saveNameField.Height()
				)

				g.saveNameField.Update(g.Keyboard.Read())
				if g.Keyboard.ClickOutCheck(x0, y0, x1, y1) {
					g.EnableKeyboard(false)
				}
			}
		},
	}
}

func (g *Gui) openSaveDialog() {
	y := g.startButton.Position().Y + g.startButton.Height() + 5
	h := g.startButton.Height()

	// Load input field for save name if the save button is pressed
	g.saveNameField = component.NewInputField(150, h, "", 500)
	g.saveNameField.SetPosition(5+g.saveButton.Position().X+g.saveButton.Width(), y)
	g.saveNameField.SetUserData(MainMenu)
	g.saveNameField.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		if g.Keyboard.GetEnabled() {
			return
		}

		g.EnableKeyboard(true)
	})
	g.Scene.Add(g.saveNameField)

	// Also load the submit save button
	g.submitButton = gui.NewButton("Save")
	g.submitButton.SetPosition(5+g.saveNameField.Position().X+g.saveNameField.Width(), y)
	g.submitButton.SetUserData(MainMenu)
	g.submitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		filename := g.Keyboard.Read()
		if len(filename) == 0 {
			return
		}

		// Save the game
		g.SaveGame(file.Touch(filename, state.SaveFileExtension))
		g.EnableKeyboard(false)
		g.Keyboard.Clear()
	})
	g.Scene.Add(g.submitButton)

	g.saveButton.SetEnabled(false)
	g.dialogOpen = true
}

func (g *Gui) EnableKeyboard(enable bool) {
	if enable {
		g.Keyboard.Enable(true)

		if g.saveNameField != nil {
			g.saveNameField.SetEnabled(false)
			g.saveNameField.Activate(true)
		}
	} else {
		g.Keyboard.Enable(false)

		if g.saveNameField != nil {
			g.saveNameField.SetEnabled(true)
			g.saveNameField.Activate(false)
		}
	}
}

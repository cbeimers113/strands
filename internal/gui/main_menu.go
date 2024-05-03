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
				views[SimulationView].open(true)
				g.closeSaveDialog()
			})
			g.Scene.Add(g.startButton)
			nextY = g.startButton.Position().Y + g.startButton.Height() + 5

			g.saveButton = gui.NewButton("Save Simulation")
			w = g.saveButton.Width()
			g.saveButton.SetPosition((float32(width)-w)/2, nextY)
			g.saveButton.SetUserData(MainMenu)
			g.saveButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.closeOpenDialog()
				g.openSaveDialog()
			})
			g.Scene.Add(g.saveButton)
			nextY = g.saveButton.Position().Y + g.saveButton.Height() + 5

			g.openButton = gui.NewButton("Open Simulation")
			w = g.openButton.Width()
			g.openButton.SetPosition((float32(width)-w)/2, nextY)
			g.openButton.SetUserData(MainMenu)
			g.openButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.openOpenDialog()
			})
			g.Scene.Add(g.openButton)
			nextY = g.openButton.Position().Y + g.openButton.Height() + 5

			g.settingsButton = gui.NewButton("Settings")
			w = g.settingsButton.Width()
			g.settingsButton.SetPosition((float32(width)-w)/2, nextY)
			g.settingsButton.SetUserData(MainMenu)
			g.settingsButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.closeSaveDialog()
				Open(ConfigMenu, true)
			})
			g.Scene.Add(g.settingsButton)
			nextY = g.settingsButton.Position().Y + g.settingsButton.Height() + 5

			g.exitButton = gui.NewButton("Exit")
			w = g.exitButton.Width()
			g.exitButton.SetPosition((float32(width)-w)/2, nextY)
			g.exitButton.SetUserData(MainMenu)
			g.exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.closeSaveDialog()
				g.App.Exit()
			})
			g.Scene.Add(g.exitButton)

			g.State.SetInMenu(true)
		},

		close: func() {
			g.closeSaveDialog()
			g.closeOpenDialog()

			g.Scene.Remove(g.startButton)
			g.Scene.Remove(g.saveButton)
			g.Scene.Remove(g.openButton)
			g.Scene.Remove(g.savesList)
			g.Scene.Remove(g.settingsButton)
			g.Scene.Remove(g.exitButton)

			g.State.SetInMenu(false)
		},

		refresh: func() {
			// Update the input field text if the dialog is opened
			if g.saveDialog != nil && g.saveDialog.Enabled() {
				var (
					x0 = g.saveDialog.Position().X
					y0 = g.saveDialog.Position().Y
					x1 = x0 + g.saveDialog.Width()
					y1 = y0 + g.saveDialog.Height()
				)

				// If there was a click event outside of the keyboard, disable it
				g.saveDialog.Update(g.Keyboard.Read())
				if g.Keyboard.ClickOutCheck(x0, y0, x1, y1) {
					g.enableKeyboard(false)
				}
			}

			// Check if a file has been slected for loading
			if g.savesList != nil && g.savesList.Selected != "" {
				g.Context.LoadGame(g.savesList.Selected)
				g.savesList.Selected = ""
				Open(SimulationView, true)
			}
		},
	}
}

func (g *Gui) openSaveDialog() {
	y := g.startButton.Position().Y + g.startButton.Height() + 5
	h := g.startButton.Height()

	// Load input field for save name if the save button is pressed
	g.saveDialog = component.NewDialog("", 5+g.saveButton.Position().X+g.saveButton.Width(), y, 150, h, int(MainMenu),

		// On click input field
		func() {
			if g.Keyboard.GetEnabled() {
				return
			}

			g.enableKeyboard(true)
		},

		// On submit...
		func() {
			filename := g.Keyboard.Read()
			if len(filename) == 0 {
				return
			}

			// Save the game
			g.SaveGame(file.Touch(filename, state.SaveFileExtension))
			Open(SimulationView, true)
		},

		// On Cancel
		func() {
			g.closeSaveDialog()
		},
	)

	g.saveDialog.SetUserData(MainMenu)
	g.saveDialog.Subscribe(gui.OnClick, func(name string, ev interface{}) {

	})
	g.Scene.Add(g.saveDialog)

	g.saveButton.SetEnabled(false)
	g.saveDialog.Open()
}

func (g *Gui) closeSaveDialog() {
	g.enableKeyboard(false)
	g.Keyboard.Clear()
	g.Scene.Remove(g.saveDialog)

	if g.saveButton != nil {
		g.saveButton.SetEnabled(true)
	}
}

func (g *Gui) openOpenDialog() {
	g.closeSaveDialog()
	g.openButton.SetEnabled(false)

	y := g.startButton.Position().Y
	h := g.exitButton.Position().Y + g.exitButton.Height() - g.startButton.Position().Y

	// Load savefile list so we can select one to open
	filepaths := state.GetSavesList()
	g.savesList = component.NewFileList(filepaths, 250, h, int(MainMenu))
	g.savesList.SetPosition(5+g.saveButton.Position().X+g.saveButton.Width(), y)
	g.savesList.SetUserData(MainMenu)
	g.Scene.Add(g.savesList)

	// Add a cancel button to go back
	g.cancelButton = gui.NewButton("Cancel")
	g.cancelButton.SetPosition(5+g.savesList.Position().X+g.savesList.GetPanel().Width(), y)
	g.cancelButton.SetUserData(MainMenu)
	g.cancelButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		g.closeOpenDialog()
	})
	g.Scene.Add(g.cancelButton)
}

func (g *Gui) closeOpenDialog() {
	g.Scene.Remove(g.savesList)
	g.Scene.Remove(g.cancelButton)

	if g.openButton != nil {
		g.openButton.SetEnabled(true)
	}
}

func (g *Gui) enableKeyboard(enable bool) {
	if enable {
		g.Keyboard.Enable(true)

		if g.saveDialog != nil {
			g.saveDialog.SetEnabled(true)
		}
	} else {
		g.Keyboard.Enable(false)

		if g.saveDialog != nil {
			g.saveDialog.SetEnabled(false)
		}
	}
}

package gui

import (
	"fmt"
	"os"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/gui/color"
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
			g.startButton.Subscribe(gui.OnCursor, func(s string, i interface{}) {
				if !g.startButton.Enabled() {
					return
				}

				g.startButton.SetColor(color.Green)
				g.startButton.Label.SetColor(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
			})
			g.Scene.Add(g.startButton)
			nextY = g.startButton.Position().Y + g.startButton.Height() + 5

			g.saveButton = gui.NewButton("Save Simulation")
			w = g.saveButton.Width()
			g.saveButton.SetPosition((float32(width)-w)/2, nextY)
			g.saveButton.SetUserData(MainMenu)
			g.saveButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.closeBrowseDialog()
				g.openSaveDialog()
			})
			g.Scene.Add(g.saveButton)
			nextY = g.saveButton.Position().Y + g.saveButton.Height() + 5

			g.browseButton = gui.NewButton("Browse Saves")
			w = g.browseButton.Width()
			g.browseButton.SetPosition((float32(width)-w)/2, nextY)
			g.browseButton.SetUserData(MainMenu)
			g.browseButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.openBrowseDialog()
			})
			g.Scene.Add(g.browseButton)
			nextY = g.browseButton.Position().Y + g.browseButton.Height() + 5

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
				g.closeBrowseDialog()

				g.startButton.SetEnabled(false)
				g.saveButton.SetEnabled(false)
				g.browseButton.SetEnabled(false)
				g.settingsButton.SetEnabled(false)
				g.exitButton.SetEnabled(false)
				if g.savesList != nil {
					g.savesList.SetEnabled(false)
				}
				if g.cancelButton != nil {
					g.cancelButton.SetEnabled(false)
				}

				g.Popup("Exit the simulation?", "Exit", func() {
					g.App.Exit()
				})
			})
			g.exitButton.Subscribe(gui.OnCursor, func(s string, i interface{}) {
				if !g.exitButton.Enabled() {
					return
				}

				g.exitButton.SetColor(color.Red)
				g.exitButton.Label.SetColor(&math32.Color{})
			})
			g.Scene.Add(g.exitButton)

			g.State.SetInMenu(true)
		},

		close: func() {
			g.closeSaveDialog()
			g.closeBrowseDialog()

			g.Scene.Remove(g.startButton)
			g.Scene.Remove(g.saveButton)
			g.Scene.Remove(g.browseButton)
			g.Scene.Remove(g.savesList)
			g.Scene.Remove(g.settingsButton)
			g.Scene.Remove(g.exitButton)
			g.Scene.Remove(g.popup)

			g.State.SetInMenu(false)
		},

		refresh: func() {
			// Update the input field text if the save dialog is opened
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

			// Check if a file has been slected for loading and the confirmation popup needs to be opened
			if g.savesList != nil && g.savesList.Selected != "" && (g.popup == nil || !g.popup.Enabled()) {
				g.openConfirmOpenPopup()
			}

			// Check if a file has been selected for deleting and the confirmation popup needs to be opened
			if g.savesList != nil && g.savesList.Deleted != "" && (g.popup == nil || !g.popup.Enabled()) {
				g.openConfirmDeletePopup()
			}
		},
	}
}

func (g *Gui) openSaveDialog() {
	x := 5 + g.saveButton.Position().X + g.saveButton.Width()
	y := g.startButton.Position().Y + g.startButton.Height() + 5
	w := float32(150)
	h := g.startButton.Height()

	// Load input field for save name if the save button is pressed
	g.saveDialog = component.NewDialog(
		"",            // Default text
		"Save",        // Submit button text
		x,             // x position on screen
		y,             // y position on screen
		w,             // width
		h,             // height
		int(MainMenu), // parent menu

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

	g.saveButton.SetEnabled(false)
	g.saveDialog.SetUserData(MainMenu)
	g.Scene.Add(g.saveDialog)
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

func (g *Gui) openBrowseDialog() {
	g.closeSaveDialog()

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
	g.cancelButton.SetPosition(5+g.savesList.Position().X+g.savesList.Width(), y)
	g.cancelButton.SetUserData(MainMenu)
	g.cancelButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		g.closeBrowseDialog()
	})
	g.Scene.Add(g.cancelButton)

	g.browseButton.SetEnabled(false)
}

func (g *Gui) closeBrowseDialog() {
	g.Scene.Remove(g.savesList)
	g.Scene.Remove(g.cancelButton)

	if g.browseButton != nil {
		g.browseButton.SetEnabled(true)
	}
}

func (g *Gui) openConfirmOpenPopup() {
	if g.savesList == nil {
		return
	}

	width, height := g.App.GetSize()
	w := float32(150) // Min width, if the text is too long it will expand
	h := float32(100)
	x := (float32(width) - w) / 2
	y := (float32(height) - h) / 2

	prompt := fmt.Sprintf("Open %s?\nUnsaved simulation will be lost.", file.Name(g.savesList.Selected))
	g.popup = component.NewPopup(
		prompt,        // Prompt label
		"Okay",        // Submit button text
		x,             // x position on screen
		y,             // y position on screen
		w,             // popup width
		h,             // popup height,
		int(MainMenu), // parent menu

		// on submit...
		func() {
			g.LoadGame(g.savesList.Selected)
			g.savesList.Selected = ""
			g.closePopup()
			Open(SimulationView, true)
		},

		// on cancel
		func() {
			// Enable the other active components in this menu
			g.savesList.Selected = ""
			g.closePopup()
			g.browseButton.SetEnabled(false)
		},
	)

	// Disable the other active components in this menu
	g.startButton.SetEnabled(false)
	g.saveButton.SetEnabled(false)
	g.savesList.SetEnabled(false)
	g.cancelButton.SetEnabled(false)
	g.settingsButton.SetEnabled(false)
	g.exitButton.SetEnabled(false)

	g.popup.SetUserData(MainMenu)
	g.Scene.Add(g.popup)
	g.popup.Open()
}

func (g *Gui) openConfirmDeletePopup() {
	if g.savesList == nil {
		return
	}

	width, height := g.App.GetSize()
	w := float32(150) // Min width, if the text is too long it will expand
	h := float32(100)
	x := (float32(width) - w) / 2
	y := (float32(height) - h) / 2

	prompt := fmt.Sprintf("Delete %s?\nFile will be lost forever!", file.Name(g.savesList.Deleted))
	g.popup = component.NewPopup(
		prompt,        // Prompt label
		"Delete",      // Submit button text
		x,             // x position on screen
		y,             // y position on screen
		w,             // popup width
		h,             // popup height,
		int(MainMenu), // parent menu

		// on submit...
		func() {
			if err := os.Remove(g.savesList.Deleted); err != nil {
				fmt.Printf("Can't delete save file [%s]: %s\n", g.savesList.Deleted, err)
			}

			g.savesList.Deleted = ""
			g.closePopup()

			// Refresh the files
			g.closeBrowseDialog()
			g.openBrowseDialog()
		},

		// on cancel
		func() {
			g.savesList.Deleted = ""
			g.closePopup()
			g.browseButton.SetEnabled(false)
		},
	)

	// Disable the other active components in this menu
	g.startButton.SetEnabled(false)
	g.saveButton.SetEnabled(false)
	g.savesList.SetEnabled(false)
	g.cancelButton.SetEnabled(false)
	g.settingsButton.SetEnabled(false)
	g.exitButton.SetEnabled(false)

	g.popup.SetUserData(MainMenu)
	g.Scene.Add(g.popup)
	g.popup.Open()
}

func (g *Gui) closePopup() {
	if g.popup == nil {
		return
	}

	g.startButton.SetEnabled(true)
	g.saveButton.SetEnabled(true)
	g.browseButton.SetEnabled(true)
	g.savesList.SetEnabled(true)
	g.cancelButton.SetEnabled(true)
	g.settingsButton.SetEnabled(true)
	g.exitButton.SetEnabled(true)

	g.popup.SetEnabled(false)
	g.Scene.Remove(g.popup)
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

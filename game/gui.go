package game

import (
	"fmt"
	"os/exec"

	"github.com/g3n/engine/gui"
)

// TODO: Find a more efficient way to implement the gui

type View = string

const MainMenu View = "main menu"
const InfoView View = "info view"
const TileContextMenu View = "tile context menu"

type ViewControls struct {
	Open  func(bool)
	Close func()
}

// Create a map of view types to their open and close functions
var Views map[View]ViewControls

// Main menu components
var startButton *gui.Button
var exitButton *gui.Button

// Info view components
var infoLabel *gui.Label

// Tile context menu components
var tileInfoLabel *gui.Label

// Load the info text
func infoText() (txt string) {
	txt = "Strands\n"

	// Retrieve version number
	ver, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err == nil {
		txt += fmt.Sprintf("Version %s", ver)
	}

	txt += "\n"
	txt += "Controls:\n"
	txt += "WASD to move\n"
	txt += "ESC to open menu\n"

	return
}

// Load the gui
func LoadGui() {
	gui.Manager().Set(Scene)
	Views = make(map[View]ViewControls)

	// Register the views with their controls
	registerMainMenu()
	registerInfoView()
	registerTileContextMenu()
}

// Refresh the gui components
func RefreshGui() {
	for _, child := range Scene.Children() {
		viewType, ok := child.GetNode().UserData().(View)

		// If this is a gui component, refresh its corresponding menu by closing and reopening it
		if ok {
			Views[viewType].Close()
			Views[viewType].Open(false)
		}
	}
}

// Register the main menu
func registerMainMenu() {
	Views[MainMenu] = ViewControls{
		Open: func(closeOthers bool) {
			if closeOthers {
				closeViews()
			}

			var width, height int = Application.GetSize()
			var w, h float32

			startButton = gui.NewButton("Enter Simulation")
			w, h = startButton.ContentWidth(), startButton.ContentHeight()
			startButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
			startButton.SetUserData(MainMenu)
			startButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				SetPaused(false)
				Views[InfoView].Open(true)
			})
			Scene.Add(startButton)

			exitButton = gui.NewButton("Exit")
			w, h = exitButton.ContentWidth(), exitButton.ContentHeight()
			exitButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*1.5)
			exitButton.SetUserData(MainMenu)
			exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				Application.Exit()
			})
			Scene.Add(exitButton)

			SetPaused(true)
		},

		Close: func() {
			Scene.Remove(startButton)
			Scene.Remove(exitButton)
		},
	}
}

// Register the info view
func registerInfoView() {
	Views[InfoView] = ViewControls{
		Open: func(closeOthers bool) {
			if closeOthers {
				closeViews()
			}

			infoLabel = gui.NewLabel(infoText())
			infoLabel.SetPosition(0, 0)
			infoLabel.SetUserData(InfoView)
			Scene.Add(infoLabel)

			SetPaused(false)
		},

		Close: func() {
			Scene.Remove(infoLabel)
		},
	}
}

// Register the tile context menu
func registerTileContextMenu() {
	Views[TileContextMenu] = ViewControls{
		Open: func(closeOthers bool) {
			if closeOthers {
				closeViews()
			}

			var width, _ int = Application.GetSize()
			var w float32

			tileInfoLabel = gui.NewLabel("<selected tile>")
			w, _ = startButton.ContentWidth(), startButton.ContentHeight()
			tileInfoLabel.SetPosition(float32(width)-w, 0)
			tileInfoLabel.SetUserData(TileContextMenu)
			Scene.Add(tileInfoLabel)
		},

		Close: func() {
			Scene.Remove(tileInfoLabel)
		},
	}
}

// Close all gui views
func closeViews() {
	for _, view := range Views {
		view.Close()
	}
}

package game

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/g3n/engine/gui"
)

type View = string

const MainMenu View = "main menu"
const SimulationView View = "simulation view"
const TileContextMenu View = "tile context menu"

type ViewControls struct {
	Open    func(bool)
	Close   func()
	Refresh func()
}

// Create a map of view types to their open and close functions
var Views map[View]ViewControls

// Main menu components
var startButton *gui.Button
var exitButton *gui.Button

// Simulation view components
var simCursor *gui.Image
var infoLabel *gui.Label
var pausedLabel *gui.Label

// Tile context menu components
var tileInfoLabel *gui.Label

// Load the info text
func infoText() (txt string) {
	txt = "Strands\n"

	// Retrieve version number
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err == nil {
		ver := strings.TrimSpace(string(out))
		txt += fmt.Sprintf("Version %s\n", ver)
	}

	txt += "\n"
	txt += "Controls:\n"
	txt += "WASD to move\n"
	txt += "ESC to open menu\n"
	txt += "Space to toggle simulation\n"
	txt += "\n"
	txt += "Left click a tile to add water\n"
	txt += "Right click a tile to try to add a plant\n"

	// Append the WAILA (what am I looking at?) data
	if LookingAt != nil {
		txt += "\n"
		txt += LookingAt.InfoString()
	}

	txt += fmt.Sprintf("\nTotal water volume=%s\n", TotalWaterVolume)

	return
}

// Update the "Simulation Running/Paused" status
func pausedStatus() string {
	return fmt.Sprintf("Simulation %s", map[bool]string{
		true:  "Paused",
		false: "Running",
	}[IsPaused])
}

// Load the gui
func LoadGui() {
	gui.Manager().Set(Scene)
	Views = make(map[View]ViewControls)

	// Register the views with their controls
	registerMainMenu()
	registerSimulationView()
	registerTileContextMenu()
}

// Reload the gui components
func ReloadGui() {
	// Iterate over active components
	for _, child := range Scene.Children() {
		// If this is a gui component, reload it by closing and reopening
		if viewType, ok := child.GetNode().UserData().(View); ok {
			Views[viewType].Close()
			Views[viewType].Open(false)
		}
	}
}

// Refresh the gui components
func RefreshGui() {
	for _, child := range Scene.Children() {
		// If this is a gui component, call its refresh method
		if viewType, ok := child.GetNode().UserData().(View); ok {
			Views[viewType].Refresh()
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
				Views[SimulationView].Open(true)
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

		Refresh: func() {},
	}
}

// Register the simulation view
func registerSimulationView() {
	Views[SimulationView] = ViewControls{
		Open: func(closeOthers bool) {
			if closeOthers {
				closeViews()
			}

			var width, height int = Application.GetSize()
			var w, h float32

			simCursor, err := gui.NewImage("res/cursor.png")
			if err == nil {
				w, h = simCursor.ContentWidth(), simCursor.ContentHeight()
				simCursor.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
				simCursor.SetUserData(MainMenu)
				Scene.Add(simCursor)
			}

			infoLabel = gui.NewLabel(infoText())
			infoLabel.SetPosition(0, 0)
			infoLabel.SetUserData(SimulationView)
			Scene.Add(infoLabel)

			pausedLabel = gui.NewLabel(pausedStatus())
			pausedLabel.SetPosition((float32(width)-pausedLabel.ContentWidth())/2, 0)
			pausedLabel.SetUserData(SimulationView)
			Scene.Add(pausedLabel)

			SetPaused(false)
		},

		Close: func() {
			Scene.Remove(simCursor)
			Scene.Remove(infoLabel)
			Scene.Remove(pausedLabel)
		},

		Refresh: func() {
			width, _ := Application.GetSize()
			infoLabel.SetText(infoText())
			pausedLabel.SetText(pausedStatus())
			pausedLabel.SetPosition((float32(width)-pausedLabel.ContentWidth())/2, 0)
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
			w, _ = tileInfoLabel.ContentWidth(), tileInfoLabel.ContentHeight()
			tileInfoLabel.SetPosition(float32(width)-w, 0)
			tileInfoLabel.SetUserData(TileContextMenu)
			Scene.Add(tileInfoLabel)

			SetPaused(true)
		},

		Close: func() {
			Scene.Remove(tileInfoLabel)
		},

		Refresh: func() {},
	}
}

// Close all gui views
func closeViews() {
	for _, view := range Views {
		view.Close()
	}
}

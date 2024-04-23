package gui

import (
	"fmt"

	"github.com/g3n/engine/gui"

	"cbeimers113/strands/internal/context"
)

type View = string

const (
	MainMenu        View = "main menu"
	SimulationView  View = "simulation view"
	TileContextMenu View = "tile context menu"
)

type viewControls struct {
	Open    func(bool)
	Close   func()
	Refresh func()
}

// Map view types to their open and close functions
var views map[View]viewControls

func init() {
	views = make(map[View]viewControls)
}

func Open(view View, closeOthers bool) {
	views[view].Open(closeOthers)
}

// GUI manager
type Gui struct {
	*context.Context

	// Main menu components
	startButton *gui.Button
	exitButton  *gui.Button

	// Simulation view components
	simCursor   *gui.Image
	infoLabel   *gui.Label
	pausedLabel *gui.Label

	// Tile context menu components
	tileInfoLabel *gui.Label
}

// New creates a new Gui
func New(ctx *context.Context) *Gui {
	g := &Gui{Context: ctx}
	gui.Manager().Set(g.Scene)

	// Register the views with their controls
	g.registerMainMenu()
	g.registerSimulationView()
	g.registerTileContextMenu()

	return g
}

// Reload the gui components
func (g *Gui) Reload() {
	// Iterate over active components
	for _, child := range g.Scene.Children() {
		// If this is a gui component, reload it by closing and reopening
		if viewType, ok := child.GetNode().UserData().(View); ok {
			views[viewType].Close()
			views[viewType].Open(false)
		}
	}
}

// Refresh the gui components
func (g *Gui) Refresh() {
	for _, child := range g.Scene.Children() {
		// If this is a gui component, call its refresh method
		if viewType, ok := child.GetNode().UserData().(View); ok {
			views[viewType].Refresh()
		}
	}
}

// Register the main menu
func (g *Gui) registerMainMenu() {
	views[MainMenu] = viewControls{
		Open: func(closeOthers bool) {
			if closeOthers {
				g.closeViews()
			}

			var width, height int = g.App.GetSize()
			var w, h float32

			g.startButton = gui.NewButton("Enter Simulation")
			w, h = g.startButton.ContentWidth(), g.startButton.ContentHeight()
			g.startButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
			g.startButton.SetUserData(MainMenu)
			g.startButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.State.SetInMenu(false)
				views[SimulationView].Open(true)
			})
			g.Scene.Add(g.startButton)

			g.exitButton = gui.NewButton("Exit")
			w, h = g.exitButton.ContentWidth(), g.exitButton.ContentHeight()
			g.exitButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*1.5)
			g.exitButton.SetUserData(MainMenu)
			g.exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.App.Exit()
			})
			g.Scene.Add(g.exitButton)

			g.State.SetInMenu(true)
		},

		Close: func() {
			g.Scene.Remove(g.startButton)
			g.Scene.Remove(g.exitButton)
		},

		Refresh: func() {},
	}
}

// Register the simulation view
func (g *Gui) registerSimulationView() {
	views[SimulationView] = viewControls{
		Open: func(closeOthers bool) {
			if closeOthers {
				g.closeViews()
			}

			var width, height int = g.App.GetSize()
			var w, h float32

			simCursor, err := gui.NewImage("res/cursor.png")
			if err == nil {
				w, h = simCursor.ContentWidth(), simCursor.ContentHeight()
				simCursor.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
				simCursor.SetUserData(MainMenu)
				g.Scene.Add(simCursor)
			}

			g.infoLabel = gui.NewLabel(g.infoText())
			g.infoLabel.SetPosition(0, 0)
			g.infoLabel.SetUserData(SimulationView)
			g.Scene.Add(g.infoLabel)

			g.pausedLabel = gui.NewLabel(g.pausedStatus())
			g.pausedLabel.SetPosition((float32(width)-g.pausedLabel.ContentWidth())/2, 0)
			g.pausedLabel.SetUserData(SimulationView)
			g.Scene.Add(g.pausedLabel)

			g.State.SetInMenu(false)
		},

		Close: func() {
			g.Scene.Remove(g.simCursor)
			g.Scene.Remove(g.infoLabel)
			g.Scene.Remove(g.pausedLabel)
		},

		Refresh: func() {
			width, _ := g.App.GetSize()
			g.infoLabel.SetText(g.infoText())
			g.pausedLabel.SetText(g.pausedStatus())
			g.pausedLabel.SetPosition((float32(width)-g.pausedLabel.ContentWidth())/2, 0)
		},
	}
}

// Register the tile context menu
func (g *Gui) registerTileContextMenu() {
	views[TileContextMenu] = viewControls{
		Open: func(closeOthers bool) {
			if closeOthers {
				g.closeViews()
			}

			var width, _ int = g.App.GetSize()
			var w float32

			g.tileInfoLabel = gui.NewLabel("<selected tile>")
			w, _ = g.tileInfoLabel.ContentWidth(), g.tileInfoLabel.ContentHeight()
			g.tileInfoLabel.SetPosition(float32(width)-w, 0)
			g.tileInfoLabel.SetUserData(TileContextMenu)
			g.Scene.Add(g.tileInfoLabel)

			g.State.SetInMenu(true)
		},

		Close: func() {
			g.Scene.Remove(g.tileInfoLabel)
		},

		Refresh: func() {},
	}
}

// Close all gui views
func (g *Gui) closeViews() {
	for _, view := range views {
		view.Close()
	}
}

// Load the info text
func (g *Gui) infoText() string {
	txt := fmt.Sprintf("Version %s\n", g.Cfg.Version)
	txt += "\n"
	txt += "Controls:\n"
	txt += "WASD to move\n"
	txt += "ESC to open menu\n"
	txt += "Space to toggle simulation\n"
	txt += "\n"
	txt += "Left click a tile to add water\n"
	txt += "Right click a tile to try to add a plant\n"

	// Append info about simulation
	txt += "\nChemical Quantities:\n"
	for name, amnt := range g.State.Quantities {
		txt += fmt.Sprintf("%s: %s\n", name, amnt.String())
	}

	// Append the WAILA (what am I looking at?) data
	if g.State.LookingAt != nil {
		txt += "\nLooking At:\n"
		txt += g.State.LookingAt.InfoString()
	}

	return txt
}

// Update the "Simulation Running/Paused" status
func (g *Gui) pausedStatus() string {
	return fmt.Sprintf("Simulation %s", map[bool]string{
		true:  "Paused",
		false: "Running",
	}[g.State.Paused()])
}

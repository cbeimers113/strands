package gui

import (
	"github.com/g3n/engine/gui"

	"cbeimers113/strands/internal/context"
)

type View int

const (
	MainMenu View = iota
	ConfigMenu
	SimulationView
	TileContextMenu

	numMenus
)

type viewControls struct {
	open    func(bool)
	close   func()
	refresh func()
}

// Map view types to their open and close functions
var views [numMenus]viewControls

// GUI manager
type Gui struct {
	*context.Context

	// Main menu components
	startButton    *gui.Button
	settingsButton *gui.Button
	exitButton     *gui.Button

	// Config menu components
	showControlsCheck *gui.CheckRadio
	mouseXSenSlider   *gui.Slider
	mouseYSenSlider   *gui.Slider
	moveSpeedSlider   *gui.Slider
	tickSpeedSlider   *gui.Slider
	dayLengthSlider   *gui.Slider
	saveButton        *gui.Button

	// Simulation view components
	simCursor   *gui.Image
	infoLabel   *gui.Label
	pausedLabel *gui.Label

	// Tile context menu components
	tileInfoLabel   *gui.Label
	plantSeedButton *gui.Button
}

// New creates a new Gui
func New(ctx *context.Context) *Gui {
	g := &Gui{Context: ctx}
	gui.Manager().Set(g.Scene)

	// Register the views with their controls
	g.registerMainMenu()
	g.registerConfigMenu()
	g.registerSimulationView()
	g.registerTileContextMenu()

	return g
}

// Open a view and optionally close other views
func Open(view View, closeOthers bool) {
	views[view].open(closeOthers)
}

// Reload the gui components
func (g *Gui) Reload() {
	// Iterate over active components
	for _, child := range g.Scene.Children() {
		// If this is a gui component, reload it by closing and reopening
		if viewType, ok := child.GetNode().UserData().(View); ok {
			views[viewType].close()
			views[viewType].open(false)
		}
	}
}

// Refresh the gui components
func (g *Gui) Refresh() {
	for _, child := range g.Scene.Children() {
		// If this is a gui component, call its refresh method
		if viewType, ok := child.GetNode().UserData().(View); ok {
			views[viewType].refresh()
		}
	}
}

// Close all gui views
func (g *Gui) closeViews() {
	for _, view := range views {
		view.close()
	}
}

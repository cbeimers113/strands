package gui

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/text"

	"cbeimers113/strands/internal/context"
	"cbeimers113/strands/internal/gui/component"
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

var (
	// Map view types to their open and close functions
	views [numMenus]viewControls

	//go:embed font/AgaveNerdFont-Regular.ttf
	fontData []byte

	//go:embed assets/icon256.png
	iconData []byte
)

// GUI manager
type Gui struct {
	*context.Context

	// Reusable popup
	popup *component.Popup

	// Main menu components
	menuLogo       *gui.Image
	startButton    *gui.Button
	newButton      *gui.Button
	settingsButton *gui.Button
	saveDialog     *component.Dialog
	browseButton   *gui.Button
	savesList      *component.FileList
	cancelButton   *gui.Button
	exitButton     *gui.Button

	// Config menu components
	showControlsCheck *gui.CheckRadio
	exitSaveCheck     *gui.CheckRadio
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

	// GUI flags
	gameStarted bool
}

// New creates a new Gui
func New(ctx *context.Context) *Gui {
	g := &Gui{Context: ctx}

	gui.SetStyleDefault(getStyle())
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

// Create a popup notification
func (g *Gui) Popup(prompt, submitText string, submitFunc func()) {
	width, height := g.App.GetSize()
	paused := g.State.Paused()
	inMenu := g.State.InMenu()

	g.popup = component.NewPopup(prompt, submitText, 0, 0, 150, 100, -1, submitFunc, func() {
		g.State.SetInMenu(inMenu)
		g.State.SetPaused(paused)
		g.popup.SetEnabled(false)
		g.Scene.Remove(g.popup)
		g.Reload()
	})

	g.popup.SetPosition((float32(width)-g.popup.Width())/2, (float32(height)-g.popup.Height())/2)
	g.popup.Open(float32(width), float32(height))
	g.Scene.Add(g.popup)
	g.State.SetInMenu(true)
	g.State.SetPaused(true)
}

// Reload the gui components
func (g *Gui) Reload() {
	// Iterate over active components
	for _, child := range g.Scene.Children() {
		if child == nil {
			continue
		}

		// If this is a gui component, reload it by closing and reopening
		if node := child.GetNode(); node != nil {
			if viewType, ok := node.UserData().(View); ok {
				views[viewType].close()
				views[viewType].open(false)
			}
		}
	}
}

// Refresh the gui components
func (g *Gui) Refresh() {
	for _, child := range g.Scene.Children() {
		// If this is a gui component, call its refresh method
		if child != nil {
			if node := child.GetNode(); node != nil {
				if data := node.UserData(); data != nil {
					if viewType, ok := data.(View); ok {
						views[viewType].refresh()
					}
				}
			}
		}
	}
}

// Close all gui views
func (g *Gui) closeViews() {
	for _, view := range views {
		view.close()
	}
}

// Set the application icon
func (g *Gui) SetIcon() {
	// Decode the embedded PNG data
	img, err := png.Decode(bytes.NewReader(iconData))
	if err != nil {
		panic("Failed to decode embedded icon data: " + err.Error())
	}

	g.Context.Win.SetIcon([]image.Image{img})
}

// Get the style of the gui
func getStyle() *gui.Style {
	var err error

	g := gui.NewDarkStyle()

	gb := g.Button
	g.Button = gui.ButtonStyles{
		Normal:   gb.Normal,
		Over:     gb.Disabled,
		Focus:    gb.Focus,
		Pressed:  gb.Over,
		Disabled: gb.Pressed,
	}

	if g.Font, err = text.NewFontFromData(fontData); err != nil {
		fmt.Printf("Couldn't load GUI font: %s\n", err)
	}

	return g
}

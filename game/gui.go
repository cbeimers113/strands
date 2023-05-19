package game

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/g3n/engine/gui"
)

// Load the info text
var infoText = func() string {
	txt := "Strands\n"

	// Retrieve version number
	out, err := exec.Command("git", "describe", "--tags", "--abbrev=0").Output()
	if err == nil {
		semver := strings.Split(strings.TrimSpace(string(out)), ".")

		if len(semver) >= 3 {
			major, _ := strconv.Atoi(semver[0])
			minor, _ := strconv.Atoi(semver[1])
			patch, _ := strconv.Atoi(semver[2])
			txt += fmt.Sprintf("Version %d.%d.%d\n", major, minor, patch)
		}
	}

	txt += "\n"
	txt += "Controls:\n"
	txt += "WASD to move\n"
	txt += "ESC to open menu\n"

	return txt
}()

type ViewType = string

const MainMenuView ViewType = "main menu"
const SimulationView ViewType = "simulation view"

// Create a map of view types to their open and close functions
var Views map[ViewType][2]func(bool)

// Main menu components
var startButton *gui.Button
var exitButton *gui.Button

// Simulation view components
var infoLabel *gui.Label

// Load the gui
func LoadGui() {
	gui.Manager().Set(Scene)
	Views = make(map[ViewType][2]func(bool))
	Views[MainMenuView] = [2]func(bool){OpenMainMenu, closeMainMenu}
	Views[SimulationView] = [2]func(bool){OpenSimulationView, closeSimulationView}
}

// Refresh the gui components
func RefreshGui() {
	for _, child := range Scene.Children() {
		viewType, ok := child.GetNode().UserData().(ViewType)

		// If this is a gui component, refresh its corresponding menu by closing and reopening it
		if ok {
			Views[viewType][1](false)
			Views[viewType][0](false)
		}
	}
}

// Open the main menu, optionally close other open views as well
func OpenMainMenu(closeOthers bool) {
	if closeOthers {
		closeMenus()
	}

	var width, height int = Application.GetSize()
	var w, h float32

	startButton = gui.NewButton("Enter Simulation")
	w, h = startButton.ContentWidth(), startButton.ContentHeight()
	startButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
	startButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		SetPaused(false)
		OpenSimulationView(true)
	})
	Scene.Add(startButton)
	startButton.SetUserData(MainMenuView)

	exitButton = gui.NewButton("Exit")
	w, h = exitButton.ContentWidth(), exitButton.ContentHeight()
	exitButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*1.5)
	exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		Application.Exit()
	})
	Scene.Add(exitButton)
	startButton.SetUserData(MainMenuView)

	SetPaused(true)
}

// Open the simulation view, optionally close other views as well
func OpenSimulationView(closeOthers bool) {
	if closeOthers {
		closeMenus()
	}

	infoLabel = gui.NewLabel(infoText)
	infoLabel.SetPosition(0, 0)
	Scene.Add(infoLabel)
	startButton.SetUserData(SimulationView)

	SetPaused(false)
}

// Close the main menu
func closeMainMenu(_ bool) {
	Scene.Remove(startButton)
	Scene.Remove(exitButton)
}

// Close the simulation view
func closeSimulationView(_ bool) {
	Scene.Remove(infoLabel)
}

// Close all gui views
func closeMenus() {
	closeMainMenu(false)
	closeSimulationView(false)
}

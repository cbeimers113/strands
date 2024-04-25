package gui

import (
	"github.com/g3n/engine/gui"
)

// Register the main menu
func (g *Gui) registerMainMenu() {
	views[MainMenu] = viewControls{
		open: func(closeOthers bool) {
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
				views[SimulationView].open(true)
			})
			g.Scene.Add(g.startButton)

			g.settingsButton = gui.NewButton("Settings")
			w, h = g.settingsButton.ContentWidth(), g.settingsButton.ContentHeight()
			g.settingsButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*1.5)
			g.settingsButton.SetUserData(MainMenu)
			g.settingsButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				Open(ConfigMenu, true)
			})
			g.Scene.Add(g.settingsButton)

			g.exitButton = gui.NewButton("Exit")
			w, h = g.exitButton.ContentWidth(), g.exitButton.ContentHeight()
			g.exitButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*3)
			g.exitButton.SetUserData(MainMenu)
			g.exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.App.Exit()
			})
			g.Scene.Add(g.exitButton)

			g.State.SetInMenu(true)
		},

		close: func() {
			g.Scene.Remove(g.startButton)
			g.Scene.Remove(g.settingsButton)
			g.Scene.Remove(g.exitButton)
		},

		refresh: func() {},
	}
}

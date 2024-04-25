package gui

import "github.com/g3n/engine/gui"

// Register the simulation view
func (g *Gui) registerSimulationView() {
	views[SimulationView] = viewControls{
		open: func(closeOthers bool) {
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

		close: func() {
			g.Scene.Remove(g.simCursor)
			g.Scene.Remove(g.infoLabel)
			g.Scene.Remove(g.pausedLabel)
			g.State.ChangeMenuCooldown(1000)
		},

		refresh: func() {
			width, _ := g.App.GetSize()
			g.infoLabel.SetText(g.infoText())
			g.pausedLabel.SetText(g.pausedStatus())
			g.pausedLabel.SetPosition((float32(width)-g.pausedLabel.ContentWidth())/2, 0)
		},
	}
}

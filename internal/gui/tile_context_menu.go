package gui

import (
	"cbeimers113/strands/internal/entity"

	"github.com/g3n/engine/gui"
)

// Register the tile context menu
func (g *Gui) registerTileContextMenu() {
	views[TileContextMenu] = viewControls{
		open: func(closeOthers bool) {
			if closeOthers {
				g.closeViews()
			}

			if g.State.LookingAt == nil {
				return
			}

			var width, height int = g.App.GetSize()
			var w, h float32

			g.tileInfoLabel = gui.NewLabel(g.State.LookingAt.InfoString())
			w, h = g.tileInfoLabel.Width()*0.75, g.tileInfoLabel.Height()
			g.tileInfoLabel.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
			g.tileInfoLabel.SetUserData(TileContextMenu)
			g.Scene.Add(g.tileInfoLabel)

			g.plantSeedButton = gui.NewButton("Plant Seed")
			w, h = g.plantSeedButton.ContentWidth(), g.plantSeedButton.ContentHeight()
			g.plantSeedButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*3)
			g.plantSeedButton.SetUserData(TileContextMenu)
			g.plantSeedButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				if tile, ok := g.State.LookingAt.(*entity.Tile); ok {
					tile.AddPlant(g.State.Entities)
					g.tileInfoLabel.SetText(g.State.LookingAt.InfoString())
				}
			})
			g.Scene.Add(g.plantSeedButton)

			g.exitButton = gui.NewButton("Close")
			w, h = g.exitButton.ContentWidth(), g.exitButton.ContentHeight()
			g.exitButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*4.5)
			g.exitButton.SetUserData(TileContextMenu)
			g.exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				Open(SimulationView, true)
			})
			g.Scene.Add(g.exitButton)

			g.State.SetInMenu(true)
		},

		close: func() {
			g.Scene.Remove(g.tileInfoLabel)
			g.Scene.Remove(g.plantSeedButton)
			g.Scene.Remove(g.exitButton)
		},

		refresh: func() {},
	}
}

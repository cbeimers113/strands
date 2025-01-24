package gui

import (
	"fmt"

	"github.com/g3n/engine/gui"

	"cbeimers113/strands/internal/entity"
	"cbeimers113/strands/internal/gui/color"
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

			var (
				width, height int = g.App.GetSize()
				w, h          float32
				nextY         float32
			)

			g.tileInfoLabel = gui.NewLabel(g.State.LookingAt.InfoString())
			w, h = g.tileInfoLabel.Width(), g.tileInfoLabel.Height()
			g.tileInfoLabel.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
			g.tileInfoLabel.SetUserData(TileContextMenu)
			g.tileInfoLabel.SetColor(color.Black)
			g.tileInfoLabel.SetBgColor4(color.Opaque)
			g.tileInfoLabel.SetPaddings(5, 5, 5, 5)
			g.Scene.Add(g.tileInfoLabel)
			nextY = g.tileInfoLabel.Position().Y + g.tileInfoLabel.Height() + 5

			g.plantSeedButton = gui.NewButton("Plant Seed")
			w = g.plantSeedButton.Width()
			g.plantSeedButton.SetPosition((float32(width)-w)/2, nextY)
			g.plantSeedButton.SetUserData(TileContextMenu)
			g.plantSeedButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				if tile, ok := g.State.LookingAt.(*entity.Tile); ok {
					planted := tile.AddPlant(g.State.Entities, g.Scene)
					g.tileInfoLabel.SetText(g.State.LookingAt.InfoString())

					if planted {
						g.Notifications.Push(fmt.Sprintf("Seed planted at (%d, %d)", tile.MapX, tile.MapZ))
					} else {
						g.Notifications.Push(fmt.Sprintf("Couldn't plant seed on %s tile at (%d, %d)", tile.Type.Name, tile.MapX, tile.MapZ))
					}
				}
			})
			g.Scene.Add(g.plantSeedButton)
			nextY = g.plantSeedButton.Position().Y + g.plantSeedButton.Height() + 5

			g.exitButton = gui.NewButton("Close")
			w = g.exitButton.Width()
			g.exitButton.SetPosition((float32(width)-w)/2, nextY)
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

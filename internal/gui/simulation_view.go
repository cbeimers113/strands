package gui

import (
	"fmt"
	"strings"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/texture"

	"cbeimers113/strands/internal/graphics"
	"cbeimers113/strands/internal/gui/color"
)

// Register the simulation view
func (g *Gui) registerSimulationView() {
	views[SimulationView] = viewControls{
		open: func(closeOthers bool) {
			if closeOthers {
				g.closeViews()
			}

			g.gameStarted = true

			var (
				width, height int = g.App.GetSize()
				w, h          float32
				cursorTex     *texture.Texture2D
				err           error
			)

			cursorTex, err = graphics.Texture(graphics.TexCursor)
			if err != nil {
				fmt.Println(err)
			} else {
				g.simCursor = gui.NewImageFromTex(cursorTex)
				w, h = g.simCursor.Width(), g.simCursor.Height()
				g.simCursor.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
				g.simCursor.SetUserData(SimulationView)
				g.Scene.Add(g.simCursor)
			}

			g.topPanel = gui.NewPanel(float32(width), 20)
			g.topPanel.SetColor4(color.Translucent)
			g.topPanel.SetUserData(SimulationView)
			g.Scene.Add(g.topPanel)

			g.moveIcon = gui.NewImageFromTex(graphics.Textures[graphics.TexWalk])
			g.moveIcon.SetUserData(SimulationView)
			g.Scene.Add(g.moveIcon)

			g.playerLabel = gui.NewLabel(g.playerPos())
			g.playerLabel.SetColor(color.Black)
			g.playerLabel.SetUserData(SimulationView)
			g.Scene.Add(g.playerLabel)

			g.clockLabel = gui.NewLabel(g.getClock())
			g.clockLabel.SetColor(color.Black)
			g.clockLabel.SetUserData(SimulationView)
			g.Scene.Add(g.clockLabel)

			g.wailaLabel = gui.NewLabel(g.getWaila())
			g.wailaLabel.SetColor(color.Black)
			g.wailaLabel.SetUserData(SimulationView)
			g.Scene.Add(g.wailaLabel)

			g.helpLabel = gui.NewLabel(helpText())
			g.helpLabel.SetColor(color.Black)
			g.helpLabel.SetBgColor4(color.Translucent)
			g.helpLabel.SetPosition(5, 25)
			g.helpLabel.SetUserData(SimulationView)
			g.helpLabel.SetPaddings(5, 5, 5, 5)
			g.Scene.Add(g.helpLabel)

			g.quantitiesLabel = gui.NewLabel(g.getQuantities())
			g.quantitiesLabel.SetColor(color.Black)
			g.quantitiesLabel.SetBgColor4(color.Translucent)
			g.quantitiesLabel.SetPosition(5, 25)
			g.quantitiesLabel.SetUserData(SimulationView)
			g.quantitiesLabel.SetPaddings(5, 5, 5, 5)
			g.Scene.Add(g.quantitiesLabel)

			g.State.SetInMenu(false)
		},

		close: func() {
			g.Scene.Remove(g.simCursor)
			g.Scene.Remove(g.topPanel)
			g.Scene.Remove(g.moveIcon)
			g.Scene.Remove(g.playerLabel)
			g.Scene.Remove(g.clockLabel)
			g.Scene.Remove(g.wailaLabel)
			g.Scene.Remove(g.helpLabel)
			g.Scene.Remove(g.quantitiesLabel)
		},

		refresh: func() {
			width, _ := g.App.GetSize()

			g.moveIcon.SetPosition(0, 0)
			g.playerLabel.SetText(g.playerPos())
			g.playerLabel.SetPosition(g.moveIcon.Width(), 1)

			g.clockLabel.SetText(g.getClock())
			g.clockLabel.SetPosition((float32(width)-g.clockLabel.Width())/2, 1)

			g.wailaLabel.SetText(g.getWaila())
			g.wailaLabel.SetPosition(float32(width)-g.wailaLabel.Width()-4, 1)

			g.helpLabel.SetVisible(g.Cfg.ShowHelp)

			g.quantitiesLabel.SetVisible(g.State.ShowChems())
			g.quantitiesLabel.SetText(g.getQuantities())
			g.quantitiesLabel.SetPosition(float32(width)-g.quantitiesLabel.Width()-5, 25)

			if g.State.FastMovement() {
				g.moveIcon.SetTexture(graphics.Textures[graphics.TexRun])
			} else {
				g.moveIcon.SetTexture(graphics.Textures[graphics.TexWalk])
			}
		},
	}
}

// Get the player's position
func (g *Gui) playerPos() string {
	p := g.Cam.Position()
	return fmt.Sprintf("(%d, %d, %d)", int32(p.X), int32(p.Y), int32(p.Z))
}

// Get the in-game clock
func (g *Gui) getClock() string {
	playPause := ""
	if g.State.Paused() {
		playPause = ""
	}

	return fmt.Sprintf("  %s | %d t/s | %s ", g.State.Clock, g.State.TPS(), playPause)
}

// Get the WAILA info
func (g *Gui) getWaila() string {
	if g.State.LookingAt != nil {
		return g.State.LookingAt.InfoString()
	}

	return ""
}

// Get the help text
func helpText() string {
	txt := "Controls:\n"
	txt += " WASD to move\n"
	txt += " Hold shift to move faster\n"
	txt += " Caps lock to toggle fast movement\n"
	txt += " Space and CTRL to go up and down\n"
	txt += " ESC to open menu\n"
	txt += " Tab to play/pause simulation\n"
	txt += " Q to toggle chemical quantities panel\n"
	txt += " Left click a tile to add 1000 L of water\n"
	txt += " Right click a tile to open the tile menu\n"
	txt += "\nThis message can be toggled in the settings menu!"

	return txt
}

// Get chemical quantities
func (g *Gui) getQuantities() string {
	txt := "Chemical Levels:\n"
	for name, amnt := range g.State.Quantities {
		txt += fmt.Sprintf("%s: %s\n", name, amnt.String())
	}

	return strings.TrimSpace(txt)
}

package gui

import (
	"fmt"

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

			g.infoLabel = gui.NewLabel(g.infoText())
			g.infoLabel.SetColor(color.Black)
			g.infoLabel.SetBgColor4(color.Opaque)
			g.infoLabel.SetPosition(5, 5)
			g.infoLabel.SetUserData(SimulationView)
			g.infoLabel.SetPaddings(5, 5, 5, 5)
			g.Scene.Add(g.infoLabel)

			g.pausedLabel = gui.NewLabel(g.pausedStatus())
			g.pausedLabel.SetColor(color.Black)
			g.pausedLabel.SetBgColor4(color.Opaque)
			g.pausedLabel.SetPosition((float32(width)-g.pausedLabel.Width())/2, 5)
			g.pausedLabel.SetUserData(SimulationView)
			g.pausedLabel.SetPaddings(5, 5, 5, 5)
			g.Scene.Add(g.pausedLabel)

			g.State.SetInMenu(false)
		},

		close: func() {
			g.Scene.Remove(g.simCursor)
			g.Scene.Remove(g.infoLabel)
			g.Scene.Remove(g.pausedLabel)
		},

		refresh: func() {
			width, _ := g.App.GetSize()
			g.infoLabel.SetText(g.infoText())
			g.pausedLabel.SetText(g.pausedStatus())
			g.pausedLabel.SetPosition((float32(width)-g.pausedLabel.Width())/2, 5)
		},
	}
}

// Load the info text
func (g *Gui) infoText() string {
	txt := fmt.Sprintf("Version %s\n", g.Version)
	txt += fmt.Sprintf("TPS: %d\n", g.State.TPS())
	txt += fmt.Sprintf("%s\n", g.State.Clock)
	txt += "\n"

	if g.Cfg.ShowHelp {
		txt += "Controls:\n"
		txt += "WASD to move\n"
		txt += "Hold shift to move faster\n"
		txt += "Caps lock to toggle fast movement\n"
		txt += "Space and CTRL to go up and down\n"
		txt += "ESC to open menu\n"
		txt += "Tab to play/pause simulation\n"
		txt += "Left click a tile to add 10 L of water\n"
		txt += "Right click a tile to open the tile menu\n"
	}

	// Append info about simulation
	txt += "\nChemical Levels:\n"
	for name, amnt := range g.State.Quantities {
		txt += fmt.Sprintf("%s: %s\n", name, amnt.String())
	}

	// Append player info
	p := g.Cam.Position()
	if g.State.FastMovement() {
		txt += "\nFast Move"
	}

	txt += fmt.Sprintf("\n(%d, %d, %d)", int32(p.X), int32(p.Y), int32(p.Z))

	// Append the WAILA (what am I looking at?) data
	if g.State.LookingAt != nil {
		txt += "\n\nLooking At:\n"
		txt += g.State.LookingAt.InfoString()
	}

	return txt
}

// Update the "Simulation Running/Paused" status
func (g *Gui) pausedStatus() string {
	suffix := "Running"
	if g.State.Paused() {
		suffix = "Paused"
	}

	return fmt.Sprintf("Simulation %s", suffix)
}

package gui

import (
	"fmt"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
)

// Register the config menu
func (g *Gui) registerConfigMenu() {
	views[ConfigMenu] = viewControls{
		open: func(closeOthers bool) {
			if closeOthers {
				g.closeViews()
			}

			var width, height int = g.App.GetSize()
			var w, h float32
			var (
				saveMouseSensX float32 = g.Cfg.Controls.MouseSensitivityX
				saveMouseSensY float32 = g.Cfg.Controls.MouseSensitivityY
				saveTickSpeed  int     = g.Cfg.Simulation.Speed
			)

			g.mouseXSenSlider = gui.NewHSlider(175, 12.5)
			w, h = g.mouseXSenSlider.ContentWidth(), g.mouseXSenSlider.ContentHeight()
			g.mouseXSenSlider.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
			g.mouseXSenSlider.SetUserData(ConfigMenu)
			g.mouseXSenSlider.SetValue(g.Cfg.Controls.MouseSensitivityX)
			g.mouseXSenSlider.SetText(g.mouseSensXLabel())
			g.mouseXSenSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Controls.MouseSensitivityX = g.mouseXSenSlider.Value()
				g.tickSpeedSlider.SetText(g.mouseSensXLabel())
			})
			g.Scene.Add(g.mouseXSenSlider)

			g.mouseYSenSlider = gui.NewHSlider(175, 12.5)
			w, h = g.mouseYSenSlider.ContentWidth(), g.mouseYSenSlider.ContentHeight()
			g.mouseYSenSlider.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*1.5)
			g.mouseYSenSlider.SetUserData(ConfigMenu)
			g.mouseYSenSlider.SetValue(g.Cfg.Controls.MouseSensitivityY)
			g.mouseYSenSlider.SetText(g.mouseSensYLabel())
			g.mouseYSenSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Controls.MouseSensitivityY = g.mouseYSenSlider.Value()
				g.tickSpeedSlider.SetText(g.mouseSensYLabel())
			})
			g.Scene.Add(g.mouseYSenSlider)

			g.tickSpeedSlider = gui.NewHSlider(175, 12.5)
			w, h = g.tickSpeedSlider.ContentWidth(), g.tickSpeedSlider.ContentHeight()
			g.tickSpeedSlider.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*3)
			g.tickSpeedSlider.SetUserData(ConfigMenu)
			g.tickSpeedSlider.SetValue(float32(g.Cfg.Simulation.Speed) / 32)
			g.tickSpeedSlider.SetText(g.tickSpeedLabel())
			g.tickSpeedSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Simulation.Speed = int(math32.Max(1, g.tickSpeedSlider.Value()*32))
				g.tickSpeedSlider.SetText(g.tickSpeedLabel())
			})
			g.Scene.Add(g.tickSpeedSlider)

			g.saveButton = gui.NewButton("Save Settings")
			w, h = g.saveButton.ContentWidth(), g.saveButton.ContentHeight()
			g.saveButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*4.5)
			g.saveButton.SetUserData(ConfigMenu)
			g.saveButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				if err := g.Cfg.Save(); err != nil {
					fmt.Println(err)
				}
				Open(MainMenu, true)
			})
			g.Scene.Add(g.saveButton)

			g.exitButton = gui.NewButton("Close")
			w, h = g.exitButton.ContentWidth(), g.exitButton.ContentHeight()
			g.exitButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*6)
			g.exitButton.SetUserData(TileContextMenu)
			g.exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.Cfg.Controls.MouseSensitivityX = saveMouseSensX
				g.Cfg.Controls.MouseSensitivityY = saveMouseSensY
				g.Cfg.Simulation.Speed = saveTickSpeed
				Open(MainMenu, true)
			})
			g.Scene.Add(g.exitButton)

			g.State.SetInMenu(true)
		},

		close: func() {
			g.Scene.Remove(g.mouseXSenSlider)
			g.Scene.Remove(g.mouseYSenSlider)
			g.Scene.Remove(g.tickSpeedSlider)
			g.Scene.Remove(g.saveButton)
			g.Scene.Remove(g.exitButton)
		},

		refresh: func() {},
	}
}

func (g *Gui) mouseSensXLabel() string {
	return fmt.Sprintf("Mouse X Sensitivity: %.2f", g.Cfg.Controls.MouseSensitivityX)
}

func (g *Gui) mouseSensYLabel() string {
	return fmt.Sprintf("Mouse Y Sensitivity: %.2f", g.Cfg.Controls.MouseSensitivityY)
}

func (g *Gui) tickSpeedLabel() string {
	return fmt.Sprintf("Ticks Per Second: %d", g.Cfg.Simulation.Speed)
}

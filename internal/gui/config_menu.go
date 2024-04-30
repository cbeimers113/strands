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

			g.showControlsCheck = gui.NewCheckBox("Show Controls")
			w, h = g.showControlsCheck.ContentWidth(), g.showControlsCheck.ContentHeight()
			g.showControlsCheck.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
			g.showControlsCheck.SetUserData(ConfigMenu)
			g.showControlsCheck.SetValue(g.Cfg.ShowHelp)
			g.showControlsCheck.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.ShowHelp = g.showControlsCheck.Value()
			})
			g.Scene.Add(g.showControlsCheck)

			g.exitSaveCheck = gui.NewCheckBox("Save on Exit")
			w, h = g.exitSaveCheck.ContentWidth(), g.exitSaveCheck.ContentHeight()
			g.exitSaveCheck.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*1.3)
			g.exitSaveCheck.SetUserData(ConfigMenu)
			g.exitSaveCheck.SetValue(g.Cfg.ExitSave)
			g.exitSaveCheck.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.ExitSave = g.exitSaveCheck.Value()
			})
			g.Scene.Add(g.exitSaveCheck)

			g.mouseXSenSlider = gui.NewHSlider(175, 12.5)
			w, h = g.mouseXSenSlider.ContentWidth(), g.mouseXSenSlider.ContentHeight()
			g.mouseXSenSlider.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*4)
			g.mouseXSenSlider.SetUserData(ConfigMenu)
			g.mouseXSenSlider.SetValue(g.Cfg.Controls.MouseSensitivityX)
			g.mouseXSenSlider.SetText(g.mouseSensXLabel())
			g.mouseXSenSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Controls.MouseSensitivityX = g.mouseXSenSlider.Value()
				g.mouseXSenSlider.SetText(g.mouseSensXLabel())
			})
			g.Scene.Add(g.mouseXSenSlider)

			g.mouseYSenSlider = gui.NewHSlider(175, 12.5)
			w, h = g.mouseYSenSlider.ContentWidth(), g.mouseYSenSlider.ContentHeight()
			g.mouseYSenSlider.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*6)
			g.mouseYSenSlider.SetUserData(ConfigMenu)
			g.mouseYSenSlider.SetValue(g.Cfg.Controls.MouseSensitivityY)
			g.mouseYSenSlider.SetText(g.mouseSensYLabel())
			g.mouseYSenSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Controls.MouseSensitivityY = g.mouseYSenSlider.Value()
				g.mouseYSenSlider.SetText(g.mouseSensYLabel())
			})
			g.Scene.Add(g.mouseYSenSlider)

			g.moveSpeedSlider = gui.NewHSlider(175, 12.5)
			w, h = g.moveSpeedSlider.ContentWidth(), g.moveSpeedSlider.ContentHeight()
			g.moveSpeedSlider.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*8)
			g.moveSpeedSlider.SetUserData(ConfigMenu)
			g.moveSpeedSlider.SetValue(g.Cfg.Controls.MoveSpeed)
			g.moveSpeedSlider.SetText(g.moveSpeedLabel())
			g.moveSpeedSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Controls.MoveSpeed = math32.Max(0.01, g.moveSpeedSlider.Value())
				g.moveSpeedSlider.SetText(g.moveSpeedLabel())
			})
			g.Scene.Add(g.moveSpeedSlider)

			g.tickSpeedSlider = gui.NewHSlider(175, 12.5)
			w, h = g.tickSpeedSlider.ContentWidth(), g.tickSpeedSlider.ContentHeight()
			g.tickSpeedSlider.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*10)
			g.tickSpeedSlider.SetUserData(ConfigMenu)
			g.tickSpeedSlider.SetValue(float32(g.Cfg.Simulation.Speed) / 32)
			g.tickSpeedSlider.SetText(g.tickSpeedLabel())
			g.tickSpeedSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Simulation.Speed = int(math32.Max(1, g.tickSpeedSlider.Value()*32))
				g.tickSpeedSlider.SetText(g.tickSpeedLabel())
			})
			g.Scene.Add(g.tickSpeedSlider)

			g.dayLengthSlider = gui.NewHSlider(175, 12.5)
			w, h = g.dayLengthSlider.ContentWidth(), g.dayLengthSlider.ContentHeight()
			g.dayLengthSlider.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*10)
			g.dayLengthSlider.SetUserData(ConfigMenu)
			g.dayLengthSlider.SetValue(float32(g.Cfg.Simulation.DayLength+1) / 30)
			g.dayLengthSlider.SetText(g.dayLengthLabel())
			g.dayLengthSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Simulation.DayLength = int(math32.Max(1, g.dayLengthSlider.Value()*30))
				g.State.Clock.SetTime(g.State.Clock.Hour, g.State.Clock.Minute)
				g.dayLengthSlider.SetText(g.dayLengthLabel())
			})
			g.Scene.Add(g.dayLengthSlider)

			g.saveButton = gui.NewButton("Save Settings")
			w, h = g.saveButton.ContentWidth(), g.saveButton.ContentHeight()
			g.saveButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*12)
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
			g.exitButton.SetPosition((float32(width)-w)/2, (float32(height)-h)/2+h*14)
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
			g.Scene.Remove(g.showControlsCheck)
			g.Scene.Remove(g.exitSaveCheck)
			g.Scene.Remove(g.mouseXSenSlider)
			g.Scene.Remove(g.mouseYSenSlider)
			g.Scene.Remove(g.moveSpeedSlider)
			g.Scene.Remove(g.tickSpeedSlider)
			g.Scene.Remove(g.dayLengthSlider)
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

func (g *Gui) moveSpeedLabel() string {
	return fmt.Sprintf("Movement Speed: %.2f", g.Cfg.Controls.MoveSpeed)
}

func (g *Gui) tickSpeedLabel() string {
	return fmt.Sprintf("Ticks Per Second: %d", g.Cfg.Simulation.Speed)
}

func (g *Gui) dayLengthLabel() string {
	suffix := "s"
	d := g.Cfg.Simulation.DayLength

	if d == 1 {
		suffix = ""
	}

	return fmt.Sprintf("Day Length: %d min%s", d, suffix)
}

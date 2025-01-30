package gui

import (
	"fmt"

	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/gui/color"
)

// Register the config menu
func (g *Gui) registerConfigMenu() {
	views[ConfigMenu] = viewControls{
		open: func(closeOthers bool) {
			if closeOthers {
				g.closeViews()
			}

			var (
				width, height  int = g.App.GetSize()
				w, h           float32
				nextY          float32
				saveMouseSensX float32 = g.Cfg.Controls.MouseSensitivityX
				saveMouseSensY float32 = g.Cfg.Controls.MouseSensitivityY
				saveTickSpeed  int     = g.Cfg.Simulation.Speed
			)

			g.showControlsCheck = gui.NewCheckBox("Show Controls")
			w, h = g.showControlsCheck.Width(), g.showControlsCheck.Height()
			g.showControlsCheck.SetPosition((float32(width)-w)/2, (float32(height)-h)/2)
			g.showControlsCheck.SetUserData(ConfigMenu)
			g.showControlsCheck.SetValue(g.Cfg.ShowHelp)
			g.showControlsCheck.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.ShowHelp = g.showControlsCheck.Value()
			})
			g.Scene.Add(g.showControlsCheck)
			nextY = g.showControlsCheck.Position().Y + g.showControlsCheck.Height() + 5

			g.exitSaveCheck = gui.NewCheckBox("Save on Exit")
			w = g.exitSaveCheck.Width()
			g.exitSaveCheck.SetPosition((float32(width)-w)/2, nextY)
			g.exitSaveCheck.SetUserData(ConfigMenu)
			g.exitSaveCheck.SetValue(g.Cfg.ExitSave)
			g.exitSaveCheck.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.ExitSave = g.exitSaveCheck.Value()
			})
			g.Scene.Add(g.exitSaveCheck)
			nextY = g.exitSaveCheck.Position().Y + g.exitSaveCheck.Height() + 10

			g.mouseXSenSlider = gui.NewHSlider(175, 12.5)
			w = g.mouseXSenSlider.Width()
			g.mouseXSenSlider.SetPosition((float32(width)-w)/2, nextY)
			g.mouseXSenSlider.SetUserData(ConfigMenu)
			g.mouseXSenSlider.SetValue(g.Cfg.Controls.MouseSensitivityX)
			g.mouseXSenSlider.SetText(g.mouseSensXLabel())
			g.mouseXSenSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Controls.MouseSensitivityX = g.mouseXSenSlider.Value()
				g.mouseXSenSlider.SetText(g.mouseSensXLabel())
			})
			g.Scene.Add(g.mouseXSenSlider)
			nextY = g.mouseXSenSlider.Position().Y + g.mouseXSenSlider.Height() + 5

			g.mouseYSenSlider = gui.NewHSlider(175, 12.5)
			w = g.mouseYSenSlider.Width()
			g.mouseYSenSlider.SetPosition((float32(width)-w)/2, nextY)
			g.mouseYSenSlider.SetUserData(ConfigMenu)
			g.mouseYSenSlider.SetValue(g.Cfg.Controls.MouseSensitivityY)
			g.mouseYSenSlider.SetText(g.mouseSensYLabel())
			g.mouseYSenSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Controls.MouseSensitivityY = g.mouseYSenSlider.Value()
				g.mouseYSenSlider.SetText(g.mouseSensYLabel())
			})
			g.Scene.Add(g.mouseYSenSlider)
			nextY = g.mouseYSenSlider.Position().Y + g.mouseYSenSlider.Height() + 5

			g.moveSpeedSlider = gui.NewHSlider(175, 12.5)
			w = g.moveSpeedSlider.Width()
			g.moveSpeedSlider.SetPosition((float32(width)-w)/2, nextY)
			g.moveSpeedSlider.SetUserData(ConfigMenu)
			g.moveSpeedSlider.SetValue(g.Cfg.Controls.MoveSpeed)
			g.moveSpeedSlider.SetText(g.moveSpeedLabel())
			g.moveSpeedSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Controls.MoveSpeed = math32.Max(0.01, g.moveSpeedSlider.Value())
				g.moveSpeedSlider.SetText(g.moveSpeedLabel())
			})
			g.Scene.Add(g.moveSpeedSlider)
			nextY = g.moveSpeedSlider.Position().Y + g.moveSpeedSlider.Height() + 5

			g.tickSpeedSlider = gui.NewHSlider(175, 12.5)
			w = g.tickSpeedSlider.Width()
			g.tickSpeedSlider.SetPosition((float32(width)-w)/2, nextY)
			g.tickSpeedSlider.SetUserData(ConfigMenu)
			g.tickSpeedSlider.SetValue(float32(g.Cfg.Simulation.Speed) / 32)
			g.tickSpeedSlider.SetText(g.tickSpeedLabel())
			g.tickSpeedSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Simulation.Speed = int(math32.Max(1, g.tickSpeedSlider.Value()*32))
				g.tickSpeedSlider.SetText(g.tickSpeedLabel())
			})
			g.Scene.Add(g.tickSpeedSlider)
			nextY = g.tickSpeedSlider.Position().Y + g.tickSpeedSlider.Height() + 5

			g.dayLengthSlider = gui.NewHSlider(175, 12.5)
			w = g.dayLengthSlider.Width()
			g.dayLengthSlider.SetPosition((float32(width)-w)/2, nextY)
			g.dayLengthSlider.SetUserData(ConfigMenu)
			g.dayLengthSlider.SetValue(float32(g.Cfg.Simulation.DayLength+1) / 30)
			g.dayLengthSlider.SetText(g.dayLengthLabel())
			g.dayLengthSlider.Subscribe(gui.OnChange, func(name string, ev interface{}) {
				g.Cfg.Simulation.DayLength = int(math32.Max(1, g.dayLengthSlider.Value()*30))
				g.State.Clock.SetTime(g.State.Clock.Hour, g.State.Clock.Minute)
				g.dayLengthSlider.SetText(g.dayLengthLabel())
			})
			g.Scene.Add(g.dayLengthSlider)
			nextY = g.dayLengthSlider.Position().Y + g.dayLengthSlider.Height() + 10

			g.saveButton = gui.NewButton("Save Settings")
			w = g.saveButton.Width()
			g.saveButton.SetPosition((float32(width)-w)/2, nextY)
			g.saveButton.SetUserData(ConfigMenu)
			g.saveButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				if err := g.Cfg.Save(); err != nil {
					fmt.Println(err)
					g.Notifications.Push(err.Error())
				} else {
					g.Notifications.Push("Settings saved")
				}

				Open(MainMenu, true)
			})
			g.saveButton.Subscribe(gui.OnCursor, func(s string, i interface{}) {
				if !g.saveButton.Enabled() {
					return
				}

				g.saveButton.SetColor(color.Green)
				g.saveButton.Label.SetColor(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
			})
			g.Scene.Add(g.saveButton)
			nextY = g.saveButton.Position().Y + g.saveButton.Height() + 5

			g.exitButton = gui.NewButton("Close")
			w = g.exitButton.Width()
			g.exitButton.SetPosition((float32(width)-w)/2, nextY)
			g.exitButton.SetUserData(TileContextMenu)
			g.exitButton.Subscribe(gui.OnClick, func(name string, ev interface{}) {
				g.Cfg.Controls.MouseSensitivityX = saveMouseSensX
				g.Cfg.Controls.MouseSensitivityY = saveMouseSensY
				g.Cfg.Simulation.Speed = saveTickSpeed
				Open(MainMenu, true)
			})
			g.exitButton.Subscribe(gui.OnCursor, func(s string, i interface{}) {
				if !g.exitButton.Enabled() {
					return
				}

				g.exitButton.SetColor(color.Red)
				g.exitButton.Label.SetColor(&math32.Color{})
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

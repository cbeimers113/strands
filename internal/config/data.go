package config

var defaultConfig = &Config{
	Name:     "Strands | Ecosystem Simulator",
	ShowHelp: false,
	ExitSave: true,
	Window: struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}{
		Width:  1200,
		Height: 800,
	},
	Simulation: struct {
		Width     int `json:"width"`
		Height    int `json:"height"`
		Depth     int `json:"depth"`
		Speed     int `json:"ticks_per_second"`
		DayLength int `json:"day_length_mins"`
	}{
		Width:     64,
		Height:    64,
		Depth:     64,
		Speed:     24,
		DayLength: 5,
	},
	Controls: struct {
		MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
		MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
		MoveSpeed         float32 `json:"move_speed"`
	}{
		MouseSensitivityX: 0.3,
		MouseSensitivityY: 0.3,
		MoveSpeed:         0.4,
	},
}

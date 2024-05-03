package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validate(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		err  error
	}{
		{
			name: "Happy path",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
		},
		{
			name: "Sad path - invalid name",
			cfg: Config{
				Name: "",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%sapplication name empty", errInvalidCfg),
		},
		{
			name: "Sad path - invalid window width",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  0,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%swindow width [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid window height",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 0,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%swindow height [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid simulation width",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     0,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%ssimulation width [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid simulation height",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    0,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%ssimulation height [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid simulation depth",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     0,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%ssimulation depth [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid simulation speed",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%ssimulation speed (TPS) [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid day length",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 0,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%ssimulation day length too small: [0 minutes]", errInvalidCfg),
		},
		{
			name: "Sad path - mouse X sensitivity too low",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%smouse X sensitivity must be between 0 and 1", errInvalidCfg),
		},
		{
			name: "Sad path - mouse Y sensitivity too high",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 10,
					MoveSpeed:         0.5,
				},
			},
			err: fmt.Errorf("%smouse Y sensitivity must be between 0 and 1", errInvalidCfg),
		},
		{
			name: "Sad path - move speed too low",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         0,
				},
			},
			err: fmt.Errorf("%smove speed must be between 0 and 1", errInvalidCfg),
		},
		{
			name: "Sad path - move speed too high",
			cfg: Config{
				Name: "Strands Test",

				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 800,
				},

				Simulation: struct {
					Width     int `json:"-"`
					Height    int `json:"-"`
					Depth     int `json:"-"`
					Speed     int `json:"ticks_per_second"`
					DayLength int `json:"day_length_mins"`
				}{
					Width:     64,
					Height:    64,
					Depth:     64,
					Speed:     60,
					DayLength: 5,
				},

				Controls: struct {
					MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
					MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
					MoveSpeed         float32 `json:"move_speed"`
				}{
					MouseSensitivityX: 0.025,
					MouseSensitivityY: 0.015,
					MoveSpeed:         2.5,
				},
			},
			err: fmt.Errorf("%smove speed must be between 0 and 1", errInvalidCfg),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.err, tt.cfg.validate())
		})
	}
}

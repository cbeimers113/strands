package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cbeimers113/strands/internal/config"
)

func Test_Load(t *testing.T) {
	good := config.Config{
		Name:     "Strands | Test",
		Version:  "0.0.0",
		ShowHelp: true,

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
			MoveSpeed:         0.25,
		},
	}

	got, err := config.Load(config.TestCfgData)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, good, *got)
}

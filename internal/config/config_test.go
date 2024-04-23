package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cbeimers113/strands/internal/config"
)

func Test_Load(t *testing.T) {
	good := config.Config{
		Name:    "Strands Test",
		Version: "0.0.0",

		Window: struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		}{
			Width:  1200,
			Height: 800,
		},

		Simulation: struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			Depth  int `json:"depth"`
			Speed  int `json:"ticks_per_second"`
		}{
			Width:  64,
			Height: 64,
			Depth:  64,
			Speed:  60,
		},

		Controls: struct {
			MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
			MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
		}{
			MouseSensitivityX: 0.025,
			MouseSensitivityY: 0.015,
		},
	}

	got, err := config.Load(config.TestCfgData)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, good, *got)
}

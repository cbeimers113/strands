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
	}

	got, err := config.Load(config.TestCfgData)
	assert.NotNil(t, got)
	assert.Equal(t, good, *got)
	assert.NoError(t, err)
}

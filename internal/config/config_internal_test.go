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
			},
		},
		{
			name: "Sad path - invalid name",
			cfg: Config{
				Name:    "",
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
			},
			err: fmt.Errorf("%sapplication name empty", errInvalidCfg),
		},
		{
			name: "Sad path - invalid version",
			cfg: Config{
				Name:    "Strands Test",
				Version: "0.0",
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
			},
			err: fmt.Errorf("%ssemantic version not provided", errInvalidCfg),
		},
		{
			name: "Sad path - invalid window width",
			cfg: Config{
				Name:    "Strands Test",
				Version: "0.0.0",
				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  0,
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
			},
			err: fmt.Errorf("%swindow width [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid window height",
			cfg: Config{
				Name:    "Strands Test",
				Version: "0.0.0",
				Window: struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				}{
					Width:  1200,
					Height: 0,
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
			},
			err: fmt.Errorf("%swindow height [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid simulation width",
			cfg: Config{
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
					Width:  0,
					Height: 64,
					Depth:  64,
					Speed:  60,
				},
			},
			err: fmt.Errorf("%ssimulation width [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid simulation height",
			cfg: Config{
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
					Height: 0,
					Depth:  64,
					Speed:  60,
				},
			},
			err: fmt.Errorf("%ssimulation height [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid simulation depth",
			cfg: Config{
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
					Depth:  0,
					Speed:  60,
				},
			},
			err: fmt.Errorf("%ssimulation depth [0] too small", errInvalidCfg),
		},
		{
			name: "Sad path - invalid simulation speed",
			cfg: Config{
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
				},
			},
			err: fmt.Errorf("%ssimulation speed (TPS) [0] too small", errInvalidCfg),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.err, tt.cfg.validate())
		})
	}
}

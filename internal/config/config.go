package config

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Window  struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"window"`

	Simulation struct {
		Width  int `json:"width"`
		Height int `json:"height"`
		Depth  int `json:"depth"`
		Speed  int `json:"ticks_per_second"`
	} `json:"simulation"`

	Controls struct {
		MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
		MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
	} `json:"controls"`
}

const (
	errInvalidCfg = "invalid config: "
)

func Load(data []byte) (*Config, error) {
	c := &Config{}
	err := json.Unmarshal(data, c)

	if err != nil {
		return nil, err
	}

	if err = c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c Config) validate() error {
	if c.Name == "" {
		return fmt.Errorf("%sapplication name empty", errInvalidCfg)
	}
	if m := regexp.MustCompile(`^[0-9].[0-9].[0-9](-snapshot)?$`); !m.MatchString(c.Version) {
		return fmt.Errorf("%ssemantic version not provided", errInvalidCfg)
	}

	if c.Window.Width <= 0 {
		return fmt.Errorf("%swindow width [%d] too small", errInvalidCfg, c.Window.Width)
	}
	if c.Window.Height <= 0 {
		return fmt.Errorf("%swindow height [%d] too small", errInvalidCfg, c.Window.Height)
	}

	if c.Simulation.Width <= 0 {
		return fmt.Errorf("%ssimulation width [%d] too small", errInvalidCfg, c.Simulation.Width)
	}
	if c.Simulation.Height <= 0 {
		return fmt.Errorf("%ssimulation height [%d] too small", errInvalidCfg, c.Simulation.Height)
	}
	if c.Simulation.Depth <= 0 {
		return fmt.Errorf("%ssimulation depth [%d] too small", errInvalidCfg, c.Simulation.Depth)
	}
	if c.Simulation.Speed < 1 {
		return fmt.Errorf("%ssimulation speed (TPS) [%d] too small", errInvalidCfg, c.Simulation.Speed)
	}

	if c.Controls.MouseSensitivityX <= 0 || c.Controls.MouseSensitivityX > 1 {
		return fmt.Errorf("%smouse X sensitivity must be between 0 and 1", errInvalidCfg)
	}
	if c.Controls.MouseSensitivityY <= 0 || c.Controls.MouseSensitivityY > 1 {
		return fmt.Errorf("%smouse Y sensitivity must be between 0 and 1", errInvalidCfg)
	}

	return nil
}

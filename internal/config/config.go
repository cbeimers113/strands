package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"cbeimers113/strands/internal/io/file"
)

type Config struct {
	Name     string `json:"name"`
	ShowHelp bool   `json:"show_controls"`
	ExitSave bool   `json:"save_on_exit"`
	Window   struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"window"`

	Simulation struct {
		Width     int `json:"-"`
		Height    int `json:"-"`
		Depth     int `json:"-"`
		Speed     int `json:"ticks_per_second"`
		DayLength int `json:"day_length_mins"`
	} `json:"simulation"`

	Controls struct {
		MouseSensitivityX float32 `json:"mouse_sensitivity_x"`
		MouseSensitivityY float32 `json:"mouse_sensitivity_y"`
		MoveSpeed         float32 `json:"move_speed"`
	} `json:"controls"`
}

const (
	Width, Height, Depth = 64, 64, 64

	errInvalidCfg = "invalid config: "
)

var configFilePath string

func Load() (*Config, error) {
	var (
		c    *Config
		data []byte
		err  error
	)

	c = &Config{}
	configFilePath = filepath.Join(file.StoragePath, "config.json")

	// Make sure a config file exists, otherwise use the default
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No config file found, loading default config")
		c = defaultConfig
		return c, nil
	} else if data, err = os.ReadFile(configFilePath); err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Use a constant world size so that save files will be compatible
	c.Simulation.Width = Width
	c.Simulation.Height = Height
	c.Simulation.Depth = Depth

	if err = c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c Config) Save() error {
	data, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		return err
	}

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func (c Config) validate() error {
	if c.Name == "" {
		return fmt.Errorf("%sapplication name empty", errInvalidCfg)
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
	if c.Simulation.DayLength < 1 {
		return fmt.Errorf("%ssimulation day length too small: [%d minutes]", errInvalidCfg, c.Simulation.DayLength)
	}

	if c.Controls.MouseSensitivityX <= 0 || c.Controls.MouseSensitivityX > 1 {
		return fmt.Errorf("%smouse X sensitivity must be between 0 and 1", errInvalidCfg)
	}
	if c.Controls.MouseSensitivityY <= 0 || c.Controls.MouseSensitivityY > 1 {
		return fmt.Errorf("%smouse Y sensitivity must be between 0 and 1", errInvalidCfg)
	}
	if c.Controls.MoveSpeed <= 0 || c.Controls.MoveSpeed > 1 {
		return fmt.Errorf("%smove speed must be between 0 and 1", errInvalidCfg)
	}

	return nil
}

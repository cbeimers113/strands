package config

import (
	_ "embed"
)

//go:embed testdata/test_cfg.json
var TestCfgData []byte

func LoadTestConfig() (*Config, error) {
	return Load(TestCfgData)
}

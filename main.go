package main

import (
	_ "embed"

	"cbeimers113/strands/internal/config"
	"cbeimers113/strands/internal/game"
)

//go:embed cfg.json
var cfgData []byte

func main() {
	cfg, err := config.Load(cfgData)
	if err != nil {
		panic(err)
	}

	g, err := game.New(cfg)
	if err != nil {
		panic(err)
	}

	g.Start()
}

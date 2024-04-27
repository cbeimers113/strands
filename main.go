package main

import (
	_ "embed"
	"os"

	"cbeimers113/strands/internal/config"
	"cbeimers113/strands/internal/game"
)


func main() {
	cfgData, err := os.ReadFile("cfg.json")
	if err != nil {
		panic(err)
	}
	
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

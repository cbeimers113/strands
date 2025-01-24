package main

import (
	_ "embed"
	"fmt"
	"regexp"

	"cbeimers113/strands/internal/config"
	"cbeimers113/strands/internal/game"
)

//go:embed .version
var Version string

func main() {
	if m := regexp.MustCompile(`^[0-9].[0-9].[0-9](-snapshot)?$`); !m.MatchString(Version) {
		fmt.Printf("Warning: non-semantic version number provided: %s\n", Version)
	}

	// TODO: remove this tag if doing a release/automated build
	Version += " snapshot"

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	g, err := game.New(cfg, Version)
	if err != nil {
		panic(err)
	}

	g.Start()
}

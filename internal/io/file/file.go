package file

import (
	"fmt"
	"os"
	"path/filepath"
)

var StoragePath string

func init() {
	// Get the path to APPDATA/Roaming
	var err error
	if StoragePath, err = os.UserConfigDir(); err != nil {
		panic(err)
	}

	StoragePath = filepath.Join(StoragePath, "strands")

	// Create the config directory if it doesn't exist
	if err = os.MkdirAll(StoragePath, os.ModePerm); err != nil {
		panic(fmt.Errorf("couldn't create config directory [%s]: %w", StoragePath, err))
	}
}

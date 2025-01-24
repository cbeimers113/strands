package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// Touch returns a full filepath to a new file name and extension, replacing spaces with underscores and removing special chars
func Touch(filename, extension string) string {
	underscored := strings.ReplaceAll(filename, " ", "_")
	fp := ""

	for _, c := range underscored {
		if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z' || ('0' <= c && c <= '9') || c == '.' || c == '_') {
			fp += string(c)
		}
	}

	return filepath.Join(StoragePath, fp+extension)
}

// Name returns the name of a file excluding the directory path
func Name(filepath string) string {
	tokens := strings.Split(filepath, "/")
	return tokens[len(tokens)-1]
}

// Exists returns whether a given file exists
func Exists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !errors.Is(err, os.ErrNotExist)
}

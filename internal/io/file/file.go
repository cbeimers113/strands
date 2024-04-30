package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flytam/filenamify"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

// Touch returns a full filepath to a new file name and extension, replacing spaces with underscores and removing illegal chars
func Touch(filename, extension string) string {
	fp, err := filenamify.Filenamify(filename, filenamify.Options{})
	if err != nil {
		fmt.Printf("Warning: couldn't create filename for %s\n", filename)
	}

	caser := cases.Title(language.English)
	fp = caser.String(fp)

	return filepath.Join(StoragePath, fp+extension)
}

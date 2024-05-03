package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/valyala/gozstd"

	"cbeimers113/strands/internal/config"
	"cbeimers113/strands/internal/entity"
	"cbeimers113/strands/internal/io/file"
)

const SaveFileExtension = ".sim"

var ExitSaveFile = file.Touch("", SaveFileExtension)

type Save struct {
	Seed  int64          `json:"seed"`
	Clock *Clock         `json:"clock"`
	Cells []*Cell        `json:"atmosphere"`
	Tiles []*entity.Tile `json:"tiles"` //plants are embedded in the tiles they're on
	// creatures []*entity.Creature `json:"creatures"`
}

func LoadSave(cfg *config.Config, filename string) (*State, []*Cell, []*entity.Tile, error) {
	var (
		data  []byte
		save  Save
		state *State
		err   error
	)

	// Check if file exists
	if _, err = os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil, fmt.Errorf("save file does not exist: [%s]", filename)
	}

	if data, err = os.ReadFile(filename); err != nil {
		return nil, nil, nil, fmt.Errorf("error reading save file [%s]: %w", filename, err)
	}

	if data, err = gozstd.Decompress(nil, data); err != nil {
		return nil, nil, nil, fmt.Errorf("error decompressing save file [%s]: %w", filename, err)
	}

	if err = json.Unmarshal(data, &save); err != nil {
		return nil, nil, nil, fmt.Errorf("error unmarshaling save file [%s]: %w", filename, err)
	}

	state = New(cfg, save.Seed)
	state.Clock = save.Clock
	state.Clock.Config = cfg

	return state, save.Cells, save.Tiles, nil
}

func StoreSave(filename string, state *State, cells []*Cell, tiles []*entity.Tile) error {
	save := Save{
		Seed:  state.Seed,
		Clock: state.Clock,
		Cells: cells,
		Tiles: tiles,
	}

	data, err := json.MarshalIndent(save, "", "	")
	if err != nil {
		return err
	}

	data = gozstd.Compress(nil, data)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	return err
}

func GetSavesList() map[string]string {
	var saves = make(map[string]string)

	filepath.WalkDir(file.StoragePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) == SaveFileExtension {
			parts := strings.Split(path, "/")
			filename := strings.TrimSuffix(parts[len(parts)-1], SaveFileExtension)

			if filename == "" {
				filename = "autosave"
			}

			saves[filename] = path
		}

		return nil
	})

	return saves
}

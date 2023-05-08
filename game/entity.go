package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/g3n/engine/core"
)

type EntityType = string
type Entity = *core.Node

const Tile EntityType = "tile"
const Plant EntityType = "plant"
const Creature EntityType = "creature"

// Metadata format
// entityType-[gene_0, gene_1, ..., gene_n]
// genes are strings or ints encoded as strings

// Generate the metadata tag for an entity
func CreateTag(entityType EntityType, data []string) (metadata string) {
	metadata = entityType + "-"

	for i, datum := range data {
		metadata += datum

		if i < len(data)-1 {
			metadata += "-"
		}
	}

	return
}

// Get an entity's type
func TypeOf(entity Entity) (entityType EntityType) {
	data := strings.Split(entity.Name(), "-")
	entityType = data[0]

	return
}

// Get an entity's metadata
func Metadata(entity Entity) (metadata []string) {
	data := strings.Split(entity.Name(), "-")
	metadata = data[1:]

	return
}

// Get a metadata entry, return empty string if not exist
func Datum(entity Entity, index int) (datum string) {
	data := Metadata(entity)

	if index >= 0 && index < len(data) {
		datum = data[index]
	}

	return
}

// Get a metadata entry as int, return 0 if not exist or failed to decode
func DatumNum(entity Entity, index int) (datum int) {
	conv, _ := strconv.ParseInt(Datum(entity, index), 10, 0)
	datum = int(conv)

	return
}

// Modify the value of a metadata entry
func SetDatum(entity Entity, index int, datum string) {
	data := Metadata(entity)

	if index >= 0 && index < len(data) {
		data[index] = datum
		entity.SetName(CreateTag(TypeOf(entity), data))
	}
}

// Modify the value of a numerical metadata entry
func SetDatumNum(entity Entity, index int, datum int) {
	SetDatum(entity, index, fmt.Sprintf("%d", datum))
}

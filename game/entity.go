package game

import (
	"strings"

	"github.com/g3n/engine/core"
)

type EntityType = string
type Entity = *core.Node
type Strand []int

const Tile EntityType = "tile"
const Plant EntityType = "plant"
const Creature EntityType = "creature"

// Return the type of this entity
func Type(entity Entity) EntityType {
	return strings.Split(entity.Name(), " ")[0]
}

// Return this entity's metadata entry at the given index
func Datum(entity Entity, index int) int {
	data, ok := entity.UserData().(Strand)

	if ok {
		return data[index]
	}

	return 0
}

// Set this entity's metadata entry at the given index, return whether it could be set
func SetDatum(entity Entity, index int, val int) (ok bool) {
	data, ok := entity.UserData().(Strand)

	if ok && index >= 0 && index < len(data) {
		data[index] = val
	}

	return
}

package game

import (
	"github.com/g3n/engine/core"
)

type EntityType = string
type Entity = *core.Node

const Tile EntityType = "tile"
const Plant EntityType = "plant"
const Creature EntityType = "creature"

type EntityData struct {
	eType    EntityType
	metadata []int
}

// Return the type of this entity
func Type(entity Entity) EntityType {
	data, ok := entity.UserData().(EntityData)

	if ok {
		return data.eType
	}

	return ""
}

// Return this entity's metadata entry at the given index
func GetDatum(entity Entity, index int) int {
	data, ok := entity.UserData().(EntityData)

	if ok {
		return data.metadata[index]
	}

	return 0
}

// Set this entity's metadata entry at the given index
func SetDatum(entity Entity, index int, val int) {
	data, ok := entity.UserData().(EntityData)

	if ok && index >= 0 && index < len(data.metadata) {
		metadata := data.metadata
		metadata[index] = val
		entity.SetUserData(EntityData{eType: data.eType, metadata: metadata})
	}
}

package game

import (
	"strings"

	"github.com/g3n/engine/core"
)

// Entity data storage model:

// The g3n Node type has a UserData field of type interface{} (the 'any' type)
// which we can set and get any data we want. Each entity type stores its own struct
// type (eg Plant entities store PlantData, Tile entities store TileData...)
// Behind the scenes of the code on this side of the engine, all Entities are
// pointers to g3n nodes.

// The g3n Node also has a Name field which is used to store the entity type.
// Entity type is checked when accessing the entity's UserData so that the proper
// struct fields are used.

type Entity = *core.Node
type EntityType = string

const Tile EntityType = "tile"
const Plant EntityType = "plant"
const Creature EntityType = "creature"

// Return the type of this entity
func Type(entity Entity) EntityType {
	return strings.Split(entity.Name(), " ")[0]
}

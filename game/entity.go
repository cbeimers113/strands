package game

import (
	"fmt"
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
func TypeOf(entity Entity) EntityType {
	return strings.Split(entity.Name(), " ")[0]
}

// Return an infostring representing this entity
func EntityInfo(entity Entity) (infoString string) {
	eType := TypeOf(entity)
	infoString = fmt.Sprintf("%s:\n", eType)

	switch eType {
	case Tile:
		tileData, ok := entity.UserData().(TileData)

		if ok {
			infoString += fmt.Sprintf("type=%s\n", tileData.Type)
		}
	case Plant:
		plantData, ok := entity.UserData().(PlantData)

		if ok {
			infoString += fmt.Sprintf("age=%d\n", plantData.Age)
			infoString += fmt.Sprintf("colour=%x\n", plantData.Colour)
		}
	case Creature:
		infoString += "not implemented yet, how are you seeing this?"
	}

	return
}

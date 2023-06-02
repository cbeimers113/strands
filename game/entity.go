package game

import (
	"fmt"
	"strconv"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
)

// Entity data storage model:

// The g3n Node type has a UserData field of type interface{} (the 'any' type)
// which we can set and get any data we want. Each entity type stores its own struct
// type (eg Plant entities store PlantData, Tile entities store TileData...)
// The entity type wraps and promotes a *core.Node

type EntityType = string
type Entity struct {
	*graphic.Mesh
	Type EntityType
}

const Tile EntityType = "tile"
const Plant EntityType = "plant"
const Creature EntityType = "creature"

// Create a new entity with the given parameters
func NewEntity(mesh *graphic.Mesh, eType ElementType) (entity *Entity) {
	entity = &Entity{
		mesh,
		eType,
	}
	entity.GetMaterial(0).GetMaterial().SetLineWidth(10)
	
	// Store the index of this entity in its name so that the entity can be found by a game object
	entity.SetName(fmt.Sprintf("%d", len(Entities)))
	Entities[len(Entities)] = entity

	return
}

// Return an infostring representing this entity
func (entity *Entity) InfoString() (infoString string) {
	eType := entity.Type
	infoString = fmt.Sprintf("%s:\n", eType)

	switch eType {
	case Tile:
		if tileData, ok := entity.UserData().(TileData); ok {
			infoString += fmt.Sprintf("type=%s\n", tileData.Type.Name)
			infoString += fmt.Sprintf("temperature=%.2f°C\n", tileData.Temperature)
			infoString += fmt.Sprintf("moisture=%.2f%%\n", tileData.Moisture)
		}
	case Plant:
		if plantData, ok := entity.UserData().(PlantData); ok {
			infoString += fmt.Sprintf("age=%d\n", plantData.Age)
			infoString += fmt.Sprintf("colour=#%06x\n", plantData.Colour)
		}
	case Creature:
		infoString += "not implemented yet, how are you seeing this?"
	}

	return
}

// Highlight or unhighlight an entity
func (entity *Entity) Highlight(highlight bool) {
	entity.GetMaterial(0).GetMaterial().SetWireframe(highlight)
}

// Get the entity associated with a node, return nil if there isn't one
func EntityOf(node *core.Node) (entity *Entity) {
	if i, err := strconv.Atoi(node.Name()); err == nil {
		entity = Entities[i]
	}

	return
}

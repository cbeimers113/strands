package game

import (
	"fmt"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
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

// Set whether an entity is highlighted
func Highlight(entity Entity, highlight bool) {
	// TODO: Find a more efficient way to do this
	// Dig out the material and modify it
	if mesh, ok := entity.GetINode().(*graphic.Mesh); ok {
		if imat := mesh.GetMaterial(0); imat != nil {
			if mat := imat.GetMaterial(); mat != nil {
				if tex, ok := Texture("highlight"); ok {
					mat.RemoveTexture(tex)
					if highlight && !mat.HasTexture(tex) {
						mat.AddTexture(tex)
					}
				}
			}
		}
	}
}

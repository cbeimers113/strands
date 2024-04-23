package entity

import (
	"fmt"

	"github.com/g3n/engine/graphic"
)

// Entity data storage model:

// The g3n Node type has a UserData field of type interface{} (the 'any' type)
// which we can set and get any data we want. Each entity type stores its own struct
// type (eg Plant entities store PlantData, Tile entities store TileData...)
// The entity type wraps and promotes a *core.Node

const (
	Tile     EntityType = "tile"
	Plant    EntityType = "plant"
	Creature EntityType = "creature"
)

type EntityType = string

type Entity struct {
	*graphic.Mesh
	Type EntityType
}

// Create a new entity with the given parameters
func New(mesh *graphic.Mesh, eType EntityType, entities map[int]*Entity) *Entity {
	e := &Entity{
		Mesh: mesh,
		Type: eType,
	}
	e.GetMaterial(0).GetMaterial().SetLineWidth(8)

	// Store the index of this entity in its name so that the entity can be found by a game object
	e.SetName(fmt.Sprintf("%d", len(entities)))
	entities[len(entities)] = e

	return e
}

// Return an infostring representing this entity
func (e Entity) InfoString() string {
	infoString := fmt.Sprintf("%s:\n", e.Type)

	switch e.Type {
	case Tile:
		if tileData, ok := e.UserData().(*TileData); ok {
			infoString += fmt.Sprintf("type=%s\n", tileData.Type.Name)
			infoString += fmt.Sprintf("temperature=%s\n", tileData.Temperature)
			infoString += fmt.Sprintf("water level=%s\n", tileData.WaterLevel)
			infoString += fmt.Sprintf("elevation=%s\n", e.getElevation())
			infoString += fmt.Sprintf("planted=%t\n", tileData.Planted)
		}
	case Plant:
		if plantData, ok := e.UserData().(*PlantData); ok {
			infoString += fmt.Sprintf("age=%d\n", plantData.Age)
			infoString += fmt.Sprintf("colour=#%06x\n", plantData.Colour)
		}
	case Creature:
		infoString += "not implemented yet, how are you seeing this?"
	}

	return infoString
}

// Highlight or unhighlight an entity
func (e *Entity) Highlight(highlight bool) {
	e.GetMaterial(0).GetMaterial().SetWireframe(highlight)
}

package entity

import (
	"fmt"

	"github.com/g3n/engine/material"
)

type Entity interface {
	Update()
	InfoString() string
	Material() *material.Material
	SetName(string)
	Dispose()
	DisposeChildren(bool)
}

// Add an entity to the entities map
func AddEntity(entity Entity, entities map[int]Entity) {
	// Store the index of this entity in its name so that the entity can be found by a game object
	entity.SetName(fmt.Sprintf("%d", len(entities)))
	entities[len(entities)] = entity

	// If this entity is a tile, store the water sub-meshes under the same entity index
	if tile, ok := entity.(*Tile); ok {
		tile.Water.SetName(tile.Name())
	}
}

// Highlight or unhighlight an entity
func Highlight(entity Entity, highlight bool) {
	entity.Material().SetWireframe(highlight)
}

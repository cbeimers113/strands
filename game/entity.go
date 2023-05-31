package game

import (
	"fmt"
	"strconv"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Entity data storage model:

// The g3n Node type has a UserData field of type interface{} (the 'any' type)
// which we can set and get any data we want. Each entity type stores its own struct
// type (eg Plant entities store PlantData, Tile entities store TileData...)
// The entity type wraps and promotes a *core.Node

type EntityType = string
type Entity struct {
	*graphic.Mesh
	
	Type     EntityType
	Collider math32.Box3
	Material *material.Material
}

const Tile EntityType = "tile"
const Plant EntityType = "plant"
const Creature EntityType = "creature"

// Create a new entity with the given parameters
func NewEntity(mesh *graphic.Mesh, eType ElementType) (entity *Entity) {
	var mat *material.Material

	if imat := mesh.GetMaterial(0); imat != nil {
		mat = imat.GetMaterial()
	}

	entity = &Entity{
		mesh,
		eType,
		mesh.BoundingBox(),
		mat,
	}

	entity.SetName(fmt.Sprintf("%d", len(Entities)))
	Entities[len(Entities)] = entity

	return
}

// Return an infostring representing this entity
func EntityInfo(entity *Entity) (infoString string) {
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

// Set whether an entity is highlighted
func Highlight(entity *Entity, highlight bool) {
	if mat := entity.Material; mat != nil {
		if tex, ok := Texture("highlight"); ok {
			mat.RemoveTexture(tex)

			if highlight && !mat.HasTexture(tex) {
				mat.AddTexture(tex)
			}
		}
	}
}

// Check if two entities are colliding
func Colliding(entity1, entity2 *Entity) (colliding bool) {
	A := entity1.Collider
	B := entity2.Collider
	colliding = A.IsIntersectionBox(&B)

	return
}

// Get the entity associated with a node, return nil if there isn't one
func EntityOf(node *core.Node) (entity *Entity) {
	if i, err := strconv.Atoi(node.Name()); err == nil {
		entity = Entities[i]
	}

	return
}

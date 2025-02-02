package entity

import (
	"fmt"
	"strconv"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/material"
)

type Entity interface {
	Update()
	Name() string
	InfoString() string
	Material() *material.Material
	SetName(string)
	GetINode() core.INode
}

// Add an entity to the entities map
func AddEntity(entity Entity, entities map[int]Entity, scene *core.Node) {
	// Store the index of this entity in its name so that the entity can be found by a game object
	entity.SetName(fmt.Sprintf("%d", len(entities)))
	entities[len(entities)] = entity
	scene.Add(entity.GetINode())
}

// Remove an entity from the entities map
func RemoveEntity(entity Entity, entities map[int]Entity, scene *core.Node) {
	// Remove the entity and shift all the entities above it down the list
	if idx, err := strconv.Atoi(entity.Name()); err == nil {
		for i := idx + 1; i < len(entities); i++ {
			entities[i-1] = entities[i]
			entities[i-1].SetName(fmt.Sprintf("%d", i-1))
		}

		scene.Remove(entity.GetINode())
		delete(entities, len(entities)-1)
	}
}

// Highlight or unhighlight an entity
func Highlight(entity Entity, highlight bool) {
	if mat := entity.Material(); mat != nil {
		mat.SetWireframe(highlight)
	}
}

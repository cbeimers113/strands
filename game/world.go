package game

import (
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
)

const Width int = 64
const Depth int = 64
const TileSize float32 = 4

var Sun *light.Ambient
var Entities map[int]Entity

// Add an entity to the entity list
func AddEntityTo(node *core.Node, entity *graphic.Mesh) {
	node.Add(entity)
	Entities[len(Entities)] = node.ChildAt(len(node.Children()) - 1).GetNode()
}

// Remove an entity from the world
func RemoveEntity(entity Entity) {
	// Recursively remove children of this entity
	for _, child := range entity.Children() {
		RemoveEntity(child.GetNode())
	}

	dropEntity(entity)
	entity.SetVisible(false)
	entity.RemoveAll(true)
	entity.Parent().GetNode().Remove(entity)
	entity.DisposeChildren(true)
	entity.Dispose()
}

// Remove an entity from the Entities list
func dropEntity(entity Entity) {
	index := -1

	// Find the entity's index in the Entities list
	for i, ent := range Entities {
		if ent == entity {
			index = i
			break
		}
	}

	// If the entity is in the Entities list, remove it and shift all the entities above it down the list
	if index >= 0 {
		for i := index + 1; i < len(Entities); i++ {
			Entities[i-1] = Entities[i]
		}
		delete(Entities, len(Entities)-1)
	}
}

// Load the world into the scene
func LoadWorld() {
	// Sun
	Sun = light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 8.0)
	Scene.Add(Sun)
	Entities = make(map[int]Entity)

	// Tiles
	pnoise := perlin.NewPerlin(1, 0.1, 2, rand.Int63())

	// Create heightmap
	var heightmap [Width][Depth]float32
	var min float32 = 1_000_000_000.0
	var max float32 = -min

	for x := 0; x < Width; x++ {
		for z := 0; z < Depth; z++ {
			height := float32(math.Abs(pnoise.Noise2D(float64(x), float64(z)) * 1000))
			heightmap[x][z] = height

			// Record min and max so that the tile types can be mapped to height range
			if height < min {
				min = height
			}

			if height > max {
				max = height
			}
		}
	}

	// Add tiles to world
	for x := 0; x < Width; x++ {
		for z := 0; z < Depth; z++ {
			// Map the heightmap value to the TileTypes array to determine tile type
			y := math32.Min(float32(len(TileTypes))*(heightmap[x][z]-min)/(max-min), float32(len(TileTypes)-1))
			tType := TileTypes[int(y)]
			AddEntityTo(Scene, NewTile(x, z, y, tType))
		}
	}
}

// Update the game world, deltaTime is time since last update in ms
func Update(deltaTime int) {
	for _, entity := range Entities {
		switch Type(entity) {
		case Plant:
			GrowPlant(entity)
		}
	}
}

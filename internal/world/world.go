package world

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/atmosphere"
	"cbeimers113/strands/internal/chem"
	"cbeimers113/strands/internal/context"
	"cbeimers113/strands/internal/entity"
)

type World struct {
	*context.Context

	sun        *light.Ambient
	tilemap    [][]*entity.Entity
	atmosphere *atmosphere.Atmosphere
}

func New(ctx *context.Context) *World {
	w := &World{
		Context: ctx,

		// Sun TODO: Sun entity type for day/night cycle
		sun:        light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 8.0),
		atmosphere: atmosphere.New(ctx),
	}

	w.Scene.Add(w.sun)
	w.createMap()

	return w
}

// Create the tilemap
func (w *World) createMap() {
	heightmap, min, max := w.makeHeightmap()
	w.makeTilemap(heightmap, min, max)
	w.assignTileNeighbourhoods()
}

// Remove an entity from the world
func (w *World) removeEntity(entity *entity.Entity) {
	index := w.entityIndex(entity)

	// If the entity is in the Entities list, remove it and shift all the entities above it down the list
	if index >= 0 {
		for i := index + 1; i < len(w.State.Entities); i++ {
			w.State.Entities[i-1] = w.State.Entities[i]
			w.State.Entities[i-1].SetName(fmt.Sprintf("%d", i-1))
		}

		delete(w.State.Entities, len(w.State.Entities)-1)
	}

	entity.DisposeChildren(true)
	entity.Dispose()
}

// Find an entity's index in the Entities list, return -1 if not found
func (w *World) entityIndex(entity *entity.Entity) int {
	for i, ent := range w.State.Entities {
		if ent == entity {
			return i
		}
	}

	return -1
}

// Check if a given coordinate is within the tilemap boundaries
func (w *World) inBounds(x, z int) bool {
	return x >= 0 && x < w.Cfg.Simulation.Width && z >= 0 && z < w.Cfg.Simulation.Depth
}

// Generate a heightmap, return the map and its min and max values
func (w *World) makeHeightmap() ([][]float32, float32, float32) {
	var heightmap = make([][]float32, w.Cfg.Simulation.Width)
	var min float32 = 1_000_000_000.0
	var max float32 = -min
	pnoise := perlin.NewPerlin(1, 0.1, 2, rand.Int63())

	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		heightmap[x] = make([]float32, w.Cfg.Simulation.Depth)

		for z := 0; z < w.Cfg.Simulation.Depth; z++ {
			height := float32(math.Abs(pnoise.Noise2D(float64(x), float64(z))))
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

	return heightmap, min, max
}

// Create a tilemap with a given heightmap specification
func (w *World) makeTilemap(heightmap [][]float32, min, max float32) {
	w.tilemap = make([][]*entity.Entity, w.Cfg.Simulation.Width)

	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		w.tilemap[x] = make([]*entity.Entity, w.Cfg.Simulation.Depth)

		for z := 0; z < w.Cfg.Simulation.Depth; z++ {
			// Map the heightmap value to the TileTypes array to determine tile type
			height := float32(len(entity.TileTypes)) * (heightmap[x][z] - min) / (max - min)
			tType := entity.Stone

			if int(height) < len(entity.TileTypes) {
				tType = entity.TileTypes[int(height)]
			}

			height /= 3

			// Each tile spawns at 22°C with 10 L of water on top of it
			tile := entity.NewTile(w.State.Entities, x, z, height, 22.0, 10, tType)
			w.Scene.Add(tile.GetINode())
			w.tilemap[x][z] = tile
		}
	}
}

// Give each tile in a tilemap a list of pointers to its neighbours
func (w *World) assignTileNeighbourhoods() {
	// Base hexmap neighbourhood offsets
	nbOffsets := [][]int{
		{1, 0},  // Right
		{1, -1}, // Top right
		{1, 1},  // Bottom right
		{-1, 0}, // Left
		{0, -1}, // Top left
		{0, 1},  // Bottom left
	}

	// Assign each tile's neighbours
	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		for z := 0; z < w.Cfg.Simulation.Depth; z++ {
			var neighbours entity.Neighbourhood
			tile := w.tilemap[x][z]

			for i, offs := range nbOffsets {
				xOffs := x + offs[0]
				zOffs := z + offs[1]

				// Stagger offsets on the x axis for every other row for "top/bottom" neighbours
				if z%2 == 0 && i%3 != 0 {
					xOffs--
				}

				if w.inBounds(xOffs, zOffs) {
					neighbours[i] = w.tilemap[xOffs][zOffs]
				}
			}

			// Load the neighbour pointers into the tile's metadata
			data, _ := tile.UserData().(*entity.TileData)
			data.Neighbours = neighbours
		}
	}
}

// Update the game world, deltaTime is time since last update in ms, return total water vol for info text
func (w *World) Update(deltaTime float32) {
	// Track quantities of various substances in the world
	waterLevel := chem.Quantity{Units: chem.Litre}

	if !w.State.Paused() {
		update_callbacks := map[entity.EntityType]func(*entity.Entity){
			entity.Plant:    entity.UpdatePlant,
			entity.Creature: entity.UpdateCreature,
		}

		// Concurrently update plants and creatures
		for _, entity := range w.State.Entities {
			if update, ok := update_callbacks[entity.Type]; ok {
				go update(entity)
				entity.Highlight(w.State.LookingAt == entity)
			}
		}
	}

	// Update the tilemap
	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		for z := 0; z < w.Cfg.Simulation.Depth; z++ {
			tile := w.tilemap[x][z]
			entity.UpdateTile(tile, w.State.Paused())
			tile.Highlight(w.State.LookingAt == tile)
		}
	}

	// Update the atmosphere
	w.atmosphere.Update(deltaTime)

	// Update water level count after updating tiles at atmosphere
	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		for z := 0; z < w.Cfg.Simulation.Depth; z++ {
			tile := w.tilemap[x][z]
			waterLevel.Value += tile.GetWaterLevel().Value
		}
	}

	w.State.Quantities[chem.Water] = waterLevel
}

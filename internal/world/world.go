package world

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/atmosphere"
	"cbeimers113/strands/internal/chem"
	"cbeimers113/strands/internal/context"
	"cbeimers113/strands/internal/entity"
)

type World struct {
	*context.Context

	light      *light.Ambient
	sun        *graphic.Mesh
	tilemap    [][]*entity.Tile
	atmosphere *atmosphere.Atmosphere
}

func New(ctx *context.Context) *World {
	w := &World{
		Context: ctx,

		light:      light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0),
		atmosphere: atmosphere.New(ctx),
	}

	geom := geometry.NewSphere(6, 12, 12)
	mat := material.NewStandard(&math32.Color{
		R: 0.2,
		G: 0.2,
		B: 0.06,
	})
	w.sun = graphic.NewMesh(geom, mat)

	w.Scene.Add(w.light)
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
func (w *World) removeEntity(entity entity.Entity) {
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
func (w *World) entityIndex(entity entity.Entity) int {
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
	w.tilemap = make([][]*entity.Tile, w.Cfg.Simulation.Width)

	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		w.tilemap[x] = make([]*entity.Tile, w.Cfg.Simulation.Depth)

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

			tile.Neighbours = neighbours
		}
	}
}

// updateSun adjusts the sun's light intensity and position based on the internal clock
func (w *World) updateSun() {
	p := w.State.Clock.Progress(w.Cfg.Simulation.DayLength)

	// Update sunlight using a fine tuned sine wave function
	i := 6*math32.Sin(2*math32.Pi*(p-0.25)) + 8
	w.light.SetIntensity(i)

	// Move sun object based on time of day
	ox, oz := float32(w.Cfg.Simulation.Width)/2-6, float32(w.Cfg.Simulation.Depth)/2-6 // Centre of map
	d := 2 * math32.Pi * (p - 0.25)
	dx := math32.Cos(d)
	dx *= float32(w.Cfg.Simulation.Width) * 1.5
	dy := math32.Sin(d)
	dy *= float32(w.Cfg.Simulation.Height) * 1.5
	w.sun.SetPosition(ox+dx, dy, oz)
	w.light.SetPosition(ox+dx, dy, oz)

	// Set opacity of the sun so we can't see it at night and shift into red at dusk/dawn
	// opacity scales up from 4am to 6 am (sunrise)
	// opacity scales down from 6pm to 8pm (sunset)
	var change bool
	var o float32

	// sunrise
	min := float32(4) / 24
	max := float32(6) / 24
	if p >= min && p <= max {
		o = (p - min) / (max - min)
		change = true
	}

	// sunset
	min = float32(18) / 24
	max = float32(20) / 24
	if p >= min && p <= max {
		o = 1 - (p-min)/(max-min)
		change = true
	}

	if change {
		if imat := w.sun.GetMaterial(0); imat != nil {
			if ms, ok := imat.(*material.Standard); ok {
				ms.SetOpacity(o)
				ms.SetColor(&math32.Color{
					R: 0.2,
					G: 0.2 * o,
					B: 0.06 * o,
				})
			}
		}
	}
}

// Update the game world, deltaTime is time since last update in ms, return total water vol for info text
func (w *World) Update(deltaTime float32) {
	// Track quantities of various substances in the world
	waterLevel := chem.Quantity{Units: chem.Litre}

	if !w.State.Paused() {
		w.updateSun()
		w.atmosphere.Update(deltaTime)

		// Concurrently update plants and creatures
		for _, e := range w.State.Entities {
			go e.Update()
			entity.Highlight(e, w.State.LookingAt == e)
		}
	}

	// Update the tilemap independant of paused state so that things like tile highlighting will still work when physics paused
	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		for z := 0; z < w.Cfg.Simulation.Depth; z++ {
			tile := w.tilemap[x][z]
			tile.PausePhysics(w.State.Paused())
			tile.Update()
			entity.Highlight(tile, w.State.LookingAt == tile)
		}
	}

	// Update water level count after updating tiles and atmosphere
	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		for z := 0; z < w.Cfg.Simulation.Depth; z++ {
			tile := w.tilemap[x][z]
			waterLevel.Value += tile.WaterLevel.Value
		}
	}

	w.State.Quantities[chem.Water] = waterLevel
}

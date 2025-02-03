package world

import (
	"math"

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
	"cbeimers113/strands/internal/state"
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

func Load(ctx *context.Context, tiles []*entity.Tile, cells []*state.Cell) *World {
	w := &World{
		Context: ctx,

		light:      light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0),
		atmosphere: atmosphere.Load(ctx, cells),
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

	width := w.Cfg.Simulation.Width
	depth := w.Cfg.Simulation.Depth

	w.State.Quantities[chem.Water] = &chem.Quantity{Units: chem.Litre}
	w.tilemap = make([][]*entity.Tile, width)

	for x := 0; x < width; x++ {
		w.tilemap[x] = make([]*entity.Tile, depth)

		for z := 0; z < depth; z++ {
			tile := tiles[x+z*depth]
			tile.Rand = w.State.Rand
			w.tilemap[x][z] = tile
			tile.Refresh(w.State.Entities, w.Scene)

			w.State.Quantities[chem.Water].Value += tile.WaterLevel.Value
		}
	}

	w.assignTileNeighbourhoods()
	return w
}

// Create the tilemap
func (w *World) createMap() {
	heightmap, min, max := w.makeHeightmap()
	w.makeTilemap(heightmap, min, max)
	w.assignTileNeighbourhoods()
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
	pnoise := perlin.NewPerlin(1, 0.1, 2, w.State.Rand.Int63())

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
	width := w.Cfg.Simulation.Width
	depth := w.Cfg.Simulation.Depth

	w.State.Quantities[chem.Water] = &chem.Quantity{Units: chem.Litre}
	w.tilemap = make([][]*entity.Tile, width)

	for x := 0; x < width; x++ {

		w.tilemap[x] = make([]*entity.Tile, depth)
		for z := 0; z < depth; z++ {
			// Map the heightmap value to the TileTypes array to determine tile type
			height := float32(len(entity.TileTypes)) * (heightmap[x][z] - min) / (max - min)
			tType := entity.Stone

			if int(height) < len(entity.TileTypes) {
				tType = entity.TileTypes[int(height)]
			}

			height /= 3

			// Each tile spawns at 22°C with 10 L of water on top of it
			tile := entity.NewTile(x, z, height, 22.0, 10, tType, w.State.Rand)
			tile.Refresh(w.State.Entities, w.Scene)
			w.tilemap[x][z] = tile

			w.State.Quantities[chem.Water].Value += tile.WaterLevel.Value
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
	srMin := float32(4) / 24
	srMax := float32(6) / 24
	if p >= srMin && p <= srMax {
		o = (p - srMin) / (srMax - srMin)
		change = true
	}

	// sunset
	ssMin := float32(18) / 24
	ssMax := float32(20) / 24
	if p >= ssMin && p <= ssMax {
		o = 1 - (p-ssMin)/(ssMax-ssMin)
		change = true
	}

	// night
	if p < srMin || p > ssMax {
		o = 0
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

// Update the game world, deltaTime is time since last update in ms
func (w *World) Update(deltaTime float32) {
	if !w.State.Paused() {
		w.updateSun()
		w.atmosphere.Update(deltaTime)

		// Update plants and creatures
		for _, e := range w.State.Entities {
			if _, isTile := e.(*entity.Tile); isTile {
				continue
			}

			e.Update()
			entity.Highlight(e, w.State.LookingAt == e)
		}
	}

	// Update the tilemap
	for x := 0; x < w.Cfg.Simulation.Width; x++ {
		for z := 0; z < w.Cfg.Simulation.Depth; z++ {
			tile := w.tilemap[x][z]
			tile.PausePhysics(w.State.Paused())
			tile.Update()
			entity.Highlight(tile, w.State.LookingAt == tile)
		}
	}
}

// GetAtmosphere returns a linear slice of Cells representing the atmosphere
func (w World) GetAtmosphere() []*state.Cell {
	return w.atmosphere.GetCells()
}

// SetAtmosphere sets the cells in the atmosphere from a linear slice of Cells
func (w *World) SetAtmosphere(cells []*state.Cell) {
	w.atmosphere.SetCells(cells)
}

// GetTiles returns a linear slice Tiles representing the map
func (w World) GetTiles() []*entity.Tile {
	width := w.Cfg.Simulation.Width
	depth := w.Cfg.Simulation.Depth
	t := make([]*entity.Tile, width*depth)

	for x := 0; x < width; x++ {
		for z := 0; z < depth; z++ {
			t[x+z*depth] = w.tilemap[x][z]
		}
	}

	return t
}

// GetTile returns the tile at x, z
func (w World) GetTile(x, z int) *entity.Tile {
	return w.tilemap[x][z]
}

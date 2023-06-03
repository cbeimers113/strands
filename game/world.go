package game

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
)

// World size in tiles
const Width int = 64
const Height int = 64
const Depth int = 64

var Sun *light.Ambient
var Entities map[int]*Entity

// Remove an entity from the world
func RemoveEntity(entity *Entity) {
	index := EntityIndex(entity)

	// If the entity is in the Entities list, remove it and shift all the entities above it down the list
	if index >= 0 {
		for i := index + 1; i < len(Entities); i++ {
			Entities[i-1] = Entities[i]
			Entities[i-1].SetName(fmt.Sprintf("%d", i-1))
		}
		delete(Entities, len(Entities)-1)
	}

	entity.DisposeChildren(true)
	entity.Dispose()
}

// Find an entity's index in the Entities list, return -1 if not found
func EntityIndex(entity *Entity) int {
	for i, ent := range Entities {
		if ent == entity {
			return i
		}
	}

	return -1
}

// Check if a given coordinate is within the tilemap boundaries
func InBounds(x, z int) (inBounds bool) {
	inBounds = x >= 0 && x < Width && z >= 0 && z < Depth

	return
}

// Generate a heightmap, return the map and its min and max values
func makeHeightmap() ([Width][Depth]float32, float32, float32) {
	var heightmap [Width][Depth]float32
	var min float32 = 1_000_000_000.0
	var max float32 = -min
	pnoise := perlin.NewPerlin(1, 0.1, 2, rand.Int63())

	for x := 0; x < Width; x++ {
		for z := 0; z < Depth; z++ {
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
func makeTilemap(heightmap [Width][Depth]float32, min, max float32) [Width][Depth]*Entity {
	var tilemap [Width][Depth]*Entity

	for x := 0; x < Width; x++ {
		for z := 0; z < Depth; z++ {
			// Map the heightmap value to the TileTypes array to determine tile type
			height := math32.Min(float32(len(TileTypes))*(heightmap[x][z]-min)/(max-min), float32(len(TileTypes)-1))
			tType := TileTypes[int(height)]
			height /= 3

			// Each tile spawns at 22°C with 1 L of water on top of it
			tile := NewTile(x, z, height, 22.0, 1.0, tType)
			Scene.Add(tile.GetINode())
			tilemap[x][z] = tile
		}
	}

	return tilemap
}

// Give each tile in a tilemap a list of pointers to its neighbours
func assignTileNeighbourhoods(tilemap [Width][Depth]*Entity) {
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
	for x := 0; x < Width; x++ {
		for z := 0; z < Depth; z++ {
			var neighbours Neighbourhood
			tile := tilemap[x][z]

			for i, offs := range nbOffsets {
				xOffs := x + offs[0]
				zOffs := z + offs[1]

				// Stagger offsets on the x axis for every other row for "top/bottom" neighbours
				if z%2 == 0 && i%3 != 0 {
					xOffs--
				}

				if InBounds(xOffs, zOffs) {
					neighbours[i] = tilemap[xOffs][zOffs]
				}
			}

			// Load the neighbour pointers into the tile's metadata
			data, _ := tile.UserData().(*TileData)
			data.Neighbours = neighbours
		}
	}
}

// Create the tilemap
func CreateMap() {
	heightmap, min, max := makeHeightmap()
	tilemap := makeTilemap(heightmap, min, max)
	assignTileNeighbourhoods(tilemap)
}

// Load the world into the scene
func LoadWorld() {
	// Sun TODO: Sun entity type for day/night cycle
	Sun = light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 8.0)
	Scene.Add(Sun)
	Entities = make(map[int]*Entity)

	CreateMap()
	CreateAtmosphere()
}

// Update the game world, deltaTime is time since last update in ms
func UpdateWorld(deltaTime float32) {
	update_callbacks := map[EntityType]func(*Entity){
		Plant:    UpdatePlant,
		Tile:     UpdateTile,
		Creature: UpdateCreature,
	}

	for _, entity := range Entities {
		if update, ok := update_callbacks[entity.Type]; ok {
			go update(entity)
			entity.Highlight(LookingAt == entity)
		}
	}
}

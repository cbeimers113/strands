package game

import (
	"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

var Sun *light.Ambient
var Entities map[string]Entity

// Add a plant with given genetics to a given tile, return whether the plant was added
func AddPlant(colour int, tile Entity) (success bool) {
	success = true

	// Check if there is already a plant here
	for _, child := range tile.Children() {
		if Type(child.GetNode()) == Plant {
			success = false
			break
		}
	}

	if success {
		geom := geometry.NewCylinder(float64(TileSize)/8, float64(TileSize)/2, 8, 8, true, true)
		mat := material.NewStandard(math32.NewColorHex(uint(colour) / 10))
		plant := graphic.NewMesh(geom, mat)
		tex, ok := Texture("stalk")

		if ok {
			mat.AddTexture(tex)
		}

		plant.SetPosition(0, plant.Scale().Y/2, 0)
		plant.SetName(Plant)
		plant.SetUserData(Strand{colour, 0})
		tile.Add(plant)
		Entities[tile.Name()] = tile.ChildAt(len(tile.Children()) - 1).GetNode()
	}

	return
}

// Spawn a hex tile of type tType at x, y
func CreateTile(x, y int, tType string) {
	geom := CreateHexagon(TileSize)
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	tile := graphic.NewMesh(geom, mat)
	tex, ok := Texture(tType)
	posX := (float32(x)+(0.5*float32(y%2))) * TileSize * math32.Sin(math32.Pi/3)
	posZ := float32(y) * TileSize * 0.75

	if ok {
		mat.AddTexture(tex)
	}

	tile.SetPosition(posX, 0, posZ)
	tile.SetRotationY(math32.Pi / 2)
	tile.SetName(Tile)
	tile.SetUserData(Strand{x, y, 0})
	Scene.Add(tile)
	Entities[tile.Name()] = Scene.ChildAt(len(Scene.Children()) - 1).GetNode()
}

// Load the world into the scene
func LoadWorld() {
	// Sun
	Sun = light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 8.0)
	Scene.Add(Sun)
	Entities = make(map[string]Entity)

	// Tiles
	pnoise := perlin.NewPerlin(1, 0.1, 2, rand.Int63())
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			// Construct the tile at x, y
			noise := int(math.Abs(pnoise.Noise2D(float64(x), float64(y))*100)) % len(TypeStrata)
			tType := TypeStrata[noise]
			CreateTile(x, y, tType)
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

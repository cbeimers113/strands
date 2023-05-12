package game

import (
	"fmt"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Tile metadata mapping
const TileX int = 0
const TileY int = 1
const TileT int = 2

const Water string = "water"
const Dirt string = "dirt"
const Grass string = "grass"
const Stone string = "stone"

// Store list of tile types ordered by spawn height
var TileTypes []string = []string{
	Water,
	Dirt,
	Grass,
	Stone,
}

// Perform an action on a tile entity on right click
func OnRightClickTile(tile Entity) {
	AddPlantTo(tile, 0x00dd05)
}

// Perform an action on a tile entity on left click
func OnLeftClickTile(tile Entity) {
	println("No left click behaviour defined for ", tile.Name())
}

// Spawn a hex tile of type tType at x, y
func NewTile(x, y int, tType string) (tile *graphic.Mesh) {
	geom := CreateHexagon(TileSize)
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	tile = graphic.NewMesh(geom, mat)
	tex, ok := Texture(tType)
	posX := (float32(x) + (0.5 * float32(y%2))) * TileSize * math32.Sin(math32.Pi/3)
	posZ := float32(y) * TileSize * 0.75

	if ok {
		mat.AddTexture(tex)
	}

	tile.SetPosition(posX, 0, posZ)
	tile.SetRotationY(math32.Pi / 2)
	tile.SetName(fmt.Sprintf("%s (%s)", Tile, tType))
	tile.SetUserData(Strand{x, y, TypeIndex(tType)})

	return
}

// Check what index of the tile types strata a type is, return -1 if invalid type
func TypeIndex(tType string) int {
	for i, t := range TileTypes {
		if t == tType {
			return i
		}
	}

	return -1
}

// Check if there's a plant on this tile
func HasPlant(tile Entity) bool {
	for _, child := range tile.Children() {
		if Type(child.GetNode()) == Plant {
			return true
		}
	}

	return false
}

// Add a plant with given genetics to a given tile, return whether the plant was added
func AddPlantTo(tile Entity, colour int) (success bool) {
	success = !HasPlant(tile)

	if success {
		AddEntityTo(tile, NewPlant(colour))
	}

	return
}

// TODO: Remove a plant from a tile
// func RemovePlant(tile Entity, plant Plant) (success bool) { return false }

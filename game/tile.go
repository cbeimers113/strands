package game

import (
	"fmt"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Tile metadata mapping
const TileX int = 0  // Position in the tilemap, not the game world
const TilyY int = 1  // These x and y values are used to access the tile's neighbourhood (TODO)
const TileT int = 2

type TileType = string

const Water TileType = "water"
const Dirt TileType = "dirt"
const Grass TileType = "grass"
const Stone TileType = "stone"

// Store list of tile types ordered by spawn height
var TileTypes []TileType = []TileType{
	Water,
	Dirt,
	Grass,
	Stone,
}

// Perform an action on a tile entity on right click
func OnRightClickTile(tile Entity) {
	AddEntityTo(tile, NewPlant(0x00dd05))
}

// Perform an action on a tile entity on left click
func OnLeftClickTile(tile Entity) {
	println("No left click behaviour defined for ", tile.Name())
}

// Spawn a hex tile of type tType at x, z (tile precision) at height y (game world precision)
func NewTile(x, z int, y float32, tType TileType) (tile *graphic.Mesh) {
	geom := CreateHexagon(TileSize)
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	tile = graphic.NewMesh(geom, mat)
	tex, ok := Texture(tType)
	posX := (float32(x) + (0.5 * float32(z%2))) * TileSize * math32.Sin(math32.Pi/3)
	posZ := float32(z) * TileSize * 0.75

	if ok {
		mat.AddTexture(tex)
	}

	tile.SetPosition(posX, 0, posZ)
	tile.SetRotationY(math32.Pi / 2)
	tile.SetName(fmt.Sprintf("%s (%s)", Tile, tType))
	tile.SetUserData(Strand{x, z, TypeIndex(tType)})

	return
}

// Check what index of the tile types strata a type is, return -1 if invalid type
func TypeIndex(tType TileType) int {
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

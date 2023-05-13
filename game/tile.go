package game

import (
	"fmt"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

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

// Which data a tile will store
type TileData struct {
	MapX int
	MapZ int
	// No need for MapY as tiles are not stacked
	// World (x, y, z) stored in tile.Position()
	Type TileType
}

// Perform an action on a tile entity on right click
func OnRightClickTile(tile Entity) {
	if !HasPlant(tile) {
		AddEntityTo(tile, NewPlant(0x00dd05))
	}
}

// Perform an action on a tile entity on left click
func OnLeftClickTile(tile Entity) {
	println("No left click behaviour defined for ", tile.Name())
}

// Spawn a hex tile of type tType at mapX, mapZ (tile precision), y (game world precision)
func NewTile(mapX, mapZ int, y float32, tType TileType) (tile *graphic.Mesh) {
	geom := CreateHexagon(TileSize)
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	tile = graphic.NewMesh(geom, mat)
	tex, ok := Texture(tType)
	x := (float32(mapX) + (0.5 * float32(mapZ%2))) * TileSize * math32.Sin(math32.Pi/3)
	z := float32(mapZ) * TileSize * 0.75

	if ok {
		mat.AddTexture(tex)
	}

	tile.SetPosition(x, y, z)
	tile.SetRotationY(math32.Pi / 2)
	tile.SetName(fmt.Sprintf("%s (%s)", Tile, tType))
	tile.SetUserData(TileData{MapX: mapX, MapZ: mapZ, Type: tType})

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

// Perform per-frame updates to a Tile
func UpdateTile(tile Entity) {

}

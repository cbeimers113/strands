package game

import (
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type TileType struct {
	Name    string
	Fertile bool
}

var Sand TileType = TileType{Name: "sand", Fertile: false}
var Dirt TileType = TileType{Name: "dirt", Fertile: true}
var Grass TileType = TileType{Name: "grass", Fertile: true}
var Stone TileType = TileType{Name: "stone", Fertile: false}

// Store list of tile types ordered by spawn height
var TileTypes []TileType = []TileType{
	Sand,
	Dirt,
	Grass,
	Stone,
}

// Represents the tiles surrounding this one
type Neighbourhood = [6]*Entity

// Which data a tile will store
type TileData struct {
	// Static properties
	MapX int
	MapZ int
	// World (x, y, z) stored in tile.Position()
	Type       TileType
	Neighbours Neighbourhood // Pointers to any neighbouring tiles

	// Dynamic properties
	Planted     bool
	Temperature float32
	Moisture    float32
}

// Perform an action on a tile entity on right click
func OnRightClickTile(tile *Entity) {
	// Try to plant a plant here
	if tileData, ok := tile.UserData().(TileData); ok && !tileData.Planted && tileData.Type.Fertile {
		tile.Add(NewRandomPlant().GetINode())
		tileData.Planted = true
		tile.SetUserData(tileData)
	}
}

// Perform an action on a tile entity on left click
func OnLeftClickTile(tile *Entity) {
	println("No left click behaviour defined for tile")
}

// Spawn a hex tile of type tType at mapX, mapZ (tile precision), y (game world precision)
func NewTile(mapX, mapZ int, y float32, tType TileType) (tile *Entity) {
	geom := CreateHexagon(TileSize, y)
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	mesh := graphic.NewMesh(geom, mat)
	x := (float32(mapX) + (0.5 * float32(mapZ%2))) * TileSize * math32.Sin(math32.Pi/3)
	z := float32(mapZ) * TileSize * 0.75

	if tex, ok := Texture(tType.Name); ok {
		mat.AddTexture(tex)
	}

	tile = NewEntity(mesh, Tile)
	tile.SetPosition(x, y, z)
	tile.SetRotationY(math32.Pi / 2)
	tile.SetUserData(TileData{MapX: mapX, MapZ: mapZ, Type: tType, Temperature: 22.0})

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

// Perform per-frame updates to a Tile
func UpdateTile(tile *Entity) {

}

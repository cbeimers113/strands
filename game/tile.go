package game

import (
	"math/rand"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

type TileType struct {
	Name      string
	Fertility float32
}

var Sand TileType = TileType{Name: "sand", Fertility: 0.05}
var Dirt TileType = TileType{Name: "dirt", Fertility: 0.33}
var Grass TileType = TileType{Name: "grass", Fertility: 0.80}
var Stone TileType = TileType{Name: "stone", Fertility: 0.00}

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
	Type TileType

	Neighbours Neighbourhood // Pointers to any neighbouring tiles
	Water      *graphic.Mesh // This is the water tile object that can exist on top of the tile

	// Dynamic properties
	Planted     bool
	Temperature *Quantity
	WaterLevel  *Quantity
}

// Perform an action on a tile entity on right click
func OnRightClickTile(tile *Entity) {
	// Try to plant a plant here
	r := rand.Float32()

	if tileData, ok := tile.UserData().(*TileData); ok && !tileData.Planted && r < tileData.Type.Fertility {
		tile.Add(NewRandomPlant().GetINode())
		tileData.Planted = true
	}
}

// Perform an action on a tile entity on left click
func OnLeftClickTile(tile *Entity) {
	println("No left click behaviour defined for tile")
}

// Spawn a hex tile of type tType at mapX, mapZ (tile precision), y (game world precision)
func NewTile(mapX, mapZ int, height, temp, waterLevel float32, tType TileType) (tile *Entity) {
	tileMesh := createTileMesh(tType.Name)
	waterMesh := createTileMesh("water")
	tileMesh.Add(waterMesh)

	x := (float32(mapX) + (0.5 * float32(mapZ%2))) * math32.Sin(math32.Pi/3)
	z := float32(mapZ) * 0.75

	tile = NewEntity(tileMesh, Tile)
	tile.SetPosition(x, height, z)
	tile.SetRotationY(math32.Pi / 2)
	tile.SetUserData(&TileData{
		MapX:        mapX,
		MapZ:        mapZ,
		Type:        tType,
		Water:       waterMesh,
		Temperature: &Quantity{temp, Celcius},
		WaterLevel:  &Quantity{waterLevel, Litre},
	})

	return
}

// Create a base tile mesh with a given texture
func createTileMesh(texture string) (tileMesh *graphic.Mesh) {
	geom := CreateHexagon()
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	tileMesh = graphic.NewMesh(geom, mat)

	if tex, ok := Texture(texture); ok {
		mat.AddTexture(tex)
	}

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

// Update the tile's water level
func (tile *Entity) updateWaterLevel(tileData *TileData) {
	waterLevel := &tileData.WaterLevel.Value
	water := tileData.Water

	// TODO: adjust water each frame

	water.SetVisible(*waterLevel > 0)
	water.SetScaleY(LitresToCubicMetres(*waterLevel))
	water.SetPositionY(DimensionsOf(water).Y)
}

// Perform per-frame updates to a Tile
func UpdateTile(tile *Entity) {
	if tileData, ok := tile.UserData().(*TileData); ok {
		tile.updateWaterLevel(tileData)
	}
}

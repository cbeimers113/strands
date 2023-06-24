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
	WaterTick  int

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
	tile.AddWater(CubicMetresToLitres(1))
}

// Spawn a hex tile of type tType at mapX, mapZ (tile precision), y (game world precision)
func NewTile(mapX, mapZ int, height, temp, waterLevel float32, tType TileType) (tile *Entity) {
	x := (float32(mapX) + (0.5 * float32(mapZ%2))) * math32.Sin(math32.Pi/3)
	z := float32(mapZ) * 0.75
	tileMesh := createTileMesh(tType.Name)
	tile = NewEntity(tileMesh, Tile)
	waterMesh := createTileMesh("water")

	TotalWaterVolume.Value += waterLevel
	waterMesh.SetName(tile.Name())
	tileMesh.Add(waterMesh)
	tile.SetPosition(x, height, z)
	tile.SetRotationY(math32.Pi / 2)
	tile.SetUserData(&TileData{
		MapX:        mapX,
		MapZ:        mapZ,
		Type:        tType,
		Water:       waterMesh,
		WaterTick:   rand.Intn(100),
		Temperature: &Quantity{temp, Celcius},
		WaterLevel:  &Quantity{waterLevel, Litre},
	})

	return
}

// Create a base tile mesh with a given texture
func createTileMesh(texture string) (tileMesh *graphic.Mesh) {
	geom := NewHexMesh()
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	mat.SetTransparent(texture == "water")
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
func (tile *Entity) updateWaterLevel() {
	if tileData, ok := tile.UserData().(*TileData); ok {
		waterLevel := tileData.WaterLevel.Value
		water := tileData.Water

		water.SetScaleY(LitresToCubicMetres(waterLevel))
		water.SetPositionY(DimensionsOf(water).Y)
		water.SetVisible(tileData.WaterLevel.Value > 0)

		// Lower water texture opacity for low water level
		if imat := water.GetMaterial(0); imat != nil {
			if ms, ok := imat.(*material.Standard); ok {
				ms.SetOpacity(math32.Min(50, waterLevel) / 60)
			}
		}
	}
}

// Add an amount of water to a tile. Add a negative amount to remove water.
// If the new volume is negative, return the amount that was lost
func (tile *Entity) AddWater(delta float32) float32 {
	if tileData, ok := tile.UserData().(*TileData); ok {
		waterLevel := &tileData.WaterLevel.Value
		backflow := -(*waterLevel + delta)
		pLevel := *waterLevel
		*waterLevel += delta
		*waterLevel = math32.Max(0, *waterLevel)
		TotalWaterVolume.Value += *waterLevel - pLevel

		return math32.Max(backflow, 0)
	}
	return 0
}

// Get the elevation of the top of the tile, including its water
func (tile *Entity) GetElevation() *Quantity {
	var elevation float32

	if tileData, ok := tile.UserData().(*TileData); ok {
		elevation = tile.Position().Y
		elevation += tileData.Water.Position().Y
		elevation += LitresToCubicMetres(tileData.WaterLevel.Value)
	}

	return &Quantity{
		Value: elevation,
		Units: Metre,
	}
}

// Perform per-frame updates to a Tile
func UpdateTile(tile *Entity) {
	if !IsPaused {
		tile.doWaterSpread()
	}

	tile.updateWaterLevel()
}

// Spread water to other tiles
func (tile *Entity) doWaterSpread() {
	if tileData, ok := tile.UserData().(*TileData); ok {
		tileData.WaterTick++

		if tileData.WaterTick > 10 && tileData.WaterLevel.Value > 0 {
			var lowerNeighbours []*Entity

			// Filter out the neighbouring tiles which aren't lower than this one
			for _, neighbour := range tileData.Neighbours {
				if neighbour == nil {
					continue
				}

				if neighbour.GetElevation().Value < tile.GetElevation().Value {
					lowerNeighbours = append(lowerNeighbours, neighbour)
				}
			}

			// Shuffle the lower neighbours to give some randomness to flow direction
			rand.Shuffle(len(lowerNeighbours), func(i, j int) {
				lowerNeighbours[i], lowerNeighbours[j] = lowerNeighbours[j], lowerNeighbours[i]
			})

			// Distribute water to lower neighbours
			for len(lowerNeighbours) > 0 {
				var neighbour *Entity

				// Find lowest neighbour
				for _, nb := range lowerNeighbours {
					if neighbour == nil || nb.GetElevation().Value < neighbour.GetElevation().Value {
						neighbour = nb
					}
				}

				// Stop when this tile is already a local minimum
				if neighbour.GetElevation().Value >= tile.GetElevation().Value {
					break
				}

				delta := CubicMetresToLitres(tile.GetElevation().Value-neighbour.GetElevation().Value) / float32(len(lowerNeighbours))

				if δ := neighbour.AddWater(delta - tile.AddWater(-delta)); δ != 0 {
					// TODO: This shouldn't ever be 0, but do something if it is
					println(δ)
				}

				// Remove neighbour from lowerNeighbours
				for i, nb := range lowerNeighbours {
					if nb == neighbour {
						lowerNeighbours[i] = lowerNeighbours[len(lowerNeighbours)-1]
						lowerNeighbours = lowerNeighbours[:len(lowerNeighbours)-1]
						break
					}
				}

				// Stop when we run out of water to spread
				if tileData.WaterLevel.Value == 0 {
					break
				}
			}

			tileData.WaterTick = 0
		}
	}
}

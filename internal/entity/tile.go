package entity

import (
	"math/rand"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/chem"
	"cbeimers113/strands/internal/graphics"
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
	Temperature *chem.Quantity
	WaterLevel  *chem.Quantity
}

// Spawn a hex tile of type tType at mapX, mapZ (tile precision), y (game world precision)
func NewTile(entities map[int]*Entity, mapX, mapZ int, height, temp, waterLevel float32, tType TileType) (tile *Entity) {
	x := (float32(mapX) + (0.5 * float32(mapZ%2))) * math32.Sin(math32.Pi/3)
	z := float32(mapZ) * 0.75
	tileMesh := createTileMesh(tType.Name)
	tile = New(tileMesh, Tile, entities)
	waterMesh := createTileMesh("water")

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
		Temperature: &chem.Quantity{Value: temp, Units: chem.Celcius},
		WaterLevel:  &chem.Quantity{Value: waterLevel, Units: chem.Litre},
	})

	return
}

// Create a base tile mesh with a given texture
func createTileMesh(texture string) (tileMesh *graphic.Mesh) {
	geom := graphics.NewHexMesh()
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	mat.SetTransparent(texture == "water")
	tileMesh = graphic.NewMesh(geom, mat)

	if tex, ok := graphics.Texture(texture); ok {
		mat.AddTexture(tex)
	}

	return
}

// Get the elevation of the top of the tile, including its water
func (t *Entity) getElevation() *chem.Quantity {
	var elevation float32

	if tileData, ok := t.UserData().(*TileData); ok {
		elevation = t.Position().Y
		elevation += tileData.Water.Position().Y
		elevation += chem.LitresToCubicMetres(tileData.WaterLevel.Value)
	}

	return &chem.Quantity{
		Value: elevation,
		Units: chem.Metre,
	}
}

// Spread water to other tiles
func (t *Entity) doWaterSpread() {
	if tileData, ok := t.UserData().(*TileData); ok {
		tileData.WaterTick++

		if tileData.WaterTick > 10 && tileData.WaterLevel.Value > 0 {
			var lowerNeighbours []*Entity

			// Filter out the neighbouring tiles which aren't lower than this one
			for _, neighbour := range tileData.Neighbours {
				if neighbour == nil {
					continue
				}

				if neighbour.getElevation().Value < t.getElevation().Value {
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
					if neighbour == nil || nb.getElevation().Value < neighbour.getElevation().Value {
						neighbour = nb
					}
				}

				// Stop when this tile is already a local minimum
				if neighbour.getElevation().Value >= t.getElevation().Value {
					break
				}

				delta := chem.CubicMetresToLitres(t.getElevation().Value-neighbour.getElevation().Value) / float32(len(lowerNeighbours))

				if δ := neighbour.AddWater(delta - t.AddWater(-delta)); δ != 0 {
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

// Update the tile's water level
func (t *Entity) updateWaterLevel() {
	if tileData, ok := t.UserData().(*TileData); ok {
		waterLevel := tileData.WaterLevel.Value
		water := tileData.Water

		water.SetScaleY(chem.LitresToCubicMetres(waterLevel))
		water.SetPositionY(graphics.DimensionsOf(water).Y)
		water.SetVisible(tileData.WaterLevel.Value > 0)

		// Lower water texture opacity for low water level
		if imat := water.GetMaterial(0); imat != nil {
			if ms, ok := imat.(*material.Standard); ok {
				ms.SetOpacity(math32.Min(50, waterLevel) / 60)
			}
		}
	}
}

// Perform per-frame updates to a Tile
func UpdateTile(t *Entity, paused bool) {
	if !paused {
		t.doWaterSpread()
	}

	t.updateWaterLevel()
}

// Get the volume of water on a tile
func (t *Entity) GetWaterLevel() *chem.Quantity {
	if tileData, ok := t.UserData().(*TileData); ok {
		return tileData.WaterLevel
	}

	return nil
}

// Add an amount of water to a tile. Add a negative amount to remove water.
// If the new volume is negative, return the amount that was lost
func (t *Entity) AddWater(delta float32) float32 {
	if tileData, ok := t.UserData().(*TileData); ok {
		waterLevel := &tileData.WaterLevel.Value
		backflow := -(*waterLevel + delta)
		*waterLevel += delta
		*waterLevel = math32.Max(0, *waterLevel)

		return math32.Max(backflow, 0)
	}
	return 0
}

// Add a plant to a tile if plantable
func (t *Entity) AddPlant(entities map[int]*Entity) {
	if tileData, ok := t.UserData().(*TileData); ok && !tileData.Planted && tileData.Type.Fertility > 0 {
		t.Add(NewRandomPlant(entities).GetINode())
		tileData.Planted = true
	}
}

package entity

import (
	"fmt"
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
type Neighbourhood = [6]*Tile

// Which data a tile will store
type Tile struct {
	*graphic.Mesh

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

	physicsEnabled bool
}

// Spawn a hex tile of type tType at mapX, mapZ (tile precision), y (game world precision)
func NewTile(entities map[int]Entity, mapX, mapZ int, height, temp, waterLevel float32, tType TileType) *Tile {
	x := (float32(mapX) + (0.5 * float32(mapZ%2))) * math32.Sin(math32.Pi/3)
	z := float32(mapZ) * 0.75
	tileMesh := createTileMesh(tType.Name)
	waterMesh := createTileMesh("water")

	tile := &Tile{
		Mesh: tileMesh,

		MapX: mapX,
		MapZ: mapZ,
		Type: tType,

		Water:     waterMesh,
		WaterTick: rand.Intn(100),

		Temperature: &chem.Quantity{Value: temp, Units: chem.Celcius},
		WaterLevel:  &chem.Quantity{Value: waterLevel, Units: chem.Litre},
	}

	waterMesh.SetName(tile.Name())
	tile.Add(waterMesh)
	tile.SetPosition(x, height, z)
	tile.SetRotationY(math32.Pi / 2)
	tile.GetMaterial(0).GetMaterial().SetLineWidth(8)

	// Add this tile to the entities list
	AddEntity(tile, entities)

	return tile
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
func (t *Tile) getElevation() *chem.Quantity {
	elevation := t.Position().Y
	elevation += t.Water.Position().Y
	elevation += chem.LitresToCubicMetres(t.WaterLevel.Value)

	return &chem.Quantity{
		Value: elevation,
		Units: chem.Metre,
	}
}

// Spread water to other tiles
func (t *Tile) doWaterSpread() {
	t.WaterTick++

	if t.WaterTick > 10 && t.WaterLevel.Value > 0 {
		var lowerNeighbours []*Tile

		// Filter out the neighbouring tiles which aren't lower than this one
		for _, neighbour := range t.Neighbours {
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
			var neighbour *Tile

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
			if t.WaterLevel.Value == 0 {
				break
			}
		}

		t.WaterTick = 0
	}
}

// Update the tile's water level
func (t *Tile) updateWaterLevel() {
	waterLevel := t.WaterLevel.Value
	water := t.Water

	water.SetScaleY(chem.LitresToCubicMetres(waterLevel))
	water.SetPositionY(graphics.DimensionsOf(water).Y)
	water.SetVisible(t.WaterLevel.Value > 0)

	// Lower water texture opacity for low water level
	if imat := water.GetMaterial(0); imat != nil {
		if ms, ok := imat.(*material.Standard); ok {
			ms.SetOpacity(math32.Min(50, waterLevel) / 60)
		}
	}
}

// Perform per-frame updates to a Tile
func (t *Tile) Update() {
	if t.physicsEnabled {
		t.doWaterSpread()
	}

	t.updateWaterLevel()
}

// Pause or unpause physics
func (t *Tile) PausePhysics(pause bool) {
	t.physicsEnabled = !pause
}

// Add an amount of water to a tile. Add a negative amount to remove water.
// If the new volume is negative, return the amount that was lost
func (t *Tile) AddWater(delta float32) float32 {
	waterLevel := &t.WaterLevel.Value
	backflow := -(*waterLevel + delta)
	*waterLevel += delta
	*waterLevel = math32.Max(0, *waterLevel)

	return math32.Max(backflow, 0)
}

// Add a plant to a tile if plantable
func (t *Tile) AddPlant(entities map[int]Entity) {
	if !t.Planted && t.Type.Fertility > 0 {
		t.Add(NewRandomPlant(entities).GetINode())
		t.Planted = true
	}
}

// Infostring returns a string representation of the tile
func (t Tile) InfoString() string {
	infoString := "Tile:\n"
	infoString += fmt.Sprintf("type=%s\n", t.Type.Name)
	infoString += fmt.Sprintf("temperature=%s\n", t.Temperature)
	infoString += fmt.Sprintf("water level=%s\n", t.WaterLevel)
	infoString += fmt.Sprintf("elevation=%s\n", t.getElevation())
	infoString += fmt.Sprintf("planted=%t\n", t.Planted)

	return infoString
}

// Material returns the tile's material
func (t Tile) Material() *material.Material {
	return t.GetMaterial(0).GetMaterial()
}

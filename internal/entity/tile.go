package entity

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/g3n/engine/core"
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

// The tile entity
type Tile struct {
	*graphic.Mesh `json:"-"`
	*rand.Rand    `json:"-"`

	// Static properties
	MapX   int      `json:"map_x"`
	WorldY float32  `json:"world_y"`
	MapZ   int      `json:"map_z"`
	Type   TileType `json:"type"`

	Neighbours Neighbourhood `json:"-"` // Pointers to any neighbouring tiles
	Water      *graphic.Mesh `json:"-"` // This is the water tile object that can exist on top of the tile
	WaterTick  int           `json:"water_tick"`

	// Dynamic properties
	Plants      []*Plant       `json:"plants"`
	Temperature *chem.Quantity `json:"temperature"`
	WaterLevel  *chem.Quantity `json:"water_level"`

	physicsEnabled bool `json:"-"`
}

// Spawn a hex tile of type tType at mapX, mapZ (tile precision), worldY (game world precision)
func NewTile(mapX, mapZ int, worldY, temp, waterLevel float32, tType TileType, rng *rand.Rand) *Tile {
	tile := &Tile{
		Rand: rng,

		MapX:   mapX,
		WorldY: worldY,
		MapZ:   mapZ,
		Type:   tType,

		WaterTick: rng.Intn(100),

		Temperature: &chem.Quantity{Value: temp, Units: chem.Celcius},
		WaterLevel:  &chem.Quantity{Value: waterLevel, Units: chem.Litre},
	}

	return tile
}

// Create a base tile mesh with a given texture
func createTileMesh(texture string) (tileMesh *graphic.Mesh) {
	geom := graphics.NewHexMesh()
	mat := material.NewStandard(math32.NewColorHex(0x111111))
	mat.SetTransparent(texture == graphics.TexWater)
	tileMesh = graphic.NewMesh(geom, mat)

	if tex, err := graphics.Texture(texture); err == nil {
		mat.AddTexture(tex)
	} else {
		fmt.Println(err)
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
		t.Rand.Shuffle(len(lowerNeighbours), func(i, j int) {
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

	if water == nil {
		return
	}

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

// Refresh the game object(s)
func (t *Tile) Refresh(entities map[int]Entity, scene *core.Node) {
	x := (float32(t.MapX) + (0.5 * float32(t.MapZ%2))) * math32.Sin(math32.Pi/3)
	z := float32(t.MapZ) * 0.75

	// Check if meshes have been created
	if t.Mesh == nil {
		t.Mesh = createTileMesh(t.Type.Name)
		AddEntity(t, entities, scene)
	} else {
		if tex, err := graphics.Texture(t.Type.Name); tex != nil {
			mat := t.GetMaterial(0).GetMaterial()

			// Remove any existing textures before adding the one that matches the tile type
			for texName := range graphics.Textures {
				if oTex, err := graphics.Texture(texName); err == nil {
					if mat.HasTexture(oTex) {
						mat.RemoveTexture(oTex)
					}
				}
			}

			mat.AddTexture(tex)
		} else {
			fmt.Printf("Couldn't get tile texture for %s: %s\n", t.Type.Name, err)
		}

		// Make sure the entities map is pointing to this tile at the specified index
		if i, err := strconv.Atoi(t.Name()); err == nil {
			entities[i] = t
		}
	}

	if t.Water == nil {
		t.Water = createTileMesh(graphics.TexWater)
		t.Water.SetName(t.Name())
		t.Add(t.Water)
	}

	t.SetPosition(x, t.WorldY, z)
	t.SetRotationY(math32.Pi / 2)
	t.GetMaterial(0).GetMaterial().SetLineWidth(8)

	// Refresh all plants
	for _, plant := range t.Plants {
		plant.Rand = t.Rand
		plant.Refresh(entities, scene)

		has := false
		for _, child := range t.Children() {
			if child == plant.GetINode() {
				has = true
			}
		}

		if !has {
			t.Add(plant)
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

// Add a plant to a tile if plantable, returns whether it was planted
func (t *Tile) AddPlant(entities map[int]Entity, scene *core.Node) bool {
	if t.Type.Fertility > 0 {
		plant := NewRandomPlant(entities, t.Rand)
		plant.Refresh(entities, scene)
		t.Add(plant)
		t.Plants = append(t.Plants, plant)
		return true
	}

	return false
}

// Infostring returns a string representation of the tile
func (t Tile) InfoString() string {
	return fmt.Sprintf(
		"%s, : %s,  : %s,  : %s,  : %d", 
		t.Type.Name,
		t.Temperature,
		t.WaterLevel,
		t.getElevation(),
		len(t.Plants),
	)
}

// Material returns the tile's material
func (t Tile) Material() *material.Material {
	if mat0 := t.GetMaterial(0); mat0 != nil {
		return mat0.GetMaterial()
	}
	return nil
}

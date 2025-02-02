package entity

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/graphics"
)

// The Plant entity
type Plant struct {
	*graphic.Mesh `json:"-"`
	*rand.Rand    `json:"-"`

	// Whole Plant
	Age    int     `json:"age"`
	Colour int     `json:"colour"`
	Height float32 `json:"height"`
	Radius float32 `json:"radius"`
	X      float32 `json:"x"`
	Z      float32 `json:"y"`
	RotX   float32 `json:"rot_x"`
	RotY   float32 `json:"rot_y"`

	// Leaves
	NumLeaves        int     `json:"num_leaves"`         // (0, ..), How many leaves the plant has. (More and bigger leaves consume more resources)
	LeafSpawnHeight  float32 `json:"leaf_spawn_height"`  // [0.1, 1], How far up the stem the leaves will appear (from top), ie: value of 0.25 means leaves will spawn on the top quarter of the stem
	AvgLeafSize      float32 `json:"avg_leaf_size"`      // (0, ..), Average size of leaf, ie: value of 1 equals the default size
	LeafSizeVariance float32 `json:"leaf_size_variance"` // [0, 1], How much the leaf sizes can vary, ie: a value of 0.5 means the leaves can be up to 50% bigger or smaller than AvgSize
	Leaves           []*graphic.Mesh
}

// Create a new plant
func NewPlant(entities map[int]Entity, colour, numLeaves int, height, radius, x, z, rotX, rotY float32, rng *rand.Rand) *Plant {
	plant := &Plant{
		Rand: rng,

		Colour:    colour,
		Height:    height,
		Radius:    radius,
		X:         x,
		Z:         z,
		RotX:      rotX,
		RotY:      rotY,
		NumLeaves: numLeaves,

		// Hard-coded starting values for leaf data
		LeafSpawnHeight:  0.5,
		AvgLeafSize:      1,
		LeafSizeVariance: 0.1,
	}

	return plant
}

// Create a new random plant
func NewRandomPlant(entities map[int]Entity, rng *rand.Rand) *Plant {
	// Random shade of green
	colour := (int(0xdd+(2*rng.Float32()-1)*0x0f) << 8)
	numLeaves := rng.Intn(5) + 1
	height := float32(1)
	radius := float32(0.125)
	x := rng.Float32()/4 - 1.0/8
	z := rng.Float32()/4 - 1.0/8
	rotX := math32.Pi * rng.Float32() / 4
	rotY := 2 * math32.Pi * rng.Float32()
	plant := NewPlant(entities, colour, numLeaves, height, radius, x, z, rotX, rotY, rng)

	return plant
}

// Create a new leaf
func NewLeaf() (mesh *graphic.Mesh) {
	geom := graphics.NewLeafMesh(2, 6, 2, 2)
	mat := material.NewStandard(math32.NewColorHex(0x101010))
	mesh = graphic.NewMesh(geom, mat)

	if tex, err := graphics.Texture(graphics.TexGrass); err == nil {
		mat.AddTexture(tex)
	} else {
		fmt.Println(err)
	}

	return
}

// Refresh reloads the in-game configuration of the plant
func (p *Plant) Refresh(entities map[int]Entity, scene *core.Node) {
	geom := geometry.NewCylinder(float64(p.Radius), float64(p.Height), 8, 8, true, true)
	mat := material.NewStandard(math32.NewColorHex(uint(p.Colour) / 10))
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetScale(0.1, 0.1, 0.1)

	if tex, err := graphics.Texture(graphics.TexStalk); err == nil {
		mat.AddTexture(tex)
	} else {
		fmt.Println(err)
	}

	if p.Mesh == nil {
		p.Mesh = mesh
		AddEntity(p, entities, scene)
	} else {
		// Make sure the entities map is pointing to this plant at the specified index
		if i, err := strconv.Atoi(p.Name()); err == nil {
			entities[i] = p
		}
	}

	p.SetPosition(p.X, mesh.Scale().Y, p.Z)
	p.SetRotation(p.RotX, p.RotY, 0)

	// Make sure each leaf exists
	if p.Leaves == nil || len(p.Leaves) != p.NumLeaves {
		p.Leaves = make([]*graphic.Mesh, p.NumLeaves)
	}

	for i := 0; i < p.NumLeaves; i++ {
		if p.Leaves[i] == nil {
			leaf := NewLeaf()
			leaf.SetName(p.Name())
			p.Add(leaf)
			p.Leaves[i] = leaf
		}

		// TODO: store leaf positions and rotation in genetics and make reloadable
		p.Leaves[i].SetScale(0.1, 0.1, 0.1)
		p.Leaves[i].SetRotation(p.Rand.Float32()*math32.Pi/12, p.Rand.Float32()*2*math32.Pi, p.Rand.Float32()*math32.Pi/12)
	}
}

// Perform per-frame updates to a plant
func (p *Plant) Update() {
	p.Age++
	scale := p.Scale()
	scale.Y = 0.1 + (0.001 * math32.Min(float32(p.Age), 10000))
	p.SetScale(scale.X, scale.Y, scale.Z)
	p.SetPosition(p.X, 0.5+p.Scale().Y/2, p.Z)
}

// Infostring returns a string representation of the plant
func (p Plant) InfoString() string {
	return fmt.Sprintf("plant,  : %dt,  : #%06x", p.Age, p.Colour)
}

// Material returns the plant's material
func (p Plant) Material() *material.Material {
	return p.GetMaterial(0).GetMaterial()
}

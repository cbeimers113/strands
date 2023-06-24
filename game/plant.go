package game

import (
	"math/rand"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Which data a Plant will store
type PlantData struct {
	// Whole Plant
	Age    int
	Colour int
	Height float32
	Radius float32
	X      float32
	Z      float32
	RotX   float32
	RotY   float32

	// Leaves
	NumLeaves        int     // (0, ..), How many leaves the plant has. (More and bigger leaves consume more resources)
	LeafSpawnHeight  float32 // [0.1, 1], How far up the stem the leaves will appear (from top), ie: value of 0.25 means leaves will spawn on the top quarter of the stem
	AvgLeafSize      float32 // (0, ..), Average size of leaf, ie: value of 1 equals the default size
	LeafSizeVariance float32 // [0, 1], How much the leaf sizes can vary, ie: a value of 0.5 means the leaves can be up to 50% bigger or smaller than AvgSize
}

// Perform action on plant entity on right click
func OnRightClickPlant(plant *Entity) {
	println("No right click behaviour defined for plant")
}

// Perform action on plant entity on left click
func OnLeftClickPlant(plant *Entity) {
	println("No left click behaviour defined for plant")
}

// Create a new plant
func NewPlant(colour, numLeaves int, height, radius, x, z, rotX, rotY float32) (plant *Entity) {
	geom := geometry.NewCylinder(float64(radius), float64(height), 8, 8, true, true)
	mat := material.NewStandard(math32.NewColorHex(uint(colour) / 10))
	mesh := graphic.NewMesh(geom, mat)
	mesh.SetScale(0.1, 0.1, 0.1)

	if tex, ok := Texture("stalk"); ok {
		mat.AddTexture(tex)
	}

	plant = NewEntity(mesh, Plant)
	plant.SetPosition(x, mesh.Scale().Y, z)
	plant.SetRotation(rotX, rotY, 0)
	plant.SetUserData(&PlantData{
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
	})

	for i := 0; i < numLeaves; i++ {
		leaf := NewLeaf()
		leaf.SetScale(0.1, 0.1, 0.1)
		leaf.SetRotation(rand.Float32() * math32.Pi / 12, rand.Float32() * 2 * math32.Pi, rand.Float32() * math32.Pi / 12)
		plant.Add(leaf)
	}

	return
}

// Create a new random plant
func NewRandomPlant() *Entity {
	// Random shade of green
	colour := (int(0xdd+(2*rand.Float32()-1)*0x0f) << 8)
	numLeaves := rand.Intn(5) + 1
	height := float32(1)
	radius := float32(0.125)
	x := rand.Float32()/4 - 1.0/8
	z := rand.Float32()/4 - 1.0/8
	rotX := math32.Pi * rand.Float32() / 4
	rotY := 2 * math32.Pi * rand.Float32()
	plant := NewPlant(colour, numLeaves, height, radius, x, z, rotX, rotY)

	return plant
}

// Create a new leaf
func NewLeaf() (mesh *graphic.Mesh) {
	geom := NewLeafMesh(2, 6, 2, 2)
	mat := material.NewStandard(math32.NewColorHex(0x101010))
	mesh = graphic.NewMesh(geom, mat)

	if tex, ok := Texture("grass"); ok {
		mat.AddTexture(tex)
	}

	return
}

// Grow the plant until maturity is reached
func (plant *Entity) growPlant(plantData *PlantData) {
	plantData.Age++

	if plantData.Age < 10000 { // TODO: Standardize "maturity" for plants
		scale := plant.Scale()
		scale.Y *= 1.001
		plant.SetScale(scale.X, scale.Y, scale.Z)
		plant.SetPosition(plantData.X, 0.5+plant.Scale().Y/2, plantData.Z)
	}
}

// Perform per-frame updates to a plant
func UpdatePlant(plant *Entity) {
	if plantData, ok := plant.UserData().(*PlantData); ok && !IsPaused {
		plant.growPlant(plantData)
	}
}

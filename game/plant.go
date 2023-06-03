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
	Age    int
	Colour int
	Height float32
	Radius float32
	X      float32
	Z      float32
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
func NewPlant(colour int, height, radius, x, z float32) (plant *Entity) {
	geom := geometry.NewCylinder(float64(radius), float64(height), 8, 8, true, true)
	mat := material.NewStandard(math32.NewColorHex(uint(colour) / 10))
	mesh := graphic.NewMesh(geom, mat)

	if tex, ok := Texture("stalk"); ok {
		mat.AddTexture(tex)
	}

	plant = NewEntity(mesh, Plant)
	plant.SetPosition(x, mesh.Scale().Y, z)
	plant.SetUserData(&PlantData{Colour: colour, Height: height, Radius: radius, X: x, Z: z})

	return
}

// Create a new random plant
func NewRandomPlant() *Entity {
	// Random shade of green
	colour := (int(0xdd+(2*rand.Float32()-1)*0x0f) << 8)
	height := (1 / 4) * (0.95 + rand.Float32()/10)
	radius := (1 / 16) * (0.95 + rand.Float32()/10)
	x := rand.Float32()/2 - 1/4
	z := rand.Float32()/2 - 1/4
	plant := NewPlant(colour, height, radius, x, z)

	return plant
}

// Grow the plant until maturity is reached
func (plant *Entity) growPlant(plantData *PlantData) {
	plantData.Age++

	if plantData.Age < 1000 { // TODO: Standardize "maturity" for plants
		scale := plant.Scale()
		scale.Y *= 1.001
		plant.SetScale(scale.X, scale.Y, scale.Z)
		plant.SetPosition(plantData.X, plant.Scale().Y/2, plantData.Z)
	}
}

// Perform per-frame updates to a plant
func UpdatePlant(plant *Entity) {
	if plantData, ok := plant.UserData().(*PlantData); ok {
		plant.growPlant(plantData)
	}
}

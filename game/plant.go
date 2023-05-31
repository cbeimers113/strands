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
	println("No right click behaviour defined for ", plant.Name())
}

// Perform action on plant entity on left click
func OnLeftClickPlant(plant *Entity) {
	println("No left click behaviour defined for ", plant.Name())
}

// Create a new plant
func NewPlant(colour int, height, radius float32) (plant *Entity) {
	geom := geometry.NewCylinder(float64(radius), float64(height), 8, 8, true, true)
	mat := material.NewStandard(math32.NewColorHex(uint(colour) / 10))
	mesh := graphic.NewMesh(geom, mat)
	x := rand.Float32()*TileSize/2 - TileSize/4
	z := rand.Float32()*TileSize/2 - TileSize/4

	if tex, ok := Texture("stalk"); ok {
		mat.AddTexture(tex)
	}

	plant = NewEntity(mesh, Plant)
	plant.SetPosition(x, mesh.Scale().Y/2, z)
	plant.SetUserData(PlantData{Colour: colour, Height: height, Radius: radius, X: x, Z: z})

	return
}

// Create a new random plant
func NewRandomPlant() *Entity {
	// Random shade of green
	colour := (int(0xdd+(2*rand.Float32()-0.5)*0x0f) << 8) // TODO: This changes the colour of the highlight texture too
	height := (TileSize / 4) * (0.95 + rand.Float32()/10)
	radius := (TileSize / 16) * (0.95 + rand.Float32()/10)

	return NewPlant(colour, height, radius)
}

// Perform per-frame updates to a plant
func UpdatePlant(plant *Entity) {
	if data, ok := plant.UserData().(PlantData); ok {
		data.Age++

		// Grow until maturity is reached
		if data.Age < 1000 {
			scale := plant.Scale()
			scale.Y *= 1.001
			plant.SetScale(scale.X, scale.Y, scale.Z)
			plant.SetPosition(data.X, plant.Scale().Y/2, data.Z)
		}

		// Update changes to the plant data
		plant.SetUserData(data)
	}
}

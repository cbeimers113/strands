package game

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

// Which data a Plant will store
type PlantData struct {
	Colour int
	Age    int
}

// Perform action on plant entity on right click
func OnRightClickPlant(plant Entity) {
	println("No right click behaviour defined for ", plant.Name())
}

// Perform action on plant entity on left click
func OnLeftClickPlant(plant Entity) {
	println("No left click behaviour defined for ", plant.Name())
}

// Create a new plant
func NewPlant(colour int) (plant *graphic.Mesh) {
	geom := geometry.NewCylinder(float64(TileSize)/8, float64(TileSize)/2, 8, 8, true, true)
	mat := material.NewStandard(math32.NewColorHex(uint(colour) / 10))
	plant = graphic.NewMesh(geom, mat)
	tex, ok := Texture("stalk")

	if ok {
		mat.AddTexture(tex)
	}

	plant.SetPosition(0, plant.Scale().Y/2, 0)
	plant.SetName(Plant)
	plant.SetUserData(PlantData{Colour: colour, Age: 0})

	return
}

// Perform per-frame updates to a plant
func UpdatePlant(plant Entity) {
	data, _ := plant.UserData().(PlantData)
	data.Age++

	// Grow until maturity is reached
	if data.Age < 1000 {
		scale := plant.Scale()
		scale.Y *= 1.001
		plant.SetScale(scale.X, scale.Y, scale.Z)
		plant.SetPosition(0, plant.Scale().Y, 0)
	}

	// Update changes to the plant data
	plant.SetUserData(data)
}

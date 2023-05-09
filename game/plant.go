package game

// Plant metadata mapping
const PlantColour int = 0
const PlantAge int = 1

// Perform action on plant entity on right click
func OnRightClickPlant(plant Entity) {
	println("No right click behaviour defined for ", plant.Name())
}

// Perform action on plant entity on left click
func OnLeftClickPlant(plant Entity) {
	println("No left click behaviour defined for ", plant.Name())
}

// Grow the plant slowly over time
func GrowPlant(plant Entity) {
	age := GetDatum(plant, PlantAge)
	age++

	// Grow until maturity is reached
	if age < 1000 {
		scale := plant.Scale()
		scale.Y *= 1.001
		plant.SetScale(scale.X, scale.Y, scale.Z)
		plant.SetPosition(0, plant.Scale().Y/2, 0)
	}

	SetDatum(plant, PlantAge, age)
}

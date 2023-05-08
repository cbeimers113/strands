package game

const Water string = "water"
const Dirt string = "dirt"
const Grass string = "grass"
const Stone string = "stone"

// Store list of tile types ordered by spawn height
var TypeStrata []string = []string{
	Water,
	Dirt,
	Grass,
	Stone,
}

// Perform an action on a tile entity on right click
func OnRightClickTile(tile Entity) {
	AddPlant(0x00dd05, tile)
}

// Perform an action on a tile entity on left click
func OnLeftClickTile(tile Entity) {
	println("No left click behaviour defined for ", tile.Name())
}

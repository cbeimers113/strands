package game

// Represents cells of chemical quantities in the atmospheres
var Atmosphere [Width][Height][Depth]map[ElementType]*Quantity

// Create the atmosphere
func CreateAtmosphere() {
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			for z := 0; z < Depth; z++ {
				Atmosphere[x][y][z] = make(map[ElementType]*Quantity)
				
				// Starting quantities
				Atmosphere[x][y][z][Water] = &Quantity{
					Value: 0.5,
					Units: Litre,
				}
			}
		}
	}
}

// Update the atmosphere
func UpdateAtmosphere(deltaTime float32) {
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			for z := 0; z < Depth; z++ {
				
			}
		}
	}
}

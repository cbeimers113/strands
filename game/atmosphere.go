package game

// Represents cells of chemical quantities and atmospheric temperature in the atmospheres
type Cell struct {
	Elements    map[ElementType]*Quantity
	Temperature float32
}

var Atmosphere [Width][Height][Depth]Cell

// Create the atmosphere
func CreateAtmosphere() {
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			for z := 0; z < Depth; z++ {
				Atmosphere[x][y][z] = Cell{
					Elements:    make(map[ElementType]*Quantity),
					Temperature: 22.0,
				}

				// Starting quantities
				Atmosphere[x][y][z].Elements[Water] = &Quantity{
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

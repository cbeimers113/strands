package atmosphere

import (
	"cbeimers113/strands/internal/chem"
	"cbeimers113/strands/internal/context"
)

// Represents cells of chemical quantities and atmospheric temperature in the atmospheres
type cell struct {
	Elements    map[chem.ElementType]*chem.Quantity
	Temperature float32
}

type Atmosphere struct {
	*context.Context

	cells [][][]cell
}

// Create the atmosphere
func New(ctx *context.Context) *Atmosphere {
	a := &Atmosphere{Context: ctx}
	a.cells = make([][][]cell, a.Cfg.Simulation.Width)

	for x := 0; x < a.Cfg.Simulation.Width; x++ {
		a.cells[x] = make([][]cell, a.Cfg.Simulation.Height)

		for y := 0; y < a.Cfg.Simulation.Height; y++ {
			a.cells[x][y] = make([]cell, a.Cfg.Simulation.Depth)

			for z := 0; z < a.Cfg.Simulation.Depth; z++ {
				a.cells[x][y][z] = cell{
					Elements:    make(map[chem.ElementType]*chem.Quantity),
					Temperature: 22.0,
				}

				// Starting quantities
				a.cells[x][y][z].Elements[chem.Water] = &chem.Quantity{
					Value: 0.5,
					Units: chem.Litre,
				}
			}
		}
	}

	return a
}

// Update the atmosphere
func (a *Atmosphere) Update(deltaTime float32) {
	for x := 0; x < a.Cfg.Simulation.Width; x++ {
		for y := 0; y < a.Cfg.Simulation.Height; y++ {
			for z := 0; z < a.Cfg.Simulation.Depth; z++ {

			}
		}
	}
}

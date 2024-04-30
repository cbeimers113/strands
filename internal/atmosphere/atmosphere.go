package atmosphere

import (
	"cbeimers113/strands/internal/chem"
	"cbeimers113/strands/internal/context"
	"cbeimers113/strands/internal/state"
)

type Atmosphere struct {
	*context.Context

	cells [][][]*state.Cell
}

// Create the atmosphere
func New(ctx *context.Context) *Atmosphere {
	a := &Atmosphere{Context: ctx}
	a.cells = make([][][]*state.Cell, a.Cfg.Simulation.Width)

	for x := 0; x < a.Cfg.Simulation.Width; x++ {
		a.cells[x] = make([][]*state.Cell, a.Cfg.Simulation.Height)

		for y := 0; y < a.Cfg.Simulation.Height; y++ {
			a.cells[x][y] = make([]*state.Cell, a.Cfg.Simulation.Depth)

			for z := 0; z < a.Cfg.Simulation.Depth; z++ {
				var t float32 = 22.0
				a.cells[x][y][z] = &state.Cell{
					Quantities:  make(map[chem.ElementType]*chem.Quantity),
					Temperature: t,
				}

				// Starting quantities
				a.cells[x][y][z].Quantities[chem.Water] = &chem.Quantity{
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
				// TODO: this
			}
		}
	}
}

// GetCells returns the atmosphere as a linear slice of cells
func (a Atmosphere) GetCells() []*state.Cell {
	w := a.Cfg.Simulation.Width
	h := a.Cfg.Simulation.Height
	d := a.Cfg.Simulation.Depth
	cells := make([]*state.Cell, w*h*d)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			for z := 0; z < d; z++ {
				cells[x+w*(y+h*z)] = a.cells[x][y][z]
			}
		}
	}

	return cells
}

// SetCells loads a slice of cells into the atmosphere
func (a *Atmosphere) SetCells(cells []*state.Cell) {
	w := a.Cfg.Simulation.Width
	h := a.Cfg.Simulation.Height
	d := a.Cfg.Simulation.Depth

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			for z := 0; z < d; z++ {
				a.cells[x][y][z] = cells[x+w*(y+h*z)]
			}
		}
	}
}

package state

import "cbeimers113/strands/internal/chem"

// Cell represents a single 3D section of chemical quantities and temperature in the atmosphere
type Cell struct {
	Quantities  map[chem.ElementType]*chem.Quantity `json:"quantities"`
	Temperature float32                             `json:"temperature"`
}

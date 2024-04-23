package chem

import "fmt"

// The various measurement units
type Unit string

const (
	Celcius Unit = "°C"
	Litre   Unit = "L"
	Gram    Unit = "g"
	Metre   Unit = "m"
)

// Represents an amount of an element
type Quantity struct {
	Value float32
	Units Unit
}

// Create a string representation of a quantity
func (q Quantity) String() (str string) {
	space := " "

	// Exception to spacing between value and unit is degrees
	if q.Units == Celcius {
		space = ""
	}

	str = fmt.Sprintf("%.2f%s%s", q.Value, space, q.Units)

	return
}

// Convert from litres to cubic metres (dimensions of one tile is 1 cubic metre)
func LitresToCubicMetres(litres float32) float32 {
	return litres / 1000
}

// Convert from cubic metres to litres (dimensions of one tile is 1 cubic metre)
func CubicMetresToLitres(cubicMetres float32) float32 {
	return cubicMetres * 1000
}

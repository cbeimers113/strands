package game

import "fmt"

// The various measurement units
type Unit string

const Celcius Unit = "°C"
const Litre Unit = "L"
const Gram Unit = "g"
const Metre Unit = "m"

// Represents an amount of an element
type Quantity struct {
	Value float32
	Units Unit
}

// Create a string representation of a quantity
func (q *Quantity) String() (str string) {
	space := " "

	// Exception to spacing between value and unit is degrees
	if q.Units == Celcius {
		space = ""
	}

	str = fmt.Sprintf("%.2f%s%s", q.Value, space, q.Units)

	return
}

// Convert from litres to cubic metres
func LitresToCubicMetres(litres float32) float32 {
	return litres / 1000
}

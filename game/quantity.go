package game

import "fmt"

// The various measurement units
type Unit string

const Celcius Unit = "°C"
const Litre Unit = "L"
const Gram Unit = "g"

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

package game

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

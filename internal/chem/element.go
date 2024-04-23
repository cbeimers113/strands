package chem

// Represents elements that exist dynamically in the world (eg. non-tile types)
type ElementType = string

const Water ElementType = "water"

var ElementTypes []ElementType = []ElementType{
	Water,
}

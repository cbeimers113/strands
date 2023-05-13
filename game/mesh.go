package game

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Generate a hexagon mesh with a given width
func CreateHexagon(width float32) (geom *geometry.Geometry) {
	geom = geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)
	uvs := math32.NewArrayF32(0, 16)

	// Vertex positioning for a hexagon made of 6 equilateral triangles
	r := width / 2
	dx := r * math32.Cos(math32.Pi/3) // This works out to r / 2
	dy := r * math32.Sin(math32.Pi/3) // This works out to r * 0.866

	// List all vertices in the hexagon
	vertices.Append(
		0.0, 0.0, 0.0, // centre 0
		-dx, 0.0, dy, // top left 1
		dx, 0.0, dy, // top right 2
		-r, 0.0, 0.0, // middle left 3
		r, 0.0, 0.0, // middle right 4
		-dx, 0.0, -dy, // bottom left 5
		dx, 0.0, -dy, // bottom right 6
	)

	// Add triangles' vertices in counter-clockwise
	indices.Append(
		0, 3, 1, // top left
		0, 1, 2, // top mid
		0, 2, 4, // top right
		0, 4, 6, // bottom right
		0, 6, 5, // bottom mid
		0, 5, 3, // bottom left
	)

	// Texture mapping, width and height are mapped to [0, 1], origin is top left
	uvs.Append(
		(r-dx)/width, 1.0, // bottom left
		0.0, 0.5, // middle left
		(r-dx)/width, 0.0, // top left
		(width-dx)/width, 0.0, // top right
		1.0, 0.5, // middle right
		(width-dx)/width, 1.0, // bottom right
	)

	geom.SetIndices(indices)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))
	return
}

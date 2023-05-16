package game

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
)

// Generate a hexagon mesh with a given width and height
func CreateHexagon(width, height float32) (geom *geometry.Geometry) {
	geom = geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)
	uvs := math32.NewArrayF32(0, 16)

	// Vertex positioning for a hexagon made of 6 equilateral triangles
	r := width / 2
	dx := r * math32.Cos(math32.Pi/3) // This works out to r / 2
	dy := r * math32.Sin(math32.Pi/3) // This works out to r * 0.866

	// List all vertices in the hexagon as though looking at the front each face
	vertices.Append(
		0.0, 0.0, 0.0, // surface: centre 0
		-dx, 0.0, dy, // surface: top left 1
		dx, 0.0, dy, // surface: top right 2
		-r, 0.0, 0.0, // surface: middle left 3
		r, 0.0, 0.0, // surface: middle right 4
		-dx, 0.0, -dy, // surface: bottom left 5
		dx, 0.0, -dy, // surface: bottom right 6
		-dx, -height, dy, // sides: top left 7
		dx, -height, dy, // sides: top right 8
		-r, -height, 0.0, // sides: middle left 9
		r, -height, 0.0, // sides: middle right 10
		-dx, -height, -dy, // sides: bottom left 11
		dx, -height, -dy, // sides: bottom right 12
	)

	// Add triangles' vertices in counter-clockwise
	indices.Append(
		0, 3, 1, // top left
		0, 1, 2, // top mid
		0, 2, 4, // top right
		0, 4, 6, // bottom right
		0, 6, 5, // bottom mid
		0, 5, 3, // bottom left
		1, 9, 7, // top left side 1
		3, 9, 1, // top left side 2
		7, 8, 2, // top side 1
		7, 2, 1, // top side 2
		8, 10, 4, // top right side 1
		8, 4, 2, // top right side 2
		10, 12, 6, // bottom right side 1
		10, 6, 4, // bottom right side 2
		12, 11, 5, // bottom side 1
		12, 5, 6, // bottom side 2
		11, 9, 3, // bottom left side 1
		5, 11, 3, // bottom left side 2
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

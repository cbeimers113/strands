package graphics

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

// Generate a 1 cubic metre hexagon mesh
func NewHexMesh() (geom *geometry.Geometry) {
	geom = geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 42)
	indices := math32.NewArrayU32(0, 72)
	uvs := math32.NewArrayF32(0, 12)

	// Vertex positioning for a hexagon made of 6 equilateral triangles
	r := float32(0.5)
	dx := r * math32.Cos(math32.Pi/3) // This works out to r / 2
	dy := r * math32.Sin(math32.Pi/3) // This works out to r * 0.866

	// List all vertices in the hexagon as though looking at the front each face
	vertices.Append(
		0.0, 0.5, 0.0, // surface: centre 0
		-dx, 0.5, dy, // surface: top left 1
		dx, 0.5, dy, // surface: top right 2
		-r, 0.5, 0.0, // surface: middle left 3
		r, 0.5, 0.0, // surface: middle right 4
		-dx, 0.5, -dy, // surface: bottom left 5
		dx, 0.5, -dy, // surface: bottom right 6
		-dx, 0.0, dy, // sides: top left 7
		dx, 0.0, dy, // sides: top right 8
		-r, 0.0, 0.0, // sides: middle left 9
		r, 0.0, 0.0, // sides: middle right 10
		-dx, 0.0, -dy, // sides: bottom left 11
		dx, 0.0, -dy, // sides: bottom right 12
		0.0, 0.0, 0.0, // bottom surface: centre 13

	)

	// Add triangles' vertices clockwise
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
		7, 9, 13, // bottom surface top left
		8, 7, 13, // bottom surface top mid
		10, 8, 13, // bottom surface top right
		12, 10, 13, // bottom surface bottom right
		11, 12, 13, // bottom surface bottom mid
		9, 11, 13, // bottom surface bottom left
	)

	// Texture mapping, width and height are mapped to [0, 1], origin is top left
	uvs.Append(
		r-dx, 1.0, // bottom left
		0.0, 0.5, // middle left
		r-dx, 0.0, // top left
		1-dx, 0.0, // top right
		1.0, 0.5, // middle right
		1-dx, 1.0, // bottom right
	)

	geom.SetIndices(indices)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	return
}

// Generate a leaf mesh with a given width, length, stem length, and centre (point on length axis where width is max)
func NewLeafMesh(w, l, s, c float32) (geom *geometry.Geometry) {
	geom = geometry.NewGeometry()
	vertices := math32.NewArrayF32(0, 16)
	indices := math32.NewArrayU32(0, 16)
	uvs := math32.NewArrayF32(0, 16)

	// Vertex positioning for a leaf shape:
	c += s
	l += s

	// List all vertices in the leaf as though looking at the front each face
	vertices.Append(
		-w/8, 0, 0, // top face: stem bottom left corner 0
		w/8, 0, 0, // top face: stem bottom right corner 1
		-w/8, 0, s, // top face: stem top left corner 2
		w/8, 0, s, // top face: stem top right corner 3
		-w/2, 0, c, // top face: left centre point 4
		0, 0, c, // top face: mid centre point 5
		w/2, 0, c, // top face: right centre point 6
		0, 0, l, // top face: tip 7

		-w/8, 0, 0, // bottom face: stem bottom left corner 8
		w/8, 0, 0, // bottom face: stem bottom right corner 9
		-w/8, 0, s, // bottom face: stem top left corner 10
		w/8, 0, s, // bottom face: stem top right corner 11
		-w/2, 0, c, // bottom face: left centre point 12
		0, 0, c, // bottom face: mid centre point 13
		w/2, 0, c, // bottom face: right centre point 14
		0, 0, l, // bottom face: tip 15
	)

	// Add triangles' vertices clockwise
	indices.Append(
		0, 2, 1, // top face: stem left
		2, 3, 1, // top face: stem right
		2, 4, 5, // top face: lower half left
		2, 5, 3, // top face: lower half mid
		3, 5, 6, // top face: lower half right
		4, 7, 5, // top face: upper half left
		5, 7, 6, // top face: upper half right

		9, 10, 8, // bottom face: stem left
		9, 11, 10, // bottom face: stem right
		13, 12, 10, // bottom face: lower half left
		11, 13, 10, // bottom face: lower half mid
		14, 13, 11, // bottom face: lower half right
		13, 15, 12, // bottom face: upper half left
		14, 15, 13, // bottom face: upper half right
	)

	// Texture mapping, width and height are mapped to [0, 1], origin is top left
	uvs.Append(
		0.0, 0.0, // top left
		0.0, 1.0, // top right
		1.0, 1.0, // bottom right
		1.0, 0.0, // bottom left
	)

	geom.SetIndices(indices)
	geom.AddVBO(gls.NewVBO(vertices).AddAttrib(gls.VertexPosition))
	geom.AddVBO(gls.NewVBO(uvs).AddAttrib(gls.VertexTexcoord))

	return
}

// Get the dimensions of a mesh
func DimensionsOf(mesh *graphic.Mesh) *math32.Vector3 {
	bb := mesh.BoundingBox()
	x := math32.Abs(bb.Max.X - bb.Min.X)
	y := math32.Abs(bb.Max.Y - bb.Min.Y)
	z := math32.Abs(bb.Max.Z - bb.Min.Z)

	return math32.NewVector3(x, y, z)
}

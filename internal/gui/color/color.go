package color

import "github.com/g3n/engine/math32"

var (
	Background = &math32.Color{R: 0.25, G: 0.25, B: 0.25}
	Focus      = &math32.Color{R: 0.45, G: 0.45, B: 0.45}
	Green      = &math32.Color{R: 0.25, G: 0.50, B: 0.25}
	Yellow     = &math32.Color{R: 0.85, G: 0.55, B: 0.00}
	Red        = &math32.Color{R: 0.50, G: 0.25, B: 0.25}
	White      = &math32.Color{R: 1.00, G: 1.00, B: 1.00}
	Black      = &math32.Color{R: 0.00, G: 0.00, B: 0.00}

	Opaque = &math32.Color4{R: 1.00, G: 1.00, B: 1.00, A: 0.75}
)

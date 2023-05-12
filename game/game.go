package game

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"
)

const SimSpeed int = 60 // Simulation update speed in TPS

var Application *app.Application
var Scene *core.Node
var Cam *camera.Camera
var Running bool

// Callback for when the window is resized, update camera to match
func onResize(evname string, ev interface{}) {
	width, height := Application.GetSize()
	Application.Gls().Viewport(0, 0, int32(width), int32(height))
	Cam.SetAspect(float32(width) / float32(height))
}

// Build the application with an empty scene
func buildApplication() {
	Textures = make(map[string]*texture.Texture2D)
	Application = app.App()
	Scene = core.NewNode()
	Scene.SetName("world")

	// Configure camera
	Cam = camera.New(1)
	Cam.SetPosition(float32(Width)*TileSize/2, TileSize, float32(Height)*TileSize/2)
	Cam.SetRotation(-45, 0, 0)
	Scene.Add(Cam)
	RegisterControls()

	// Refresh display
	Application.Subscribe(window.OnWindowSize, onResize)
	Application.Gls().ClearColor(0, 0, 0, 1.0)
	onResize("", nil)
}

// Run the application
func Run() {
	if Running {
		return
	}

	buildApplication()
	LoadGUI()
	LoadWorld()
	Running = true

	// Update every n ms so that SimSpeed updates happen per second
	tickThreshold := 1000 / SimSpeed
	timer := 0

	Application.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		Application.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(Scene, Cam)

		// TPS counter
		timer += int(deltaTime / 1000000)

		if timer >= tickThreshold {
			Update(timer)
			timer = 0
		}
	})
}

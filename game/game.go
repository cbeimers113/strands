package game

import (
	"math/rand"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const SimSpeed int = 60 // Simulation update speed in TPS

var Application *app.Application
var Scene *core.Node
var Cam *camera.Camera
var Win *window.GlfwWindow

// This determines if everything in the simulation is frozen, including the player
var IsFrozen bool

// This determines if the simulation physics are paused, but the player can still interact with the simulation
var IsPaused bool

// Set whether the game is paused
func SetPaused(paused bool) {
	IsFrozen = paused

	if Win != nil {
		switch paused {
		case true:
			Win.SetInputMode(glfw.CursorMode, int(window.CursorNormal))
		case false:
			Win.SetInputMode(glfw.CursorMode, int(window.CursorDisabled))
		}
	}
}

// Callback for when the window is resized, update camera and gui to match
func onResize(evname string, ev interface{}) {
	width, height := Application.GetSize()
	Application.Gls().Viewport(0, 0, int32(width), int32(height))
	Cam.SetAspect(float32(width) / float32(height))
	ReloadGui()
}

// Run the application
func Run() {
	Textures = make(map[string]*texture.Texture2D)
	Application = app.App()
	Win, _ = Application.IWindow.(*window.GlfwWindow)
	Scene = core.NewNode()

	// Configure camera
	Cam = camera.New(1)
	Cam.SetPosition(float32(Width)/2, 1, float32(Depth)/2)
	Cam.SetRotation(0, 0, 0)
	Scene.Add(Cam)
	RegisterControls()

	// Load game
	LoadGui()
	LoadWorld()
	Views[MainMenu].Open(true)

	// Refresh display
	Application.Subscribe(window.OnWindowSize, onResize)
	Application.Gls().ClearColor(0, 0, 0, 1.0)
	onResize("", nil)

	// Update every n ms so that SimSpeed updates happen per second
	var tickThreshold float32 = 1000 / float32(SimSpeed)
	var deltaTime float32 = 0

	// Seed the PRNG
	rand.Seed(time.Now().UnixNano())

	Application.Run(func(renderer *renderer.Renderer, duration time.Duration) {
		Application.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(Scene, Cam)

		if !IsFrozen {
			// TPS counter
			deltaTime += float32(duration.Milliseconds())
			if deltaTime >= tickThreshold {
				UpdateWorld(deltaTime)
				UpdateAtmosphere(deltaTime)
				UpdatePlayer(deltaTime)
				RefreshGui()
				deltaTime = 0
			}
		}
	})
}

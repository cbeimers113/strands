package game

import (
	"errors"
	"math/rand"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"github.com/go-gl/glfw/v3.3/glfw"

	"cbeimers113/strands/internal/config"
	"cbeimers113/strands/internal/context"
	"cbeimers113/strands/internal/gui"
	"cbeimers113/strands/internal/input"
	"cbeimers113/strands/internal/player"
	"cbeimers113/strands/internal/state"
	"cbeimers113/strands/internal/world"
)

type Game struct {
	*context.Context

	gui    *gui.Gui
	iman   *input.InputManager
	world  *world.World
	player *player.Player
}

func New(cfg *config.Config) (*Game, error) {
	ctx := &context.Context{
		App:   app.App(),
		Scene: core.NewNode(),
		Cam:   camera.New(1),
		Cfg:   cfg,
		State: state.New(),
	}

	g := &Game{
		Context: ctx,
		gui:     gui.New(ctx),
		world:   world.New(ctx),
		player:  player.New(ctx),
		iman:    input.New(ctx),
	}

	// Create window
	var ok bool
	if g.Win, ok = g.App.IWindow.(*window.GlfwWindow); !ok {
		return nil, errors.New("cannot instantiate window")
	}

	// Configure camera
	g.Cam.SetPosition(float32(cfg.Simulation.Width)/2, 10, float32(cfg.Simulation.Depth)/2)
	g.Cam.SetRotation(0, 0, 0)
	g.Scene.Add(g.Cam)

	return g, nil
}

// Start the game
func (g *Game) Start() {
	// Open main menu on game start
	gui.Open(gui.MainMenu, true)

	// Refresh display
	g.App.Subscribe(window.OnWindowSize, g.onResize)
	g.App.Gls().ClearColor(0, 0, 0, 1.0)
	g.Win.SetSize(g.Cfg.Window.Width, g.Cfg.Window.Height)
	g.Win.SetTitle(g.Cfg.Name)
	g.onResize("", nil)

	// Update every n ms so that n updates happen per second
	var tickThreshold float32 = 1000 / float32(g.Cfg.Simulation.Speed)
	var deltaTime float32 = 0

	// Seed the PRNG
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Start the main loop
	g.App.Run(func(renderer *renderer.Renderer, duration time.Duration) {
		deltaTime += float32(duration.Milliseconds())
		g.App.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(g.Scene, g.Cam)

		if g.State.InMenu() {
			g.Win.SetInputMode(glfw.CursorMode, int(window.CursorNormal))
		} else {
			g.Win.SetInputMode(glfw.CursorMode, int(window.CursorDisabled))

			if deltaTime >= tickThreshold {
				g.world.Update(deltaTime)
				g.iman.Update(g.player)
				g.player.Update(deltaTime)
				g.gui.Refresh()
				deltaTime = 0
			}
		}
	})
}

// Callback for when the window is resized, update camera and gui to match
func (g *Game) onResize(evname string, ev interface{}) {
	width, height := g.App.GetSize()
	g.App.Gls().Viewport(0, 0, int32(width), int32(height))
	g.Cam.SetAspect(float32(width) / float32(height))
	g.gui.Reload()
}

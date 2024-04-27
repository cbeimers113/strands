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
		State: state.New(cfg),
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

	// Configure app
	g.App.Subscribe(window.OnWindowSize, g.onResize)
	g.App.Gls().ClearColor(0, 0, 0, 1.0)
	g.Win.SetSize(g.Cfg.Window.Width, g.Cfg.Window.Height)
	g.Win.SetTitle(g.Cfg.Name)
	g.onResize("", nil)

	// Update world every n ms so that n updates happen per second
	var deltaTimeWorld float32 = 0
	// Update the player controller 60 times per second to keep it smooth no matter the sim speed
	var deltaTimeController float32 = 0

	var lastTick time.Time
	var tps int

	// Seed the PRNG
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Start the main loop
	g.App.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		g.App.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(g.Scene, g.Cam)

		if g.State.InMenu() {
			g.Win.SetInputMode(glfw.CursorMode, int(window.CursorNormal))
			deltaTimeController = 0
			deltaTimeWorld = 0
		} else {
			g.Win.SetInputMode(glfw.CursorMode, int(window.CursorDisabled))
			deltaTimeWorld += float32(deltaTime.Milliseconds())
			deltaTimeController += float32(deltaTime.Milliseconds())

			// Update the controller on a fixed interval
			if deltaTimeController >= 1000/60 {
				g.iman.Update(g.player)
				g.player.Update(deltaTimeController)
				g.gui.Refresh()
				deltaTimeController = 0
			}

			// Update the world at a dynamic configurable rate
			if deltaTimeWorld >= 1000/float32(g.Cfg.Simulation.Speed) {
				g.world.Update(deltaTimeWorld)

				if !g.State.Paused() {
					g.State.Clock.Update(deltaTimeWorld)
				}

				deltaTimeWorld = 0
				tps++
			}

			// Count how many ticks happened since one second ago, update TPS count
			if time.Since(lastTick) >= 1*time.Second {
				g.State.SetTPS(tps)
				tps = 0
				lastTick = time.Now()
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

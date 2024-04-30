package game

import (
	_ "embed"
	"errors"
	"fmt"
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
	"cbeimers113/strands/internal/io/input_manager"
	"cbeimers113/strands/internal/io/keyboard"
	"cbeimers113/strands/internal/player"
	"cbeimers113/strands/internal/state"
	"cbeimers113/strands/internal/world"
)

type Game struct {
	*context.Context

	gui    *gui.Gui
	iman   *input_manager.InputManager
	world  *world.World
	player *player.Player
}

func New(cfg *config.Config, version string) (*Game, error) {
	ctx := &context.Context{
		Version:  version,
		App:      app.App(),
		Scene:    core.NewNode(),
		Cam:      camera.New(1),
		Cfg:      cfg,
		State:    state.New(cfg, time.Now().UnixNano()),
		Keyboard: keyboard.New(),
	}

	g := &Game{
		Context: ctx,
		gui:     gui.New(ctx),
		world:   world.New(ctx),
		player:  player.New(ctx),
		iman:    input_manager.New(ctx),
	}

	// Create window
	var ok bool
	if g.Win, ok = g.App.IWindow.(*window.GlfwWindow); !ok {
		return nil, errors.New("cannot instantiate window")
	}
	g.Win.SetSize(g.Cfg.Window.Width, g.Cfg.Window.Height)
	g.Win.SetTitle(g.Cfg.Name)

	// Configure app
	g.App.Subscribe(window.OnWindowSize, g.onResize)
	g.App.Subscribe(app.OnExit, g.onExit)
	g.App.Gls().ClearColor(0, 0, 0, 1.0)
	g.onResize("", nil)

	// Configure camera
	g.Cam.SetPosition(float32(cfg.Simulation.Width)/2, 10, float32(cfg.Simulation.Depth)/2)
	g.Cam.SetRotation(0, 0, 0)
	g.Scene.Add(g.Cam)

	// If exit save is enabled, load the exit save file
	if g.Cfg.ExitSave {
		g.LoadGame(state.ExitSaveFile)
	}

	return g, nil
}

// Load a save file
func (g *Game) LoadGame(filename string) {
	g.Win.SetTitle(fmt.Sprintf("Loading Simulation [%s]", filename))
	st, cells, tiles, err := state.LoadSave(g.Cfg, filename)
	if err != nil {
		fmt.Printf("Couldn't load save file [%s]: %s\n", filename, err) // TODO: popup with warning
		return
	}

	g.State = st
	g.world.SetAtmosphere(cells)
	g.world.SetTiles(tiles)
	g.Win.SetTitle(g.Cfg.Name)
}

// Save a save file
func (g Game) SaveGame(filename string) {
	g.Win.SetTitle(fmt.Sprintf("Saving Simulation [%s]", filename))
	if err := state.StoreSave(filename, g.State, g.world.GetAtmosphere(), g.world.GetTiles()); err != nil {
		fmt.Printf("Couldn't create save file [%s]: %s", filename, err)
	}
	g.Win.SetTitle(g.Cfg.Name)
}

// Start the game
func (g *Game) Start() {
	// Open main menu on game start
	gui.Open(gui.MainMenu, true)

	// Update world every n ms so that n updates happen per second
	var deltaTimeWorld float32 = 0
	// Update the player controller 60 times per second to keep it smooth no matter the sim speed
	var deltaTimeController float32 = 0

	var lastTick time.Time
	var tps int

	// Start the main loop
	g.App.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		g.App.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(g.Scene, g.Cam)

		if g.State.InMenu() {
			g.Win.SetInputMode(glfw.CursorMode, int(window.CursorNormal))
			deltaTimeController = 0
			deltaTimeWorld = 0
			g.gui.Refresh()
		} else {
			g.Win.SetInputMode(glfw.CursorMode, int(window.CursorDisabled))
			deltaTimeWorld += float32(deltaTime.Milliseconds())
			deltaTimeController += float32(deltaTime.Milliseconds())

			// Update the controller at a fixed interval
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

		// Poll for file load/save actions
		if g.SaveFile != "" {
			g.SaveGame(g.SaveFile)
			g.SaveFile = ""
		}
		if g.LoadFile != "" {
			g.LoadGame(g.LoadFile)
			g.LoadFile = ""
		}
	})
}

// Callback for when the app is closed
func (g Game) onExit(string, interface{}) {
	g.Cfg.Save()

	if g.Cfg.ExitSave {
		g.SaveGame(state.ExitSaveFile)
	}
}

// Callback for when the window is resized, update camera and gui to match
func (g *Game) onResize(string, interface{}) {
	width, height := g.App.GetSize()
	g.App.Gls().Viewport(0, 0, int32(width), int32(height))
	g.Cam.SetAspect(float32(width) / float32(height))
	g.gui.Reload()
}

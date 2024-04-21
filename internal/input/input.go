package input

import (
	"github.com/g3n/engine/window"

	"cbeimers113/strands/internal/context"
	"cbeimers113/strands/internal/entity"
	"cbeimers113/strands/internal/gui"
	"cbeimers113/strands/internal/player"
)

type InputManager struct {
	*context.Context

	prevMX float32
	prevMY float32

	// Δx and y, player movement
	dx float32
	dz float32

	// player looking x and y
	lx float32
	ly float32
}

func New(ctx *context.Context) *InputManager {
	i := &InputManager{Context: ctx}

	// Register the controls with the game application
	i.App.Subscribe(window.OnKeyDown, i.KeyDown)
	i.App.Subscribe(window.OnKeyUp, i.KeyUp)
	i.App.Subscribe(window.OnKeyRepeat, i.KeyHold)
	i.App.Subscribe(window.OnMouseDown, i.MouseDown)
	i.App.Subscribe(window.OnCursor, i.MouseMove)

	return i
}

func (i *InputManager) Update(player *player.Player) {
	player.Move(i.dx, i.dz)
	player.Look(i.lx, i.ly)

	i.dx = 0
	i.dz = 0
	i.lx = 0
	i.ly = 0
}

// Handle key down events for the game
func (i *InputManager) KeyDown(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	case window.KeyEscape:
		if i.State.InMenu() {
			gui.Open(gui.SimulationView, true)
		} else {
			gui.Open(gui.MainMenu, true)
		}
	case window.KeyS:
		i.dz = 1
	case window.KeyW:
		i.dz = -1
	case window.KeyD:
		i.dx = 1
	case window.KeyA:
		i.dx = -1
	case window.KeySpace:
		i.State.SetPaused(!i.State.Paused())
	}
}

// Handle key up events for the game
func (i *InputManager) KeyUp(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	case window.KeyS:
		i.dz = 0
	case window.KeyW:
		i.dz = 0
	case window.KeyD:
		i.dx = 0
	case window.KeyA:
		i.dx = 0
	}
}

// Handle key hold events for the game
func (i *InputManager) KeyHold(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	}
}

// Handle mouse click events for the game
func (i *InputManager) MouseDown(evname string, ev interface{}) {
	me := ev.(*window.MouseEvent)

	if !i.State.InMenu() && i.State.LookingAt != nil {
		switch i.State.LookingAt.Type {
		case entity.Tile:
			switch me.Button {
			case window.MouseButton1:
				entity.OnLeftClickTile(i.State.LookingAt)
			case window.MouseButton2:
				entity.OnRightClickTile(i.State.LookingAt, i.State.Entities)
			}
		case entity.Plant:
			switch me.Button {
			case window.MouseButton1:
				entity.OnLeftClickPlant(i.State.LookingAt)
			case window.MouseButton2:
				entity.OnRightClickPlant(i.State.LookingAt)
			}
		case entity.Creature:
			switch me.Button {
			case window.MouseButton1:
				entity.OnLeftClickCreature(i.State.LookingAt)
			case window.MouseButton2:
				entity.OnRightClickCreature(i.State.LookingAt)
			}
		default:
			println("No action defined for button ", me.Button, " on ", i.State.LookingAt.Type)
		}
	}
}

// Handle mouse movement events for the game
func (i *InputManager) MouseMove(evname string, ev interface{}) {
	me := ev.(*window.CursorEvent)
	mx := me.Xpos
	my := me.Ypos

	i.lx = i.prevMX - mx
	i.ly = my - i.prevMY

	i.prevMX = mx
	i.prevMY = my
}

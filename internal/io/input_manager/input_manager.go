package input_manager

import (
	"fmt"

	"github.com/g3n/engine/window"

	"cbeimers113/strands/internal/chem"
	"cbeimers113/strands/internal/context"
	"cbeimers113/strands/internal/entity"
	"cbeimers113/strands/internal/gui"
	"cbeimers113/strands/internal/player"
)

type InputManager struct {
	*context.Context

	// Δx and y, player movement
	dx     float32
	dz     float32
	prevMX float32
	prevMY float32
	shift  bool
	ctrl   bool

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
	var mf float32 = i.Cfg.Controls.MoveSpeed
	if i.shift {
		mf *= 4
	}

	player.Move(i.dx*mf, i.dz*mf)
	player.Look(i.lx, i.ly)

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
		i.dz = 0.01
	case window.KeyW:
		i.dz = -0.01
	case window.KeyD:
		i.dx = 0.01
	case window.KeyA:
		i.dx = -0.01
	case window.KeyLeftShift:
		i.shift = true
	case window.KeyCapsLock:
		i.shift = !i.shift
	case window.KeyLeftControl:
		i.ctrl = true
	case window.KeySpace:
		i.State.SetPaused(!i.State.Paused())

		// Debugging keyboard: toggle keyboard on
	case window.KeyTab:
		i.Keyboard.Enable(!i.Keyboard.GetEnabled())
	}

	i.Keyboard.Shift(i.shift)
	i.Keyboard.Ctrl(i.ctrl)
	i.Keyboard.Input(kev.Key)
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
	case window.KeyLeftShift:
		i.shift = false
	case window.KeyLeftControl:
		i.ctrl = false
	}
}

// Handle key hold events for the game
func (i *InputManager) KeyHold(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)
	i.Keyboard.Shift(i.shift)
	i.Keyboard.Ctrl(i.ctrl)
	i.Keyboard.Input(kev.Key)
}

// Handle mouse click events for the game
func (i *InputManager) MouseDown(evname string, ev interface{}) {
	me := ev.(*window.MouseEvent)

	if !i.State.InMenu() && i.State.LookingAt != nil {
		switch i.State.LookingAt.(type) {

		case *entity.Tile:
			tile := i.State.LookingAt.(*entity.Tile)

			switch me.Button {
			case window.MouseButton1:
				tile.AddWater(chem.CubicMetresToLitres(1))
			case window.MouseButton2:
				gui.Open(gui.TileContextMenu, false)
			}

		case *entity.Plant:
			// plant := i.State.LookingAt.(*entity.Plant)

			switch me.Button {
			case window.MouseButton1:
				// Left click: ?
			case window.MouseButton2:
				// Right click: TODO: plant context menu
			}

		// case *entity.Creature:
		// 	creature := i.State.LookingAt.(*entity.Creature)

		// 	switch me.Button {
		// 	case window.MouseButton1:
		// 		// On left click creature
		// 	case window.MouseButton2:
		// 		// On right click creature
		// 	}

		default:
			fmt.Printf("No action defined for button %+v on %+v\n", me.Button, i.State.LookingAt)
		}
	} else if i.State.InMenu() {
		i.Keyboard.RegisterMouseEvent(me.Xpos, me.Ypos)
	}
}

// Handle mouse movement events for the game
func (i *InputManager) MouseMove(evname string, ev interface{}) {
	me := ev.(*window.CursorEvent)
	mx := me.Xpos
	my := me.Ypos

	if !i.State.InMenu() {
		dx := i.prevMX - mx
		dy := my - i.prevMY
		i.lx += dx / 1000
		i.ly += dy / 1000
	} else {
		i.lx = 0
		i.ly = 0
	}

	i.prevMX = mx
	i.prevMY = my
}

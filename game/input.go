package game

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

var prevMX float32
var prevMY float32

// Handle key down events for the game
func KeyDown(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	case window.KeyEscape:
		if IsPaused {
			Views[SimulationView].Open(true)
		} else {
			Views[MainMenu].Open(true)
		}
	case window.KeyS:
		PlayerMoveZ = 1
	case window.KeyW:
		PlayerMoveZ = -1
	case window.KeyD:
		PlayerMoveX = 1
	case window.KeyA:
		PlayerMoveX = -1
	}
}

// Handle key up events for the game
func KeyUp(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	case window.KeyS:
		PlayerMoveZ = 0
	case window.KeyW:
		PlayerMoveZ = 0
	case window.KeyD:
		PlayerMoveX = 0
	case window.KeyA:
		PlayerMoveX = 0
	}
}

// Handle key hold events for the game
func KeyHold(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	}
}

// Handle mouse click events for the game
func MouseDown(evname string, ev interface{}) {
	me := ev.(*window.MouseEvent)
	r := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	r.SetFromCamera(Cam, 0, 0)
	i := r.IntersectObject(Scene, true)

	var object *core.Node

	// If we hit something, trigger any necessary callbacks
	if len(i) != 0 {
		object = i[0].Object.GetNode()
		entity := EntityOf(object)

		if entity != nil {
			switch entity.Type {
			case Tile:
				switch me.Button {
				case window.MouseButton1:
					OnLeftClickTile(entity)
				case window.MouseButton2:
					OnRightClickTile(entity)
				}
			case Plant:
				switch me.Button {
				case window.MouseButton1:
					OnLeftClickPlant(entity)
				case window.MouseButton2:
					OnRightClickPlant(entity)
				}
			case Creature:
				switch me.Button {
				case window.MouseButton1:
					OnLeftClickCreature(entity)
				case window.MouseButton2:
					OnRightClickCreature(entity)
				}
			default:
				println("No action defined for button ", me.Button, " on ", entity.Type)
			}
		}
	}
}

// Handle mouse movement events for the game
func MouseMove(evname string, ev interface{}) {
	me := ev.(*window.CursorEvent)
	mx := me.Xpos
	my := me.Ypos
	PlayerLookX = prevMX - mx
	PlayerLookY = my - prevMY
	prevMX = mx
	prevMY = my
}

// Register the controls with the game application
func RegisterControls() {
	Application.Subscribe(window.OnKeyDown, KeyDown)
	Application.Subscribe(window.OnKeyUp, KeyUp)
	Application.Subscribe(window.OnKeyRepeat, KeyHold)
	Application.Subscribe(window.OnMouseDown, MouseDown)
	Application.Subscribe(window.OnCursor, MouseMove)
}

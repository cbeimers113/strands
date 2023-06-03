package game

import (
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

	if !IsPaused && LookingAt != nil {
		switch LookingAt.Type {
		case Tile:
			switch me.Button {
			case window.MouseButton1:
				OnLeftClickTile(LookingAt)
			case window.MouseButton2:
				OnRightClickTile(LookingAt)
			}
		case Plant:
			switch me.Button {
			case window.MouseButton1:
				OnLeftClickPlant(LookingAt)
			case window.MouseButton2:
				OnRightClickPlant(LookingAt)
			}
		case Creature:
			switch me.Button {
			case window.MouseButton1:
				OnLeftClickCreature(LookingAt)
			case window.MouseButton2:
				OnRightClickCreature(LookingAt)
			}
		default:
			println("No action defined for button ", me.Button, " on ", LookingAt.Type)
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

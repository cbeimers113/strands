package game

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

// Handle key down events for the game
func KeyDown(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
	}
}

// Handle key up events for the game
func KeyUp(evname string, ev interface{}) {
	kev := ev.(*window.KeyEvent)

	switch kev.Key {
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

	w, h := Application.GetSize()
	r := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	r.SetFromCamera(Cam, (-.5+me.Xpos/float32(w))*2.0, (.5-me.Ypos/float32(h))*2.0)
	i := r.IntersectObject(Scene, true)

	var object Entity

	// If we hit something, trigger any necessary callbacks
	if len(i) != 0 {
		object = i[0].Object.GetNode()

		switch Type(object) {
		case Tile:
			switch me.Button {
			case window.MouseButton1:
				OnLeftClickTile(object)
			case window.MouseButton2:
				OnRightClickTile(object)
			}
		case Plant:
			switch me.Button {
			case window.MouseButton1:
				OnLeftClickPlant(object)
			case window.MouseButton2:
				OnRightClickPlant(object)
			}
		case Creature:
			switch me.Button {
			case window.MouseButton1:
				OnLeftClickCreature(object)
			case window.MouseButton2:
				OnRightClickCreature(object)
			}
		default:
			println("No action defined for button ", me.Button, " on ", Type(object))
		}
	}
}

// Register the controls with the game application
func RegisterControls() {
	Application.Subscribe(window.OnKeyDown, KeyDown)
	Application.Subscribe(window.OnKeyUp, KeyUp)
	Application.Subscribe(window.OnKeyRepeat, KeyHold)
	Application.Subscribe(window.OnMouseDown, MouseDown)
	camera.NewOrbitControl(Cam).SetTarget(*math32.NewVector3(float32(Width)*TileSize/2, 0, float32(Depth)*TileSize/2))
}

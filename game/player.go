package game

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/math32"
)

const lookSensitivityX float32 = 0.025
const lookSensitivityY float32 = 0.015

// Movement offset inpus
var PlayerLookX float32
var PlayerLookY float32
var PlayerMoveX float32
var PlayerMoveZ float32

// What entity the player is looking at
var LookingAt *Entity

// What tile is currently selected by the player
var SelectedTile *Entity

// The rotation offset for looking around
var rotationY float32
var rotationX float32

// Update the player
func UpdatePlayer(deltaTime float32) {
	// Looking
	θ := Cam.Rotation().Y
	rotationX = θ + PlayerLookX*lookSensitivityX
	rotationY += PlayerLookY * lookSensitivityY
	rotationY = math32.Clamp(rotationY, -math32.Pi/2, math32.Pi/2)

	Cam.SetRotation(-rotationY, rotationX, 0)
	Look()

	PlayerLookX = 0
	PlayerLookY = 0

	// Movement
	Cam.TranslateX(PlayerMoveX)
	Cam.TranslateZ(PlayerMoveZ)
}

// Update which entity the player is looking at
func Look() {
	r := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	r.SetFromCamera(Cam, 0, 0)
	i := r.IntersectObject(Scene, true)

	LookingAt = nil
	var object *core.Node

	// If we hit something, set the "looking at" entity to it
	if len(i) != 0 {
		object = i[0].Object.GetNode()

		if entity := EntityOf(object); entity != nil {
			LookingAt = entity
		}
	}
}

package game

import (
	"github.com/g3n/engine/math32"
)

const lookSensitivityX float32 = 0.025
const lookSensitivityY float32 = 0.0373

var PlayerLookX float32
var PlayerLookY float32
var PlayerMoveX float32
var PlayerMoveZ float32

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

	PlayerLookX = 0
	PlayerLookY = 0

	// Movement
	Cam.TranslateX(PlayerMoveX)
	Cam.TranslateZ(PlayerMoveZ)
}

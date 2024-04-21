package player

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/context"
)

const (
	lookSensitivityX float32 = 0.025
	lookSensitivityY float32 = 0.015
)

type Player struct {
	*context.Context

	// Movement offset inpus
	lookX float32
	lookY float32
	moveX float32
	moveZ float32

	// The rotation offset for looking around
	rotationY float32
	rotationX float32
}

func New(ctx *context.Context) *Player {
	return &Player{Context: ctx}
}

// Update the player
func (p *Player) Update(deltaTime float32) {
	// Looking
	θ := p.Cam.Rotation().Y
	p.rotationX = θ + p.lookX*lookSensitivityX
	p.rotationY += p.lookY * lookSensitivityY
	p.rotationY = math32.Clamp(p.rotationY, -math32.Pi/2, math32.Pi/2)

	p.Cam.SetRotation(-p.rotationY, p.rotationX, 0)

	// Update which entity the player is looking at
	r := collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	r.SetFromCamera(p.Cam, 0, 0)
	i := r.IntersectObject(p.Scene, true)

	p.State.LookingAt = nil
	var object *core.Node

	// If we hit something, set the "looking at" entity to it
	if len(i) != 0 {
		object = i[0].Object.GetNode()

		if entity := p.State.EntityOf(object); entity != nil {
			p.State.LookingAt = entity
		}
	}

	p.lookX = 0
	p.lookY = 0

	// Movement
	p.Cam.TranslateX(p.moveX)
	p.Cam.TranslateZ(p.moveZ)
}

func (p *Player) Move(dx, dz float32) {
	p.moveX = dx
	p.moveZ = dz
}

func (p *Player) Look(lx, ly float32) {
	p.lookX = lx
	p.lookY = ly
}

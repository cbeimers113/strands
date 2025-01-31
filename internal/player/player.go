package player

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/math32"

	"cbeimers113/strands/internal/context"
)

const (
	// How far from the edge of the map we can be
	mapBound float32 = 50

	// Keep us this far from the map bounds to prevent getting stuck
	mapBoundPadding float32 = 0.1
)

type Player struct {
	*context.Context

	// Movement offset inputs
	lookX float32
	lookY float32
	moveX float32
	moveZ float32
	moveY float32

	// The rotation offset for looking around
	rotationY float32
	rotationX float32
}

func New(ctx *context.Context) *Player {
	return &Player{Context: ctx}
}

// Update the player
func (p *Player) Update(deltaTime float32, maxX, maxZ float32, centre math32.Vector3) {
	// Looking
	p.rotationX = p.Cam.Rotation().Y
	p.rotationX += p.lookX * p.Cfg.Controls.MouseSensitivityX / 10 * deltaTime
	p.rotationY += p.lookY * p.Cfg.Controls.MouseSensitivityY / 10 * deltaTime
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

	// Movement: move on each axis separately
	p.Cam.TranslateX(p.moveX * deltaTime)
	p.Cam.TranslateZ(p.moveZ * deltaTime)

	// Y translation should be relative to the world and not the camera's rotation
	p.Cam.SetPositionY(p.Cam.Position().Y + p.moveY*deltaTime)

	pos := p.Cam.Position()

	// Check new position, if we are OOB on any axis, put us back in bounds
	boundMin := -mapBound + mapBoundPadding
	boundMax := mapBound - mapBoundPadding
	if pos.X < -mapBound || pos.X > maxX+mapBound {
		nx := math32.Clamp(pos.X, boundMin, maxX+boundMax)
		p.Cam.SetPositionX(nx)
	}
	if pos.Y < -mapBound/2 || pos.Y > mapBound*2 {
		ny := math32.Clamp(pos.Y, boundMin/2, boundMax*2)
		p.Cam.SetPositionY(ny)
	}
	if pos.Z < -mapBound || pos.Z > maxZ+mapBound {
		nz := math32.Clamp(pos.Z, boundMin, maxZ+boundMax)
		p.Cam.SetPositionZ(nz)
	}

	// Player look deceleration
	p.lookX *= 0.75
	p.lookY *= 0.75

	// Player movement deceleration
	p.moveX *= 0.75
	p.moveZ *= 0.75
}

func (p *Player) MoveHorizontal(dx, dz float32) {
	p.moveX += dx
	p.moveZ += dz
}

func (p *Player) MoveVertical(dy float32) {
	p.moveY = dy
}

func (p *Player) Look(lx, ly float32) {
	p.lookX += lx
	p.lookY += ly
}

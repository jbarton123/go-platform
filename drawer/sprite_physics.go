package drawer

import (
	"github.com/faiface/pixel"
)

type Physics interface {
	Update(dt float64, ctrl pixel.Vec, platforms []Platform)
}

type SpritePhysics struct {
	Gravity   float64
	RunSpeed  float64
	JumpSpeed float64
	Rect      pixel.Rect
	Vel       pixel.Vec
	Ground    bool
}

func NewSpritePhysicsUpdater(gravity, runSpeed, jumpSpeed float64, rect pixel.Rect, vel pixel.Vec, ground bool) *SpritePhysics {
	return &SpritePhysics{
		Gravity:   gravity,
		RunSpeed:  runSpeed,
		JumpSpeed: jumpSpeed,
		Rect:      rect,
		Vel:       vel,
		Ground:    ground,
	}
}

func (gp *SpritePhysics) Update(dt float64, ctrl pixel.Vec, platforms []Platform) {
	// apply controls
	switch {
	case ctrl.X < 0:
		gp.Vel.X = -gp.RunSpeed
	case ctrl.X > 0:
		gp.Vel.X = +gp.RunSpeed
	default:
		gp.Vel.X = 0
	}

	// apply gravity and velocity
	gp.Vel.Y += gp.Gravity * dt
	gp.Rect = gp.Rect.Moved(gp.Vel.Scaled(dt))

	// check collisions against each platform
	gp.Ground = false
	if gp.Vel.Y <= 0 {
		for _, p := range platforms {
			if gp.Rect.Max.X <= p.Rect.Min.X || gp.Rect.Min.X >= p.Rect.Max.X {
				continue
			}
			if gp.Rect.Min.Y > p.Rect.Max.Y || gp.Rect.Min.Y < p.Rect.Max.Y+gp.Vel.Y*dt {
				continue
			}
			gp.Vel.Y = 0
			gp.Rect = gp.Rect.Moved(pixel.V(0, p.Rect.Max.Y-gp.Rect.Min.Y))
			gp.Ground = true
		}
	}

	// jump if on the ground and the player wants to jump
	if gp.Ground && ctrl.Y > 0 {
		gp.Vel.Y = gp.JumpSpeed
	}
}

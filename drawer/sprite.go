package drawer

import (
	"github.com/faiface/pixel"
	"math"
)

const (
	idle AnimState = iota
	running
	jumping
)

type AnimState int

type SpriteDrawer interface {
	Draw(target pixel.Target, physics *SpritePhysics)
	Update(dt float64, physics *SpritePhysics)
}

type spriteAnimation struct {
	sheet   pixel.Picture
	anims   map[string][]pixel.Rect
	rate    float64
	state   AnimState
	counter float64
	dir     float64
	frame   pixel.Rect
	sprite  *pixel.Sprite
}

func NewSpriteDrawer(sheet pixel.Picture, anims map[string][]pixel.Rect, rate float64, state AnimState, counter, dir float64, frame pixel.Rect) SpriteDrawer {
	return &spriteAnimation{
		sheet:   sheet,
		anims:   anims,
		rate:    rate,
		state:   state,
		counter: counter,
		dir:     dir,
		frame:   frame,
	}
}

func (ga *spriteAnimation) Draw(target pixel.Target, physics *SpritePhysics) {
	if ga.sprite == nil {
		ga.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}

	// draw the correct frame with the correct position and direction
	ga.sprite.Set(ga.sheet, ga.frame)
	ga.sprite.Draw(target, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			physics.Rect.W()/ga.sprite.Frame().W(),
			physics.Rect.H()/ga.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(-ga.dir, 1)).
		Moved(physics.Rect.Center()),
	)
}

func (ga *spriteAnimation) Update(dt float64, physics *SpritePhysics) {
	ga.counter += dt

	// determine the new animation state
	var newState AnimState
	switch {
	case !physics.Ground:
		newState = jumping
	case physics.Vel.Len() == 0:
		newState = idle
	case physics.Vel.Len() > 0:
		newState = running
	}

	// reset the time counter if the state changed
	if ga.state != newState {
		ga.state = newState
		ga.counter = 0
	}

	// determine the correct animation frame
	switch ga.state {
	case idle:
		ga.frame = ga.anims["Front"][0]
	case running:
		i := int(math.Floor(ga.counter / ga.rate))
		ga.frame = ga.anims["Run"][i%len(ga.anims["Run"])]
	case jumping:
		speed := physics.Vel.Y
		i := int((-speed/physics.JumpSpeed + 1) / 2 * float64(len(ga.anims["Jump"])))
		if i < 0 {
			i = 0
		}
		if i >= len(ga.anims["Jump"]) {
			i = len(ga.anims["Jump"]) - 1
		}
		ga.frame = ga.anims["Jump"][i]
	}

	// set the facing direction of the gopher
	if physics.Vel.X != 0 {
		if physics.Vel.X > 0 {
			ga.dir = +1
		} else {
			ga.dir = -1
		}
	}
}

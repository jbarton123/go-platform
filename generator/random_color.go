package generator

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

type randomcolor struct {
}

func NewRandomColor() *randomcolor {
	return &randomcolor{}
}

func (*randomcolor) Generate() pixel.RGBA {
again:
	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	length := math.Sqrt(r*r + g*g + b*b)

	if length == 0 {
		goto again
	}

	return pixel.RGB(r/length, g/length, b/length)
}

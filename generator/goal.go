package generator

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type GoalGenerator interface {
	Generate(imd *imdraw.IMDraw)
	Update(dt float64, randomColor pixel.RGBA)
}

type goal struct {
	pos    pixel.Vec
	radius float64
	step   float64

	counter float64
	cols    map[int]pixel.RGBA
}

func NewGoalGenerator(pos pixel.Vec, radius float64, step float64, counter float64, cols map[int]pixel.RGBA) GoalGenerator {
	return &goal{
		pos:     pos,
		radius:  radius,
		step:    step,
		counter: counter,
		cols:    cols,
	}
}

func (g *goal) Generate(imd *imdraw.IMDraw) {
	for i := len(g.cols) - 1; i >= 0; i-- {
		imd.Color = g.cols[i]
		imd.Push(g.pos)
		imd.Circle(float64(i+1)*g.radius/float64(len(g.cols)), 0)
	}
}

func (g *goal) Update(dt float64, randomColor pixel.RGBA) {
	g.counter += dt

	for g.counter > g.step {
		g.counter -= g.step
		for i := len(g.cols) - 2; i >= 0; i-- {
			g.cols[i+1] = g.cols[i]
		}
		g.cols[0] = randomColor
	}
}
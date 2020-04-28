package drawer

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
)

type PlatformDrawer interface {
	Draw(imd *imdraw.IMDraw)
}

type Platform struct {
	Rect  pixel.Rect
	Color color.Color
}

func (p *Platform) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.Color
	imd.Push(p.Rect.Min, p.Rect.Max)
	imd.Rectangle(0)
}

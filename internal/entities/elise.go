package entities

import (
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

const eliseWidth = 20.0
const eliseHeight = 100.0

type Elise struct {
	pos pixel.Vec
}

func (elise Elise) HitBox() EntityHitBox {
	min := elise.pos.Add(pixel.V(-eliseWidth/2, 0))
	max := elise.pos.Add(pixel.V(eliseWidth/2, eliseHeight))
	return EntityHitBox{
		Entity: 2,
		HitBox: pixel.Rect{min, max},
	}
}

func MakeElise(position pixel.Vec) Elise {
	return Elise{pos: position}
}

func (elise Elise) Tick(gameTimeMs float64) Entity {
	return elise
}

func (elise Elise) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(elise.pos.Add(pixel.V(0, eliseHeight/2)),
		draw.Image(*imageMap, internal.IEliseWalk1))
}

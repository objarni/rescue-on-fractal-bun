package entities

import (
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

type Elise struct {
	pos pixel.Vec
}

func MakeElise(position pixel.Vec) Elise {
	return Elise{pos: position}
}

func (elise Elise) Tick(gameTimeMs float64) Entity {
	return elise
}

func (elise Elise) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(elise.pos,
		draw.Image(*imageMap, internal.IEliseWalk1))
}

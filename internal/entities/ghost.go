package entities

import (
	"github.com/faiface/pixel"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

type Ghost struct {
	pos      pixel.Vec
	baseLine float64
}

func (ghost Ghost) HitBox() EntityHitBox {
	return EntityHitBox{
		Entity: 1,
		HitBox: pixel.Rect{},
	}
}

func (ghost Ghost) Tick(gameTimeMs float64) Entity {
	return Ghost{
		pos:      internal.V(ghost.pos.X, ghost.baseLine+math.Sin(gameTimeMs/300.0)*50),
		baseLine: ghost.baseLine,
	}
}

func (ghost Ghost) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(ghost.pos,
		draw.Image(*imageMap, internal.IGhost))
}

func MakeGhost(position pixel.Vec) Entity {
	return Ghost{pos: position, baseLine: position.Y}
}

/* notes ghost/elise behaviour
when ghost does overlap elise
when ghost does not overlap elise
when ghost overlaps light box
*/

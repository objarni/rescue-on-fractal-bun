package entities

import (
	"github.com/faiface/pixel"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

const ghostWidth = 50
const ghostHeight = 125

type Ghost struct {
	pos        pixel.Vec
	baseLine   float64
	dirX       float64
	minX, maxX float64
}

func (ghost Ghost) Handle(_ EventBox) Entity {
	return ghost
}

func (ghost Ghost) HitBox() pixel.Rect {
	min := ghost.pos.Add(pixel.V(-ghostWidth/2, 0))
	max := ghost.pos.Add(pixel.V(ghostWidth/2, ghostHeight))
	return pixel.Rect{Min: min, Max: max}
}

func (ghost Ghost) Tick(gameTimeMs float64, receiver EventBoxReceiver) Entity {
	receiver.AddEventBox(EventBox{
		Event: "DAMAGE",
		Box:   ghost.HitBox(),
	})
	ghost.pos.X += ghost.dirX
	if ghost.pos.X > ghost.maxX {
		ghost.dirX = -1
	}
	if ghost.pos.X < ghost.minX {
		ghost.dirX = 1
	}
	ghost.pos = internal.V(ghost.pos.X, ghost.baseLine+math.Sin(gameTimeMs/300.0)*50)
	return ghost
}

func (ghost Ghost) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(ghost.pos.Add(pixel.V(0, ghostHeight/2)),
		draw.Image(*imageMap, internal.IGhost))
}

func MakeGhost(area pixel.Rect) Entity {
	startPos := area.Center()
	return Ghost{
		pos:      startPos,
		baseLine: startPos.Y,
		dirX:     1,
		minX:     area.Min.X,
		maxX:     area.Max.X,
	}
}

/* notes ghost/elise behaviour
when ghost does overlap elise
when ghost does not overlap elise
when ghost overlaps light box
*/

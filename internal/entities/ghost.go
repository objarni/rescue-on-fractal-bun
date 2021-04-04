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
	gameTimeMs float64
}

func (ghost Ghost) Handle(_ EventBox) Entity {
	return ghost
}

func (ghost Ghost) HitBox() pixel.Rect {
	min := ghost.pos.Add(pixel.V(-ghostWidth/2, 0))
	max := ghost.pos.Add(pixel.V(ghostWidth/2, ghostHeight))
	rect := pixel.Rect{min, max}
	return rect
}

func (ghost Ghost) Tick(receiver EventBoxReceiver) Entity {
	receiver.AddEventBox(EventBox{
		Event: "DAMAGE",
		Box:   ghost.HitBox(),
	})
	ghost.gameTimeMs += internal.TickTimeMs
	ghost.pos = internal.V(ghost.pos.X, ghost.baseLine+math.Sin(ghost.gameTimeMs/300.0)*50)
	return ghost
}

func (ghost Ghost) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(ghost.pos.Add(pixel.V(0, ghostHeight/2)),
		draw.Image(*imageMap, internal.IGhost))
}

func MakeGhost(position pixel.Vec) Entity {
	return Ghost{pos: position, baseLine: position.Y, gameTimeMs: 0}
}

/* notes ghost/elise behaviour
when ghost does overlap elise
when ghost does not overlap elise
when ghost overlaps light box
*/

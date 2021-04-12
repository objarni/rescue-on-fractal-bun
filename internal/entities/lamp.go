package entities

import (
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

const lampWidth = 95
const lampHeight = 300

type Lamp struct {
	pos pixel.Vec
}

func (lamp Lamp) Handle(EventBox) Entity {
	return lamp
}

func (lamp Lamp) HitBox() pixel.Rect {
	min := lamp.pos.Add(pixel.V(-lampWidth/2, 0))
	max := lamp.pos.Add(pixel.V(lampWidth/2, lampHeight))
	rect := pixel.Rect{Min: min, Max: max}
	return rect
}

func (lamp Lamp) Tick(_ float64, _ EventBoxReceiver) Entity {
	return lamp
}

func (lamp Lamp) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(lamp.pos.Add(pixel.V(0, lampHeight/2)),
		draw.Image(*imageMap, internal.IStreetLight))
}

func MakeLamp(area pixel.Rect) Entity {
	x := area.Center().X
	y := area.Min.Y
	return Lamp{pos: pixel.V(x, y)}
}

/* notes lamp/elise behaviour
when lamp does overlap elise
when lamp does not overlap elise
when lamp overlaps light box
*/

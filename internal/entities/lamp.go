package entities

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	d "objarni/rescue-on-fractal-bun/internal/draw"
	"objarni/rescue-on-fractal-bun/internal/events"
)

const lampWidth = 95
const lampHeight = 300

type Lamp struct {
	pos pixel.Vec
	on  bool
}

func (lamp Lamp) String() string {
	panic("implement me")
}

func (lamp Lamp) Handle(eb EventBox) Entity {
	if eb.Event == events.ButtonPressed {
		lamp.on = !lamp.on
	}
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

func (lamp Lamp) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	image := d.Image(*imageMap, internal.IStreetLight)
	if !lamp.on {
		image = d.Color(colornames.Black, image)
	}
	return d.Moved(lamp.pos.Add(pixel.V(0, lampHeight/2)),
		image)
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

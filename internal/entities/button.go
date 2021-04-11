package entities

import (
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

const buttonWidth = 20
const buttonHeight = 50

type Button struct {
	pos pixel.Vec
}

func (button Button) Handle(EventBox) Entity {
	return button
}

func (button Button) HitBox() pixel.Rect {
	min := button.pos.Add(pixel.V(-buttonWidth/2, 0))
	max := button.pos.Add(pixel.V(buttonWidth/2, buttonHeight))
	rect := pixel.Rect{Min: min, Max: max}
	return rect
}

func (button Button) Tick(_ float64, _ EventBoxReceiver) Entity {
	return button
}

func (button Button) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(button.pos.Add(pixel.V(0, buttonHeight/2)),
		draw.Image(*imageMap, internal.IButton))
}

func MakeButton(area pixel.Rect) Entity {
	x := area.Center().X
	y := area.Min.Y
	return Button{pos: pixel.V(x, y)}
}

/* notes button/elise behaviour
when button does overlap elise
when button does not overlap elise
when button overlaps light box
*/

package entities

import (
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

const buttonWidth = 20
const buttonHeight = 50

type Button struct {
	pos     px.Vec
	pressed bool
}

func (button Button) Handle(eb EventBox) Entity {
	if eb.Event == "ACTION" {
		button.pressed = true
	}
	return button
}

func (button Button) HitBox() px.Rect {
	min := button.pos.Add(px.V(-buttonWidth/2, 0))
	max := button.pos.Add(px.V(buttonWidth/2, buttonHeight))
	rect := px.Rect{Min: min, Max: max}
	return rect
}

func (button Button) Tick(_ float64, ebr EventBoxReceiver) Entity {
	if button.pressed {
		button.pressed = false
		ebr.AddEventBox(EventBox{
			Event: "BUTTON_PRESSED",
			Box:   button.HitBox().Resized(button.HitBox().Center(), px.V(500, 500)),
		})
	}
	return button
}

func (button Button) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(button.pos.Add(px.V(0, buttonHeight/2)),
		draw.Image(*imageMap, internal.IButton))
}

func MakeButton(area px.Rect) Entity {
	x := area.Center().X
	y := area.Min.Y
	return Button{pos: px.V(x, y)}
}

/* notes button/elise behaviour
when button does overlap elise
when button does not overlap elise
when button overlaps light box
*/

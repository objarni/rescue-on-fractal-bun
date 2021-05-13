package entities

import (
	px "github.com/faiface/pixel"
	d "github.com/objarni/pixelop"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/events"
)

const buttonWidth = 20
const buttonHeight = 50

type Button struct {
	pos     px.Vec
	pressed bool
}

func (button Button) String() string {
	panic("implement me")
}

func (button Button) Handle(eb EventBox) Entity {
	if eb.Event == events.Action {
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
			Event: events.ButtonPressed,
			Box:   button.HitBox().Resized(button.HitBox().Center(), px.V(500, 500)),
		})
	}
	return button
}

func (button Button) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	return d.Moved(button.pos.Add(px.V(0, buttonHeight/2)),
		d.Image((*imageMap)[internal.IButton], internal.IButton.String()))
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

package entities

import (
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

type EventBoxReceiver interface {
	AddEventBox(EventBox)
}

type Entity interface {
	GfxOp(imageMap *internal.ImageMap) draw.WinOp
	Tick(ebReceiver EventBoxReceiver) Entity
	HitBoxRect() pixel.Rect
	Handle(eb EventBox) Entity
}

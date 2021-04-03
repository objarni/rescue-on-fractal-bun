package entities

import (
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

type EventBoxReceiver interface {
	AddEventBox(EventBox)
}

type Entity interface {
	GfxOp(imageMap *internal.ImageMap) draw.WinOp
	Tick(ebReceiver EventBoxReceiver) Entity
	HitBox() EntityHitBox
	Handle(eb EventBox) Entity
}

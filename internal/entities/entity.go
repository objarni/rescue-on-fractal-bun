package entities

import (
	"github.com/faiface/pixel"
	d "github.com/objarni/pixelop"
	"objarni/rescue-on-fractal-bun/internal"
)

type EventBoxReceiver interface {
	AddEventBox(EventBox)
}

type Entity interface {
	GfxOp(imageMap *internal.ImageMap) d.WinOp
	// Tick TODO: how can entities communicate 'remove me' to scene?
	Tick(gameTimeMs float64, ebReceiver EventBoxReceiver) Entity
	HitBox() pixel.Rect
	HitBoxes() []pixel.Rect
	Handle(eb EventBox) Entity
	String() string
}

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
	// TODO: nil from Tick means remove entity from canvas/world
	Tick(gameTimeMs float64, ebReceiver EventBoxReceiver) Entity
	HitBox() pixel.Rect
	HitBoxes() []pixel.Rect
	Handle(eb EventBox) Entity
	String() string
}

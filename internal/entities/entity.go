package entities

import (
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

type Entity interface {
	GfxOp(imageMap *internal.ImageMap) draw.WinOp
	// TODO: remove gameTimeMs, ticks are constant!
	Tick(gameTimeMs float64) Entity
	HitBox() EntityHitBox
}

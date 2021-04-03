package entities

import (
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

type Entity interface {
	GfxOp(imageMap *internal.ImageMap) draw.WinOp
	Tick() Entity
	HitBox() EntityHitBox
}

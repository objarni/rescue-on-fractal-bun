package entities

import (
	px "github.com/faiface/pixel"
	d "github.com/objarni/pixelop"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/printers"
)

type LevelBoundary struct {
	areaRect px.Rect
}

func (levelBoundary LevelBoundary) HitBoxes() []px.Rect {
	panic("implement me")
}

func (levelBoundary LevelBoundary) String() string {
	return "LevelBoundary: " + printers.PrintRect(levelBoundary.areaRect)
}

func (levelBoundary LevelBoundary) Handle(_ EventBox) Entity {
	return levelBoundary
}

func (levelBoundary LevelBoundary) HitBox() px.Rect {
	return px.ZR
}

func (levelBoundary LevelBoundary) Tick(_ float64, _ EventBoxReceiver) Entity {
	return levelBoundary
}

func (levelBoundary LevelBoundary) GfxOp(_ *internal.ImageMap) d.WinOp {
	return d.OpSequence()
}

func MakeLevelBoundary(areaRect px.Rect) Entity {
	return LevelBoundary{
		areaRect: areaRect,
	}
}

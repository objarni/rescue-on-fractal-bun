package entities

import (
	"fmt"
	px "github.com/faiface/pixel"
	d "github.com/objarni/pixelop"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/events"
	"objarni/rescue-on-fractal-bun/internal/printers"
)

type LevelBoundary struct {
	areaRect px.Rect
}

func (lb LevelBoundary) HitBoxes() []px.Rect {
	panic("implement me")
}

func (lb LevelBoundary) String() string {
	return "LevelBoundary: " + printers.PrintRect(lb.areaRect)
}

func (lb LevelBoundary) Handle(_ EventBox) Entity {
	return lb
}

func (lb LevelBoundary) HitBox() px.Rect {
	return px.ZR
}

func (lb LevelBoundary) bottomLeft() px.Vec {
	return lb.areaRect.Min
}

func (lb LevelBoundary) topLeft() px.Vec {
	rect := lb.areaRect
	return px.V(rect.Min.X, rect.Max.Y)
}

func (lb LevelBoundary) topRight() px.Vec {
	return lb.areaRect.Max
}

func (lb LevelBoundary) bottomRight() px.Vec {
	rect := lb.areaRect
	return px.V(rect.Max.X, rect.Min.Y)
}

func (lb LevelBoundary) Tick(_ float64, ebr EventBoxReceiver) Entity {
	for ix, r := range lb.borderRectangles() {
		if r.Area() <= 0 {
			panic(fmt.Sprintf("border rectangle without area %d,%v", ix, r))
		}
		ebr.AddEventBox(EventBox{
			Event: events.Wall,
			Box:   r,
		})
	}
	return lb
}

func (lb LevelBoundary) GfxOp(_ *internal.ImageMap) d.WinOp {
	rSeq := d.Colored(
		colornames.Aquamarine,
		d.ImdOpSequence(
			d.Rectangle(lb.bottomLeft(), lb.topLeft(), 2),
			d.Rectangle(lb.topLeft(), lb.topRight(), 2),
			d.Rectangle(lb.topRight(), lb.bottomRight(), 2),
			d.Rectangle(lb.bottomRight(), lb.bottomLeft(), 2),
		),
	)
	return d.OpSequence(d.ToWinOp(rSeq))
}

func (lb LevelBoundary) borderRectangles() []px.Rect {
	thickness := 50.0
	return []px.Rect{
		{Min: lb.bottomLeft().Add(px.V(-thickness, 0)), Max: lb.topLeft()},
		{Min: lb.topLeft(), Max: lb.topRight().Add(px.V(0, thickness))},
		{Max: lb.topRight().Add(px.V(thickness, 0)), Min: lb.bottomRight()},
		{Min: lb.bottomLeft().Add(px.V(0, -thickness)), Max: lb.bottomRight()},
	}
}

func MakeLevelBoundary(areaRect px.Rect) Entity {
	return LevelBoundary{
		areaRect: areaRect,
	}
}

package entities

import (
	px "github.com/faiface/pixel"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	d "objarni/rescue-on-fractal-bun/internal/draw"
	"objarni/rescue-on-fractal-bun/internal/events"
)

const ghostWidth = 50
const ghostHeight = 125

type Ghost struct {
	pos         px.Vec
	baseLine    float64
	dirX        float64
	minX, maxX  float64
	curveHeight float64
}

func (ghost Ghost) String() string {
	panic("implement me")
}

func (ghost Ghost) Handle(_ EventBox) Entity {
	return ghost
}

func (ghost Ghost) HitBox() px.Rect {
	min := ghost.pos.Add(px.V(-ghostWidth/2, 0))
	max := ghost.pos.Add(px.V(ghostWidth/2, ghostHeight))
	return px.Rect{Min: min, Max: max}
}

func (ghost Ghost) Tick(gameTimeMs float64, receiver EventBoxReceiver) Entity {
	receiver.AddEventBox(EventBox{
		Event: events.Damage,
		Box:   ghost.HitBox(),
	})
	ghost.pos.X += ghost.dirX
	if ghost.pos.X > ghost.maxX {
		ghost.dirX = -1
	}
	if ghost.pos.X < ghost.minX {
		ghost.dirX = 1
	}
	yPos := ghost.baseLine + math.Sin(gameTimeMs/300.0)*ghost.curveHeight
	ghost.pos = internal.V(ghost.pos.X, yPos)
	return ghost
}

func (ghost Ghost) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	image := d.Image(*imageMap, internal.IGhost)
	if ghost.dirX < 0 {
		image = d.Mirrored(image)
	}
	return d.Moved(ghost.pos.Add(px.V(0, ghostHeight/2)),
		image)
}

func MakeGhost(area px.Rect) Entity {
	startPos := area.Center()
	return Ghost{
		pos:         startPos,
		baseLine:    startPos.Y,
		dirX:        1,
		minX:        area.Min.X,
		maxX:        area.Max.X,
		curveHeight: (area.Max.Y - area.Min.Y) / 2,
	}
}

/* notes ghost/elise behaviour
when ghost does overlap elise
when ghost does not overlap elise
when ghost overlaps light box
*/

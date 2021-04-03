package entities

import (
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

const eliseWidth = 20.0
const eliseHeight = 100.0

type Elise struct {
	pos                       pixel.Vec
	leftPressed, rightPressed bool
}

func (elise Elise) HitBox() EntityHitBox {
	min := elise.pos.Add(pixel.V(-eliseWidth/2, 0))
	max := elise.pos.Add(pixel.V(eliseWidth/2, eliseHeight))
	return EntityHitBox{
		Entity: 2,
		HitBox: pixel.Rect{min, max},
	}
}

func MakeElise(position pixel.Vec) Elise {
	elise := Elise{pos: position}
	return elise
}

func (elise Elise) Tick(gameTimeMs float64) Elise {
	eliseMoveSpeed := 1.2
	if elise.leftPressed && !elise.rightPressed {
		elise.pos = elise.pos.Add(internal.V(-eliseMoveSpeed, 0))
	}
	if !elise.leftPressed && elise.rightPressed {
		elise.pos = elise.pos.Add(internal.V(eliseMoveSpeed, 0))
	}
	return elise
}

func (elise Elise) HandleKeyDown(key internal.ControlKey) Elise {
	if key == internal.Left {
		elise.leftPressed = true
	}
	if key == internal.Right {
		elise.rightPressed = true
	}
	return elise
}

func (elise Elise) HandleKeyUp(key internal.ControlKey) Elise {
	if key == internal.Left {
		elise.leftPressed = false
	}
	if key == internal.Right {
		elise.rightPressed = false
	}
	return elise
}

func (elise Elise) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	return draw.Moved(elise.pos.Add(pixel.V(0, eliseHeight/2)),
		draw.Image(*imageMap, internal.IEliseCrouch))
}

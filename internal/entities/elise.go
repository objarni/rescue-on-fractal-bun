package entities

import (
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

const eliseWidth = 20.0
const eliseHeight = 100.0

type Elise struct {
	Pos                       pixel.Vec
	leftPressed, rightPressed bool
	gameTimeMs                int
}

func (elise Elise) HitBox() pixel.Rect {
	min := elise.Pos.Add(pixel.V(-eliseWidth/2, 0))
	max := elise.Pos.Add(pixel.V(eliseWidth/2, eliseHeight))
	rect := pixel.Rect{min, max}
	return rect
}

func MakeElise(position pixel.Vec) Elise {
	elise := Elise{Pos: position}
	return elise
}

func (elise Elise) Tick() Elise {
	elise.gameTimeMs += 5
	eliseMoveSpeed := 1.2
	if elise.leftPressed && !elise.rightPressed {
		elise.Pos = elise.Pos.Add(internal.V(-eliseMoveSpeed, 0))
	}
	if !elise.leftPressed && elise.rightPressed {
		elise.Pos = elise.Pos.Add(internal.V(eliseMoveSpeed, 0))
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
	image := EliseWalkFrame(float64(elise.gameTimeMs)/1000.0, 10)
	return draw.Moved(elise.Pos.Add(pixel.V(0, eliseHeight/2)),
		draw.Image(*imageMap, image))
}

func (elise Elise) Handle(eb EventBox) Elise {
	if eb.Event == "DAMAGE" {
		elise.Pos = elise.Pos.Add(pixel.V(5, 0))
	}
	if eb.Event == "LEFT_DOWN" {
		elise.leftPressed = true
	}
	if eb.Event == "RIGHT_DOWN" {
		elise.rightPressed = true
	}
	return elise
}

var frames = [...]internal.Image{
	internal.IEliseWalk6,
	internal.IEliseWalk5,
	internal.IEliseWalk4,
	internal.IEliseWalk3,
	internal.IEliseWalk2,
	internal.IEliseWalk1,
}

func EliseWalkFrame(gameTimeS float64, targetFPS int) internal.Image {
	var eliseAnimation = internal.Animation{
		Frames:    6,
		TargetFPS: targetFPS,
	}
	frame := eliseAnimation.FrameAtTime(gameTimeS)
	image := frames[frame]
	return image
}

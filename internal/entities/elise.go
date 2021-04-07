package entities

import (
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

const eliseWidth = 20.0
const eliseHeight = 100.0

type Elise struct {
	Pos                       px.Vec
	leftPressed, rightPressed bool
	gameTimeMs                int
	flip                      bool
	actionDown                bool
}

func (elise Elise) HitBox() px.Rect {
	min := elise.Pos.Add(px.V(-eliseWidth/2, 0))
	max := elise.Pos.Add(px.V(eliseWidth/2, eliseHeight))
	rect := px.Rect{Min: min, Max: max}
	return rect
}

func MakeElise(position px.Vec) Entity {
	return Elise{Pos: position}
}

func (elise Elise) Tick(eb EventBoxReceiver) Entity {
	elise.gameTimeMs += 5
	eliseMoveSpeed := 1.2
	if elise.leftPressed && !elise.rightPressed {
		elise.flip = true
		elise.Pos = elise.Pos.Add(internal.V(-eliseMoveSpeed, 0))
	}
	if !elise.leftPressed && elise.rightPressed {
		elise.flip = false
		elise.Pos = elise.Pos.Add(internal.V(eliseMoveSpeed, 0))
	}
	if elise.actionDown {
		elise.actionDown = false
		hitBox := elise.HitBox()
		eb.AddEventBox(EventBox{
			Event: "ACTION",
			Box:   hitBox.Resized(hitBox.Center(), px.V(40, 40)),
		})
	}
	return elise
}

func (elise Elise) GfxOp(imageMap *internal.ImageMap) draw.WinOp {
	image := EliseWalkFrame(float64(elise.gameTimeMs)/1000.0, 10)
	imgOp := draw.Image(*imageMap, image)
	if elise.flip {
		imgOp = draw.Mirrored(imgOp)
	}
	return draw.Moved(elise.Pos.Add(px.V(0, eliseHeight/2)), imgOp)
}

func (elise Elise) Handle(eb EventBox) Entity {
	if eb.Event == "DAMAGE" {
		elise.Pos = elise.Pos.Add(px.V(5, 0))
	}
	if eb.Event == "LEFT_DOWN" {
		elise.leftPressed = true
	}
	if eb.Event == "RIGHT_DOWN" {
		elise.rightPressed = true
	}
	if eb.Event == "LEFT_UP" {
		elise.leftPressed = false
	}
	if eb.Event == "RIGHT_UP" {
		elise.rightPressed = false
	}
	if eb.Event == "ACTION_DOWN" {
		elise.actionDown = true
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

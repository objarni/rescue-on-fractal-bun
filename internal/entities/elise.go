package entities

import (
	"fmt"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	d "objarni/rescue-on-fractal-bun/internal/draw"
	"objarni/rescue-on-fractal-bun/internal/events"
	pr "objarni/rescue-on-fractal-bun/internal/printers"
	"strings"
)

const eliseWidth = 20.0
const eliseHeight = 100.0

type Elise struct {
	Pos, Vel                  px.Vec
	leftPressed, rightPressed bool
	gameTimeMs                float64
	flip                      bool
	actionDown                bool
}

func (elise Elise) String() string {
	state := fmt.Sprintf("Elise %v", "standing")
	hb := fmt.Sprintf("HitBox %v", pr.PrintRect(elise.HitBox()))
	pos := fmt.Sprintf("Pos: %v", pr.PrintVec(elise.Pos))
	vel := fmt.Sprintf("Vel: %v", pr.PrintVec(elise.Vel))
	facing := "right"
	if elise.flip {
		facing = "left"
	}
	facing = "Facing " + facing
	all := []string{state, hb, pos, vel, facing}
	return strings.Join(all, "\n") + "\n"
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

func (elise Elise) Tick(gameTimeMs float64, eventBoxReceiver EventBoxReceiver) Entity {
	elise.gameTimeMs = gameTimeMs
	eliseMoveSpeed := 1.2
	eliseGravity := -0.1
	if elise.leftPressed && !elise.rightPressed {
		elise.flip = true
		elise.Pos = elise.Pos.Add(internal.V(-eliseMoveSpeed, 0))
	}
	if !elise.leftPressed && elise.rightPressed {
		elise.flip = false
		elise.Pos = elise.Pos.Add(internal.V(eliseMoveSpeed, 0))
	}
	elise.Vel = elise.Vel.Add(px.V(0, eliseGravity))
	elise.Pos = elise.Pos.Add(elise.Vel)
	if elise.actionDown {
		elise.actionDown = false
		hitBox := elise.HitBox()
		eventBoxReceiver.AddEventBox(EventBox{
			Event: events.Action,
			Box:   hitBox.Resized(hitBox.Center(), px.V(40, 40)),
		})
	}
	return elise
}

func (elise Elise) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	image := internal.IEliseWalk2
	if elise.rightPressed || elise.leftPressed {
		image = EliseWalkFrame(elise.gameTimeMs/1000.0, 10)
	}
	imgOp := d.Image(*imageMap, image)
	if elise.flip {
		imgOp = d.Mirrored(imgOp)
	}
	return d.Moved(elise.Pos.Add(px.V(0, eliseHeight/2)), imgOp)
}

func (elise Elise) Handle(eb EventBox) Entity {
	if eb.Event == events.Damage {
		elise.Pos = elise.Pos.Add(px.V(5, 0))
	}
	if eb.Event == events.KeyLeftDown {
		elise.leftPressed = true
	}
	if eb.Event == events.KeyRightDown {
		elise.rightPressed = true
	}
	if eb.Event == events.KeyLeftUp {
		elise.leftPressed = false
	}
	if eb.Event == events.KeyRightUp {
		elise.rightPressed = false
	}
	if eb.Event == events.KeyActionDown {
		elise.actionDown = true
	}
	if eb.Event == events.Wall {
		overlap := eb.Box.Intersect(elise.HitBox())
		elise.Pos = elise.Pos.Add(px.V(0, overlap.H()))
		elise.Vel = px.ZV
	}
	return elise
}

var frames = [...]internal.Image{
	internal.IEliseJump7,
	internal.IEliseJump6,
	internal.IEliseJump5,
	internal.IEliseJump4,
	internal.IEliseJump3,
	internal.IEliseJump2,
	internal.IEliseJump1,
}

var walkFrames = [...]internal.Image{
	internal.IEliseWalk6,
	internal.IEliseWalk5,
	internal.IEliseWalk4,
	internal.IEliseWalk3,
	internal.IEliseWalk2,
	internal.IEliseWalk1,
}

func EliseWalkFrame(gameTimeS float64, targetFPS int) internal.Image {
	var eliseAnimation = internal.Animation{
		Frames:    len(walkFrames),
		TargetFPS: targetFPS,
	}
	frame := eliseAnimation.FrameAtTime(gameTimeS)
	image := walkFrames[frame]
	return image
}

package entities

import (
	"fmt"
	px "github.com/faiface/pixel"
	"math"
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
	facingLeft                bool
	actionDown                bool
	jumping                   bool
}

func MakeElise(position px.Vec) Entity {
	return Elise{Pos: position}
}

func (elise Elise) String() string {
	generalState := "standing"
	if elise.jumping {
		generalState = "jumping"
	}
	state := fmt.Sprintf("Elise %v", generalState)
	hb := fmt.Sprintf("HitBox %v", pr.PrintRect(elise.HitBox()))
	pos := fmt.Sprintf("Pos: %v", pr.PrintVec(elise.Pos))
	vel := fmt.Sprintf("Vel: %v", pr.PrintVec(elise.Vel))
	facing := "right"
	if elise.facingLeft {
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

func (elise Elise) Tick(gameTimeMs float64, eventBoxReceiver EventBoxReceiver) Entity {
	elise.gameTimeMs = gameTimeMs

	// Movement constants
	eliseWalkAcceleration := 0.1
	maxHorisontalSpeed := 1.2
	eliseGravity := -0.1

	// Horisontal acceleration
	directionSign := 0.0
	if elise.leftPressed && !elise.rightPressed {
		elise.facingLeft = true
		directionSign = -1.0
		//if math.Abs(elise.Vel.X) < maxHorisontalSpeed {
		//	elise.Vel = elise.Vel.Add(px.V(eliseWalkAcceleration, 0))
		//}
	}
	if !elise.leftPressed && elise.rightPressed {
		elise.facingLeft = false
		directionSign = 1.0
		//if math.Abs(elise.Vel.X) < maxHorisontalSpeed {
		//	elise.Vel = elise.Vel.Add(px.V(-eliseWalkAcceleration, 0))
		//}
	}
	if math.Abs(elise.Vel.X) < maxHorisontalSpeed {
		elise.Vel = elise.Vel.Add(px.V(eliseWalkAcceleration*directionSign, 0))
	}

	// Vertical acceleration
	elise.Vel = elise.Vel.Add(px.V(0, eliseGravity))

	// Update position from velocity
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
	if eb.Event == events.KeyJumpDown && !elise.jumping {
		elise.jumping = true
		elise.Vel = px.V(0, 100)
		elise.Pos = elise.Pos.Add(px.V(0, 30)) // Get out of any ground tile!
	}
	if eb.Event == events.Wall {
		overlap := eb.Box.Intersect(elise.HitBox())
		elise.jumping = false
		h := overlap.H()
		w := overlap.W()
		if h > w {
			overlapCenterX := overlap.Center().X
			if overlapCenterX > elise.Pos.X {
				elise.Pos = elise.Pos.Sub(px.V(w, 0))
			} else {
				elise.Pos = elise.Pos.Add(px.V(w, 0))
			}
			elise.Vel = px.V(0, elise.Vel.Y)
		} else {
			elise.Pos = elise.Pos.Add(px.V(0, h))
			elise.Vel = px.V(elise.Vel.X, 0)
		}
	}
	return elise
}

func (elise Elise) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	image := internal.IEliseWalk2
	if elise.rightPressed || elise.leftPressed {
		image = EliseWalkFrame(elise.gameTimeMs/1000.0, 10)
	}
	imgOp := d.Image(*imageMap, image)
	if elise.facingLeft {
		imgOp = d.Mirrored(imgOp)
	}
	return d.Moved(elise.Pos.Add(px.V(0, eliseHeight/2)), imgOp)
}

var _ = [...]internal.Image{
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

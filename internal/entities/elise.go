package entities

import (
	"fmt"
	px "github.com/faiface/pixel"
	d "github.com/objarni/pixelop"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/events"
	pr "objarni/rescue-on-fractal-bun/internal/printers"
	"strings"
)

const eliseWidth = 20.0
const eliseHeight = 100.0

type EliseState int

const (
	EliseStanding = iota
	EliseWalking
)

func (state EliseState) String() string {
	return []string{
		"Standing",
		"Walking",
	}[state]
}

type Elise struct {
	Pos, Vel                  px.Vec
	leftPressed, rightPressed bool
	gameTimeMs                float64
	facingLeft                bool
	actionDown                bool
	jumping                   int // impulse must last more than 1 tick
	state                     EliseState
}

func (elise Elise) HitBoxes() []px.Rect {
	panic("implement me")
}

func MakeElise(position px.Vec) Entity {
	return Elise{Pos: position, state: EliseStanding}
}

func (elise Elise) String() string {
	generalState := elise.state.String()
	state := fmt.Sprintf("Elise %v", generalState)
	vel := fmt.Sprintf("Vel: %v", pr.PrintVec(elise.Vel))
	gfx := fmt.Sprintf("Gfx:\n%s", elise.GfxOp(nil).String())
	facing := "right"
	if elise.facingLeft {
		facing = "left"
	}
	facing = "Facing " + facing
	all := []string{state, vel, facing, gfx}
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
	eliseWalkDeceleration := 0.04 // When gliding, but not pressing a direction
	maxHorisontalSpeed := 1.4
	eliseGravity := -0.1
	haltSpeed := 0.1

	// Horisontal acceleration
	horisontalSpeed := math.Abs(elise.Vel.X)
	direction := 0.0
	if elise.Vel.X > 0 {
		direction = 1.0
	}
	if elise.Vel.X < 0 {
		direction = -1.0
	}
	intention := 0.0
	if elise.leftPressed && !elise.rightPressed {
		elise.state = EliseWalking
		elise.facingLeft = true
		intention = -1.0
	} else if !elise.leftPressed && elise.rightPressed {
		elise.state = EliseWalking
		elise.facingLeft = false
		intention = 1.0
	} else {
		elise.state = EliseStanding
	}
	if intention == 0 {
		// Slow down Elise until halt when no movement intention

		// Slow enough to halt?
		if horisontalSpeed > haltSpeed {
			elise.Vel = elise.Vel.Add(px.V(-direction*eliseWalkDeceleration, 0))
		} else {
			elise.Vel = px.V(0, elise.Vel.Y)
		}
	}
	if math.Abs(elise.Vel.X) < maxHorisontalSpeed {
		elise.Vel = elise.Vel.Add(px.V(eliseWalkAcceleration*intention, 0))
	}

	// Vertical acceleration
	if elise.jumping > 0 {
		elise.jumping -= 1
		elise.Vel = elise.Vel.Add(px.V(0, 2))
		elise.Pos = elise.Pos.Add(px.V(0, 1)) // Get out of any ground tile!
	}
	isGrounded := elise.Pos.Y == 0
	if isGrounded {
		elise.Vel = px.V(elise.Vel.X, 0)
	} else {
		elise.Vel = elise.Vel.Add(px.V(0, eliseGravity))
	}

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
	if eb.Event == events.KeyJumpDown && elise.jumping == 0 {
		elise.jumping = 2
	}
	if eb.Event == events.Wall {
		overlap := eb.Box.Intersect(elise.HitBox())
		h := overlap.H()
		w := overlap.W()
		if h > w { // hit wall horisontally
			overlapCenterX := overlap.Center().X
			sign := 1.0
			if overlapCenterX > elise.Pos.X {
				sign = -1.0
			}
			elise.Pos = elise.Pos.Add(px.V(sign*w, 0))
			elise.Vel = px.V(0, elise.Vel.Y)
		} else { // hit wall vertically (ground)
			elise.Pos = elise.Pos.Add(px.V(0, h))
			elise.Vel = px.V(elise.Vel.X, 0)
		}
	}
	return elise
}

func (elise Elise) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	image := internal.IEliseWalk2
	if elise.state == EliseWalking {
		image = EliseWalkFrame(elise.gameTimeMs/1000.0, 10)
	}
	var sprite *px.Sprite = nil
	if imageMap != nil {
		sprite = (*imageMap)[image]
	}
	imgOp := d.Image(sprite, image.String())
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

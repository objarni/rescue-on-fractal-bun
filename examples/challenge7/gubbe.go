package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"objarni/rescue-on-fractal-bun/internal"
)

const MAXVELOCITY = 2
const ACCELERATION = 0.03
const DECCELERATION = 0.95

type State int

const (
	Standing State = iota
	Kicking
	Walking
)

type Looking int

const (
	Right Looking = iota
	Left
)

func (looking Looking) String() string {
	return [...]string{"Right", "Left"}[looking]
}

type Image int

const (
	WalkRight1 Image = iota
	WalkRight2
	StandingRight
	KickRight
)

func (image Image) String() string {
	return [...]string{
		"WalkRight1",
		"WalkRight2",
		"StandingRight",
		"KickRight",
	}[image]
}

func (state State) String() string {
	return [...]string{"Standing", "Kicking", "Walking"}[state]
}

type Controls struct {
	left, right, kick bool
}

type Gubbe struct {
	state    State
	looking  Looking
	image    Image
	counter  int
	pos      pixel.Vec
	vel      pixel.Vec
	acc      pixel.Vec
	controls Controls
}

func (gubbe *Gubbe) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Left {
		gubbe.controls.left = true
	}
	if key == internal.Right {
		gubbe.controls.right = true
	}
	return gubbe
}

func (gubbe *Gubbe) HandleKeyUp(key internal.ControlKey) internal.Thing {
	if key == internal.Left {
		gubbe.controls.left = false
	}
	if key == internal.Right {
		gubbe.controls.right = false
	}
	return gubbe
}

func (gubbe *Gubbe) Render(win *pixelgl.Window) {

}

func (gubbe *Gubbe) Tick() bool {
	// STATE DEPENDENT BEHAVIOR
	switch gubbe.state {
	case Standing:
		if gubbe.controls.right && !gubbe.controls.left {
			initWalking(gubbe, Right)
		}
		if gubbe.controls.left && !gubbe.controls.right {
			initWalking(gubbe, Left)
		}
		gubbe.vel = gubbe.vel.Scaled(DECCELERATION)
	case Walking:
		gubbe.counter++
		if gubbe.counter%30 == 0 {
			if gubbe.image == WalkRight1 {
				gubbe.image = WalkRight2
			} else {
				gubbe.image = WalkRight1
			}
		}
		// No directions, or ambigious orders = initStanding still
		if gubbe.controls.left == gubbe.controls.right {
			initStanding(gubbe)
		} else if gubbe.controls.right && gubbe.looking == Left {
			initWalking(gubbe, Right)
		} else if gubbe.controls.left && gubbe.looking == Right {
			initWalking(gubbe, Left)
		}
	case Kicking:
	}

	// STATE INDEPENDENT BEHAVIOR
	gubbe.vel = gubbe.vel.Add(gubbe.acc)
	if gubbe.vel.Len() > MAXVELOCITY {
		gubbe.vel = gubbe.vel.Unit().Scaled(MAXVELOCITY)
	}
	gubbe.pos = gubbe.pos.Add(gubbe.vel)

	return true
}

func tickGubbe(gubbe *Gubbe, controls Controls) {
	if controls.left {
		gubbe.HandleKeyDown(internal.Left)
	} else {
		gubbe.HandleKeyUp(internal.Left)
	}
	if controls.right {
		gubbe.HandleKeyDown(internal.Right)
	} else {
		gubbe.HandleKeyUp(internal.Right)
	}
	gubbe.Tick()
}

func initStanding(g *Gubbe) {
	g.state = Standing
	g.image = StandingRight
	g.acc = pixel.ZV
	g.counter = 0
}

func initWalking(g *Gubbe, looking Looking) {
	g.state = Walking
	g.looking = looking
	g.image = WalkRight1
	g.acc = pixel.Vec{X: ACCELERATION, Y: 0}
	if looking == Left {
		g.acc.X = -g.acc.X
	}
	g.counter = 0
}

func MakeGubbe(pos pixel.Vec) Gubbe {
	return Gubbe{
		state:    Standing,
		looking:  Right,
		image:    StandingRight,
		counter:  0,
		pos:      pos,
		vel:      pixel.ZV,
		acc:      pixel.ZV,
		controls: Controls{false, false, false},
	}
}

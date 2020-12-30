package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const MAXVELOCITY = 2
const ACCELERATION = 0.03
const DECCELERATION = 0.95

type Gubbe struct {
	state   State
	looking Looking
	image   Image
	counter int
	pos     pixel.Vec
	vel     pixel.Vec
	acc     pixel.Vec
}

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

func stepGubbe(g *Gubbe, controls Controls) {

	// STATE DEPENDENT BEHAVIOR

	switch g.state {
	case Standing:
		if controls.right && !controls.left {
			initWalking(g, Right)
		}
		if controls.left && !controls.right {
			initWalking(g, Left)
		}
		g.vel = g.vel.Scaled(DECCELERATION)
	case Walking:
		g.counter++
		if g.counter%30 == 0 {
			if g.image == WalkRight1 {
				g.image = WalkRight2
			} else {
				g.image = WalkRight1
			}
		}
		// No directions, or ambigious orders = initStanding still
		if controls.left == controls.right {
			initStanding(g)
		} else if controls.right && g.looking == Left {
			initWalking(g, Right)
		} else if controls.left && g.looking == Right {
			initWalking(g, Left)
		}
	case Kicking:
	}

	// STATE INDEPENDENT BEHAVIOR
	g.vel = g.vel.Add(g.acc)
	if g.vel.Len() > MAXVELOCITY {
		g.vel = g.vel.Unit().Scaled(MAXVELOCITY)
	}
	g.pos = g.pos.Add(g.vel)
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

func MakeGubbe(win *pixelgl.Window) Gubbe {
	return Gubbe{
		state:   Standing,
		looking: Right,
		pos:     win.Bounds().Center().Add(pixel.Vec{0, -screenheight / 4}),
		vel:     pixel.ZV,
		acc:     pixel.ZV,
	}
}

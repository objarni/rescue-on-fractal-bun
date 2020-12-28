package main

import (
	"github.com/faiface/pixel"
)

const MAXVELOCITY = 2

type Gubbe struct {
	state        State
	looking      Looking
	image        Image
	stepTillSwap int
	pos          pixel.Vec
	vel          pixel.Vec
	acc          pixel.Vec
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

	// CONTROL BEHAVIOR

	switch g.state {
	case Standing:
		if controls.right {
			g.state = Walking
			g.looking = Right
			g.image = WalkRight1
			g.acc = pixel.Vec{X: 1, Y: 0}
			g.stepTillSwap = 10
		}
		if controls.left {
			g.state = Walking
			g.looking = Left
			g.image = WalkRight1
			g.acc = pixel.Vec{X: -1, Y: 0}
			g.stepTillSwap = 10
		}
		g.vel = g.vel.Scaled(0.9)
	case Walking:
		// Switch image every step
		g.stepTillSwap--
		if g.stepTillSwap == 0 {
			g.stepTillSwap = 10
			if g.image == WalkRight1 {
				g.image = WalkRight2
			} else {
				g.image = WalkRight1
			}
		}
		if controls.right {
			g.state = Walking
			g.looking = Right
			g.image = WalkRight1
			g.acc = pixel.Vec{X: 1, Y: 0}
		} else if controls.left {
			g.state = Walking
			g.looking = Left
			g.image = WalkRight1
			g.acc = pixel.Vec{X: -1, Y: 0}
		} else {
			g.state = Standing
			g.image = StandingRight
			g.acc = pixel.ZV
		}
	case Kicking:
	}

	// Update position & velocity
	g.vel = g.vel.Add(g.acc)
	// Cap velocity to 5 pixels per step
	if g.vel.Len() > MAXVELOCITY {
		g.vel = g.vel.Unit().Scaled(MAXVELOCITY)
	}
	g.pos = g.pos.Add(g.vel)
}
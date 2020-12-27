package main

import (
	"github.com/faiface/pixel"
)

const speed = 50 // Pixels per second

type Gubbe struct {
	state       State
	looking     Looking
	timeInState float64
	position    pixel.Vec
}

type Looking int

const (
	Right Looking = iota
	Left
)

func (looking Looking) String() string {
	return [...]string{"Right", "Left"}[looking]
}

type State int

const (
	Standing State = iota
	Kicking
	Walking
)

func (state State) String() string {
	return [...]string{"Standing", "Kicking", "Walking"}[state]
}

type Controls struct {
	left, right, kick bool
}

func updateGubbe(g *Gubbe, deltaSeconds float64, controls Controls) {
	if g.state == Standing {
		g.state = Walking
		if controls.left {
			g.looking = Left
		}
		if controls.left && controls.right {
			g.state = Standing
			g.looking = Right
		}
		if controls.kick {
			g.state = Kicking
		}
	}
	if g.state == Walking {
		var dir float64 = 1
		if g.looking == Left {
			dir = -1
		}
		g.position = g.position.Add(pixel.Vec{deltaSeconds * dir * speed, 0})
	}
}

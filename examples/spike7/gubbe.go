package main

import (
	"github.com/faiface/pixel"
	"time"
)

func updateGubbe(g *Gubbe, delta float64, controls Controls) {
	g.state = Walking
	if controls.left {
		g.looking = Left
	}
	if controls.left && controls.right {
		g.state = Standing
		g.looking = Right
	}
}

type Gubbe struct {
	state           State
	looking         Looking
	stateChangeTime time.Time // For animation purposes
	position        pixel.Vec
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

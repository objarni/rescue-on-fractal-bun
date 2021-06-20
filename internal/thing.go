package internal

import (
	"github.com/faiface/pixel/pixelgl"
)

// TickMs How many milliseconds is a tick?
const TickMs = 5

type Controllable interface {
	HandleKeyDown(key ControlKey) Thing
	HandleKeyUp(key ControlKey) Thing
}

type Animated interface {
	Render(win *pixelgl.Window)
	Tick() bool
}

type Thing interface {
	Controllable
	Animated
}

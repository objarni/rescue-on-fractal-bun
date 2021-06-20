package internal

import (
	"github.com/faiface/pixel/pixelgl"
)

// TickMs How many milliseconds is a tick?
const TickMs = 5

type Thing interface {
	HandleKeyDown(key ControlKey) Thing
	HandleKeyUp(key ControlKey) Thing
	Render(win *pixelgl.Window)
	Tick() bool
	WantToExitProgram() bool
}

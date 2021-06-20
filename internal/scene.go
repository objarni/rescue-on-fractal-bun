package internal

import (
	"github.com/faiface/pixel/pixelgl"
)

// TickMs How many milliseconds is a tick?
const TickMs = 5

type Scene interface {
	HandleKeyDown(key ControlKey) Scene
	HandleKeyUp(key ControlKey) Scene
	Render(win *pixelgl.Window)
	Tick() bool
	WantToExitProgram() bool
}

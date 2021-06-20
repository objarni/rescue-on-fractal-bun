package scenes

import (
	"github.com/faiface/pixel/pixelgl"
	"objarni/rescue-on-fractal-bun/internal"
)

// TickMs How many milliseconds is a tick?
const TickMs = 5

type Scene interface {
	WantToExitProgram() bool
	Tick() Scene
	Render(win *pixelgl.Window)
	HandleKeyDown(key internal.ControlKey) Scene
	HandleKeyUp(key internal.ControlKey) Scene
}

package internal

import (
	"github.com/faiface/pixel/pixelgl"
)

type Scene interface {
	HandleKeyDown(key ControlKey) Scene

	HandleKeyUp(key ControlKey) Scene

	// @thought: generalize to 'rescue draw operations'?
	// itches:
	//   - text atlas recreation in every scene
	//   - copy-paste textbox code e.g. scale = 2
	Render(win *pixelgl.Window)
}

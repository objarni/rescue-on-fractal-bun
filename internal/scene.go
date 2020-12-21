package internal

import (
	"github.com/faiface/pixel/pixelgl"
)

type Scene interface {
	HandleKeyDown(key ControlKey) Scene
	HandleKeyUp(key ControlKey) Scene
	Render(win *pixelgl.Window)
}

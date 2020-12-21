package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"objarni/rescue-on-fractal-bun/internal"
)

func run() {

	scene := internal.MakeMenuScene()

	cfg := pixelgl.WindowConfig{
		Title:  "Rescue on fractal bun (work title)",
		Bounds: pixel.R(0, 0, 800, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	controllerMap := make(map[pixelgl.Button]internal.ControlKey)
	controllerMap[pixelgl.KeyUp] = internal.Up
	controllerMap[pixelgl.KeyDown] = internal.Down

	for !win.Closed() {

		// Escape closes main window unconditionally
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		for key, control := range controllerMap {
			// Hmm. Just Pressed/Released APIs is 'key repeat' - problem?
			if win.JustPressed(key) {
				scene.HandleKeyDown(control)
			}
			if win.JustReleased(key) {
				scene.HandleKeyUp(control)
			}
		}
		scene.Render(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

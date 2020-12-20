package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
)

func run() {

	scene := internal.MenuScene{}

	cfg := pixelgl.WindowConfig{
		Title:  "Rescue on fractal bun (work title)",
		Bounds: pixel.R(0, 0, 800, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		// TODO: DRY - make a 'dictionary' from pixelgl key to ControllerKey
		if win.JustPressed(pixelgl.KeyDown) {
			scene.HandleKeyDown(internal.Down)
		}
		if win.JustReleased(pixelgl.KeyDown) {
			scene.HandleKeyUp(internal.Down)
		}
		if win.JustPressed(pixelgl.KeyUp) {
			scene.HandleKeyDown(internal.Up)
		}
		if win.JustReleased(pixelgl.KeyUp) {
			scene.HandleKeyUp(internal.Up)
		}
		//scene.Render()
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"objarni/rescue-on-fractal-bun/internal"
	"time"
)

func run() {

	var scene internal.Scene = internal.MakeMenuScene()

	cfg := pixelgl.WindowConfig{
		Title:  "Rescue on fractal bun (work title)",
		Bounds: pixel.R(0, 0, 800, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	err = speaker.Init(beep.SampleRate(22050), 2000)

	controllerMap := make(map[pixelgl.Button]internal.ControlKey)
	controllerMap[pixelgl.KeyUp] = internal.Up
	controllerMap[pixelgl.KeyDown] = internal.Down
	controllerMap[pixelgl.KeyLeft] = internal.Left
	controllerMap[pixelgl.KeyRight] = internal.Right
	controllerMap[pixelgl.KeySpace] = internal.Jump
	controllerMap[pixelgl.KeyRightControl] = internal.Action

	for !win.Closed() {

		// Escape closes main window unconditionally
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		for key, control := range controllerMap {
			// Hmm. Just Pressed/Released APIs is 'key repeat' at least on win - problem?
			if win.JustPressed(key) {
				fmt.Println("pressed: ", key)
				scene = scene.HandleKeyDown(control)
			}
			if win.JustReleased(key) {
				fmt.Println("released: ", key)
				scene = scene.HandleKeyUp(control)
			}
			if scene == nil {
				win.SetClosed(true)
				continue
			}
		}
		scene.Render(win)
		win.Update()
		time.Sleep(time.Millisecond * 5)
	}
}

func main() {
	pixelgl.Run(run)
}
package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
	"time"
)

func run() {

	var scene internal.Thing = scenes.MakeMenuScene()

	cfg := pixelgl.WindowConfig{
		Title:  "Rescue on fractal bun (work title)",
		Bounds: pixel.R(0, 0, 800, 600),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	err = speaker.Init(beep.SampleRate(22050), 1000)

	keyMap := make(map[pixelgl.Button]internal.ControlKey)
	keyMap[pixelgl.KeyUp] = internal.Up
	keyMap[pixelgl.KeyDown] = internal.Down
	keyMap[pixelgl.KeyLeft] = internal.Left
	keyMap[pixelgl.KeyRight] = internal.Right
	keyMap[pixelgl.KeySpace] = internal.Jump
	keyMap[pixelgl.KeyRightControl] = internal.Action

	padMap := make(map[pixelgl.GamepadButton]internal.ControlKey)
	padMap[pixelgl.ButtonDpadUp] = internal.Up
	padMap[pixelgl.ButtonDpadDown] = internal.Down
	padMap[pixelgl.ButtonDpadLeft] = internal.Left
	padMap[pixelgl.ButtonDpadRight] = internal.Right
	padMap[pixelgl.ButtonA] = internal.Jump
	padMap[pixelgl.ButtonB] = internal.Action

	for !win.Closed() {

		// Escape closes main window unconditionally
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		// Keyboard control
		for key, control := range keyMap {
			// Hmm. Just Pressed/Released APIs is 'key repeat' at least on win - problem?
			if win.JustPressed(key) {
				fmt.Println("pressed: ", key)
				scene = scene.HandleKeyDown(control)
			}
			if win.JustReleased(key) {
				fmt.Println("released: ", key)
				scene = scene.HandleKeyUp(control)
			}
		}

		// Gamepad control
		for pad, control := range padMap {
			// @remember: do we want to check all joysticks not just 1?
			if win.JoystickJustPressed(pixelgl.Joystick1, pad) {
				fmt.Println("pressed: ", pad)
				scene = scene.HandleKeyDown(control)
			}
			if win.JoystickJustReleased(pixelgl.Joystick1, pad) {
				fmt.Println("released: ", pad)
				scene = scene.HandleKeyUp(control)
			}
		}
		if !scene.Tick() {
			win.SetClosed(true)
			continue
		}

		scene.Render(win)
		win.Update()
		time.Sleep(time.Millisecond * 5)
	}
}

func main() {
	pixelgl.Run(run)
}

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

		if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonDpadLeft) {
			scene = scene.HandleKeyDown(internal.Left)
		}
		if win.JoystickJustReleased(pixelgl.Joystick1, pixelgl.ButtonDpadLeft) {
			scene = scene.HandleKeyUp(internal.Left)
		}
		if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonDpadRight) {
			scene = scene.HandleKeyDown(internal.Right)
		}
		if win.JoystickJustReleased(pixelgl.Joystick1, pixelgl.ButtonDpadRight) {
			scene = scene.HandleKeyUp(internal.Right)
		}
		if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonDpadUp) {
			scene = scene.HandleKeyDown(internal.Up)
		}
		if win.JoystickJustReleased(pixelgl.Joystick1, pixelgl.ButtonDpadUp) {
			scene = scene.HandleKeyUp(internal.Up)
		}
		if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonDpadDown) {
			scene = scene.HandleKeyDown(internal.Down)
		}
		if win.JoystickJustReleased(pixelgl.Joystick1, pixelgl.ButtonDpadDown) {
			fmt.Println(pixelgl.ButtonDpadDown)
			scene = scene.HandleKeyUp(internal.Down)
		}
		if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonA) {
			scene = scene.HandleKeyDown(internal.Jump)
		}
		if win.JoystickJustReleased(pixelgl.Joystick1, pixelgl.ButtonA) {
			fmt.Println(pixelgl.ButtonDpadDown)
			scene = scene.HandleKeyUp(internal.Jump)
		}
		if win.JoystickJustPressed(pixelgl.Joystick1, pixelgl.ButtonB) {
			scene = scene.HandleKeyDown(internal.Action)
		}
		if win.JoystickJustReleased(pixelgl.Joystick1, pixelgl.ButtonB) {
			fmt.Println(pixelgl.ButtonDpadDown)
			scene = scene.HandleKeyUp(internal.Action)
		}

		//if win.JoystickPresent(pixelgl.Joystick1) {
		//	fmt.Println("js1:", win.JoystickName(pixelgl.Joystick1))
		//}

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
			if !scene.Tick() {
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

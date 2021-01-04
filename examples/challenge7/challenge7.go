package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	"time"
)

const screenwidth = 800
const screenheight = 600

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:    "Kick the ball",
		Bounds:   pixel.R(0, 0, screenwidth, screenheight),
		Position: pixel.Vec{X: screenwidth / 2, Y: screenheight / 2},
	})
	internal.PanicIfError(err)
	err = speaker.Init(
		beep.SampleRate(44100),
		2000,
	)
	internal.PanicIfError(err)

	var prevtime = time.Now()

	var rest float64 = 0
	keyTranslation := map[pixelgl.Button]internal.ControlKey{
		pixelgl.KeyLeft:         internal.Left,
		pixelgl.KeyRight:        internal.Right,
		pixelgl.KeySpace:        internal.Jump,
		pixelgl.KeyRightControl: internal.Action,
	}
	cfg := TryReadCfgFrom("json/challenge7.json", Config{})

	var scene internal.Thing = MakeStartScene(&cfg)

	for !win.Closed() {
		// Janitor
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		if win.JustPressed(pixelgl.KeyNumLock) {
			cfg = TryReadCfgFrom("json/challenge7.json", cfg)
		}

		// Compute time deltaMs
		var now = time.Now()
		var deltaMs = now.Sub(prevtime).Seconds() * 1000
		prevtime = now
		deltaMs += rest

		// Keyboard input
		for key, value := range keyTranslation {
			if win.JustPressed(key) {
				scene = scene.HandleKeyDown(value)
			}
			if win.JustReleased(key) {
				scene = scene.HandleKeyUp(value)
			}
		}

		// Update entities
		steps := int(math.Floor(deltaMs / 5))
		rest = deltaMs - float64(steps*5)
		for i := 0; i < steps; i++ {
			scene.Tick()
		}

		// Render
		scene.Render(win)

		// Window/OS
		win.Update()
		time.Sleep(time.Millisecond * 5)
	}
}

func main() {
	pixelgl.Run(run)
}

package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
	"os"
	"time"
	"unicode"
)

func run() {
	configFile := "json/rescue.json"
	cfg := scenes.TryReadCfgFrom(configFile, scenes.Config{})
	info, err := os.Stat(configFile)
	internal.PanicIfError(err)
	cfgTime := info.ModTime()

	// Load resources
	res := loadResources()

	// Initial scene - depends on --level cmd line arg!
	var scene internal.Thing
	if len(os.Args) == 1 {
		scene = scenes.MakeMenuScene(&cfg, &res)
	} else {
		argsWithoutProg := os.Args[1:]
		if argsWithoutProg[0] == "--level" {
			// TODO: parameterize on level arg!
			levelName := argsWithoutProg[1]
			fmt.Println("Loading level:", levelName)
			scene = scenes.MakeLevelScene(&cfg, &res)
		} else {
			fmt.Printf("Unknown cmd.line arg: %v\n", argsWithoutProg)
			panic("Don't understand cmd.line")
		}
	}

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Rescue on fractal bun (work title)",
		Bounds: pixel.R(0, 0, 800, 600),
	})
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

	fpsCounter := 0

	for !win.Closed() {

		// Escape closes main window unconditionally
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		// Tweak system
		info, err := os.Stat(configFile)
		internal.PanicIfError(err)
		if cfgTime != info.ModTime() {
			cfgTime = info.ModTime()
			cfg = scenes.TryReadCfgFrom(configFile, cfg)
		}

		// Keyboard control
		for key, control := range keyMap {
			// Hmm. Just Pressed/Released APIs is 'key repeat' at least on win - problem?
			if win.JustPressed(key) {
				//fmt.Println("pressed: ", key)
				scene = scene.HandleKeyDown(control)
			}
			if win.JustReleased(key) {
				//fmt.Println("released: ", key)
				scene = scene.HandleKeyUp(control)
			}
		}

		// Gamepad control
		for pad, control := range padMap {
			// @remember: do we want to check all joysticks not just 1?
			if win.JoystickJustPressed(pixelgl.Joystick1, pad) {
				//fmt.Println("pressed: ", pad)
				scene = scene.HandleKeyDown(control)
			}
			if win.JoystickJustReleased(pixelgl.Joystick1, pad) {
				//fmt.Println("released: ", pad)
				scene = scene.HandleKeyUp(control)
			}
		}
		if !scene.Tick() {
			win.SetClosed(true)
			continue
		}

		start := time.Now()
		scene.Render(win)
		win.Update()
		duration := time.Since(start)
		// Only update FPS every 10th tick
		fpsCounter++
		if fpsCounter == 40 {
			res.FPS = computeFPS(duration)
			fpsCounter = 0
		}

		time.Sleep(time.Millisecond * 5)
	}
}

func computeFPS(renderTime time.Duration) float64 {
	return 1.0 / renderTime.Seconds()
}

func loadResources() scenes.Resources {
	res := scenes.Resources{}
	face := internal.LoadTTFForSure("assets/Font.ttf", 32)
	res.Atlas = text.NewAtlas(face, text.RangeTable(unicode.Latin), text.ASCII)
	res.Ghost = internal.LoadSpriteForSure("assets/TGhost.png")
	res.MapPoint = internal.LoadSpriteForSure("assets/TMapPoint.png")
	res.PlayerStanding = internal.LoadSpriteForSure("assets/TStanding.png")
	res.Blip = internal.LoadWavForSure("assets/Bounce.wav")
	return res
}

func main() {
	pixelgl.Run(run)
}

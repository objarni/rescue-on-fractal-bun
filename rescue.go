package main

import (
	"flag"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
	"os"
	"runtime/pprof"
	"time"
	"unicode"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var startlevel = flag.String("level", "", "start at specified game level")

func run() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	cfg := scenes.TryReadCfgFrom(internal.ConfigFile, scenes.Config{})
	info, err := os.Stat(internal.ConfigFile)
	internal.PanicIfError(err)
	cfgTime := info.ModTime()

	// Load resources
	res := loadResources()

	// Initial scene - depends on --level cmd line arg!
	var scene internal.Thing
	if *startlevel != "" {
		levelName := startlevel
		fmt.Println("Loading level:", *levelName)
		scene = scenes.MakeLevelScene(&cfg, &res, *levelName)
	} else {
		scene = scenes.MakeMenuScene(&cfg, &res)
	}

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Rescue",
		Bounds: pixel.R(0, 0, internal.ScreenWidth, internal.ScreenHeight),
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

	accumulatedMs := 0.0
	timePrev := time.Now()

	for !win.Closed() {

		// Escape closes main window unconditionally
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		// Tweak system
		info, err := os.Stat(internal.ConfigFile)
		internal.PanicIfError(err)
		if cfgTime != info.ModTime() {
			cfgTime = info.ModTime()
			fmt.Println("Reading ", internal.ConfigFile, " at ", time.Now().Format(time.Stamp))
			cfg = scenes.TryReadCfgFrom(internal.ConfigFile, cfg)
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

		timeNow := time.Now()
		deltaMs := 1000.0 * timeNow.Sub(timePrev).Seconds()
		timePrev = timeNow
		var steps int
		accumulatedMs, steps = gameTimeSteps(accumulatedMs, deltaMs)
		for i := 0; i < steps; i++ {
			if !scene.Tick() {
				win.SetClosed(true)
			}
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

func loadResources() internal.Resources {
	face := internal.LoadTTFForSure("assets/Font.ttf", 32)
	levels := map[string]internal.Level{
		"GhostForest":   internal.LoadLevel("assets/levels/GhostForest.tmx"),
		"ForestOpening": internal.LoadLevel("assets/levels/ForestOpening.tmx"),
	}
	mapSigns := internal.BuildMapSignArray(levels)

	res := internal.Resources{
		Atlas:    text.NewAtlas(face, text.RangeTable(unicode.Latin), text.ASCII),
		Blip:     internal.LoadWavForSure("assets/Bounce.wav"),
		FPS:      0,
		MapSigns: mapSigns,
		Levels:   levels,
	}
	res.ImageMap = map[internal.Image]*pixel.Sprite{
		internal.IMap:                  internal.LoadSpriteForSure("assets/TMap.png"),
		internal.IGhost:                internal.LoadSpriteForSure("assets/TGhost.png"),
		internal.IMapSymbol:            internal.LoadSpriteForSure("assets/THeadsup.png"),
		internal.ISignPost:             internal.LoadSpriteForSure("assets/TMapPoint.png"),
		internal.ITemporaryPlayerImage: internal.LoadSpriteForSure("assets/TEliseWalk1.png"),
		internal.IEliseWalk1:           internal.LoadSpriteForSure("assets/TEliseWalk1.png"),
		internal.IEliseWalk2:           internal.LoadSpriteForSure("assets/TEliseWalk2.png"),
		internal.IEliseWalk3:           internal.LoadSpriteForSure("assets/TEliseWalk3.png"),
		internal.IEliseWalk4:           internal.LoadSpriteForSure("assets/TEliseWalk4.png"),
		internal.IEliseWalk5:           internal.LoadSpriteForSure("assets/TEliseWalk5.png"),
		internal.IEliseWalk6:           internal.LoadSpriteForSure("assets/TEliseWalk6.png"),
		internal.IEliseCrouch:          internal.LoadSpriteForSure("assets/TEliseCrouch.png"),
		internal.IButton:               internal.LoadSpriteForSure("assets/TButton.png"),
		internal.IStreetLight:          internal.LoadSpriteForSure("assets/TStreetLight.png"),
	}
	if len(res.ImageMap) < int(internal.AfterLastImage) {
		panic("Expect one image loaded per map item")
	}
	return res
}

func main() {
	pixelgl.Run(run)
}

func gameTimeSteps(accumulated float64, deltaMs float64) (float64, int) {
	// How many whole ticks can we step?
	accumulated += deltaMs
	steps := int(accumulated / 5) // 200 logical frames per second
	accumulated -= float64(steps * 5)
	return accumulated, steps
}

package main

import (
	"flag"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
	"objarni/rescue-on-fractal-bun/internal/tweaking"
	"os"
	"runtime/pprof"
	"time"
	"unicode"
)

var cpuProfile = flag.String("cpuProfile", "", "write cpu profile to file")
var startLevel = flag.String("level", "", "start at specified game level")

func run() {

	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	cfg := tweaking.TryReadCfgFrom(internal.ConfigFile, tweaking.Config{})
	info, err := os.Stat(internal.ConfigFile)
	internal.PanicIfError(err)
	cfgTime := info.ModTime()

	// Load resources
	res := loadResources()

	// Initial scene - depends on --level cmd line arg!
	var scene internal.Scene
	startLevelName := *startLevel
	if startLevelName != "" {
		fmt.Println("Loading level:", startLevelName)
		var pos px.Vec
		for _, mapSign := range res.MapSigns {
			if mapSign.LevelName == startLevelName {
				pos = mapSign.LevelPos
				break
			}
		}
		scene = scenes.MakeLevelScene(&cfg, &res, startLevelName, pos)
	} else {
		scene = scenes.MakeMenuScene(&cfg, &res)
	}

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Rescue",
		Bounds: px.R(0, 0, internal.ScreenWidth, internal.ScreenHeight),
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

	keyMap[pixelgl.KeyW] = internal.Up
	keyMap[pixelgl.KeyS] = internal.Down
	keyMap[pixelgl.KeyA] = internal.Left
	keyMap[pixelgl.KeyD] = internal.Right
	keyMap[pixelgl.KeyEnter] = internal.Action

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
		cfg = maybeReloadConfig(cfgTime, cfg)

		// Keyboard control
		for key, control := range keyMap {
			// TODO: Hmm. Just Pressed/Released APIs is 'key repeat' at least on win - problem?
			if win.JustPressed(key) {
				scene = scene.HandleKeyDown(control)
			}
			if win.JustReleased(key) {
				scene = scene.HandleKeyUp(control)
			}
		}

		// Gamepad control
		for pad, control := range padMap {
			// TODO: do we want to check all joysticks not just 1?
			if win.JoystickJustPressed(pixelgl.Joystick1, pad) {
				scene = scene.HandleKeyDown(control)
			}
			if win.JoystickJustReleased(pixelgl.Joystick1, pad) {
				scene = scene.HandleKeyUp(control)
			}
		}
		timeNow := time.Now()
		deltaMs := 1000.0 * timeNow.Sub(timePrev).Seconds()
		timePrev = timeNow
		var steps int
		accumulatedMs, steps = gameTimeSteps(accumulatedMs, deltaMs, cfg.LengthOfTickInMS)
		for i := 0; i < steps; i++ {
			if !scene.Tick() {
				win.SetClosed(true)
			}
		}

		start := time.Now()
		scene.Render(win)

		if isAnalogJoystickDisplaced(win) {
			sprite := res.ImageMap[internal.IUseDPAD]
			sprite.Draw(win, px.IM.Moved(win.Bounds().Center()))
		}

		win.Update()
		duration := time.Since(start)

		// Only update DisplayFrameMs every 10th tick
		fpsCounter++
		if fpsCounter == 40 {
			res.FPS = computeFPS(duration)
			fpsCounter = 0
		}

		time.Sleep(time.Millisecond * 5)
	}
}

func maybeReloadConfig(cfgTime time.Time, cfg tweaking.Config) tweaking.Config {
	info, err := os.Stat(internal.ConfigFile)
	if err != nil {
		fmt.Printf("could not stat config file skipping config %v\n", internal.ConfigFile)
	} else {
		if cfgTime != info.ModTime() {
			cfgTime = info.ModTime()
			fmt.Println("Reading ", internal.ConfigFile, " at ", time.Now().Format(time.Stamp))
			cfg = tweaking.TryReadCfgFrom(internal.ConfigFile, cfg)
		}
	}
	return cfg
}

func isAnalogJoystickDisplaced(win *pixelgl.Window) bool {
	analogXAxis := win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftX)
	analogYAxis := win.JoystickAxis(pixelgl.Joystick1, pixelgl.AxisLeftY)
	displaced := false
	if math.Abs(analogYAxis) > 0.5 || math.Abs(analogXAxis) > 0.5 {
		displaced = true
	}
	return displaced
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
		Atlas:       text.NewAtlas(face, text.RangeTable(unicode.Latin), text.ASCII),
		Blip:        internal.LoadWavForSure("assets/Bounce.wav"),
		ButtonClick: internal.LoadWavForSure("assets/SButtonClick.wav"),
		RobotMove:   internal.LoadWavForSure("assets/SRobotMove.wav"),
		FPS:         0,
		MapSigns:    mapSigns,
		Levels:      levels,
	}
	res.ImageMap = map[internal.Image]*px.Sprite{
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
		internal.IEliseJump1:           internal.LoadSpriteForSure("assets/TEliseJumping1.png"),
		internal.IEliseJump2:           internal.LoadSpriteForSure("assets/TEliseJumping2.png"),
		internal.IEliseJump3:           internal.LoadSpriteForSure("assets/TEliseJumping3.png"),
		internal.IEliseJump4:           internal.LoadSpriteForSure("assets/TEliseJumping4.png"),
		internal.IEliseJump5:           internal.LoadSpriteForSure("assets/TEliseJumping5.png"),
		internal.IEliseJump6:           internal.LoadSpriteForSure("assets/TEliseJumping6.png"),
		internal.IEliseJump7:           internal.LoadSpriteForSure("assets/TEliseJumping7.png"),
		internal.IButton:               internal.LoadSpriteForSure("assets/TButton.png"),
		internal.IStreetLight:          internal.LoadSpriteForSure("assets/TStreetLight.png"),
		internal.ISpider:               internal.LoadSpriteForSure("assets/TSpider.png"),
		internal.IUseDPAD:              internal.LoadSpriteForSure("assets/TUseDPAD.png"),
		internal.IRobot1:               internal.LoadSpriteForSure("assets/TRobot1.png"),
	}
	if len(res.ImageMap) < int(internal.AfterLastImage) {
		panic("Expect one image loaded per map item")
	}
	return res
}

func main() {
	pixelgl.Run(run)
}

func gameTimeSteps(accumulated float64, deltaMs float64, lengthOfTickInMS float64) (float64, int) {
	// How many whole ticks can we step?
	accumulated += deltaMs
	steps := int(accumulated / lengthOfTickInMS) // 200 logical frames per second
	accumulated -= float64(steps) * lengthOfTickInMS
	return accumulated, steps
}

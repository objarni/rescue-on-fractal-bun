package main

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	"time"
)

const screenwidth = 800
const screenheight = 600

func run() {
	cfg := pixelgl.WindowConfig{
		Title:    "Kick the ball",
		Bounds:   pixel.R(0, 0, screenwidth, screenheight),
		Position: pixel.Vec{X: screenwidth / 2, Y: screenheight / 2},
	}
	win, err := pixelgl.NewWindow(cfg)
	internal.PanicIfError(err)
	err = speaker.Init(
		beep.SampleRate(44100),
		2000,
	) //done := make(chan bool)
	internal.PanicIfError(err)

	var scene internal.Thing = MakeStartScene()
	var prevtime = time.Now()

	var rest float64 = 0
	keyTranslation := map[pixelgl.Button]internal.ControlKey{
		pixelgl.KeyLeft:         internal.Left,
		pixelgl.KeyRight:        internal.Right,
		pixelgl.KeySpace:        internal.Jump,
		pixelgl.KeyRightControl: internal.Action,
	}
	for !win.Closed() {
		// Janitor
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		//config, err = TryReadCfgFrom("json/challenge7.json", config)
		internal.PanicIfError(err)

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

type Config struct {
	Gravity float64
	SpeedX  float64
	StartX  float64
	StartY  float64
}

func TryReadCfgFrom(filename string, defaultCfg Config) (Config, error) {
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile error, defaulting")
		return defaultCfg, err
	}
	var cfg = defaultCfg
	err = json.Unmarshal(byteArray, &cfg)
	if err != nil {
		fmt.Println("JSON parse error, defaulting. JSON was: " + string(byteArray))
		return defaultCfg, err
	}
	return cfg, nil
}

func main() {
	pixelgl.Run(run)
}

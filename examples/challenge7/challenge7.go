package main

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	"time"
)

const screenwidth = 800
const screenheight = 600

type StartScene struct {
}

type GameScene struct {
	ball  internal.Thing
	gubbe internal.Thing
}

func (startScene *StartScene) HandleKeyUp(_ internal.ControlKey) internal.Thing {
	return startScene
}

func (startScene *StartScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		var config Config
		config, err := TryReadCfgFrom("json/challenge7.json", config)
		internal.PanicIfError(err)

		var gubbeStandingRightSprite = internal.LoadSpriteForSure("assets/TStanding.png")
		var gubbeWalkingRightSprite1 = internal.LoadSpriteForSure("assets/TWalking-1.png")
		var gubbeWalkingRightSprite2 = internal.LoadSpriteForSure("assets/TWalking-2.png")
		gubbeImage2Sprite := map[Image]*pixel.Sprite{
			WalkRight1:    gubbeWalkingRightSprite1,
			WalkRight2:    gubbeWalkingRightSprite2,
			StandingRight: gubbeStandingRightSprite,
		}
		var gubbe = MakeGubbe(pixel.Vec{X: 100, Y: 150}, gubbeImage2Sprite)
		var scene internal.Thing = &GameScene{
			ball:  MakeBall(config),
			gubbe: &gubbe,
		}

		return scene
	}
	return startScene
}

func (startScene *StartScene) Tick() bool {
	return true
}

func (startScene *StartScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Green)
}

func (gameScene *GameScene) HandleKeyUp(key internal.ControlKey) internal.Thing {
	gameScene.gubbe.HandleKeyUp(key)
	return gameScene
}

func (gameScene *GameScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	gameScene.gubbe.HandleKeyDown(key)
	return gameScene
}

func (gameScene *GameScene) Tick() bool {
	gameScene.gubbe.Tick()
	gameScene.ball.Tick()
	return true
}

func (gameScene *GameScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Lightskyblue)
	drawGround(win)
	gameScene.gubbe.Render(win)
	gameScene.ball.Render(win)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:    "Kick the ball",
		Bounds:   pixel.R(0, 0, screenwidth, screenheight),
		Position: pixel.Vec{X: screenwidth / 2, Y: screenheight / 2},
	}
	win, err := pixelgl.NewWindow(cfg)
	internal.PanicIfError(err)
	err = speaker.Init(
		beep.SampleRate(22050),
		2000,
	) //done := make(chan bool)
	internal.PanicIfError(err)

	var scene internal.Thing = &StartScene{}
	var prevtime = time.Now()

	var rest float64 = 0
	keyTranslation := map[pixelgl.Button]internal.ControlKey{
		pixelgl.KeyLeft:  internal.Left,
		pixelgl.KeyRight: internal.Right,
		pixelgl.KeySpace: internal.Jump,
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

func drawGround(win *pixelgl.Window) {
	var imd = imdraw.New(nil)
	win.SetSmooth(true)
	imd.Clear()
	imd.Color = colornames.Darkgreen
	imd.Push(pixel.ZV)
	imd.Color = colornames.Lightgreen
	imd.Push(pixel.Vec{X: screenwidth, Y: 75})
	imd.Rectangle(0)
	imd.Draw(win)
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

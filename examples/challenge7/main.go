package main

import (
	"encoding/json"
	"fmt"
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

func run() {
	cfg := pixelgl.WindowConfig{
		Title:    "Kick the ball",
		Bounds:   pixel.R(0, 0, screenwidth, screenheight),
		Position: pixel.Vec{X: screenwidth / 2, Y: screenheight / 2},
	}
	win, err := pixelgl.NewWindow(cfg)
	internal.PanicIfError(err)

	var imd = imdraw.New(nil)
	win.SetSmooth(true)
	var config Config
	var prevtime = time.Now()

	var gubbeStandingRightSprite = internal.LoadSpriteForSure("assets/TStanding.png")
	var gubbeWalkingRightSprite1 = internal.LoadSpriteForSure("assets/TWalking-1.png")
	var gubbeWalkingRightSprite2 = internal.LoadSpriteForSure("assets/TWalking-2.png")
	gubbeImage2Sprite := map[Image]*pixel.Sprite{
		WalkRight1:    gubbeWalkingRightSprite1,
		WalkRight2:    gubbeWalkingRightSprite2,
		StandingRight: gubbeStandingRightSprite,
	}

	config, err = TryReadCfgFrom("json/challenge7.json", config)
	internal.PanicIfError(err)

	var ball = MakeBall(config)
	var gubbe = MakeGubbe(win)

	var rest float64 = 0
	for !win.Closed() {
		// Janitor
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		config, err = TryReadCfgFrom("json/challenge7.json", config)
		internal.PanicIfError(err)

		// Compute time deltaMs
		var now = time.Now()
		var deltaMs = now.Sub(prevtime).Seconds() * 1000
		prevtime = now
		deltaMs += rest

		// Update entities
		steps := int(math.Floor(deltaMs / 5))
		rest = deltaMs - float64(steps*5)
		for i := 0; i < steps; i++ {
			controls := Controls{
				left:  win.Pressed(pixelgl.KeyLeft),
				right: win.Pressed(pixelgl.KeyRight),
				kick:  win.Pressed(pixelgl.KeySpace),
			}
			stepGubbe(&gubbe, controls)
			ball.Tick()
		}

		// Render
		win.Clear(colornames.Lightskyblue)
		drawGround(imd, win)
		drawGubbe(gubbe, gubbeImage2Sprite, win)
		ball.Render(win)

		// Window/OS
		win.Update()
		time.Sleep(time.Millisecond * 5)
	}
}

func drawGubbe(gubbe Gubbe, gubbeImage2Sprite map[Image]*pixel.Sprite, win *pixelgl.Window) {
	mx := pixel.IM.Scaled(pixel.ZV, 1)
	mx = mx.Moved(gubbe.pos)
	if gubbe.looking == Left {
		mx = mx.ScaledXY(
			gubbe.pos,
			pixel.Vec{X: -1, Y: 1},
		)
	}
	gubbeSprite := gubbeImage2Sprite[gubbe.image]
	gubbeSprite.Draw(win, mx)
}

func drawGround(imd *imdraw.IMDraw, win *pixelgl.Window) {
	imd.Clear()
	imd.Push(pixel.ZV)
	imd.Color = colornames.Darkgreen
	imd.Push(pixel.Vec{X: screenwidth, Y: 75})
	imd.Color = colornames.Lightgreen
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

// rita blårektangel
// animera med tidsstämpel (millisekunder t.ex.)

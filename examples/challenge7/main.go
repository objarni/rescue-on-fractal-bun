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
const ballradius = 50

type Ball struct {
	Pos pixel.Vec
	Vel pixel.Vec
	Rot float64
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:    "Kick the ball",
		Bounds:   pixel.R(0, 0, screenwidth, screenheight),
		Position: pixel.Vec{X: screenwidth / 2, Y: screenheight / 2},
	}
	win, err := pixelgl.NewWindow(cfg)
	internal.PanicIfError(err)

	internal.PanicIfError(err)

	var imd = imdraw.New(nil)
	win.SetSmooth(true)
	var config Config
	var prevtime = time.Now()
	var rot float64 = 0

	var ballSprite = internal.LoadSpriteForSure("assets/Ball.png")

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

	var ballState = Ball{
		Pos: pixel.Vec{X: config.StartX, Y: config.StartY},
		Vel: pixel.Vec{X: config.SpeedX, Y: 0}, Rot: 0,
	}

	var gubbe = Gubbe{
		state:   Standing,
		looking: Right,
		pos:     win.Bounds().Center().Add(pixel.Vec{0, -screenheight / 4}),
		vel:     pixel.ZV,
		acc:     pixel.ZV,
	}
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
		updateBall(rot, deltaMs/1000, &ballState, config)
		steps := int(math.Floor(deltaMs / 5))
		rest = deltaMs - float64(steps*5)
		for i := 0; i < steps; i++ {
			controls := Controls{
				left:  win.Pressed(pixelgl.KeyLeft),
				right: win.Pressed(pixelgl.KeyRight),
				kick:  win.Pressed(pixelgl.KeySpace),
			}
			stepGubbe(&gubbe, controls)
		}

		// Render
		win.Clear(colornames.Lightskyblue)
		imd.Clear()
		imd.Push(pixel.Vec{0, 0})
		imd.Color = colornames.Darkgreen
		imd.Push(pixel.Vec{screenwidth, 75})
		imd.Color = colornames.Lightgreen
		imd.Rectangle(0)
		imd.Draw(win)
		mx := pixel.IM.Scaled(pixel.ZV, 1)
		mx = mx.Moved(gubbe.pos)
		if gubbe.looking == Left {
			mx = mx.ScaledXY(gubbe.pos, pixel.Vec{-1, 1})
		}
		gubbeSprite := gubbeImage2Sprite[gubbe.image]
		gubbeSprite.Draw(win, mx)
		drawBall(ballState, ballSprite, win)

		// Window/OS
		win.Update()
		time.Sleep(time.Millisecond * 5)
	}
}

func updateBall(rot float64, delta float64, ballState *Ball, config Config) float64 {
	rot += delta
	ballState.Pos = ballState.Pos.Add(
		ballState.Vel.Scaled(delta))
	ballState.Vel = ballState.Vel.Add(
		pixel.Vec{X: 0, Y: -config.Gravity * delta})
	ballState.Rot -= delta * ballState.Vel.X / 40
	if ballState.Pos.X < ballradius {
		ballState.Pos.X = ballradius + 1
		ballState.Vel = ballState.Vel.ScaledXY(pixel.Vec{X: -1, Y: 1})
	}
	if ballState.Pos.X > screenwidth-ballradius {
		ballState.Pos.X = screenwidth - ballradius - 1
		ballState.Vel = ballState.Vel.ScaledXY(pixel.Vec{X: -1, Y: 1})
	}
	if ballState.Pos.Y < ballradius*1.5 {
		ballState.Pos.Y = ballradius*1.5 + 1
		ballState.Vel = ballState.Vel.ScaledXY(pixel.Vec{X: 1, Y: -1})
	}
	return rot
}

func drawBall(ballState Ball, ballSprite *pixel.Sprite, win *pixelgl.Window) {
	mx := pixel.IM
	mx = mx.Moved(ballState.Pos)
	mx = mx.Rotated(ballState.Pos, ballState.Rot)
	ballSprite.Draw(win, mx)
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

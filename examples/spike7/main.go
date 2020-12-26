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
		Position: pixel.Vec{X: 300, Y: 300},
	}
	win, err := pixelgl.NewWindow(cfg)
	internal.PanicIfError(err)

	internal.PanicIfError(err)

	var imd = imdraw.New(nil)
	imd.Color = colornames.Darkslateblue
	win.SetSmooth(true)
	var config Config
	var prevtime = time.Now()
	var rot float64 = 0

	var ballSprite = internal.LoadSpriteForSure("assets/Ball.png")

	var gubbeStandingRightSprite = internal.LoadSpriteForSure("assets/Standing.jpg")
	//var gubbeWalkingRightSprite1 = internal.LoadSpriteForSure("assets/Walk-1.jpg")
	//var gubbeWalkingRightSprite2 = internal.LoadSpriteForSure("assets/Walk-2.jpg")
	//
	config, err = TryReadCfgFrom("json/spike6.json", config)
	internal.PanicIfError(err)

	var ballState = Ball{
		Pos: pixel.Vec{X: config.StartX, Y: config.StartY},
		Vel: pixel.Vec{X: config.SpeedX, Y: 0}, Rot: 0,
	}

	var gubbe = Gubbe{
		state:           Standing,
		looking:         Right,
		stateChangeTime: time.Time{},
		position:        win.Bounds().Center(),
	}

	for !win.Closed() {
		// Janitor
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		config, err = TryReadCfgFrom("json/spike7.json", config)
		internal.PanicIfError(err)

		// Compute time delta
		var now = time.Now()
		var delta = now.Sub(prevtime).Seconds()
		prevtime = now

		// Update entities
		updateBall(rot, delta, &ballState, config)
		updateGubbe(&gubbe, delta, Controls{})

		// Render
		win.Clear(colornames.Lightskyblue)
		drawBall(ballState, ballSprite, win)
		if gubbe.looking == Right {
			if gubbe.state == Standing {
				mx := pixel.IM.Scaled(pixel.ZV, 0.1)
				mx = mx.Moved(gubbe.position)
				gubbeStandingRightSprite.Draw(win, mx)
			}
		}

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
	if ballState.Pos.Y < ballradius {
		ballState.Pos.Y = ballradius + 1
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

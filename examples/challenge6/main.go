package main

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
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
		Title:    "Bouncing ball",
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
	var ballSprite, err2 = internal.LoadSprite("assets/Ball.png")
	config, err = TryReadCfgFrom("json/challenge6.json", config)
	internal.PanicIfError(err)
	var ballState = Ball{
		Pos: pixel.Vec{X: config.StartX, Y: config.StartY},
		Vel: pixel.Vec{X: config.SpeedX, Y: 0}, Rot: 0,
	}
	internal.PanicIfError(err2)
	for !win.Closed() {
		win.Clear(colornames.Lightskyblue)

		var now = time.Now()
		var delta = now.Sub(prevtime).Seconds()
		prevtime = now

		config, err = TryReadCfgFrom("json/challenge6.json", config)
		internal.PanicIfError(err)

		rot += delta
		mx := pixel.IM
		mx = mx.Moved(ballState.Pos)
		mx = mx.Rotated(ballState.Pos, ballState.Rot)
		ballSprite.Draw(win, mx)

		// Update ball
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

		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

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

// rita blårektangel
// animera med tidsstämpel (millisekunder t.ex.)

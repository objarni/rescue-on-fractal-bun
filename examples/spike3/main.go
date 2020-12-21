package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	time2 "time"
)

const screenwidth = 500
const boxwidth = 30
const speed = 200

func run() {
	cfg := pixelgl.WindowConfig{
		Title:    "Animated rectangle",
		Bounds:   pixel.R(0, 0, screenwidth, 300),
		Position: pixel.Vec{X: 300, Y: 300},
	}
	win, err := pixelgl.NewWindow(cfg)
	failOnError(err)

	var x float64 = 0
	imd := imdraw.New(nil)
	dir := 1.0

	time := time2.Now()
	for !win.Closed() {
		now := time2.Now()
		delta := now.Sub(time).Seconds()
		time = now

		x = x + delta*speed*dir
		if x > screenwidth-boxwidth {
			x = screenwidth - boxwidth
			dir = -1
		}
		if x < 0 {
			x = 0
			dir = 1
		}

		win.Clear(colornames.Blue)
		imd.Clear()
		imd.Push(pixel.Vec{X: x, Y: 20})
		imd.Push(pixel.Vec{X: x + 30, Y: 40})
		imd.Rectangle(0)
		imd.Draw(win)
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		win.Update()
	}
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pixelgl.Run(run)
}

// rita blårektangel
// animera med tidsstämpel (millisekunder t.ex.)

package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:    "Push A, S and D to play drums!",
		Bounds:   pixel.R(0, 0, 400, 300),
		Position: pixel.Vec{X: 500, Y: 500},
	}
	win, err := pixelgl.NewWindow(cfg)
	failOnError(err)

	imd := imdraw.New(nil)
	imd.Push(pixel.Vec{X: 10, Y: 20})
	imd.Push(pixel.Vec{X: 30, Y: 40})
	imd.Rectangle(0)

	for !win.Closed() {
		win.Clear(colornames.Blue)
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

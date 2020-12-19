package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 500), basicAtlas)

	var x float64 = 0
	var y float64 = 0

	for !win.Closed() {
		win.Clear(colornames.Black)

		drawHelloWorldAt(basicTxt, x, y, win)

		win.Update()
	}
}

func drawHelloWorldAt(basicTxt *text.Text, x float64, y float64, win *pixelgl.Window) {
	basicTxt.Clear()
	basicTxt.Orig = pixel.V(x, y)
	fmt.Fprintln(basicTxt, "Hello, text!")
	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
}

func main() {
	pixelgl.Run(run)
}

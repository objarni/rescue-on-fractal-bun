package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Asd struct {
	A bool
	S bool
	D bool
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Push A, S and D to play drums!",
		Bounds: pixel.R(0, 0, 400, 300),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	//basicTxt := text.New(pixel.V(100, 500), basicAtlas)

	asd := Asd{false, false, false}

	for !win.Closed() {
		win.Clear(colornames.Blue)
		asd.A = win.Pressed(pixelgl.KeyA)
		asd.S = win.Pressed(pixelgl.KeyS)
		asd.D = win.Pressed(pixelgl.KeyD)
		fmt.Println(asd)
		win.Update()
	}
}

//func drawHelloWorldAt(basicTxt *text.Text, pos Pos, win *pixelgl.Window) {
//	basicTxt.Clear()
//	basicTxt.Orig = pixel.V(pos.X, pos.Y)
//	_, _ = fmt.Fprintln(basicTxt, "Hello, text!")
//	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
//}

func main() {
	pixelgl.Run(run)
}

package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"io/ioutil"
	"strconv"
)

//type Pos struct {
//	X int `json:"x"`
//	Y int `json:"y"`
//}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Move hello world by changing pos.json",
		Bounds: pixel.R(0, 0, 400, 300),
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

		// Start with reading the position every frame;
		// Really resource intensive but it's a start!
		x = float64(tryReadIntFrom("x.txt", int(x)))
		y = float64(tryReadIntFrom("y.txt", int(y)))

		drawHelloWorldAt(basicTxt, x, y, win)

		win.Update()
	}
}

func tryReadIntFrom(fileName string, valueOnError int) int {
	content, err := ioutil.ReadFile(fileName)
	fmt.Printf("Maybe Read: %s\n", content)

	if err != nil {
		fmt.Printf("Was an error, defaulting\n")
		return valueOnError
	}

	i, err := strconv.Atoi(string(content))
	if err != nil {
		fmt.Printf("no parse, defaulting\n")
		return valueOnError
	}
	fmt.Printf("parsed: %d\n", i)
	return i
}

func drawHelloWorldAt(basicTxt *text.Text, x float64, y float64, win *pixelgl.Window) {
	basicTxt.Clear()
	basicTxt.Orig = pixel.V(x, y)
	_, _ = fmt.Fprintln(basicTxt, "Hello, text!")
	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
}

func main() {
	pixelgl.Run(run)
}

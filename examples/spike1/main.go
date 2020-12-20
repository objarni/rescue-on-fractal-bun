package main

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"io/ioutil"
)

type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

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

	var pos Pos = Pos{200, 200}

	for !win.Closed() {
		win.Clear(colornames.Black)
		pos = TryReadPosFrom("pos.json", pos)
		drawHelloWorldAt(basicTxt, pos, win)
		win.Update()
	}
}

func TryReadPosFrom(filename string, defaultPos Pos) Pos {
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile error, defaulting")
		return defaultPos
	}
	var pos Pos
	err = json.Unmarshal(byteArray, &pos)
	if err != nil {
		fmt.Println("JSON parse error, defaulting. JSON was: " + string(byteArray))
		return defaultPos
	}
	return pos
}
func drawHelloWorldAt(basicTxt *text.Text, pos Pos, win *pixelgl.Window) {
	basicTxt.Clear()
	basicTxt.Orig = pixel.V(pos.X, pos.Y)
	_, _ = fmt.Fprintln(basicTxt, "Hello, text!")
	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
}

func main() {
	pixelgl.Run(run)
}

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
	"time"
)

type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Goroutine protocol
//  run <- posReader
// .. But to avoid blocking on read, run does read
//  with a select which defaults to 'do nothing'.
//  Then posReader basically sleeps x amount time,
//  then tries to read pos.json and on success
//  sends a Pos value over the channel 'positions'.
//  The channel can be unbuffered since run
//  reads next frame.

func readPosition(position chan<- Pos) {
	for {
		pos, err := TryReadPosFrom("json/challenge1.json", Pos{0, 0})
		if err == nil {
			position <- pos
		}
		time.Sleep(time.Second)
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Move hello world by changing spike1.json",
		Bounds: pixel.R(0, 0, 400, 300),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 500), basicAtlas)

	var pos = Pos{200, 200}

	positions := make(chan Pos)
	go readPosition(positions)

	for !win.Closed() {
		win.Clear(colornames.Black)

		select {
		case newPos := <-positions:
			pos = newPos
		default:
		}

		drawHelloWorldAt(basicTxt, pos, win)
		win.Update()
		time.Sleep(time.Millisecond)
	}
}

func TryReadPosFrom(filename string, defaultPos Pos) (Pos, error) {
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		//fmt.Println("ReadFile error, defaulting")
		return defaultPos, err
	}
	var pos Pos
	err = json.Unmarshal(byteArray, &pos)
	if err != nil {
		//fmt.Println("JSON parse error, defaulting. JSON was: " + string(byteArray))
		return defaultPos, err
	}
	return pos, nil
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

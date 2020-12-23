package main

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"os"
	"time"
)

const screenwidth = 800
const screenheight = 600
const speed = 200
const boxwidth = 30

func run() {
	cfg := pixelgl.WindowConfig{
		Title:    "Animated rectangle",
		Bounds:   pixel.R(0, 0, screenwidth, screenheight),
		Position: pixel.Vec{X: 300, Y: 300},
	}
	win, err := pixelgl.NewWindow(cfg)
	panicOnError(err)

	backgroundSprite, err := loadSprite("assets/MapSketch.jpg")
	panicOnError(err)

	var x float64 = 0
	var imd = imdraw.New(nil)
	var dir = 1.0
	var config Config
	var prevtime = time.Now()
	for !win.Closed() {
		var now = time.Now()
		var delta = now.Sub(prevtime).Seconds()
		prevtime = now

		x = x + delta*speed*dir
		if x > screenwidth-boxwidth {
			x = screenwidth - boxwidth
			dir = -1
		}
		if x < 0 {
			x = 0
			dir = 1
		}

		config, err = TryReadCfgFrom("json/spike4.json", config)
		panicOnError(err)

		win.Clear(colornames.Blue)
		backgroundSprite.Draw(win, pixel.IM.
			Scaled(pixel.ZV, config.Scale).
			Moved(win.Bounds().Center()),
		)

		imd.Clear()
		imd.Push(pixel.Vec{X: x, Y: 20})
		imd.Push(pixel.Vec{X: x + 30, Y: 40})
		imd.Rectangle(0)
		imd.Draw(win)

		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		win.Update()
		time.Sleep(time.Millisecond * 5)
	}
}

type Config struct {
	Scale float64
}

func TryReadCfgFrom(filename string, defaultCfg Config) (Config, error) {
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile error, defaulting")
		return defaultCfg, err
	}
	var cfg Config = defaultCfg
	err = json.Unmarshal(byteArray, &cfg)
	if err != nil {
		fmt.Println("JSON parse error, defaulting. JSON was: " + string(byteArray))
		return defaultCfg, err
	}
	return cfg, nil
}

func loadSprite(path string) (*pixel.Sprite, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	dim := img.Bounds().Max
	frame := pixel.R(0, 0, float64(dim.X), float64(dim.Y))
	sprite := pixel.NewSprite(pic, frame)
	return sprite, nil
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pixelgl.Run(run)
}

// rita blårektangel
// animera med tidsstämpel (millisekunder t.ex.)

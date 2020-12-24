package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

const screenwidth = 800
const screenheight = 600

type Drop struct {
	X float64
	Y float64
	Z float64
}

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

	drops := list.New()
	drops.PushBack(&Drop{40, 40, 40})

	var x float64 = 0
	var imd = imdraw.New(nil)
	imd.Color = colornames.Darkslateblue
	var dir = 1.0
	var config Config
	var prevtime = time.Now()
	for !win.Closed() {
		var now = time.Now()
		var delta = now.Sub(prevtime).Seconds()
		prevtime = now

		x = x + delta*config.DropSpeed*dir
		if x > float64(screenwidth-config.DropLength) {
			x = float64(screenwidth - config.DropLength)
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

		// Om det finns färre än max antal droppar,
		// lägg till en ovanför skärmen enl. config
		// fast med slump-x och slump-z
		// när en droppe har z < 0 ta bort den
		// varje frame flytta i y-led enligt dt
		imd.Clear()
		for drop := drops.Front(); drop != nil; {
			var d *Drop = drop.Value.(*Drop)
			updateDrop(d, delta, config)

			imd.Push(pixel.Vec{X: d.X, Y: d.Y})
			imd.Push(pixel.Vec{X: d.X + 1, Y: d.Y + config.DropLength})
			imd.Rectangle(0)

			if d.Z < 0 {
				remove := drop
				drop = drop.Next()
				drops.Remove(remove)
			} else {
				drop = drop.Next()
			}

		}

		if drops.Len() < config.DropMaxCount {
			drops.PushBack(&Drop{
				float64(rand.Intn(screenwidth)),
				float64(rand.Intn(screenheight)),
				float64(rand.Intn(config.DropMaxLife)),
			})
		}
		fmt.Println(drops.Len())

		imd.Draw(win)

		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		win.Update()
		time.Sleep(time.Millisecond * 5)
	}
}

func updateDrop(d *Drop, delta float64, config Config) {
	d.Y -= delta * config.DropSpeed
	d.Z -= delta * config.DropSpeed
}

type Config struct {
	Scale        float64
	DropLength   float64
	DropSpeed    float64
	DropMaxCount int
	DropMaxLife  int
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

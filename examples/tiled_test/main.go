package main

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
	_ "image/png"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "TilePix",
		Bounds: pixel.R(0, 0, 320, 320),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Load and initialise the map.
	m, err := tilepix.ReadFile("examples/tiled_test/Level2.tmx")
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.Aqua)

		//Draw all layers to the window.
		if err := m.DrawAll(win, color.White, pixel.IM); err != nil {
			panic(err)
		}

		//sprite := m.TileLayers[0].Tileset.Sprite
		//sprite.Draw(win, pixel.IM.Moved(pixel.Vec{300, 150}))
		//_ = m.DrawAll(win, colornames.Green, pixel.IM)
		//m.TileLayers[0].Draw(win)
		win.Update()
	}
}

type Config struct {
	Gravity float64
	SpeedX  float64
	StartX  float64
	StartY  float64
}

func main() {
	pixelgl.Run(run)
}

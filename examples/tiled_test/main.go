package main

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	_ "image/png"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "TilePix",
		Bounds: pixel.R(0, 0, 640, 320),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Load and initialise the map.
	m, err := tilepix.ReadFile("assets/levels/GhostForest.tmx")
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(color.White)

		// Draw all layers to the window.
		if err := m.DrawAll(win, color.White, pixel.IM); err != nil {
			panic(err)
		}

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

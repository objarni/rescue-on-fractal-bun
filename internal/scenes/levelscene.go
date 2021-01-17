package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
)

type LevelScene struct {
	cfg       *Config
	playerPos pixel.Vec
	level     Level
}

type Level struct {
	width, height int32
	clearColor    color.RGBA
	mapPoints     []MapPoint
}

type MapPoint struct {
	pos        pixel.Vec
	discovered bool
	mapTarget  string
}

func MakeLevelScene(cfg *Config) *LevelScene {
	return &LevelScene{
		cfg:       cfg,
		playerPos: pixel.Vec{3500, 600},
		level: Level{
			width:      5000,
			height:     768,
			clearColor: colornames.Blue900,
		},
	}
}

func (scene *LevelScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	return MakeMenuScene(scene.cfg)
}

func (scene *LevelScene) HandleKeyUp(key internal.ControlKey) internal.Thing {
	return scene
}

func (scene *LevelScene) Render(win *pixelgl.Window) {
	win.Clear(scene.level.clearColor)
}

func (scene *LevelScene) Tick() bool {
	return true
}

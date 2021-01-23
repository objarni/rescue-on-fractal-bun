package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
)

type LevelScene struct {
	cfg                       *Config
	playerPos                 pixel.Vec
	leftPressed, rightPressed bool
	level                     Level
}

type Level struct {
	width, height float64
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
		playerPos: pixel.Vec{X: 3000, Y: 60},
		level: Level{
			width:      3500,
			height:     768,
			clearColor: colornames.Blue900,
		},
	}
}

func (scene *LevelScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Left {
		scene.leftPressed = true
	}
	if key == internal.Right {
		scene.rightPressed = true
	}
	if key == internal.Action {
		return MakeMapScene(scene.cfg, "Korsningen")
	}
	return scene
}

func (scene *LevelScene) HandleKeyUp(key internal.ControlKey) internal.Thing {
	if key == internal.Left {
		scene.leftPressed = false
	}
	if key == internal.Right {
		scene.rightPressed = false
	}
	return scene
}

func (scene *LevelScene) Render(win *pixelgl.Window) {
	// Clear screen
	win.Clear(colornames.Black)

	// Camera transform
	imd := imdraw.New(nil)
	pw := internal.PlayerWidth
	ph := internal.PlayerHeight
	halfScreen := v(internal.ScreenWidth/2, internal.ScreenHeight/2)
	playerHead := v(0, ph)
	cam := scene.playerPos.Sub(halfScreen).Add(playerHead)
	imd.SetMatrix(pixel.IM.Moved(cam.Scaled(-1)))

	// Level backdrop
	imd.Color = scene.level.clearColor
	imd.Push(v(0, 0))
	imd.Push(v(scene.level.width, scene.level.height))
	imd.Rectangle(0)

	// Player
	imd.Color = colornames.Brown200
	px := scene.playerPos.X
	py := scene.playerPos.Y
	plBL := v(px-pw/2, py)
	plTR := v(px+pw/2, py+ph)
	imd.Push(plBL)
	imd.Push(plTR)
	imd.Rectangle(0)
	imd.Draw(win)
}

func (scene *LevelScene) Tick() bool {
	if scene.leftPressed && !scene.rightPressed {
		scene.playerPos = scene.playerPos.Add(v(-5, 0))
	}
	if !scene.leftPressed && scene.rightPressed {
		scene.playerPos = scene.playerPos.Add(v(5, 0))
	}

	return true
}

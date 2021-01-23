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

// TODO: discovered should probably be stored somewhere else
type MapPoint struct {
	pos        pixel.Vec
	discovered bool
	mapTarget  string
}

func MakeLevelScene(cfg *Config) *LevelScene {
	mapPoints := []MapPoint{
		{
			pos:        pixel.Vec{3400, 60},
			discovered: true,
			mapTarget:  "Hembyn",
		},
		{
			pos:        pixel.Vec{300, 60},
			discovered: false,
			mapTarget:  "Korsningen",
		},
	}
	return &LevelScene{
		cfg:       cfg,
		playerPos: pixel.Vec{X: 3000, Y: 60},
		level: Level{
			width:      3500,
			height:     768,
			clearColor: colornames.Blue900,
			mapPoints:  mapPoints,
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
	win.Clear(colornames.Black)
	imd := scene.cameraTransform()
	scene.drawBackdrop(imd)
	scene.drawMapPoints(win, imd)
	scene.drawPlayer(win, imd)
	imd.Draw(win)
}

func (scene *LevelScene) cameraTransform() *imdraw.IMDraw {
	//	start := time.Now()
	imd := imdraw.New(nil)
	//	duration := time.Since(start)
	//	fmt.Println(duration)
	halfScreen := v(internal.ScreenWidth/2, internal.ScreenHeight/2)
	playerHead := v(0, internal.PlayerHeight)
	cam := scene.playerPos.Sub(halfScreen).Add(playerHead)
	imd.SetMatrix(pixel.IM.Moved(cam.Scaled(-1)))
	return imd
}

func (scene *LevelScene) drawPlayer(win *pixelgl.Window, imd *imdraw.IMDraw) {
	playerWidthHalf := internal.PlayerWidth / 2
	playerBottomLeft := scene.playerPos.Sub(v(playerWidthHalf, 0))
	playerTopRight := scene.playerPos.Add(v(playerWidthHalf, internal.PlayerHeight))

	imd.Color = colornames.Brown200
	imd.Push(playerBottomLeft)
	imd.Push(playerTopRight)
	imd.Rectangle(0)
	imd.Color = colornames.Brown100
	imd.Push(playerBottomLeft)
	imd.Push(playerTopRight)
	imd.Rectangle(2)
}

func (scene *LevelScene) drawBackdrop(imd *imdraw.IMDraw) {
	imd.Color = scene.level.clearColor
	imd.Push(v(0, 0))
	imd.Push(v(scene.level.width, scene.level.height))
	imd.Rectangle(0)
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

func (scene *LevelScene) drawMapPoints(win *pixelgl.Window, imd *imdraw.IMDraw) {
	for _, mapPoint := range scene.level.mapPoints {
		if mapPoint.discovered {
			imd.Color = colornames.Green800
		} else {
			imd.Color = colornames.Orange500
		}
		imd.Push(mapPoint.pos)
		imd.Push(mapPoint.pos.Add(v(10, 10)))
		imd.Rectangle(0)
	}
}

/*
type MapPoint struct {
	pos        pixel.Vec
	discovered bool
	mapTarget  string
}
*/

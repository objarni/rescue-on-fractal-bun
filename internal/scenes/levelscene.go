package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"objarni/rescue-on-fractal-bun/internal"
)

type LevelScene struct {
	cfg                       *Config
	res                       *Resources
	playerPos                 pixel.Vec
	leftPressed, rightPressed bool
	level                     internal.Level
	ghost                     *pixel.Sprite // TODO: move sprites to res.structure
}

func MakeLevelScene(cfg *Config, res *Resources) *LevelScene {
	mapPoints := []internal.MapPoint{
		{
			Pos:        pixel.Vec{X: 3400, Y: 60},
			Discovered: true,
			Location:   "Hembyn",
		},
		{
			Pos:        pixel.Vec{X: 300, Y: 60},
			Discovered: false,
			Location:   "Korsningen",
		},
	}
	ghost := internal.LoadSpriteForSure("assets/TGhost.png")
	return &LevelScene{
		cfg:       cfg,
		res:       res,
		playerPos: pixel.Vec{X: 3000, Y: 60},
		level: internal.Level{
			Width:      3500,
			Height:     768,
			ClearColor: colornames.Blue900,
			MapPoints:  mapPoints,
		},
		ghost: ghost,
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
		return MakeMapScene(scene.cfg, scene.res, "Korsningen")
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
	camMx := scene.cameraMatrix()
	for i := 0; i < int(scene.level.Width); i += 500 {
		scene.ghost.Draw(win, pixel.IM.Moved(v(float64(i), 200)).Chained(camMx))
	}
}

func (scene *LevelScene) cameraTransform() *imdraw.IMDraw {
	//	start := time.Now()
	imd := imdraw.New(nil)
	//	duration := time.Since(start)
	//	fmt.Println(duration)
	cameraMatrix := scene.cameraMatrix()
	imd.SetMatrix(cameraMatrix)
	return imd
}

func (scene *LevelScene) cameraMatrix() pixel.Matrix {
	halfScreen := v(internal.ScreenWidth/2, internal.ScreenHeight/2)
	playerHead := v(0, internal.PlayerHeight)
	cam := scene.playerPos.Sub(halfScreen).Add(playerHead)
	cameraMatrix := pixel.IM.Moved(cam.Scaled(-1))
	return cameraMatrix
}

func (scene *LevelScene) drawPlayer(_ *pixelgl.Window, imd *imdraw.IMDraw) {
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
	imd.Color = scene.level.ClearColor
	imd.Push(v(0, 0))
	imd.Push(v(float64(scene.level.Width), float64(scene.level.Height)))
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

func (scene *LevelScene) drawMapPoints(_ *pixelgl.Window, imd *imdraw.IMDraw) {
	for _, mapPoint := range scene.level.MapPoints {
		if mapPoint.Discovered {
			imd.Color = colornames.Green800
		} else {
			imd.Color = colornames.Orange500
		}
		imd.Push(mapPoint.Pos)
		imd.Push(mapPoint.Pos.Add(v(10, 10)))
		imd.Rectangle(0)
	}
}

/*
type MapPoint struct {
	Pos        pixel.Vec
	Discovered bool
	Location  string
}
*/

package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"objarni/rescue-on-fractal-bun/internal"
)

type LevelScene struct {
	cfg                       *Config
	res                       *Resources
	playerPos                 pixel.Vec
	leftPressed, rightPressed bool
	level                     internal.Level
}

func MakeLevelScene(cfg *Config, res *Resources) *LevelScene {
	level := internal.LoadLevel("assets/levels/GhostForest.tmx")
	pos := level.MapPoints[0].Pos
	return &LevelScene{
		cfg:       cfg,
		res:       res,
		playerPos: pos,
		level:     level,
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
	// Clear screen
	win.Clear(colornames.Black)
	camMx := scene.cameraMatrix()
	win.SetMatrix(camMx)
	imd := imdraw.New(nil)

	// map rectangle. TODO: remove when player cannot see past limits!
	// Then, just clear screen to map background color
	scene.drawBackdrop(imd)

	layers := scene.level.TilepixMap.TileLayers
	_ = layers[0].Draw(win) // Background
	_ = layers[1].Draw(win) // Platforms
	_ = layers[2].Draw(win) // Walls

	// Draw objects
	scene.drawMapPoints(win, imd)
	scene.drawPlayer(win)
	//imd.Draw(win)
	for i := 0; i < scene.level.Width; i += 500 {
		scene.res.Ghost.Draw(win,
			pixel.IM.Moved(v(float64(i), 200)))
	}

	_ = layers[3].Draw(win) // Foreground

	// Heads-up display

	// FPS
	win.SetMatrix(pixel.IM)
	tb := text.New(pixel.V(0, 0), scene.res.Atlas)
	_, _ = fmt.Fprintf(tb, "FPS=%1.1f", scene.res.FPS)
	tb.DrawColorMask(win, pixel.IM, colornames.White)
}

func (scene *LevelScene) cameraMatrix() pixel.Matrix {
	halfScreen := v(internal.ScreenWidth/2, internal.ScreenHeight/2)
	playerHead := v(0, internal.PlayerHeight)
	cam := scene.playerPos.Sub(halfScreen).Add(playerHead)
	cameraMatrix := pixel.IM.Moved(cam.Scaled(-1))
	return cameraMatrix
}

func (scene *LevelScene) drawPlayer(win *pixelgl.Window) {
	scene.res.PlayerStanding.Draw(win, pixel.IM.Moved(scene.playerPos))
}

func (scene *LevelScene) drawBackdrop(imd *imdraw.IMDraw) {
	imd.Color = scene.level.ClearColor
	imd.Push(v(0, 0))
	imd.Push(v(float64(scene.level.Width*32), float64(scene.level.Height*32)))
	imd.Rectangle(0)
}

func (scene *LevelScene) Tick() bool {
	if scene.leftPressed && !scene.rightPressed {
		scene.playerPos = scene.playerPos.Add(v(-scene.cfg.LevelSceneMoveSpeed, 0))
	}
	if !scene.leftPressed && scene.rightPressed {
		scene.playerPos = scene.playerPos.Add(v(scene.cfg.LevelSceneMoveSpeed, 0))
	}
	return true
}

func (scene *LevelScene) drawMapPoints(win *pixelgl.Window, imd *imdraw.IMDraw) {
	for _, mapPoint := range scene.level.MapPoints {
		alignVec := v(0, scene.res.MapPoint.Frame().Center().Y)
		scene.res.MapPoint.Draw(
			win,
			pixel.IM.Moved(mapPoint.Pos).Moved(alignVec))
		if mapPoint.Discovered {
			imd.Color = colornames.Green800
		} else {
			imd.Color = colornames.Orange500
		}
		imd.Push(mapPoint.Pos.Sub(v(4, 4)))
		imd.Push(mapPoint.Pos.Add(v(4, 4)))
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

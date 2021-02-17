package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
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
	pos := level.MapSigns[0].Pos
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
	win.Clear(colornames.Yellow50)
	win.SetMatrix(scene.cameraMatrix())

	// Level backdrop
	// TODO: remove when player cannot see past limits!
	// Then, just clear screen to map background color
	imd := imdraw.New(nil)
	scene.backdropGfx().Render(imd)
	imd.Draw(win)

	layers := scene.level.TilepixMap.TileLayers
	_ = layers[0].Draw(win) // Background
	_ = layers[1].Draw(win) // Platforms
	_ = layers[2].Draw(win) // Walls

	// Draw objects
	scene.drawMapPoints(win)
	scene.drawPlayer(win)
	for i := 0; i < scene.level.Width; i += 500 {
		scene.res.Ghost.Draw(win,
			pixel.IM.Moved(v(float64(i), 200)))
	}

	_ = layers[3].Draw(win) // Foreground

	scene.drawHeadsUpDisplay(win)
}

func (scene *LevelScene) drawHeadsUpDisplay(win *pixelgl.Window) {
	win.SetMatrix(pixel.IM)
	if scene.isMapSignClose() {
		scene.res.InLevelHeadsUp.DrawColorMask(
			win,
			pixel.IM.Moved(scene.res.InLevelHeadsUp.Frame().Center()),
			colornames.GreenA400)
	} else {
		scene.res.InLevelHeadsUp.Draw(
			win,
			pixel.IM.Moved(scene.res.InLevelHeadsUp.Frame().Center()))
	}

	// FPS
	tb := text.New(pixel.V(0, 0), scene.res.Atlas)
	_, _ = fmt.Fprintf(tb, "FPS=%1.1f", scene.res.FPS)
	tb.DrawColorMask(win, pixel.IM, colornames.Brown800)
}

func (scene *LevelScene) isMapSignClose() bool {
	for _, mapSign := range scene.level.MapSigns {
		if scene.playerPos.Sub(mapSign.Pos).Len() < 10 {
			return true
		}
	}
	return false
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

func (scene *LevelScene) backdropGfx() draw.ImdOp {
	widthPixels := scene.level.Width * 32
	heightPixels := scene.level.Height * 32
	rectangle := draw.Rectangle(draw.C(0, 0), draw.C(widthPixels, heightPixels), 0)
	return draw.Colored(scene.level.ClearColor, rectangle)
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

func (scene *LevelScene) drawMapPoints(win *pixelgl.Window) {
	for _, mapPoint := range scene.level.MapSigns {
		alignVec := v(0, scene.res.MapPoint.Frame().Center().Y)
		scene.res.MapPoint.Draw(
			win,
			pixel.IM.Moved(mapPoint.Pos).Moved(alignVec))
	}
}

/*
Windows rit-API:er som används i Rescue:
image.Draw(win, mx)
image.DrawColorMask(win, mx, color)
text.DrawColorMask(win, mx, color)
text.Draw(win, pixel.IM)  # används inte men finns i API:et!
layer.Draw(win)
imd.Draw(win)

Samtliga påverkas av win.Matrix (sätts med win.SetMatrix).

Förenkling: skulle kunna använda identitetsmatris i de
anrop som har mx parameter, för att istället _alltid_ använda
win.Matrix.

Har även hittat en "SetColorMask" i win; detta betyder
att jag kan unifiera till att bara använda Draw()-anrop,
och därmed flytta ut denna data/kunskap till modellen, så
att det finns allmänna Color operationer att beskriva grafiken med.
Det blir då desto viktigare att dessa "resettas" efter ritoperationer
eftersom de annars kommer spilla över i t.ex. image eller layer ritning
(antar jag).


*/

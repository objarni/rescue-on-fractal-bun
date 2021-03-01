package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

type LevelScene struct {
	cfg                       *Config
	res                       *internal.Resources
	playerPos                 pixel.Vec
	leftPressed, rightPressed bool
	level                     internal.Level
}

func MakeLevelScene(cfg *Config, res *internal.Resources) *LevelScene {
	level := internal.LoadLevel("assets/levels/GhostForest.tmx")
	pos := level.SignPosts[0].Pos
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
		if scene.isMapSignClose() {
			mapSign := scene.closestMapSign()
			return MakeMapScene(scene.cfg, scene.res, mapSign.Location)
		} else {
			scene.playerPos = scene.playerPos.Add(v(10, 0))
		}
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

	// Level backdrop
	// TODO: remove when player cannot see past limits!
	// Then, just clear screen to map background color
	movedBackdrop := draw.Moved(scene.cameraVector(), draw.ToWinOp(scene.backdropGfx()))
	movedBackdrop.Render(pixel.IM, win)
	//draw.ToWinOp(scene.backdropGfx()).Render(win)

	movedBackground := draw.Moved(scene.cameraVector(), draw.TileLayer(scene.level.TilepixMap, "Background"))
	movedBackground.Render(pixel.IM, win)

	movedPlatforms := draw.Moved(scene.cameraVector(), draw.TileLayer(scene.level.TilepixMap, "Platforms"))
	movedPlatforms.Render(pixel.IM, win)

	movedWalls := draw.Moved(scene.cameraVector(), draw.TileLayer(scene.level.TilepixMap, "Walls"))
	movedWalls.Render(pixel.IM, win)

	win.SetMatrix(scene.cameraMatrix())

	// Draw objects
	scene.drawSignPosts(win)
	win.SetMatrix(scene.cameraMatrix())
	scene.drawPlayer(win)

	// IGhost
	win.SetMatrix(pixel.IM)
	ghostPos := v(float64(0), 200)
	ghostSprite := draw.Moved(scene.cameraVector(), draw.Moved(ghostPos, draw.Image(scene.res.ImageMap, internal.IGhost)))
	ghostSprite.Render(pixel.IM, win)
	win.SetMatrix(scene.cameraMatrix())

	movedForeground := draw.Moved(scene.cameraVector(), draw.TileLayer(scene.level.TilepixMap, "Foreground"))
	movedForeground.Render(pixel.IM, win)

	scene.drawHeadsUpDisplay(win)
}

func (scene *LevelScene) drawHeadsUpDisplay(win *pixelgl.Window) {
	// TODO: crop this screen-sized image and translate it in position
	// (coloring only works now since it's the only image used in headsup!)
	mapSymbolCenter := scene.res.ImageMap[internal.IMapSymbol].Frame().Center()
	op := draw.Moved(mapSymbolCenter, draw.Image(scene.res.ImageMap, internal.IMapSymbol))
	if scene.isMapSignClose() {
		op = draw.Color(colornames.GreenA400, op)
	}
	op.Render(pixel.IM, win)

	// FPS
	win.SetMatrix(pixel.IM)

	tb := text.New(pixel.V(0, 0), scene.res.Atlas)
	_, _ = fmt.Fprintf(tb, "FPS=%1.1f", scene.res.FPS)
	tb.DrawColorMask(win, pixel.IM, colornames.Brown800)
}

func (scene *LevelScene) isMapSignClose() bool {
	sign := scene.closestMapSign()
	return scene.playerPos.Sub(sign.Pos).Len() < 10
}

func (scene *LevelScene) closestMapSign() internal.SignPost {
	closestMapPoint := internal.SignPost{}
	shortestDistance := 1000000000.0
	for _, mapSign := range scene.level.SignPosts {
		mapSignDistance := scene.playerPos.Sub(mapSign.Pos).Len()
		if mapSignDistance < shortestDistance {
			closestMapPoint = mapSign
			shortestDistance = mapSignDistance
		}
	}
	return closestMapPoint
}

func (scene *LevelScene) cameraMatrix() pixel.Matrix {
	reversed := scene.cameraVector()
	cameraMatrix := pixel.IM.Moved(reversed)
	return cameraMatrix
}

func (scene *LevelScene) cameraVector() pixel.Vec {
	halfScreen := v(internal.ScreenWidth/2, internal.ScreenHeight/2)
	playerHead := v(0, internal.PlayerHeight)
	cam := scene.playerPos.Sub(halfScreen).Add(playerHead)
	reversed := cam.Scaled(-1)
	return reversed
}

func (scene *LevelScene) drawPlayer(win *pixelgl.Window) {
	draw.Moved(
		scene.playerPos.Add(scene.cameraVector()),
		draw.Image(scene.res.ImageMap, internal.ITemporaryPlayerImage)).Render(pixel.IM, win)
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

func (scene *LevelScene) drawSignPosts(win *pixelgl.Window) {
	for _, mapPoint := range scene.level.SignPosts {
		alignVec := v(0, scene.res.ImageMap[internal.ISignPost].Frame().Center().Y)
		signPostOp := draw.Moved(scene.cameraVector(), draw.Moved(mapPoint.Pos, draw.Moved(alignVec,
			draw.Image(scene.res.ImageMap, internal.ISignPost))))
		signPostOp.Render(pixel.IM, win)
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
Det blir då desto viktigare att dessa "resets" efter rit operationer
eftersom de annars kommer spilla över i t.ex. image eller layer ritning
(antar jag).


*/

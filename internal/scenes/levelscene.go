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

func MakeLevelScene(cfg *Config, res *internal.Resources, levelName string) *LevelScene {
	level := res.Levels[levelName]
	//level := internal.LoadLevel("assets/levels/GhostForest.tmx")
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
	// TODO: if player cannot see past level limists,
	// this clear is not needed (camera like WonderBoy)
	win.Clear(colornames.Yellow50)

	gfx := draw.OpSequence(
		draw.Moved(
			scene.cameraVector(),
			draw.OpSequence(
				draw.ToWinOp(scene.backdropOp()),
				draw.TileLayer(scene.level.TilepixMap, "Background"),
				draw.TileLayer(scene.level.TilepixMap, "Platforms"),
				draw.TileLayer(scene.level.TilepixMap, "Walls"),
				scene.signPostsOp(),
				scene.playerOp(),
				scene.ghostOp(),
				draw.TileLayer(scene.level.TilepixMap, "Foreground"),
			),
		),
		scene.mapSymbolOp(),
	)
	gfx.Render(pixel.IM, win)

	scene.drawFPS(win)
}

/* gfxOp start */
func (scene *LevelScene) ghostOp() draw.WinOp {
	return draw.Moved(scene.cameraVector(),
		draw.Moved(v(float64(0), 200),
			draw.Image(scene.res.ImageMap, internal.IGhost)))
}

func (scene *LevelScene) mapSymbolOp() draw.WinOp {
	mapSymbolCenter := scene.res.ImageMap[internal.IMapSymbol].Frame().Center()
	op := draw.Moved(mapSymbolCenter, draw.Image(scene.res.ImageMap, internal.IMapSymbol))
	if scene.isMapSignClose() {
		op = draw.Color(colornames.GreenA400, op)
	}
	return op
}

func (scene *LevelScene) playerOp() draw.WinOp {
	playerOp := draw.Moved(
		scene.playerPos,
		draw.Image(scene.res.ImageMap, internal.ITemporaryPlayerImage))
	return playerOp
}

func (scene *LevelScene) backdropOp() draw.ImdOp {
	widthPixels := scene.level.Width * 32
	heightPixels := scene.level.Height * 32
	rectangle := draw.Rectangle(draw.C(0, 0), draw.C(widthPixels, heightPixels), 0)
	return draw.Colored(scene.level.ClearColor, rectangle)
}

func (scene *LevelScene) signPostsOp() draw.WinOp {
	ops := []draw.WinOp{}
	for _, mapPoint := range scene.level.SignPosts {
		alignVec := v(0, scene.res.ImageMap[internal.ISignPost].Frame().Center().Y)
		signPostOp := draw.Moved(mapPoint.Pos, draw.Moved(alignVec,
			draw.Image(scene.res.ImageMap, internal.ISignPost)))
		ops = append(ops, signPostOp)
	}
	return draw.OpSequence(ops...)
}

/* gfxOp stop */

func (scene *LevelScene) drawFPS(win *pixelgl.Window) {
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
	// Potential: ClosestPoint could take an array of objects implementing
	// 'WithPoint' interface, and we only define anon func here
	getPoint := func(mp internal.SignPost) pixel.Vec { return mp.Pos }
	points := []pixel.Vec{}
	for _, val := range scene.level.SignPosts {
		point := getPoint(val)
		points = append(points, point)
	}
	return scene.level.SignPosts[internal.ClosestPoint(scene.playerPos, points)]
}

func (scene *LevelScene) cameraVector() pixel.Vec {
	halfScreen := v(internal.ScreenWidth/2, internal.ScreenHeight/2)
	playerHead := v(0, internal.PlayerHeight)
	cam := scene.playerPos.Sub(halfScreen).Add(playerHead)
	reversed := cam.Scaled(-1)
	return reversed
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

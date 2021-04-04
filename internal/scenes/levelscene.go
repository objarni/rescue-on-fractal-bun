package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
	"objarni/rescue-on-fractal-bun/internal/entities"
)

type LevelScene struct {
	cfg                       *Config
	res                       *internal.Resources
	playerPos                 pixel.Vec
	leftPressed, rightPressed bool
	level                     internal.Level
	timeMs                    float64
	entities                  []entities.Entity
	elise                     entities.Elise
	entityCanvas              entities.EntityCanvas
}

func MakeLevelScene(cfg *Config, res *internal.Resources, levelName string) *LevelScene {
	level := res.Levels[levelName]
	pos := level.SignPosts[0].Pos.Add(pixel.V(0, -50))
	elise := entities.MakeElise(pos)
	return &LevelScene{
		cfg:          cfg,
		res:          res,
		playerPos:    pos,
		level:        level,
		timeMs:       0,
		entities:     []entities.Entity{entities.MakeGhost(internal.V(2000, 150))},
		elise:        elise,
		entityCanvas: entities.MakeEntityCanvas(),
	}
}

func (scene *LevelScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	scene.elise = scene.elise.HandleKeyDown(key)
	if key == internal.Left {
		scene.leftPressed = true
	}
	if key == internal.Right {
		scene.rightPressed = true
	}
	if key == internal.Action {
		if scene.isMapSignClose() {
			mapSign := scene.closestMapSign()
			return MakeMapScene(scene.cfg, scene.res, mapSign.Text)
		} else {
			scene.playerPos = scene.playerPos.Add(internal.V(10, 0))
		}
	}
	return scene
}

func (scene *LevelScene) HandleKeyUp(key internal.ControlKey) internal.Thing {
	scene.elise = scene.elise.HandleKeyUp(key)
	if key == internal.Left {
		scene.leftPressed = false
	}
	if key == internal.Right {
		scene.rightPressed = false
	}
	return scene
}

func (scene *LevelScene) Render(win *pixelgl.Window) {
	// TODO: if player cannot see past level limits,
	// this clear is not needed (camera like WonderBoy)
	win.Clear(colornames.Yellow50)

	gfx := draw.OpSequence(
		draw.Moved(
			scene.cameraVector(),
			draw.OpSequence(
				draw.ToWinOp(scene.backdropOp()),
				draw.Color(colornames.Black, draw.TileLayer(scene.level.TilepixMap, "Background")),
				draw.Color(colornames.Black, draw.TileLayer(scene.level.TilepixMap, "Platforms")),
				draw.Color(colornames.Black, draw.TileLayer(scene.level.TilepixMap, "Walls")),
				draw.Color(colornames.Black, scene.signPostsOp()),
				draw.TileLayer(scene.level.TilepixMap, "Objects"),
				scene.entitiesOp(scene.timeMs),
				draw.Color(colornames.Black, draw.TileLayer(scene.level.TilepixMap, "Foreground")),
				scene.debugGfx(),
			),
		),
		scene.mapSymbolOp(),
	)
	gfx.Render(pixel.IM, win)

	//scene.drawFPS(win)
}

func (scene *LevelScene) debugGfx() draw.WinOp {
	rectangles := []draw.ImdOp{}
	for _, entity := range scene.entities {
		r := entity.HitBox().HitBox
		rectangles = append(rectangles, rectDrawOp(r))
	}
	rectangles = append(rectangles, rectDrawOp(scene.elise.HitBox().HitBox))
	color := draw.Colored(colornames.White, draw.ImdOpSequence(rectangles...))
	return draw.OpSequence(draw.ToWinOp(color))
}

func rectDrawOp(r pixel.Rect) draw.ImdOp {
	return draw.Rectangle(C(r.Min), C(r.Max), 2)
}

func (scene *LevelScene) entityOp() draw.WinOp {
	return draw.OpSequence(scene.entities[0].GfxOp(&scene.res.ImageMap),
		scene.elise.GfxOp(&scene.res.ImageMap))
}

func (scene *LevelScene) mapSymbolOp() draw.WinOp {
	mapSymbolCenter := scene.res.ImageMap[internal.IMapSymbol].Frame().Center()
	op := draw.Moved(mapSymbolCenter, draw.Image(scene.res.ImageMap, internal.IMapSymbol))
	if scene.isMapSignClose() {
		op = draw.Color(colornames.GreenA400, op)
	}
	return op
}

var frames = [...]internal.Image{
	internal.IEliseWalk6,
	internal.IEliseWalk5,
	internal.IEliseWalk4,
	internal.IEliseWalk3,
	internal.IEliseWalk2,
	internal.IEliseWalk1,
}

func (scene *LevelScene) playerOp(gameTimeS float64) draw.WinOp {
	image := eliseWalkFrame(gameTimeS, scene.cfg.LevelSceneEliseFPS)
	return draw.Moved(scene.playerPos,
		draw.Image(scene.res.ImageMap, image))
}

func eliseWalkFrame(gameTimeS float64, targetFPS int) internal.Image {
	var eliseAnimation = internal.Animation{
		Frames:    6,
		TargetFPS: targetFPS,
	}
	frame := eliseAnimation.FrameAtTime(gameTimeS)
	image := frames[frame]
	return image
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
		alignVec := internal.V(0, scene.res.ImageMap[internal.ISignPost].Frame().Center().Y)
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
	halfScreen := internal.V(internal.ScreenWidth/2, internal.ScreenHeight/2)
	playerHead := internal.V(0, internal.PlayerHeight)
	cam := scene.elise.Pos.Sub(halfScreen).Add(playerHead)
	reversed := cam.Scaled(-1)
	return reversed
}

func (scene *LevelScene) Tick() bool {
	// Handle event boxes from previous tick first of all
	scene.entityCanvas.Consequences(func(eb entities.EventBox, box entities.EntityHitBox) {
		id := box.Entity
		fmt.Printf("%v will handle %v\n", id, eb.Event)
		// Elise?
		if id == -1 {
			scene.elise = scene.elise.Handle(eb)
		} else {
			scene.entities[id] = scene.entities[id].Handle(eb)
		}
		// Otherwise, ordinary entity
	})
	// Reset the canvas
	scene.entityCanvas = entities.MakeEntityCanvas()
	scene.timeMs += 5.0
	scene.elise = scene.elise.Tick()
	scene.entityCanvas.AddEntityHitBox(scene.elise.HitBox())
	for i := range scene.entities {
		scene.entities[i] = scene.entities[i].Tick(&scene.entityCanvas)
		scene.entityCanvas.AddEntityHitBox(scene.entities[i].HitBox())
	}
	return true
}

func (scene *LevelScene) entitiesOp(timeMs float64) draw.WinOp {
	return draw.OpSequence(scene.playerOp(timeMs/1000.0), scene.entityOp())
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

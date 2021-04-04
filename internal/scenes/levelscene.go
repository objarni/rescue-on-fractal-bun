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
	cfg          *Config
	res          *internal.Resources
	level        internal.Level
	timeMs       float64
	entities     []entities.Entity
	elise        entities.Elise
	entityCanvas entities.EntityCanvas
}

func MakeLevelScene(cfg *Config, res *internal.Resources, levelName string) *LevelScene {
	level := res.Levels[levelName]
	pos := level.SignPosts[0].Pos.Add(pixel.V(0, -50))
	elise := entities.MakeElise(pos)
	return &LevelScene{
		cfg:          cfg,
		res:          res,
		level:        level,
		timeMs:       0,
		entities:     []entities.Entity{entities.MakeGhost(internal.V(2000, 150))},
		elise:        elise,
		entityCanvas: entities.MakeEntityCanvas(),
	}
}

func (scene *LevelScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	scene.elise = scene.elise.HandleKeyDown(key)
	if key == internal.Action {
		if scene.isMapSignClose() {
			mapSign := scene.closestMapSign()
			return MakeMapScene(scene.cfg, scene.res, mapSign.Text)
		}
	}
	return scene
}

func (scene *LevelScene) HandleKeyUp(key internal.ControlKey) internal.Thing {
	scene.elise = scene.elise.HandleKeyUp(key)
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
				scene.entityOp(),
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
	rectangles := make([]draw.ImdOp, 0)
	for _, entity := range scene.entities {
		ghost := entity
		rect := ghost.HitBoxRect()
		box4 := entities.EntityHitBox{
			Entity: 0,
			HitBox: rect,
		}
		r := box4.HitBox
		rectangles = append(rectangles, rectDrawOp(r))
	}
	elise := scene.elise
	rect := elise.HitBoxRect()
	box3 := entities.EntityHitBox{
		Entity: -1, // TODO: entities need not know their ID, levelscenes concern
		HitBox: rect,
	}
	rectangles = append(rectangles, rectDrawOp(box3.HitBox))
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

func (scene *LevelScene) backdropOp() draw.ImdOp {
	widthPixels := scene.level.Width * 32
	heightPixels := scene.level.Height * 32
	rectangle := draw.Rectangle(draw.C(0, 0), draw.C(widthPixels, heightPixels), 0)
	return draw.Colored(scene.level.ClearColor, rectangle)
}

func (scene *LevelScene) signPostsOp() draw.WinOp {
	ops := make([]draw.WinOp, 0)
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
	return scene.elise.Pos.Sub(sign.Pos).Len() < 75
}

func (scene *LevelScene) closestMapSign() internal.SignPost {
	// Potential: ClosestPoint could take an array of objects implementing
	// 'WithPoint' interface, and we only define anon func here
	getPoint := func(mp internal.SignPost) pixel.Vec { return mp.Pos }
	points := make([]pixel.Vec, 0)
	for _, val := range scene.level.SignPosts {
		point := getPoint(val)
		points = append(points, point)
	}
	return scene.level.SignPosts[internal.ClosestPoint(scene.elise.Pos, points)]
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
	elise := scene.elise
	rect := elise.HitBoxRect()
	box3 := entities.EntityHitBox{
		Entity: -1, // TODO: entities need not know their ID, levelscenes concern
		HitBox: rect,
	}
	scene.entityCanvas.AddEntityHitBox(box3)
	for i := range scene.entities {
		scene.entities[i] = scene.entities[i].Tick(&scene.entityCanvas)
		ghost := scene.entities[i]
		rect := ghost.HitBoxRect()
		box4 := entities.EntityHitBox{
			Entity: 0,
			HitBox: rect,
		}
		scene.entityCanvas.AddEntityHitBox(box4)
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

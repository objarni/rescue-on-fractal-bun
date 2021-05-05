package scenes

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	d "objarni/rescue-on-fractal-bun/internal/draw"
	"objarni/rescue-on-fractal-bun/internal/entities"
	"objarni/rescue-on-fractal-bun/internal/events"
)

type LevelScene struct {
	cfg          *Config
	res          *internal.Resources
	level        internal.Level
	gameTimeMs   float64
	entities     []entities.Entity
	entityCanvas entities.EntityCanvas
	buttonClick  *beep.Buffer
	robotMoveSfx *beep.Buffer
}

func MakeLevelScene(cfg *Config, res *internal.Resources, levelName string,
	playerPos px.Vec) *LevelScene {
	level := res.Levels[levelName]
	pos := playerPos
	return &LevelScene{
		cfg:          cfg,
		res:          res,
		level:        level,
		gameTimeMs:   0,
		entities:     SpawnEntities(pos, level),
		entityCanvas: entities.MakeEntityCanvas(),
		buttonClick:  res.ButtonClick,
		robotMoveSfx: res.RobotMove,
	}
}

func SpawnEntities(pos px.Vec, level internal.Level) []entities.Entity {
	elise := entities.MakeElise(pos)
	es := []entities.Entity{elise}
	for _, esp := range level.EntitySpawnPoints {
		if esp.EntityType == "Ghost" {
			es = append(es, entities.MakeGhost(esp.SpawnAt))
		} else if esp.EntityType == "Button" {
			es = append(es, entities.MakeButton(esp.SpawnAt))
		} else if esp.EntityType == "Lamp" {
			es = append(es, entities.MakeLamp(esp.SpawnAt))
		} else if esp.EntityType == "Spider" {
			es = append(es, entities.MakeSpider(esp.SpawnAt))
		} else if esp.EntityType == "Robot" {
			es = append(es, entities.MakeRobot(esp.SpawnAt))
		} else {
			panic("Unknown entity type: " + esp.EntityType)
		}
	}
	return es
}

func (scene *LevelScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	event := events.NoEvent
	if key == internal.Left {
		event = events.KeyLeftDown
	}
	if key == internal.Right {
		event = events.KeyRightDown
	}
	if key == internal.Action {
		event = events.KeyActionDown
	}
	if event != events.NoEvent {
		scene.entities[0] = scene.entities[0].Handle(entities.EventBox{
			Event: event,
			Box:   px.Rect{},
		})
	}
	if key == internal.Action {
		if scene.isMapSignClose() {
			mapSign := scene.closestMapSign()
			return MakeMapScene(scene.cfg, scene.res, mapSign.Text)
		}
	}
	return scene
}

func (scene *LevelScene) HandleKeyUp(key internal.ControlKey) internal.Thing {
	event := events.NoEvent
	if key == internal.Left {
		event = events.KeyLeftUp
	}
	if key == internal.Right {
		event = events.KeyRightUp
	}
	if event != events.NoEvent {
		scene.entities[0] = scene.entities[0].Handle(entities.EventBox{
			Event: event,
			Box:   px.Rect{},
		})
	}
	return scene
}

func (scene *LevelScene) Render(win *pixelgl.Window) {
	// TODO: if player cannot see past level limits,
	// this clear is not needed (camera like WonderBoy)
	win.Clear(colornames.Yellow50)

	gfx := d.OpSequence(
		d.Moved(
			scene.cameraVector(),
			d.OpSequence(
				d.ToWinOp(scene.backdropOp()),
				d.Color(colornames.White, d.TileLayer(scene.level.TilepixMap, "Background")),
				d.Color(colornames.White, d.TileLayer(scene.level.TilepixMap, "Platforms")),
				d.Color(colornames.White, d.TileLayer(scene.level.TilepixMap, "Walls")),
				d.Color(colornames.White, scene.signPostsOp()),
				scene.entityOp(),
				d.Color(colornames.White, d.TileLayer(scene.level.TilepixMap, "Foreground")),
				scene.debugGfx(),
			),
		),
		scene.mapSymbolOp(),
	)
	gfx.Render(px.IM, win.Canvas())

	// Spike
	//win.Clear(colornames.Blue200)
	//win.SetComposeMethod(px.ComposeOver)
	//imd := imdraw.New(nil)
	//imd.Color = colornames.Yellow600
	//imd.Push(px.V(0, 0))
	//imd.Push(px.V(100, 500))
	//imd.Rectangle(float64(0))
	//imd.Draw(win)
	//win.SetComposeMethod(px.ComposeXor)
	//gfx.Render(px.IM, win)
	//win.SetComposeMethod(px.ComposeOver)
}

func (scene *LevelScene) debugGfx() d.WinOp {
	if !scene.cfg.ShowEntityCanvas {
		return d.ToWinOp(d.Nothing())
	}
	hitBoxes := make([]d.ImdOp, 0)
	for _, entity := range scene.entities {
		hitBoxes = append(hitBoxes, rectDrawOp(entity.HitBox()))
	}
	hitBoxesOp := d.Colored(colornames.White, d.ImdOpSequence(hitBoxes...))

	eventBoxes := make([]d.ImdOp, 0)
	for _, eventbox := range scene.entityCanvas.EventBoxes {
		eventBoxes = append(eventBoxes, rectDrawOp(eventbox.Box))
	}
	eventBoxesOp := d.Colored(colornames.RedA700, d.ImdOpSequence(eventBoxes...))

	return d.OpSequence(d.ToWinOp(hitBoxesOp), d.ToWinOp(eventBoxesOp))
}

func rectDrawOp(r px.Rect) d.ImdOp {
	return d.Rectangle(C(r.Min), C(r.Max), 2)
}

func (scene *LevelScene) entityOp() d.WinOp {
	ops := make([]d.WinOp, 0)
	for _, entity := range scene.entities {
		ops = append(ops, entity.GfxOp(&scene.res.ImageMap))
	}
	return d.OpSequence(ops...)
}

func (scene *LevelScene) mapSymbolOp() d.WinOp {
	mapSymbolCenter := scene.res.ImageMap[internal.IMapSymbol].Frame().Center()
	op := d.Moved(mapSymbolCenter, d.Image(scene.res.ImageMap, internal.IMapSymbol))
	if scene.isMapSignClose() {
		op = d.Color(colornames.GreenA400, op)
	}
	return op
}

func (scene *LevelScene) backdropOp() d.ImdOp {
	widthPixels := scene.level.Width * 32
	heightPixels := scene.level.Height * 32
	rectangle := d.Rectangle(d.C(0, 0), d.C(float64(widthPixels), float64(heightPixels)), 0)
	return d.Colored(scene.level.ClearColor, rectangle)
}

func (scene *LevelScene) signPostsOp() d.WinOp {
	ops := make([]d.WinOp, 0)
	for _, mapPoint := range scene.level.SignPosts {
		alignVec := internal.V(0, scene.res.ImageMap[internal.ISignPost].Frame().Center().Y)
		signPostOp := d.Moved(mapPoint.Pos, d.Moved(alignVec,
			d.Image(scene.res.ImageMap, internal.ISignPost)))
		ops = append(ops, signPostOp)
	}
	return d.OpSequence(ops...)
}

/* gfxOp stop */

func (scene *LevelScene) drawFPS(win *pixelgl.Window) {
	win.SetMatrix(px.IM)
	tb := text.New(px.V(0, 0), scene.res.Atlas)
	_, _ = fmt.Fprintf(tb, "FPS=%1.1f", scene.res.FPS)
	tb.DrawColorMask(win, px.IM, colornames.Brown800)
}

func (scene *LevelScene) isMapSignClose() bool {
	sign := scene.closestMapSign()
	return scene.elisePos().Sub(sign.Pos).Len() < 75
}

func (scene *LevelScene) elisePos() px.Vec {
	return scene.entities[0].HitBox().Center()
}

func (scene *LevelScene) closestMapSign() internal.SignPost {
	// Potential: ClosestPoint could take an array of objects implementing
	// 'WithPoint' interface, and we only define anon func here
	getPoint := func(mp internal.SignPost) px.Vec { return mp.Pos }
	points := make([]px.Vec, 0)
	for _, val := range scene.level.SignPosts {
		point := getPoint(val)
		points = append(points, point)
	}
	return scene.level.SignPosts[internal.ClosestPoint(scene.elisePos(), points)]
}

func (scene *LevelScene) cameraVector() px.Vec {
	halfScreen := internal.V(internal.ScreenWidth/2, internal.ScreenHeight/2)
	playerHead := internal.V(0, internal.PlayerHeight)
	cam := scene.elisePos().Sub(halfScreen).Add(playerHead)
	reversed := cam.Scaled(-1)
	return reversed
}

func (scene *LevelScene) Tick() bool {
	// Handle event boxes from previous tick first of all
	scene.entityCanvas.Consequences(func(eb entities.EventBox, box entities.EntityHitBox) {
		id := box.Entity
		scene.entities[id] = scene.entities[id].Handle(eb)
	})

	// Play relevant sound effects (based on events)
	for _, eb := range scene.entityCanvas.EventBoxes {
		if eb.Event == events.ButtonPressed {
			streamer := scene.buttonClick.Streamer(0, scene.buttonClick.Len())
			speaker.Play(streamer)
		}
		if eb.Event == events.RobotMove {
			streamer := scene.robotMoveSfx.Streamer(0, scene.robotMoveSfx.Len())
			speaker.Play(streamer)
		}
	}

	// Reset the canvas
	scene.entityCanvas = entities.MakeEntityCanvas()
	scene.gameTimeMs += internal.TickTimeMs
	for i := range scene.entities {
		scene.entities[i] = scene.entities[i].Tick(scene.gameTimeMs, &scene.entityCanvas)
		scene.entityCanvas.AddEntityHitBox(entities.EntityHitBox{
			Entity: i,
			HitBox: scene.entities[i].HitBox(),
		})
	}

	// Add all walls
	// TODO: Refactor this to data structure popuplated
	// when loading level (and cover with approval test!)
	walls := scene.level.TilepixMap.GetTileLayerByName("Walls")
	tiles := walls.DecodedTiles
	for y := 0; y < scene.level.TilepixMap.Height; y++ {
		for x := 0; x < scene.level.TilepixMap.Width; x++ {
			ix := y*scene.level.Width + x
			tile := tiles[ix]
			if tile.IsNil() {
				continue
			}
			pos := tile.Position(ix, walls.Tileset)
			scene.entityCanvas.AddEventBox(entities.EventBox{
				Event: events.Wall,
				Box: px.R(
					pos.X-16,
					pos.Y-16,
					pos.X+16,
					pos.Y+16,
				),
			})
		}
	}

	return true
}

/*
Windows rit-API:er som används i Rescue:
image.Draw(win, mx)
image.DrawColorMask(win, mx, color)
text.DrawColorMask(win, mx, color)
text.Draw(win, px.IM)  # används inte men finns i API:et!
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

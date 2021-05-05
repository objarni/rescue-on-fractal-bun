package scenes

import (
	"fmt"
	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
	d "objarni/rescue-on-fractal-bun/internal/draw"
)

/*
SignPosts innehåller följande information:
MapPos vec   - position på karten för MapSignen
LevelPos vec - position på leveln den är kopplad till
Level string - vilken level den är kopplad till

MapPos lagras i MapScene (som har kunskap om detta).
De två senare lagras i Tiled maps (SignPost objects).

En idé är att konstruera alla SignPosts vid start av
spelet, för att på så sätt
 - validera att alla levels laddas korrekt
   (level loading kontrollerar att layers finns och
   är ordnade i rätt ordning)
 - spara all SignPosts info

Vill börja med att validera laddning av level enligt
ovan, för att sedan bygga upp SignPosts datastrukturen.
*/

type MapScene struct {
	cfg            *Config
	res            *internal.Resources
	hairCrossPos   px.Vec
	hairCrossVel   px.Vec
	highlightTimer int
	atMapSign      int
}

func MakeMapScene(cfg *Config, res *internal.Resources, mapSignName string) *MapScene {
	mapSignIx := lookupMapSignWithText(mapSignName, res.MapSigns)
	return &MapScene{
		cfg:          cfg,
		res:          res,
		hairCrossPos: res.MapSigns[mapSignIx].MapPos,
		hairCrossVel: px.ZV,
		atMapSign:    mapSignIx,
	}
}

func lookupMapSignWithText(mapSignText string, mapSigns []internal.MapSign) int {
	for ix, mapSign := range mapSigns {
		if mapSign.Text == mapSignText {
			return ix
		}
	}
	panic("error: could not find MapSign with Text=" + mapSignText)
}

func (scene *MapScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		mapSignIx := scene.FindClosestMapSign()
		if mapSignIx != -1 {
			mapSign := scene.res.MapSigns[mapSignIx]
			return MakeLevelScene(scene.cfg,
				scene.res,
				mapSign.LevelName,
				mapSign.LevelPos)
		}
	}
	if key == internal.Left {
		scene.hairCrossVel.X -= 1
	}
	if key == internal.Right {
		scene.hairCrossVel.X += 1
	}
	if key == internal.Down {
		scene.hairCrossVel.Y -= 1
	}
	if key == internal.Up {
		scene.hairCrossVel.Y += 1
	}
	return scene
}

func (scene *MapScene) HandleKeyUp(key internal.ControlKey) internal.Thing {
	if key == internal.Left {
		scene.hairCrossVel.X += 1
	}
	if key == internal.Right {
		scene.hairCrossVel.X -= 1
	}
	if key == internal.Down {
		scene.hairCrossVel.Y += 1
	}
	if key == internal.Up {
		scene.hairCrossVel.Y -= 1
	}
	return scene
}

func (scene *MapScene) Render(win *pixelgl.Window) {
	sceneGfxOp := scene.MapSceneWinOp()
	context := d.Context{Transform: px.IM}
	sceneGfxOp.DrawTo(win.Canvas(), context)
	drawMapSignTexts(win, scene)
}

func (scene *MapScene) MapSceneWinOp() d.WinOp {
	return d.OpSequence(
		d.Moved(px.R(0, 0, internal.ScreenWidth, internal.ScreenHeight).Center(),
			d.Image(scene.res.ImageMap, internal.IMap)),
		d.ToWinOp(d.ImdOpSequence(scene.mapSignsGfx(), scene.crossHairGfx())),
	)
}

func drawMapSignTexts(win *pixelgl.Window, scene *MapScene) {
	mapSignIx := FindNearMapSign(scene.hairCrossPos, scene.res.MapSigns, scene.cfg.MapSceneTargetLocMaxDistance)
	mapSignName := mapSignNameFromIx(mapSignIx)
	tb := text.New(px.ZV, scene.res.Atlas)
	d.Text(
		fmt.Sprintf("Här är du: %s\n", "Hembyn"),
		fmt.Sprintf("Gå till? %s", mapSignName),
	).Render(tb)
	textPosition := px.V(
		float64(scene.cfg.MapSceneLocationTextX),
		float64(scene.cfg.MapSceneLocationTextY))
	// Observation: tb ankras i första radens nedre vänstra hörn.
	// Med två rader blir det alltså riktigt rörigt att positionera!
	// För att kunna beräkna positionering, behövs:
	// 1) Resultatparagrafens totala bredd
	// 2) En textrads höjd
	// Finns det något sätt som inte inbegriper "beräkna textdim i förväg"?
	//win.SetComposeMethod(px.ComposeRin)
	tb.DrawColorMask(win, px.IM.Moved(textPosition), colornames.Black)
}

func mapSignNameFromIx(mapSignIx int) string {
	if mapSignIx == 0 {
		return "Hembyn"
	}
	if mapSignIx == 1 {
		return "Korsningen"
	}
	if mapSignIx == 2 {
		return "Skogen"
	}
	return "-"
}

func (scene *MapScene) mapSignsGfx() d.ImdOp {
	return scene.levelEntrances().
		Then(scene.currentMapSignOp()).
		Then(scene.crossHairsOp())
}

func (scene *MapScene) crossHairsOp() d.ImdOp {
	closestMapSignIx := scene.FindClosestMapSign()
	if closestMapSignIx > -1 {
		pos := scene.res.MapSigns[closestMapSignIx].MapPos
		radius := scene.targetLocCircleRadius()
		circle := d.Circle(radius, C(pos), scene.circleThickness())
		operation := d.Colored(colornames.Red, circle)
		return operation
	}
	return d.Nothing()
}

func C(v px.Vec) px.Vec {
	return d.C(v.X, v.Y)
}

func (scene *MapScene) FindClosestMapSign() int {
	return FindNearMapSign(scene.hairCrossPos, scene.res.MapSigns, scene.locMaxDistance())
}

func (scene *MapScene) currentMapSignOp() d.ImdOp {
	blink := scene.cfg.MapSceneBlinkSpeed
	if scene.highlightTimer/blink%2 == 0 {
		pos := scene.res.MapSigns[scene.atMapSign].MapPos
		circle := d.Circle(scene.currentLocCircleRadius(), C(pos), scene.circleThickness())
		return d.Colored(colornames.Green, circle)
	}
	return d.Nothing()
}

func (scene *MapScene) levelEntrances() d.ImdSequence {
	sequence := d.ImdOpSequence()
	for _, mapSign := range scene.res.MapSigns {
		pos := mapSign.MapPos
		operation := d.Colored(
			colornames.Darkslateblue,
			d.Circle(scene.locCircleRadius(), C(pos), scene.circleThickness()))
		sequence = sequence.Then(operation)
	}
	return sequence
}

func (scene *MapScene) circleThickness() int {
	return 3
}

func (scene *MapScene) targetLocCircleRadius() int {
	return scene.cfg.MapSceneTargetLocCircleRadius
}

func (scene *MapScene) locMaxDistance() int {
	return scene.cfg.MapSceneTargetLocMaxDistance
}

func (scene *MapScene) currentLocCircleRadius() int {
	return scene.cfg.MapSceneCurrentLocCircleRadius
}

func (scene *MapScene) locCircleRadius() int {
	return scene.cfg.MapSceneLocCircleRadius
}

func (scene *MapScene) crossHairGfx() d.ImdOp {
	h := scene.hairCrossPos
	thickness := 2
	vertical := d.Line(d.C(h.X, 0), d.C(h.X, 600), thickness)
	horisontal := d.Line(d.C(0, h.Y), d.C(800, h.Y), thickness)
	transparentPink := color.RGBA{R: 255, A: 32}
	return d.Colored(transparentPink, d.ImdOpSequence(vertical, horisontal))
}

func (scene *MapScene) Tick() bool {
	scene.hairCrossPos = scene.hairCrossPos.Add(scene.hairCrossVel.Scaled(scene.cfg.MapSceneCrosshairSpeed))
	if scene.hairCrossPos.X < 0 {
		scene.hairCrossPos.X = 0
	}
	if scene.hairCrossPos.X > 799 {
		scene.hairCrossPos.X = 799
	}
	if scene.hairCrossPos.Y < 0 {
		scene.hairCrossPos.Y = 0
	}
	if scene.hairCrossPos.Y > 599 {
		scene.hairCrossPos.Y = 599
	}
	scene.highlightTimer += 1
	return true
}

func FindNearMapSign(vec px.Vec, mapSigns []internal.MapSign, maxDist int) int {
	points := make([]px.Vec, 0)
	getPoint := func(mp internal.MapSign) px.Vec { return mp.MapPos }
	// Potential: ClosestPoint could take an array of objects implementing
	// 'WithPoint' interface, and we only define anon func here
	for _, val := range mapSigns {
		point := getPoint(val)
		points = append(points, point)
	}
	closest := internal.ClosestPoint(vec, points)
	if internal.Distance(vec, points[closest]) > float64(maxDist) {
		return -1
	}
	return closest
}

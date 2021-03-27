package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
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
	hairCrossPos   pixel.Vec
	hairCrossVel   pixel.Vec
	playerLocIx    int
	highlightTimer int
}

func MakeMapScene(cfg *Config, res *internal.Resources, mapSignName string) *MapScene {
	mapSignIx := mapSignIndexFromName(mapSignName)
	return &MapScene{
		cfg:          cfg,
		res:          res,
		hairCrossPos: res.MapSigns[mapSignIx].MapPos,
		hairCrossVel: pixel.ZV,
		playerLocIx:  mapSignIx,
	}
}

func (scene *MapScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		mapSignIx := scene.FindClosestMapSign()
		if mapSignIx != -1 {
			levelName := scene.res.MapSigns[mapSignIx].LevelName
			// TODO: spawn at right mapsign of level, not just 'mapsign 0'
			return MakeLevelScene(scene.cfg, scene.res, levelName)
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
	sceneGfxOp.Render(pixel.IM, win)
	drawMapSignTexts(win, scene)
}

func (scene *MapScene) MapSceneWinOp() draw.WinOp {
	lineOps := draw.ToWinOp(draw.ImdOpSequence(scene.mapSignsGfx(), scene.crossHairGfx()))
	mapOp := draw.Moved(pixel.Rect{Min: v(0, 0), Max: v(internal.ScreenWidth, internal.ScreenHeight)}.Center(),
		draw.Image(scene.res.ImageMap, internal.IMap))
	sceneGfxOp := draw.OpSequence(mapOp, lineOps)
	return sceneGfxOp
}

func drawMapSignTexts(win *pixelgl.Window, scene *MapScene) {
	mapSignIx := FindNearMapSign(scene.hairCrossPos, scene.res.MapSigns, scene.cfg.MapSceneTargetLocMaxDistance)
	mapSignName := mapSignNameFromIx(mapSignIx)
	tb := text.New(pixel.ZV, scene.res.Atlas)
	draw.Text(
		fmt.Sprintf("Här är du: %s\n", "Hembyn"),
		fmt.Sprintf("Gå till? %s", mapSignName),
	).Render(tb)
	textPosition := pixel.V(
		float64(scene.cfg.MapSceneLocationTextX),
		float64(scene.cfg.MapSceneLocationTextY))
	// Observation: tb ankras i första radens nedre vänstra hörn.
	// Med två rader blir det alltså riktigt rörigt att positionera!
	// För att kunna beräkna positionering, behövs:
	// 1) Resultatparagrafens totala bredd
	// 2) En textrads höjd
	// Finns det något sätt som inte inbegriper "beräkna textdim i förväg"?
	//win.SetComposeMethod(pixel.ComposeRin)
	tb.DrawColorMask(win, pixel.IM.Moved(textPosition), colornames.Black)
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

func mapSignIndexFromName(locationName string) int {
	if locationName == "Hembyn" {
		return 0
	}
	if locationName == "Korsningen" {
		return 1
	}
	if locationName == "Skogen" {
		return 2
	}
	panic(fmt.Sprintf("Unknown location name: %v", locationName))
}

func (scene *MapScene) mapSignsGfx() draw.ImdOp {
	return scene.levelEntrances().
		Then(scene.currentLocation()).
		Then(scene.crossHairLocation())
}

func (scene *MapScene) crossHairLocation() draw.ImdOp {
	closestMapSignIx := scene.FindClosestMapSign()
	if closestMapSignIx > -1 {
		pos := scene.res.MapSigns[closestMapSignIx].MapPos
		radius := scene.targetLocCircleRadius()
		circle := draw.Circle(radius, C(pos), scene.circleThickness())
		operation := draw.Colored(colornames.Red, circle)
		return operation
	}
	return draw.Nothing()
}

func C(v pixel.Vec) draw.Coordinate {
	return draw.C(int(v.X), int(v.Y))
}

func (scene *MapScene) FindClosestMapSign() int {
	return FindNearMapSign(scene.hairCrossPos, scene.res.MapSigns, scene.locMaxDistance())
}

func (scene *MapScene) currentLocation() draw.ImdOp {
	blink := scene.cfg.MapSceneBlinkSpeed
	if scene.highlightTimer/blink%2 == 0 {
		// TODO: remove playerCurrIx in favor of 'currentMapSign' or something
		pos := scene.res.MapSigns[0].MapPos
		circle := draw.Circle(scene.currentLocCircleRadius(), C(pos), scene.circleThickness())
		return draw.Colored(colornames.Green, circle)
	}
	return draw.Nothing()
}

func (scene *MapScene) levelEntrances() draw.ImdSequence {
	sequence := draw.ImdOpSequence()
	for _, mapSign := range scene.res.MapSigns {
		pos := mapSign.MapPos
		operation := draw.Colored(
			colornames.Darkslateblue,
			draw.Circle(scene.locCircleRadius(), C(pos), scene.circleThickness()))
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

func (scene *MapScene) crossHairGfx() draw.ImdOp {
	h := scene.hairCrossPos
	thickness := 2
	vertical := draw.Line(draw.C(int(h.X), 0), draw.C(int(h.X), 600), thickness)
	horisontal := draw.Line(draw.C(0, int(h.Y)), draw.C(800, int(h.Y)), thickness)
	transparentPink := color.RGBA{R: 255, A: 32}
	return draw.Colored(transparentPink, draw.ImdOpSequence(vertical, horisontal))
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

func v(x float64, y float64) pixel.Vec {
	return pixel.Vec{X: x, Y: y}
}

func FindNearMapSign(vec pixel.Vec, mapSigns []internal.MapSign, maxDist int) int {
	points := []pixel.Vec{}
	getPoint := func(mp internal.MapSign) pixel.Vec { return mp.MapPos }
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

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

type MapPoint struct {
	position   pixel.Vec
	discovered bool
}

type MapScene struct {
	cfg            *Config
	res            *internal.Resources
	mapImage       *pixel.Sprite
	mapPoints      []MapPoint
	hairCrossPos   pixel.Vec
	hairCrossVel   pixel.Vec
	playerLocIx    int
	highlightTimer int
	mapSigns       []internal.MapSign // All map signs in game
}

func MakeMapScene(cfg *Config, res *internal.Resources, locationName string) *MapScene {
	locations := []MapPoint{
		{
			position:   pixel.Vec{X: 246, Y: 109},
			discovered: true,
		},
		{
			position:   pixel.Vec{X: 355, Y: 235},
			discovered: true,
		},
		{
			position:   pixel.Vec{X: 299, Y: 375},
			discovered: false,
		},
	}
	locationIx := locationIxFromName(locationName)
	return &MapScene{
		cfg:          cfg,
		res:          res,
		mapImage:     internal.LoadSpriteForSure("assets/TMap.png"),
		hairCrossPos: locations[locationIx].position,
		hairCrossVel: pixel.ZV,
		playerLocIx:  locationIx,
		mapPoints:    locations,
	}
}

func (scene *MapScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		locationIx := scene.FindClosestLocation()
		if locationIx != -1 {
			levelName := "GhostForest" // TODO: initialize mapSigns on game boot (it's nil!)
			if scene.mapSigns != nil {
				levelName = scene.mapSigns[locationIx].LevelName
			}
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
	lineOps := draw.ToWinOp(draw.Sequence(scene.locationsGfx(), scene.crossHairGfx()))
	mapOp := draw.Moved(win.Bounds().Center(),
		draw.Image(scene.res.ImageMap, internal.IMap))
	sceneGfxOp := draw.OpSequence(mapOp, lineOps)
	sceneGfxOp.Render(pixel.IM, win)

	drawLocationTexts(win, scene)
}

func drawLocationTexts(win *pixelgl.Window, scene *MapScene) {
	locationIx := FindNearLocation(scene.hairCrossPos, scene.mapPoints, scene.cfg.MapSceneTargetLocMaxDistance)
	locationName := locationNameFromIx(locationIx)
	tb := text.New(pixel.ZV, scene.res.Atlas)
	draw.Text(
		fmt.Sprintf("Här är du: %s\n", "Hembyn"),
		fmt.Sprintf("Gå till? %s", locationName),
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
	tb.DrawColorMask(win, pixel.IM.Moved(textPosition), colornames.Black)
}

func locationNameFromIx(locationIx int) string {
	if locationIx == 0 {
		return "Hembyn"
	}
	if locationIx == 1 {
		return "Korsningen"
	}
	if locationIx == 2 {
		return "Skogen"
	}
	return "-"
}

func locationIxFromName(locationName string) int {
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

func (scene *MapScene) locationsGfx() draw.ImdOp {
	return scene.levelEntrances().
		Then(scene.currentLocation()).
		Then(scene.crossHairLocation())
}

func (scene *MapScene) crossHairLocation() draw.ImdOp {
	closestLocation := scene.FindClosestLocation()
	if closestLocation > -1 {
		pos := scene.mapPoints[closestLocation].position
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

func (scene *MapScene) FindClosestLocation() int {
	return FindNearLocation(scene.hairCrossPos, scene.mapPoints, scene.locMaxDistance())
}

func (scene *MapScene) currentLocation() draw.ImdOp {
	blink := scene.cfg.MapSceneBlinkSpeed
	if scene.highlightTimer/blink%2 == 0 {
		loc := scene.mapPoints[scene.playerLocIx]
		pos := loc.position
		circle := draw.Circle(scene.currentLocCircleRadius(), C(pos), scene.circleThickness())
		return draw.Colored(colornames.Green, circle)
	}
	return draw.Nothing()
}

func (scene *MapScene) levelEntrances() draw.ImdSequence {
	sequence := draw.Sequence()
	for _, loc := range scene.mapPoints {
		pos := loc.position
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
	return draw.Colored(transparentPink, draw.Sequence(vertical, horisontal))
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

func FindNearLocation(vec pixel.Vec, locations []MapPoint, maxDist int) int {
	points := []pixel.Vec{}
	getPoint := func(mp MapPoint) pixel.Vec { return mp.position }
	// Potential: ClosestPoint could take an array of objects implementing
	// 'WithPoint' interface, and we only define anon func here
	for _, val := range locations {
		point := getPoint(val)
		points = append(points, point)
	}
	closest := internal.ClosestPoint(vec, points)
	if internal.Distance(vec, points[closest]) > float64(maxDist) {
		return -1
	}
	return closest
}

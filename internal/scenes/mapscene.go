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
På kartan kan man besöka "mapPoints".
Man befinner sig alltid på precis en location (current).
Current visas tydligt visuellt i kartscenen, kanske
står dess namn någonstans också?
När man går in på en location, hamnar man på en "kartpunkt"
i en bana. Man kan alltid gå tillbaka till kartan
via en kartpunkt. Om man på en bana hittar en ny kartpunkt,
d.v.s en kartpunkt som inte besökts förut, dyker den upp
på kartan (kanske ett ljud spelas upp en text visas för
att demonstrera att en ny kartpunkt hittats).
Kraftindikatorn är en liten tjej som börjar glad, men
blir surare och surare. Till slut, när det blir jättearg
min, flyttas man tillbaka till senaste kartpunkten på
banan, alltså den kartpunkt som är associerad med current
location.

En location har alltså följande egenskaper:
  - position på kartan, x,y
  - om den är "Discovered" eller ej; en bool.
    (hidden eller visible kan man se det som)
Förutom detta sparas en pekare till current location
"någonstans".

*/

type MapPoint struct {
	position   pixel.Vec
	discovered bool
}

type MapScene struct {
	cfg            *Config
	res            *Resources
	mapImage       *pixel.Sprite
	mapPoints      []MapPoint
	hairCrossPos   pixel.Vec
	hairCrossVel   pixel.Vec
	playerLocIx    int
	highlightTimer int
}

func MakeMapScene(cfg *Config, res *Resources, locationName string) *MapScene {
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
		// TODO: load level scene of target location, with location
		// as starting point (map point)
		// (if there is a location close enough to crossHair, that is)
		return MakeLevelScene(scene.cfg, scene.res)
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
	scene.mapImage.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	seq := draw.Sequence(scene.locationsGfx(), scene.crossHairGfx())
	draw.ToWinOp(seq).Render(win)

	drawLocationTexts(win, scene)
}

func drawLocationTexts(win *pixelgl.Window, scene *MapScene) {
	locationIx := FindClosestLocation(scene.hairCrossPos, scene.mapPoints, scene.cfg.MapSceneTargetLocMaxDistance)
	locationName := locationNameFromIx(locationIx)
	tb := text.New(pixel.ZV, scene.res.Atlas)
	draw.Text(
		fmt.Sprintf("Här är du: %s\n", "Hembyn"),
		fmt.Sprintf("Gå till? %s", locationName),
	).Render(tb)
	//_, _ = fmt.Fprintf(tb, "Här är du: %s\n", "Hembyn")
	//_, _ = fmt.Fprintf(tb, "Gå till? %s", locationName)
	tb.DrawColorMask(win, pixel.IM.Moved(pixel.V(
		float64(scene.cfg.MapSceneLocationTextX),
		float64(scene.cfg.MapSceneLocationTextY))), colornames.Black)
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
	return FindClosestLocation(scene.hairCrossPos, scene.mapPoints, scene.locMaxDistance())
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

func FindClosestLocation(vec pixel.Vec, locations []MapPoint, maxDist int) int {
	closest := -1
	closestDist := -1.0
	for ix, val := range locations {
		d := distance(vec, val.position)
		if closest == -1 || d < closestDist {
			closest = ix
			closestDist = d
		}
	}
	if closestDist > float64(maxDist) {
		return -1
	}
	return closest
}

func distance(a pixel.Vec, b pixel.Vec) float64 {
	return a.Sub(b).Len()
}

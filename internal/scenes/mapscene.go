package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
)

/*
På kartan kan man besöka "locations".
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

type Location struct {
	position   pixel.Vec
	discovered bool
}

type MapScene struct {
	cfg            *Config
	res            *Resources
	mapImage       *pixel.Sprite
	locations      []Location
	hairCrossPos   pixel.Vec
	hairCrossVel   pixel.Vec
	playerLocIx    int
	highlightTimer int
}

func MakeMapScene(cfg *Config, res *Resources, locationName string) *MapScene {
	locations := []Location{
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
		locations:    locations,
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
	imd := imdraw.New(nil)
	drawLocations(scene, imd)
	drawCrossHair(win, scene, imd)
	drawLocationTexts(win, scene)
}

func drawLocationTexts(win *pixelgl.Window, scene *MapScene) {
	locationIx := FindClosestLocation(scene.hairCrossPos, scene.locations, scene.cfg.MapSceneTargetLocMaxDistance)
	locationName := locationNameFromIx(locationIx)
	tb := text.New(pixel.V(
		float64(scene.cfg.MapSceneLocationTextX),
		float64(scene.cfg.MapSceneLocationTextY)),
		scene.res.Atlas)
	_, _ = fmt.Fprintf(tb, "Här är du: %s\n", "Hembyn")
	_, _ = fmt.Fprintf(tb, "Gå till? %s", locationName)
	tb.DrawColorMask(win, pixel.IM, colornames.Black)
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

func drawLocations(scene *MapScene, imd *imdraw.IMDraw) *imdraw.IMDraw {
	for _, loc := range scene.locations {
		vec := loc.position
		drawCircle(imd, colornames.Darkslateblue, vec,
			scene.cfg.MapSceneLocCircleRadius)
	}
	blink := scene.cfg.MapSceneBlinkSpeed
	if scene.highlightTimer/blink%2 == 0 {
		drawCircle(
			imd, colornames.Green,
			scene.locations[scene.playerLocIx].position,
			scene.cfg.MapSceneCurrentLocCircleRadius,
		)
	}
	ix := FindClosestLocation(scene.hairCrossPos, scene.locations, scene.cfg.MapSceneTargetLocMaxDistance)
	if ix > -1 {
		drawCircle(
			imd, colornames.Red,
			scene.locations[ix].position,
			scene.cfg.MapSceneTargetLocCircleRadius,
		)
	}

	return imd
}

func drawCrossHair(win *pixelgl.Window, scene *MapScene, imd *imdraw.IMDraw) {
	imd.Color = pixel.RGBA{R: 1, A: 0.15}
	h := scene.hairCrossPos
	imd.Push(v(h.X, 0))
	imd.Push(v(h.X, 600))
	imd.Rectangle(2)
	imd.Push(v(0, h.Y))
	imd.Push(v(800, h.Y))
	imd.Rectangle(2)
	imd.Draw(win)
}

func drawCircle(
	imd *imdraw.IMDraw,
	color color.RGBA,
	vec pixel.Vec,
	radius int,
) {
	imd.Color = color
	imd.Push(vec)
	imd.Circle(float64(radius), 3)
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

func FindClosestLocation(vec pixel.Vec, locations []Location, maxDist int) int {
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

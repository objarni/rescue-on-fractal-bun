package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image/color"
	"io/ioutil"
	"objarni/rescue-on-fractal-bun/internal"
	"os"
	"unicode"
)

/*
På kartan kan man besöka "locations".
Man befinner sig alltid på precis en location (current).
Current visas tydligt visuellt i kartscenen, kanske
står dess namn någonstans också?
När man går in på en location, hamnar man på en "kartpunkt"
i en bana. Man kan alltid exita tillbaka till kartan
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
  - om den är "discovered" eller ej; en bool.
    (hidden eller visible kan man se det som)
Förutom detta sparas en pekare till current location
"någonstans".

*/

type Location struct {
	position   pixel.Vec
	discovered bool
}

type MapScene struct {
	cfg             *Config
	mapImage        *pixel.Sprite
	atlas           *text.Atlas
	locations       []Location
	heroPos         pixel.Vec
	heroVel         pixel.Vec
	highlight       int
	hightlightTimer int
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func MakeMapScene(cfg *Config) *MapScene {
	// @remember: atlas sent in to render, shared by all scenes?
	face, err := loadTTF("assets/Font.ttf", 32)
	internal.PanicIfError(err)
	//face := basicfont.Face7x13
	atlas := text.NewAtlas(face, text.RangeTable(unicode.Latin), text.ASCII)
	locs := []Location{
		{
			position:   pixel.Vec{X: 246, Y: 486},
			discovered: true,
		},
		{
			position:   pixel.Vec{X: 355, Y: 368},
			discovered: true,
		},
		{
			position:   pixel.Vec{X: 500, Y: 500},
			discovered: false,
		},
	}
	return &MapScene{
		cfg:       cfg,
		atlas:     atlas,
		mapImage:  internal.LoadSpriteForSure("assets/TMap.png"),
		heroPos:   pixel.Vec{50, 50},
		heroVel:   pixel.ZV,
		highlight: 0,
		locations: locs,
	}
}

func (scene *MapScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		return MakeMenuScene(scene.cfg)
	}
	if key == internal.Left {
		scene.heroVel.X -= 1
	}
	if key == internal.Right {
		scene.heroVel.X += 1
	}
	if key == internal.Down {
		scene.heroVel.Y -= 1
	}
	if key == internal.Up {
		scene.heroVel.Y += 1
	}
	return scene
}

func (scene *MapScene) HandleKeyUp(key internal.ControlKey) internal.Thing {
	if key == internal.Left {
		scene.heroVel.X += 1
	}
	if key == internal.Right {
		scene.heroVel.X -= 1
	}
	if key == internal.Down {
		scene.heroVel.Y += 1
	}
	if key == internal.Up {
		scene.heroVel.Y -= 1
	}
	return scene
}

func (scene *MapScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Antiquewhite)

	// Background
	scene.mapImage.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	// Level locations
	var imd = imdraw.New(nil)
	for _, loc := range scene.locations {
		vec := loc.position
		drawCircle(imd, colornames.Darkslateblue, vec,
			scene.cfg.MapSceneLocCircleRadius)
	}
	blink := scene.cfg.MapSceneBlinkSpeed
	if scene.highlight != -1 && scene.hightlightTimer/blink%2 == 0 {
		drawCircle(
			imd, colornames.Green,
			scene.locations[scene.highlight].position,
			scene.cfg.MapSceneCurrentLocCircleRadius,
		)
	}

	// Hero position
	drawCrosshair(win, imd, scene)

	// Text
	tb := text.New(pixel.V(
		float64(scene.cfg.MapSceneLocationTextX),
		float64(scene.cfg.MapSceneLocationTextY)),
		scene.atlas)
	_, _ = fmt.Fprintf(tb, "Här är du: %s\n", "Hembyn")
	_, _ = fmt.Fprintf(tb, "Gå till? %s", "Korsningen")
	tb.DrawColorMask(win, pixel.IM.Scaled(tb.Orig, 1), colornames.Black)
}

func drawCrosshair(win *pixelgl.Window, imd *imdraw.IMDraw, scene *MapScene) {
	imd.Color = pixel.RGBA{1, 0, 0, 0.15}
	h := scene.heroPos
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
	imd.Push(v(0, 600).Add(vec.ScaledXY(v(1, -1))))
	imd.Circle(float64(radius), 3)
}

func (scene *MapScene) Tick() bool {
	scene.heroPos = scene.heroPos.Add(scene.heroVel.Scaled(scene.cfg.MapSceneCrosshairSpeed))
	if scene.heroPos.X < 0 {
		scene.heroPos.X = 0
	}
	if scene.heroPos.X > 799 {
		scene.heroPos.X = 799
	}
	if scene.heroPos.Y < 0 {
		scene.heroPos.Y = 0
	}
	if scene.heroPos.Y > 599 {
		scene.heroPos.Y = 599
	}
	scene.hightlightTimer += 1
	return true
}

func v(x float64, y float64) pixel.Vec {
	return pixel.Vec{X: x, Y: y}
}

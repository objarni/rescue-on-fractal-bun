package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"objarni/rescue-on-fractal-bun/internal"
)

type MapScene struct {
	textbox     *text.Text
	mapImage    *pixel.Sprite
	levelCoords []pixel.Vec
	heroPos     pixel.Vec
	heroVel     pixel.Vec
}

func MakeMapScene() *MapScene {
	// @thought: atlas sent in to render, shared by all scenes?
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	levelCoords := []pixel.Vec{
		{X: 246, Y: 486},
		{X: 355, Y: 368},
	}
	return &MapScene{
		textbox:     text.New(pixel.V(0, 0), atlas),
		mapImage:    internal.LoadSpriteForSure("assets/TMap.png"),
		levelCoords: levelCoords,
		heroPos:     pixel.Vec{50, 50},
		heroVel:     pixel.ZV,
	}
}

func (scene *MapScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		return MakeMenuScene()
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
	imd.Color = colornames.Darkslateblue
	for _, vec := range scene.levelCoords {
		imd.Push(v(0, 600).Add(vec.ScaledXY(v(1, -1))))
		imd.Circle(12, 3)
	}

	// Hero position
	imd.Color = pixel.RGBA{1, 0, 0, 0.15}
	h := scene.heroPos
	imd.Push(v(h.X, 0))
	imd.Push(v(h.X, 600))
	imd.Rectangle(2)
	imd.Push(v(0, h.Y))
	imd.Push(v(800, h.Y))
	imd.Rectangle(2)
	imd.Draw(win)

	// Text
	tb := scene.textbox
	tb.Clear()
	_, _ = fmt.Fprintln(tb, "Map scene")
	tb.DrawColorMask(win, pixel.IM.Scaled(tb.Orig, 2), colornames.Black)
}

func (scene *MapScene) Tick() bool {
	scene.heroPos = scene.heroPos.Add(scene.heroVel.Scaled(1))
	return true
}

func v(x float64, y float64) pixel.Vec {
	return pixel.Vec{X: x, Y: y}
}

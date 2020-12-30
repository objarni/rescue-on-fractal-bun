package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"objarni/rescue-on-fractal-bun/internal"
)

type MapScene struct {
	textbox  *text.Text
	mapImage *pixel.Sprite
}

func MakeMapScene() MapScene {
	// @thought: atlas sent in to render, shared by all scenes?
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	return MapScene{
		textbox:  text.New(pixel.V(0, 0), atlas),
		mapImage: internal.LoadSpriteForSure("assets/TMap.png"),
	}
}

func (scene MapScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		return MakeMenuScene()
	}
	return scene
}

func (scene MapScene) HandleKeyUp(_ internal.ControlKey) internal.Thing {
	return scene
}

func (scene MapScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Antiquewhite)
	scene.mapImage.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
	tb := scene.textbox
	tb.Clear()
	_, _ = fmt.Fprintln(tb, "Map scene")
	tb.DrawColorMask(win, pixel.IM.Scaled(tb.Orig, 2), colornames.Black)
}

func (scene MapScene) Tick() bool {
	return true
}

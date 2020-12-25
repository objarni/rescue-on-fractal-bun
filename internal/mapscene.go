package internal

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type MapScene struct {
	textbox *text.Text
}

func MakeMapScene() MapScene {
	// @thought: atlas sent in to render, shared by all scenes?
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	return MapScene{
		textbox: text.New(pixel.V(0, 0), atlas),
	}
}

func (scene MapScene) HandleKeyDown(key ControlKey) Scene {
	if key == Jump {
		return MakeMenuScene()
	}
	return scene
}

func (scene MapScene) HandleKeyUp(_ ControlKey) Scene {
	return scene
}

func (scene MapScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Antiquewhite)
	tb := scene.textbox
	tb.Clear()
	_, _ = fmt.Fprintln(tb, "Map scene")
	tb.DrawColorMask(win, pixel.IM.Scaled(tb.Orig, 2), colornames.Black)
}

// Vill ha: drawText(x, y, txt, color)
// .. där x,y är angivet i procent, och betyder centerpunkt för texten

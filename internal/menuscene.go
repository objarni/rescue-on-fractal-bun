package internal

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type MenuItem int

const (
	Play MenuItem = iota
	Quit
)

type MenuScene struct {
	currentitem MenuItem
	textbox     *text.Text
}

func MakeMenuScene() MenuScene {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	return MenuScene{
		currentitem: Play,
		textbox:     text.New(pixel.V(0, 0), atlas),
	}
}

func (menuScene MenuScene) HandleKeyDown(key ControlKey) Scene {
	if key == Action || key == Jump {
		return MakeMapScene()
	}
	return menuScene
}

func (menuScene MenuScene) HandleKeyUp(key ControlKey) Scene {
	fmt.Println("menu key up: " + key.String())
	return menuScene
}

func (menuScene MenuScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Aliceblue)
	tb := menuScene.textbox
	tb.Clear()
	tb.Orig = pixel.V(200, 200)
	_, _ = fmt.Fprintln(tb, "Spela!")
	_, _ = fmt.Fprintln(tb, "Avsluta")
	tb.DrawColorMask(win, pixel.IM.Scaled(tb.Orig, 2), colornames.Black)
}

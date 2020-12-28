package internal

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
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
	currentitem     MenuItem
	textbox         *text.Text
	itemSwitchSound *beep.Buffer
}

func MakeMenuScene() MenuScene {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	err, _, switchSound := LoadWav("assets/MenuPointerMoved.wav")
	PanicIfError(err)
	return MenuScene{
		currentitem:     Play,
		itemSwitchSound: switchSound,
		textbox:         text.New(pixel.V(0, 0), atlas),
	}
}

func (menuScene MenuScene) HandleKeyDown(key ControlKey) Scene {
	if key == Jump {
		if menuScene.currentitem == Play {
			return MakeMapScene()
		} else {
			return nil
		}
	}
	if key == Down || key == Up {
		streamer := menuScene.itemSwitchSound.Streamer(0, menuScene.itemSwitchSound.Len())
		speaker.Play(streamer)
		menuScene.currentitem = (menuScene.currentitem + 1) % 2
	}
	return menuScene
}

func (menuScene MenuScene) HandleKeyUp(_ ControlKey) Scene {
	return menuScene
}

func (menuScene MenuScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Aliceblue)
	tb := menuScene.textbox
	tb.Clear()
	tb.Orig = pixel.V(300, 300)
	playItem := "  Spela!"
	if menuScene.currentitem == Play {
		playItem = "* Spela!"
	}
	_, _ = fmt.Fprintln(tb, playItem)

	quitItem := "  Avsluta"
	if menuScene.currentitem == Quit {
		quitItem = "* Avsluta"
	}
	_, _ = fmt.Fprintln(tb, quitItem)
	tb.DrawColorMask(win, pixel.IM.Scaled(tb.Orig, 2), colornames.Black)
}

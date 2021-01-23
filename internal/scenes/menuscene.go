package scenes

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"objarni/rescue-on-fractal-bun/internal"
)

type MenuItem int

const (
	Play MenuItem = iota
	Quit
)

type MenuScene struct {
	cfg             *Config
	currentitem     MenuItem
	textbox         *text.Text
	itemSwitchSound *beep.Buffer
	quit            bool
}

func MakeMenuScene(config *Config) *MenuScene {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	err, _, switchSound := internal.LoadWav("assets/MenuPointerMoved.wav")
	internal.PanicIfError(err)
	return &MenuScene{
		currentitem:     Play,
		itemSwitchSound: switchSound,
		textbox:         text.New(pixel.V(0, 0), atlas),
		quit:            false,
		cfg:             config,
	}
}

func (menuScene *MenuScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		if menuScene.currentitem == Play {
			return MakeMapScene(menuScene.cfg, "Hembyn")
		} else {
			menuScene.quit = true
		}
	}
	if key == internal.Down || key == internal.Up {
		streamer := menuScene.itemSwitchSound.Streamer(0, menuScene.itemSwitchSound.Len())
		speaker.Play(streamer)
		menuScene.currentitem = (menuScene.currentitem + 1) % 2
	}
	return menuScene
}

func (menuScene *MenuScene) HandleKeyUp(_ internal.ControlKey) internal.Thing {
	return menuScene
}

func (menuScene *MenuScene) Render(win *pixelgl.Window) {
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

func (menuScene *MenuScene) Tick() bool {
	return !menuScene.quit
}

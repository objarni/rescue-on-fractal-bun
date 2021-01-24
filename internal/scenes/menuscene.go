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
	res             *Resources
	currentItem     MenuItem
	textbox         *text.Text
	itemSwitchSound *beep.Buffer
	quit            bool
}

func MakeMenuScene(config *Config, res *Resources) *MenuScene {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	err, _, switchSound := internal.LoadWav("assets/MenuPointerMoved.wav")
	internal.PanicIfError(err)
	return &MenuScene{
		cfg:             config,
		res:             res,
		currentItem:     Play,
		textbox:         text.New(pixel.V(0, 0), atlas),
		itemSwitchSound: switchSound,
		quit:            false,
	}
}

func (menuScene *MenuScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		if menuScene.currentItem == Play {
			return MakeMapScene(menuScene.cfg, menuScene.res, "Hembyn")
		} else {
			menuScene.quit = true
		}
	}
	if key == internal.Down || key == internal.Up {
		streamer := menuScene.itemSwitchSound.Streamer(0, menuScene.itemSwitchSound.Len())
		speaker.Play(streamer)
		menuScene.currentItem = (menuScene.currentItem + 1) % 2
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
	if menuScene.currentItem == Play {
		playItem = "* Spela!"
	}
	_, _ = fmt.Fprintln(tb, playItem)

	quitItem := "  Avsluta"
	if menuScene.currentItem == Quit {
		quitItem = "* Avsluta"
	}
	_, _ = fmt.Fprintln(tb, quitItem)
	tb.DrawColorMask(win, pixel.IM.Scaled(tb.Orig, 2), colornames.Black)
}

func (menuScene *MenuScene) Tick() bool {
	return !menuScene.quit
}

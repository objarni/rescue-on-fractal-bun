package scenes

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
)

type MenuItem int

const (
	Play MenuItem = iota
	Quit
)

type MenuScene struct {
	cfg             *Config
	res             *internal.Resources
	currentItem     MenuItem
	textbox         *text.Text
	itemSwitchSound *beep.Buffer
	quit            bool
}

func MakeMenuScene(config *Config, res *internal.Resources) *MenuScene {
	return &MenuScene{
		cfg:             config,
		res:             res,
		currentItem:     Play,
		itemSwitchSound: res.Blip,
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
	tb := text.New(pixel.ZV, menuScene.res.Atlas)
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
	tbCenter := tb.Bounds().Center().Scaled(2)
	translation := win.Bounds().Center().Sub(tbCenter)
	tb.DrawColorMask(win, pixel.IM.Scaled(pixel.ZV, 2).Moved(translation), colornames.Black)
}

func (menuScene *MenuScene) Tick() bool {
	return !menuScene.quit
}

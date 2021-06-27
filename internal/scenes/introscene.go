package scenes

import (
	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/tweaking"
)

// TODO: Use animation structure/functions instead
//  of re-doing it custom made here in intro scene!
type IntroScene struct {
	gif              internal.GifData
	cfg              *tweaking.Config
	res              *internal.Resources
	frame            int
	frameDisplayTime int
}

func (introScene *IntroScene) WantToExitProgram() bool {
	return false
}

func MakeIntroScene(config *tweaking.Config, res *internal.Resources) *IntroScene {
	return &IntroScene{
		gif:   internal.LoadGif("testdata/example3.gif"),
		cfg:   config,
		res:   res,
		frame: 0,
	}
}

func (introScene *IntroScene) HandleKeyDown(_ internal.ControlKey) Scene {
	return MakeMapScene(introScene.cfg, introScene.res, "Hembyn")
}

func (introScene *IntroScene) HandleKeyUp(_ internal.ControlKey) Scene {
	return introScene
}

func (introScene *IntroScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.White)
	sprite := introScene.gif.Frames[introScene.frame]
	sprite.Draw(win,
		px.IM.Moved(
			sprite.Frame().Center()).Scaled(
			px.V(400, 300), 0.1))
}

func (introScene *IntroScene) Tick() Scene {
	// Switch frame?
	introScene.frameDisplayTime += internal.TickTimeMs
	if introScene.frameDisplayTime > introScene.gif.DisplayFrameMs {
		introScene.frameDisplayTime -= introScene.gif.DisplayFrameMs
		introScene.frame = introScene.frame + 1
		if introScene.frame >= introScene.gif.FrameCount {
			introScene.frame = 0
		}
	}
	return introScene
}

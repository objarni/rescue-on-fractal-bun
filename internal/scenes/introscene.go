package scenes

import (
	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/tweaking"
)

type IntroScene struct {
	gif              internal.GifData
	cfg              *tweaking.Config
	res              *internal.Resources
	frame            int
	frameDisplayTime int
}

func MakeIntroScene(config *tweaking.Config, res *internal.Resources) *IntroScene {
	return &IntroScene{
		gif:   internal.LoadGif("testdata/example.gif"),
		cfg:   config,
		res:   res,
		frame: 0,
	}
}

func (introScene *IntroScene) HandleKeyDown(_ internal.ControlKey) internal.Thing {
	return MakeMenuScene(introScene.cfg, introScene.res)
}

func (introScene *IntroScene) HandleKeyUp(_ internal.ControlKey) internal.Thing {
	return introScene
}

func (introScene *IntroScene) Render(win *pixelgl.Window) {
	//win.Clear(colornames.Black)
	sprite := introScene.gif.Frames[introScene.frame]
	sprite.Draw(win,
		px.IM.Moved(
			sprite.Frame().Center()).Scaled(
			px.V(400, 300), 0.1))
}

// Tick TODO: Why is Tick still mutating?
func (introScene *IntroScene) Tick() bool {
	// Switch frame?
	introScene.frameDisplayTime += internal.TickTimeMs
	if introScene.frameDisplayTime > introScene.gif.DisplayFrameMs {
		introScene.frameDisplayTime -= introScene.gif.DisplayFrameMs
		introScene.frame = introScene.frame + 1
		if introScene.frame >= introScene.gif.FrameCount {
			introScene.frame = 0
		}
	}
	return true
}

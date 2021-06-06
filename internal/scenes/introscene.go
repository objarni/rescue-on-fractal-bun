package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/tweaking"
)

type IntroScene struct {
	gif   internal.GifData
	cfg   *tweaking.Config
	res   *internal.Resources
	frame int
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
	win.Clear(colornames.Black)
	sprite := internal.SpriteFromImage(*introScene.gif.Images[introScene.frame])
	sprite.Draw(win, pixel.IM.Moved(sprite.Frame().Center()))
}

// Tick TODO: Why is Tick still mutating?
func (introScene *IntroScene) Tick() bool {
	introScene.frame = introScene.frame + 1
	if introScene.frame >= introScene.gif.FrameCount {
		introScene.frame = 0
	}
	return true
}

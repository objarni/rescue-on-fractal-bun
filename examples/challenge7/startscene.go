package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"objarni/rescue-on-fractal-bun/internal"
)

func (startScene *StartScene) HandleKeyUp(_ internal.ControlKey) internal.Thing {
	return startScene
}

func (startScene *StartScene) HandleKeyDown(key internal.ControlKey) internal.Thing {
	if key == internal.Jump {
		var config Config
		config, err := TryReadCfgFrom("json/challenge7.json", config)
		internal.PanicIfError(err)

		var gubbeStandingRightSprite = internal.LoadSpriteForSure("assets/TStanding.png")
		var gubbeWalkingRightSprite1 = internal.LoadSpriteForSure("assets/TWalking-1.png")
		var gubbeWalkingRightSprite2 = internal.LoadSpriteForSure("assets/TWalking-2.png")
		var gubbeKickingRightSprite = internal.LoadSpriteForSure("assets/TWalking-1.png")
		gubbeImage2Sprite := map[Image]*pixel.Sprite{
			WalkRight1:    gubbeWalkingRightSprite1,
			WalkRight2:    gubbeWalkingRightSprite2,
			StandingRight: gubbeStandingRightSprite,
			KickRight:     gubbeKickingRightSprite,
		}
		var gubbe = MakeGubbe(pixel.Vec{X: 100, Y: 150}, gubbeImage2Sprite)
		var scene internal.Thing = &GameScene{
			ball:  MakeBall(config),
			gubbe: gubbe,
		}

		return scene
	}
	return startScene
}

func (startScene *StartScene) Tick() bool {
	return true
}

func (startScene *StartScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Green)
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(0, 0), basicAtlas)
	basicTxt.Clear()
	_, _ = fmt.Fprintln(basicTxt, "PRESS SPACE TO PLAY")
	tbh := basicTxt.Bounds().Size().Scaled(0.5 * 5)
	fmt.Print(tbh)
	pos := win.Bounds().Center().Add(pixel.Vec{X: -tbh.X, Y: tbh.Y})
	basicTxt.Draw(win, pixel.IM.Moved(pos).Scaled(pos, 5))
}

/*
y
^
|
|
|
+---------> x
TEXTBOX


*/
func MakeStartScene() *StartScene {
	return &StartScene{}
}

type StartScene struct {
}

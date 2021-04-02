package scenes

import (
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"testing"
)

func Test_initialRender(t *testing.T) {
	cfg := TryReadCfgFrom("../../"+internal.ConfigFile, Config{})
	res := internal.Resources{
		ImageMap: map[internal.Image]*pixel.Sprite{},
		MapSigns: []internal.MapSign{{
			MapPos:    pixel.Vec{},
			LevelName: "",
			LevelPos:  pixel.Vec{},
		}},
	}
	mapScene := MakeMapScene(&cfg, &res, "Hembyn")
	op := mapScene.MapSceneWinOp()
	approvals.VerifyString(t, op.String())
}

package scenes_test

import (
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
	"objarni/rescue-on-fractal-bun/internal/tweaking"
	"objarni/rescue-on-fractal-bun/tests"
	"testing"
)

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

func Test_initialRender(t *testing.T) {
	cfg := tweaking.TryReadCfgFrom("../../"+internal.ConfigFile, tweaking.Config{})
	res := internal.Resources{
		ImageMap: map[internal.Image]*pixel.Sprite{},
		MapSigns: []internal.MapSign{{
			MapPos:    pixel.Vec{},
			LevelName: "",
			LevelPos:  pixel.Vec{},
			Text:      "Hembyn",
		}},
	}
	mapScene := scenes.MakeMapScene(&cfg, &res, "Hembyn")
	op := mapScene.MapSceneWinOp()
	approvals.VerifyString(t, op.String()+"\n")
}

package tests

import (
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
	"testing"
)

func Test_initialRender(t *testing.T) {
	cfg := scenes.TryReadCfgFrom("../"+internal.ConfigFile, scenes.Config{})
	res := internal.Resources{
		ImageMap: map[internal.Image]*pixel.Sprite{},
	}
	mapScene := scenes.MakeMapScene(&cfg, &res, "Skogen")
	op := mapScene.MapSceneWinOp()
	approvals.VerifyString(t, op.String())
}

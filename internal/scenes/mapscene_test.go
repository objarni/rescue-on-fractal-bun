package scenes

import (
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"testing"
	"time"
)

func Test_initialRender(t *testing.T) {
	cfg := TryReadCfgFrom("../"+internal.ConfigFile, Config{})
	res := internal.Resources{
		ImageMap: map[internal.Image]*pixel.Sprite{},
	}
	mapScene := MakeMapScene(&cfg, &res, "Skogen")
	op := mapScene.MapSceneWinOp()

	approvals.VerifyString(t, op.String())
	duration, err := time.ParseDuration("3500ms")
	if err != nil {
		panic(err)
	}
	time.Sleep(duration)
}

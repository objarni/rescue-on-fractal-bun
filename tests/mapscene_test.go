package tests

import (
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	r := approvals.UseReporter(NewBCompare())
	defer r.Close()
	os.Exit(m.Run())
}

func Test_initialRender(t *testing.T) {
	cfg := scenes.TryReadCfgFrom("../"+internal.ConfigFile, scenes.Config{})
	res := internal.Resources{
		ImageMap: map[internal.Image]*pixel.Sprite{},
	}
	mapScene := scenes.MakeMapScene(&cfg, &res, "Skogen")
	op := mapScene.MapSceneWinOp()

	approvals.VerifyString(t, op.String())
	duration, err := time.ParseDuration("3500ms")
	if err != nil {
		panic(err)
	}
	time.Sleep(duration)
}

package entities

import (
	approvals "github.com/approvals/go-approval-tests"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/tests"
	"testing"
)

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

func Test_levelBoundaryInitialState(t *testing.T) {
	levelBoundary := MakeLevelBoundary(px.R(0, 1, 2, 3))
	approvals.VerifyString(t, levelBoundary.String()+"\n")
}

func Test_levelBoundaryGfx(t *testing.T) {
	levelBoundary := MakeLevelBoundary(px.R(0, 1, 2, 3))
	approvals.VerifyString(t, levelBoundary.GfxOp(nil).String()+"\n")
}

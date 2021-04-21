package entities

import (
	approvals "github.com/approvals/go-approval-tests"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal/events"
	"objarni/rescue-on-fractal-bun/tests"
	"testing"
)

/*
- #initial state
- walking right
- walking left
- jumping when standing still
- jumping when pressing right
- clicking action when standing still
- clicking action when mid-air
- taking damage
*/

func Test_eliseInitial(t *testing.T) {
	elise := MakeElise(px.V(10, 20))
	approvals.VerifyString(t, elise.String())
}

func Test_pressingLeft(t *testing.T) {
	elise := MakeElise(px.V(10, 20))
	elise = elise.Handle(EventBox{
		Event: events.KeyLeftDown,
		Box:   px.Rect{},
	})
	elise = elise.Tick(0, nil)
	approvals.VerifyString(t, elise.String())
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

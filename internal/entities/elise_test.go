package entities

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal/events"
	"objarni/rescue-on-fractal-bun/tests"
	"strings"
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
	var state = printElise(elise)
	approvals.VerifyString(t, state)
}

func Test_pressingLeft(t *testing.T) {
	elise := MakeElise(px.V(10, 20))
	elise = elise.Handle(EventBox{
		Event: events.KeyLeftDown,
		Box:   px.Rect{},
	})
	elise = elise.Tick(0, nil)
	var state = printElise(elise)
	approvals.VerifyString(t, state)
}

func printElise(elise Entity) string {
	state := fmt.Sprintf("Elise %v", "standing")
	hb := fmt.Sprintf("HitBox %v", printRect(elise.HitBox()))
	facing := fmt.Sprintf("Facing right")
	all := []string{state, hb, facing}
	return strings.Join(all, "\n") + "\n"
}

func printRect(box px.Rect) interface{} {
	return fmt.Sprintf("[%1.0f,%1.0f->%1.0f,%1.0f]",
		box.Min.X, box.Min.Y, box.Max.X, box.Max.Y)
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

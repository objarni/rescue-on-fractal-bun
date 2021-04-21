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
   - #walking right
   - #pressing left
   - clicking action when standing still
   - taking damage
   - jumping when standing still
   - jumping when pressing right
   - clicking action when mid-air
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

func Test_walkingRight(t *testing.T) {
	elise := MakeElise(px.V(10, 20))
	before := "Before:\n" + elise.String()
	elise = elise.Handle(EventBox{
		Event: events.KeyRightDown,
		Box:   px.Rect{},
	})
	elise = passTime(4, elise)
	after := "After:\n" + elise.String()
	approvals.VerifyString(t, before+after)
}
func Test_actionWhenStanding(t *testing.T) {
	elise := MakeElise(px.V(10, 20))
	before := "** Before **\n" + elise.String()
	elise = elise.Handle(EventBox{
		Event: events.KeyActionDown,
		Box:   px.Rect{},
	})
	receiver := MakeEntityCanvas()
	elise = elise.Tick(0, &receiver)
	after := "\n** After **\n" + elise.String()
	canvas := "\n** Canvas **\n" + receiver.String()
	approvals.VerifyString(t, before+after+canvas)
}

func passTime(ticks int, elise Entity) Entity {
	for i := 0; i < ticks; i++ {
		elise = elise.Tick(0, nil)
	}
	return elise
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

// TODO: for improved readability, want some kind of 'manuscript' that is played up by helper function

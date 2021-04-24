package entities

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal/events"
	"objarni/rescue-on-fractal-bun/tests"
	"testing"
)

/*
   + initial state
   + walking right
   + pressing left
   + clicking action when standing still
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
	box := EventBox{
		Event: events.KeyLeftDown,
		Box:   px.R(1, 1, 2, 2),
	}
	approvals.VerifyString(t, simulate(box, 1))
}

func Test_actionWhenStanding(t *testing.T) {
	result := simulate(EventBox{
		Event: events.KeyActionDown,
		Box:   px.Rect{},
	}, 1)
	approvals.VerifyString(t, result)
}

func Test_walkingRight(t *testing.T) {
	ticks := 4
	box := EventBox{
		Event: events.KeyRightDown,
		Box:   px.Rect{},
	}
	approvals.VerifyString(t, simulate(box, ticks))
}

func Test_takingDamage(t *testing.T) {
	result := simulate(EventBox{
		Event: events.Damage,
		Box: px.Rect{
			Min: px.V(-10, -10),
			Max: px.V(10, 10),
		},
	}, 1)
	approvals.VerifyString(t, result)
}

func simulate(box EventBox, ticks int) string {
	elise := MakeElise(px.V(0, 0))
	canvasSequence := ""
	var entityCanvas EntityCanvas
	for ix, _ := range make([]int, ticks) {
		entityCanvas = MakeEntityCanvas()
		entityCanvas.AddEntityHitBox(EntityHitBox{
			Entity: 0,
			HitBox: elise.HitBox(),
		})
		entityCanvas.AddEventBox(box)
		entityCanvas.Consequences(func(eb EventBox, ehb EntityHitBox) {
			elise = elise.Handle(eb)
		})
		elise = elise.Tick(0, &entityCanvas)
		canvasSequence += printCanvasTick(ix, entityCanvas)
	}
	canvasSequence += printCanvasTick(ticks, entityCanvas) + "\n"
	scenario := fmt.Sprintf(
		"\n** Scenario **\nEventBox: %v\nTicks: %v\n",
		box,
		ticks,
	)
	endState := "\n** Elise end state **\n" + elise.String()
	canvas := "\n** Canvas states **\n" + canvasSequence
	return scenario + endState + canvas
}

func printCanvasTick(ticks int, entityCanvas EntityCanvas) string {
	return fmt.Sprintf(" * Tick %v *\n", ticks) + entityCanvas.String()
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

// TODO: for improved readability, want some kind of 'manuscript' that is played up by helper function

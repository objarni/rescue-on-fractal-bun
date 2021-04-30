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
   + nothing beneath = falling
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

var rectOverlappingElise = px.R(-100000, -100000, 100000, 100000)

func Test_pressingLeft(t *testing.T) {
	box := EventBox{
		Event: events.KeyLeftDown,
		Box:   rectOverlappingElise,
	}
	approvals.VerifyString(t, simulate(box, 1))
}

//func Test_falling(t *testing.T) {
//	box := EventBox{
//		Event: events.NoEvent,
//		Box:   rectOverlappingElise,
//	}
//	approvals.VerifyString(t, simulate(box, 10))
//}

func Test_actionWhenStanding(t *testing.T) {
	result := simulate(EventBox{
		Event: events.KeyActionDown,
		Box:   rectOverlappingElise,
	}, 1)
	approvals.VerifyString(t, result)
}

func Test_walkingRight(t *testing.T) {
	ticks := 4
	box := EventBox{
		Event: events.KeyRightDown,
		Box:   rectOverlappingElise,
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
	scenario := fmt.Sprintf(
		"** Scenario **\n"+
			"Event: %v\n"+
			"Tick count: %v\n"+
			"Elise start state:\n%v\n",
		box.String(),
		ticks,
		elise.String(),
	)
	simulationLog := ""
	var entityCanvas EntityCanvas
	for ix := range make([]int, ticks) {
		entityCanvas = FillCanvas(box, entityCanvas, elise)
		entityCanvas.Consequences(func(eb EventBox, ehb EntityHitBox) {
			elise = elise.Handle(eb)
		})
		elise = elise.Tick(0, &entityCanvas)
		simulationLog += printCanvasTick(ix+1, entityCanvas)
	}
	endState := "Elise end state:\n" + elise.String()
	canvas := "\n** Simulation **\n\n" + simulationLog
	return scenario + endState + canvas
}

func FillCanvas(box EventBox, entityCanvas EntityCanvas, elise Entity) EntityCanvas {
	entityCanvas = MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{
		Entity: 0,
		HitBox: elise.HitBox(),
	})
	entityCanvas.AddEventBox(box)
	return entityCanvas
}

func printCanvasTick(ticks int, entityCanvas EntityCanvas) string {
	return fmt.Sprintf(" * Tick %v *\n%v\n", ticks, entityCanvas.String())
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

// TODO: for improved readability, want some kind of 'manuscript' that is played up by helper function

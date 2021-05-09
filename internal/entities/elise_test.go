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
   + nothing beneath = falling
   + initial state
   + walking right
   + pressing left
   + clicking action when standing still
   + taking damage
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
	approvals.VerifyString(t, simulate([]EventBox{box}, 1, 0))
}

func Test_falling(t *testing.T) {
	box := EventBox{
		Event: events.NoEvent,
		Box:   rectOverlappingElise,
	}
	approvals.VerifyString(t, simulate([]EventBox{box}, 10, -100))
}

func Test_actionWhenStanding(t *testing.T) {
	box := EventBox{
		Event: events.KeyActionDown,
		Box:   rectOverlappingElise,
	}
	result := simulate([]EventBox{box}, 1, 0)
	approvals.VerifyString(t, result)
}

func Test_walkingRight(t *testing.T) {
	ticks := 4
	box := EventBox{
		Event: events.KeyRightDown,
		Box:   rectOverlappingElise,
	}
	approvals.VerifyString(t, simulate([]EventBox{box}, ticks, 0))
}

func Test_walkingRightIntoWall(t *testing.T) {
	ticks := 4
	box := EventBox{
		Event: events.KeyRightDown,
		Box:   rectOverlappingElise,
	}
	wallBox := EventBox{
		Event: events.Wall,
		Box:   px.R(11, 0, 21, 10),
	}
	approvals.VerifyString(t, simulate([]EventBox{box, wallBox}, ticks, 0))
}

func Test_takingDamage(t *testing.T) {
	box := EventBox{
		Event: events.Damage,
		Box: px.Rect{
			Min: px.V(-10, -10),
			Max: px.V(10, 10),
		},
	}
	approvals.VerifyString(t, simulate([]EventBox{box}, 1, 0))
}

func Test_startOfJumpWhenStandingStill(t *testing.T) {
	box := EventBox{
		Event: events.KeyJumpDown,
		Box:   rectOverlappingElise,
	}
	approvals.VerifyString(t, simulate([]EventBox{box}, 1, 0))
}

func Test_fullJumpStandingStill(t *testing.T) {
	box := EventBox{
		Event: events.KeyJumpDown,
		Box:   rectOverlappingElise,
	}
	approvals.VerifyString(t, simulate([]EventBox{box}, 20, 0))
}

func simulate(boxes []EventBox, ticks int, groundHeight int) string {
	elise := MakeElise(px.V(0, 0))

	es := make([]string, 0)
	for _, box := range boxes {
		es = append(es, box.String())
	}

	scenario := fmt.Sprintf(
		"*** Scenario ***\n"+
			"* Events:\n%v\n\n"+
			"Ground height: %v\n"+
			"Tick count: %v\n"+
			"\n* Elise start state:\n%v\n",
		strings.Join(es, "\n"),
		groundHeight,
		ticks,
		elise.String(),
	)
	simulationLog := ""
	var entityCanvas EntityCanvas
	for ix := range make([]int, ticks) {
		entityCanvas = FillCanvas(boxes, entityCanvas, elise, groundHeight)
		entityCanvas.Consequences(func(eb EventBox, ehb EntityHitBox) {
			elise = elise.Handle(eb)
		})
		elise = elise.Tick(0, &entityCanvas)
		simulationLog += printCanvasTick(ix+1, entityCanvas)
	}
	endState := "* Elise end state:\n" + elise.String()
	canvas := "\n\n*** Simulation ***\n\n" + simulationLog
	return scenario + endState + canvas
}

func FillCanvas(boxes []EventBox, entityCanvas EntityCanvas, elise Entity, groundHeight int) EntityCanvas {
	entityCanvas = MakeEntityCanvas()
	entityCanvas.AddEntityHitBox(EntityHitBox{
		Entity: 0,
		HitBox: elise.HitBox(),
	})

	entityCanvas.AddEventBox(EventBox{
		Event: events.Wall,
		Box: px.R(
			-1000,
			float64(groundHeight-1000),
			1000,
			float64(groundHeight),
		),
	})

	for _, box := range boxes {
		entityCanvas.AddEventBox(box)
	}

	return entityCanvas
}

func printCanvasTick(ticks int, entityCanvas EntityCanvas) string {
	return fmt.Sprintf(" * Tick %v *\n%v\n", ticks, entityCanvas.String())
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

// TODO: for improved readability, want some kind of 'manuscript' that is played up by helper function

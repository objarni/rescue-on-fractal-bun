package main

import (
	"fmt"
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/faiface/pixel"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	r := approvals.UseReporter(reporters.NewIntelliJReporter())
	code := m.Run()
	err := r.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func MakeGubbe() Gubbe {
	return Gubbe{
		state:   Standing,
		looking: Right,
		image:   StandingRight,
		pos:     pixel.Vec{},
		vel:     pixel.Vec{},
		acc:     pixel.Vec{},
	}
}

func simulateSteps(gubbe *Gubbe, steps int, controls Controls) string {
	result := ""
	for step := 0; step < steps; step++ {
		stepGubbe(gubbe, controls)
		result += fmt.Sprintf("Step %02d\n%s\n%s\n",
			step, printControls(controls), printGubbe(*gubbe))
	}
	return result
}

func toScenarioName(testName string) string {
	return strings.Replace(testName, "Test", "Scenario: ", 1)
}

func printGubbe(gubbe Gubbe) string {
	return fmt.Sprintf("State: %s\nLooking: %s\nImage: %s\nPos: %s\nVel: %s\nAcc: %s\n",
		gubbe.state.String(),
		gubbe.looking.String(),
		gubbe.image.String(),
		gubbe.pos,
		gubbe.vel,
		gubbe.acc,
	)
}

func printControls(controls Controls) string {
	pad := ""
	pad += checkPressed(controls.left, "<")
	pad += checkPressed(controls.kick, "^")
	pad += checkPressed(controls.right, ">")
	return "Pad: " + pad
}

func checkPressed(pressed bool, symbol string) string {
	if pressed {
		return symbol
	} else {
		return "."
	}
}

//func TestWalkingLeftForOneSecond(t *testing.T) {
//	gubbe := MakeGubbe()
//	for i := 0; i<10; i++ {
//		stepGubbe(&gubbe, Controls{left: true, right: false, kick: false})
//	}
//	approvals.VerifyString(t, printGubbe(gubbe))
//}
//
//func TestKicking(t *testing.T) {
//	gubbe := MakeGubbe()
//	for i := 0; i<10; i++ {
//		stepGubbe(&gubbe, Controls{left: false, right: false, kick: true})
//	}
//	approvals.VerifyString(t, printGubbe(gubbe))
//}
//
//func TestTryingToMoveWhileKicking(t *testing.T) {
//	gubbe := MakeGubbe()
//	stepGubbe(&gubbe, Controls{left: false, right: false, kick: true})
//	for i := 0; i<3; i++ {
//		stepGubbe(&gubbe, Controls{left: false, right: true, kick: false})
//	}
//	approvals.VerifyString(t, printGubbe(gubbe))
//}

/*
Issues with current approach
- does not capture 'time passing' e.g animation
- only captures end frame, no frame inbetween
- does not capture actual frame swaps, only gubbe state
- test do not describe timings for pressed keys; only on/off
- not really happy with the types for gubbe, they are too
  dissimilar from game scenes/states when actually they
  work from same 'things': key presses, time passing, and
  rendering.

Ideas
The update 'logs' states with timestamps, and updates 100 times
 per second, so giving a deltaSecond of 0.25 logs the following
 timestamps: 0.0, 0.1. What to do about the 0.05 left over?
 For now: assume main loop handles it, so that the stepGubbe
 functionality actually gets an integer number of 'animationFrames'
 to calcularte.
*/

var pressRight = Controls{left: false, right: true, kick: false}
var pressLeft = Controls{left: true, right: false, kick: false}
var pressNothing = Controls{left: false, right: false, kick: false}

func TestWalkingRight10Steps(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := MakeGubbe()
	result += simulateSteps(&gubbe, 10, pressRight)
	approvals.VerifyString(t, result)
}

func TestWalkingLeft10Steps(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := MakeGubbe()
	result += simulateSteps(&gubbe, 10, pressLeft)
	approvals.VerifyString(t, result)
}

func TestWalkingLeftThenBackAgain(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := MakeGubbe()
	result += simulateSteps(&gubbe, 5, pressLeft)
	result += simulateSteps(&gubbe, 5, pressRight)
	approvals.VerifyString(t, result)
}

func TestWalkingLeft5StepsThenStop(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := MakeGubbe()
	result += simulateSteps(&gubbe, 5, pressLeft)
	result += simulateSteps(&gubbe, 5, pressNothing)
	approvals.VerifyString(t, result)
}

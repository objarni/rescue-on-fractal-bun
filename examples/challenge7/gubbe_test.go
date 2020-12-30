package main

import (
	"fmt"
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/faiface/pixel"
	"os"
	"testing"
	"unicode"
)

var step int = 0

func TestMain(m *testing.M) {
	// This code is run before all tests
	r := approvals.UseReporter(reporters.NewIntelliJReporter())
	step = 0
	code := m.Run()
	err := r.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func makeGubbe() Gubbe {
	return Gubbe{
		state:   Standing,
		looking: Right,
		image:   StandingRight,
		pos:     pixel.ZV,
		vel:     pixel.ZV,
		acc:     pixel.ZV,
	}
}

func simulateSteps(gubbe *Gubbe, steps int, controls Controls) string {
	result := ""
	for i := 0; i < steps; i++ {
		stepGubbe(gubbe, controls)
		result += fmt.Sprintf("Step %02d\n%s\n%s\n",
			0, printControls(controls), printGubbe(*gubbe))
		step += 1
	}
	return result
}

func toScenarioName(testName string) string {
	name := ""
	for _, ch := range testName[4:] {
		if unicode.IsUpper(ch) || unicode.IsNumber(ch) {
			name += " "
		}
		name += string(ch)
	}
	return "Scenario: " + name
}

func printGubbe(gubbe Gubbe) string {
	return fmt.Sprintf("State: %s\nLooking: %s\nImage: %s\nPos: %s\nVel: %s\nAcc: %s\n",
		gubbe.state.String(),
		gubbe.looking.String(),
		gubbe.image.String(),
		printVec(gubbe.pos),
		printVec(gubbe.vel),
		printVec(gubbe.acc),
	)
}

func printVec(vec pixel.Vec) string {
	return fmt.Sprintf("(%1.1f, %1.1f)", vec.X, vec.Y)
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

/*
Issues with current approach
- not really happy with the types for gubbe, they are too
  dissimilar from game scenes/states when actually they
  work from same 'things': key presses, time passing, and
  rendering
- step is always 00 in approval files
*/

var pressRight = Controls{left: false, right: true, kick: false}
var pressLeft = Controls{left: true, right: false, kick: false}
var pressNothing = Controls{left: false, right: false, kick: false}

func TestWalkingRight5Steps(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := makeGubbe()
	result += simulateSteps(&gubbe, 10, pressRight)
	approvals.VerifyString(t, result)
}

func TestWalkingLeft5Steps(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := makeGubbe()
	result += simulateSteps(&gubbe, 10, pressLeft)
	approvals.VerifyString(t, result)
}

func TestWalkingLeftThenBackAgain(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := makeGubbe()
	result += simulateSteps(&gubbe, 5, pressLeft)
	result += simulateSteps(&gubbe, 5, pressRight)
	approvals.VerifyString(t, result)
}

func TestWalkingLeft5StepsThenStop(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := makeGubbe()
	result += simulateSteps(&gubbe, 5, pressLeft)
	result += simulateSteps(&gubbe, 5, pressNothing)
	approvals.VerifyString(t, result)
}

func TestLeftAndRightMeansStandStill(t *testing.T) {
	result := toScenarioName(t.Name()) + "\n"
	gubbe := makeGubbe()
	result += simulateSteps(&gubbe, 5, Controls{left: true, right: true, kick: false})
	approvals.VerifyString(t, result)
}

package main

import (
	"fmt"
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"os"
	"testing"
	"unicode"
)

var globalStepVariable int = 0

func TestMain(m *testing.M) {
	// This code is run before all tests
	r := approvals.UseReporter(reporters.NewIntelliJReporter())
	code := m.Run()
	err := r.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func simulateSteps(gubbe *Gubbe, simulationSteps int, controls Controls) string {
	result := ""
	for i := 0; i < simulationSteps; i++ {
		if controls.left {
			gubbe.HandleKeyDown(internal.Left)
		} else {
			gubbe.HandleKeyUp(internal.Left)
		}
		if controls.right {
			gubbe.HandleKeyDown(internal.Right)
		} else {
			gubbe.HandleKeyUp(internal.Right)
		}
		if controls.kick {
			gubbe.HandleKeyDown(internal.Action)
		} else {
			gubbe.HandleKeyUp(internal.Action)
		}
		gubbe.Tick()
		result += fmt.Sprintf("Step %02d\n%s\n%s\n",
			globalStepVariable, printControls(controls), printGubbe(*gubbe))
		globalStepVariable += 1
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
	pad += checkPressed(controls.kick, "K")
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
- globalStepVariable is always 00 in approval files
*/

var pressRight = Controls{left: false, right: true, kick: false}
var pressLeft = Controls{left: true, right: false, kick: false}
var pressNothing = Controls{left: false, right: false, kick: false}

func initGubbe() Gubbe {
	return MakeGubbe(pixel.ZV, nil)
}

func TestWalkingRight5Steps(t *testing.T) {
	globalStepVariable = 0
	result := toScenarioName(t.Name()) + "\n"
	gubbe := initGubbe()
	result += simulateSteps(&gubbe, 10, pressRight)
	approvals.VerifyString(t, result)
}

func TestWalkingLeft5Steps(t *testing.T) {
	globalStepVariable = 0
	result := toScenarioName(t.Name()) + "\n"
	gubbe := initGubbe()
	result += simulateSteps(&gubbe, 10, pressLeft)
	approvals.VerifyString(t, result)
}

func TestWalkingLeftThenBackAgain(t *testing.T) {
	globalStepVariable = 0
	result := toScenarioName(t.Name()) + "\n"
	gubbe := initGubbe()
	result += simulateSteps(&gubbe, 5, pressLeft)
	result += simulateSteps(&gubbe, 5, pressRight)
	approvals.VerifyString(t, result)
}

func TestWalkingLeft5StepsThenStop(t *testing.T) {
	globalStepVariable = 0
	result := toScenarioName(t.Name()) + "\n"
	gubbe := initGubbe()
	result += simulateSteps(&gubbe, 5, pressLeft)
	result += simulateSteps(&gubbe, 5, pressNothing)
	approvals.VerifyString(t, result)
}

func TestLeftAndRightMeansStandStill(t *testing.T) {
	globalStepVariable = 0
	result := toScenarioName(t.Name()) + "\n"
	gubbe := initGubbe()
	result += simulateSteps(&gubbe, 5, Controls{left: true, right: true, kick: false})
	approvals.VerifyString(t, result)
}

func TestLongKick(t *testing.T) {
	globalStepVariable = 0
	result := toScenarioName(t.Name()) + "\n"
	gubbe := initGubbe()
	result += simulateSteps(&gubbe, 15, Controls{left: false, right: false, kick: true})
	approvals.VerifyString(t, result)
}

func TestQuickKick(t *testing.T) {
	globalStepVariable = 0
	result := toScenarioName(t.Name()) + "\n"
	gubbe := initGubbe()
	result += simulateSteps(&gubbe, 20, Controls{left: false, right: false, kick: true})
	result += simulateSteps(&gubbe, 10, Controls{left: false, right: false, kick: false})
	approvals.VerifyString(t, result)
}

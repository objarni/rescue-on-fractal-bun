package main

import (
	"fmt"
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/faiface/pixel"
	"os"
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

func printGubbe(gubbe Gubbe) string {
	return fmt.Sprintf("State: %s\nLooking: %s\nPosition: %s\n",
		gubbe.state.String(),
		gubbe.looking.String(),
		gubbe.position)
}

func MakeGubbe() Gubbe {
	return Gubbe{
		state:       Standing,
		looking:     Right,
		timeInState: 0,
		position:    pixel.Vec{},
	}
}

func TestWalkingRightForOneSecond(t *testing.T) {
	gubbe := MakeGubbe()
	updateGubbe(&gubbe, 1, Controls{left: false, right: true, kick: false})
	approvals.VerifyString(t, printGubbe(gubbe))
}

func TestWalkingLeftForOneSecond(t *testing.T) {
	gubbe := MakeGubbe()
	updateGubbe(&gubbe, 1, Controls{left: true, right: false, kick: false})
	approvals.VerifyString(t, printGubbe(gubbe))
}

func TestPressingBothLeftAndRightMeansStandingStill(t *testing.T) {
	gubbe := MakeGubbe()
	updateGubbe(&gubbe, 0.1, Controls{left: true, right: true, kick: false})
	approvals.VerifyString(t, printGubbe(gubbe))
}

func TestKicking(t *testing.T) {
	gubbe := MakeGubbe()
	updateGubbe(&gubbe, 0.1, Controls{left: false, right: false, kick: true})
	approvals.VerifyString(t, printGubbe(gubbe))
}

func TestTryingToMoveWhileKicking(t *testing.T) {
	gubbe := MakeGubbe()
	updateGubbe(&gubbe, 0.1, Controls{left: false, right: false, kick: true})
	updateGubbe(&gubbe, 0.3, Controls{left: false, right: true, kick: false})
	approvals.VerifyString(t, printGubbe(gubbe))
}

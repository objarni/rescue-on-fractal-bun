package main

import (
	"fmt"
	"github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/faiface/pixel"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	r := approvals.UseReporter(reporters.NewIntelliJReporter())
	defer r.Close()
	code := m.Run()
	os.Exit(code)
}

func printGubbe(gubbe Gubbe) string {
	return fmt.Sprintf("State: %s\nLooking: %s\n",
		gubbe.state.String(), gubbe.looking.String())
}

func MakeGubbe() Gubbe {
	return Gubbe{
		state:           Standing,
		looking:         Right,
		stateChangeTime: time.Time{},
		position:        pixel.Vec{},
	}
}

func TestWalkingRight(t *testing.T) {
	gubbe := MakeGubbe()
	updateGubbe(&gubbe, 0.1, Controls{left: false, right: true, kick: false})
	approvals.VerifyString(t, printGubbe(gubbe))
}

func TestWalkingLeft(t *testing.T) {
	gubbe := MakeGubbe()
	updateGubbe(&gubbe, 0.1, Controls{left: true, right: false, kick: false})
	approvals.VerifyString(t, printGubbe(gubbe))
}

func TestPressingBothLeftAndRightMeansStandingStill(t *testing.T) {
	gubbe := MakeGubbe()
	updateGubbe(&gubbe, 0.1, Controls{left: true, right: true, kick: false})
	approvals.VerifyString(t, printGubbe(gubbe))
}

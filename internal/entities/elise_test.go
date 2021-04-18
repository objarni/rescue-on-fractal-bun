package entities

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/tests"
	"strings"
	"testing"
)

/*
- initial state
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

func printElise(elise Entity) string {
	state := fmt.Sprintf("Elise %v", "standing")
	hb := fmt.Sprintf("HitBox %v", elise.HitBox())
	facing := fmt.Sprintf("Facing right")
	all := []string{state, hb, facing}
	return strings.Join(all, "\n") + "\n"
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

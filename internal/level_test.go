package internal

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/tests"
	"sort"
	"strings"
	"testing"
)

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

func TestMapSignBuilder(t *testing.T) {
	var levels = map[string]Level{
		"Hembyn": {
			SignPosts: []SignPost{
				{
					Pos:  pixel.Vec{X: 100, Y: 10},
					Text: "Hembyn",
				},
				{
					Pos:  pixel.Vec{X: 1000, Y: 10},
					Text: "Skogen",
				},
			},
		},
		"Korsningen": {
			SignPosts: []SignPost{
				{
					Pos:  pixel.Vec{X: 100, Y: 10},
					Text: "Hembyn",
				},
			},
		},
	}

	approvals.VerifyString(t, mapSignsToString(BuildMapSignArray(levels)))
}

func mapSignsToString(signs []MapSign) string {
	descriptions := make([]string, 0)
	printVec := func(vec pixel.Vec) string {
		return fmt.Sprintf("<%v, %v>", vec.X, vec.Y)
	}
	for ix, sign := range signs {
		mapSignDescription := ""
		mapSignDescription += fmt.Sprintf("MapSign %v:", ix+1)
		mapSignDescription += fmt.Sprintf(" Position on map %v", printVec(sign.MapPos))
		mapSignDescription += fmt.Sprintf(" Links to %v %v\n", sign.LevelName, printVec(sign.LevelPos))
		descriptions = append(descriptions, mapSignDescription)
	}
	fmt.Print(descriptions)
	var sorted = descriptions
	sort.Strings(sorted)
	fmt.Print(sorted)
	return strings.Join(sorted, "\n")
}

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

func TestBuildMapSignArray(t *testing.T) {
	var levels = map[string]Level{
		"GhostForest": {
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
		"ForestOpening": {
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
		mapSignDescription += fmt.Sprintf(" Link between %v '%v' %v", sign.LevelName, sign.Text, printVec(sign.LevelPos))
		mapSignDescription += fmt.Sprintf(" and map pos %v", printVec(sign.MapPos))
		descriptions = append(descriptions, mapSignDescription)
	}
	var sorted = descriptions
	sort.Strings(sorted)
	return strings.Join(sorted, "\n") + "\n"
}

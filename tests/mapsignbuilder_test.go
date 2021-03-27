package tests

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"sort"
	"strings"
	"testing"
)

func TestMapSignBuilder(t *testing.T) {
	var levels = map[string]internal.Level{
		"Hembyn": {
			SignPosts: []internal.SignPost{
				{
					Pos:      pixel.Vec{100, 10},
					Location: "Hembyn",
				},
				{
					Pos:      pixel.Vec{1000, 10},
					Location: "Skogen",
				},
			},
		},
		"Korsningen": {
			SignPosts: []internal.SignPost{
				{
					Pos:      pixel.Vec{100, 10},
					Location: "Hembyn",
				},
			},
		},
	}

	approvals.VerifyString(t, mapSignsToString(internal.BuildMapSignArray(levels)))
}

func mapSignsToString(signs []internal.MapSign) string {
	descriptions := []string{}
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
	sort.Strings(descriptions)
	return strings.Join(descriptions, "\n")
}

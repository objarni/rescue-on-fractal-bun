package internal

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"sort"
	"strings"
	"testing"
)

func TestMapSignBuilder(t *testing.T) {
	var levels = map[string]Level{
		"Hembyn": {
			SignPosts: []SignPost{
				{
					Pos:  pixel.Vec{100, 10},
					Text: "Hembyn",
				},
				{
					Pos:  pixel.Vec{1000, 10},
					Text: "Skogen",
				},
			},
		},
		"Korsningen": {
			SignPosts: []SignPost{
				{
					Pos:  pixel.Vec{100, 10},
					Text: "Hembyn",
				},
			},
		},
	}

	approvals.VerifyString(t, mapSignsToString(BuildMapSignArray(levels)))
}

func mapSignsToString(signs []MapSign) string {
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

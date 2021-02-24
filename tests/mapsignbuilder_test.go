package tests

import (
	"fmt"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
)

func Example_loadFromTwoSimpleLevels() {
	var levels = map[string]internal.Level{
		"Hembyn": {
			SignPosts: []internal.SignPost{
				{
					Pos:      pixel.Vec{100, 10},
					Location: "Hembyn",
				},
				{
					Pos:      pixel.Vec{1000, 10},
					Location: "Asbest",
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

	printMapSigns(internal.BuildMapSignArray(levels))
	// Output:
	// MapSign 1:
	//   Position on map <44,55>
	//   Links to Asbestzonen <100, 10>
	// MapSign 2:
	//   Position on map <144,155>
	//   Links to Hembyn <1000, 10>
	// MapSign 3:
	//   Position on map <244,255>
	//   Links to Korsningen <678, 11>
}

func printMapSigns(signs []internal.MapSign) {
	printVec := func(vec pixel.Vec) string {
		return fmt.Sprintf("<%v, %v>", vec.X, vec.Y)
	}
	for ix, sign := range signs {
		fmt.Printf("MapSign %v:\n", ix+1)
		lPos := sign.LevelPos
		fmt.Printf("  Position on map %v\n", printVec(sign.MapPos))
		fmt.Printf("  Links to %v %v\n", sign.LevelName, printVec(lPos))
	}
}

package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
)

func LocAt(x float64, y float64) internal.MapSign {
	return internal.MapSign{
		MapPos: pixel.Vec{X: x, Y: y},
	}
}

func ExampleFindNearLocation_tooFarAwayMeansNotClose() {
	ix := FindNearMapSign(
		pixel.Vec{X: 99, Y: 99},
		[]internal.MapSign{
			LocAt(1000, 1000),
		},
		30,
	)
	fmt.Println("Index of closest location:", ix)
	// Output:
	// Index of closest location: -1
}

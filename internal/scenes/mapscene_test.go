package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
)

func LocAt(x float64, y float64) MapPoint {
	return MapPoint{
		position:   pixel.Vec{X: x, Y: y},
		discovered: false,
	}
}

func ExampleFindNearLocation_tooFarAwayMeansNotClose() {
	ix := FindNearMapSign(
		pixel.Vec{X: 99, Y: 99},
		[]MapPoint{
			LocAt(1000, 1000),
		},
		30,
	)
	fmt.Println("Index of closest location:", ix)
	// Output:
	// Index of closest location: -1
}

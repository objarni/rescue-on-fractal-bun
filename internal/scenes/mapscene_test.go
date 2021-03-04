package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
)

func Example_findClosestLocation_SinglePoint() {
	ix := FindNearLocation(
		pixel.Vec{X: 50, Y: 50},
		[]MapPoint{
			LocAt(50, 50),
		},
		1000,
	)
	fmt.Printf("Index of closest location: %v", ix)
	// Output:
	// Index of closest location: 0
}

func LocAt(x float64, y float64) MapPoint {
	return MapPoint{
		position:   pixel.Vec{X: x, Y: y},
		discovered: false,
	}
}

func Example_findClosestLocation_TwoPoints() {
	ix := FindNearLocation(
		pixel.Vec{},
		[]MapPoint{
			LocAt(100, 100),
			LocAt(50, 50),
		},
		1000,
	)
	fmt.Printf("Index of closest location: %v", ix)
	// Output:
	// Index of closest location: 1
}

func Example_findClosestLocation_TwoPoints_first_is_closest() {
	ix := FindNearLocation(
		pixel.Vec{X: 99, Y: 99},
		[]MapPoint{
			LocAt(100, 100),
			LocAt(50, 50),
		},
		1000,
	)
	fmt.Printf("Index of closest location: %v", ix)
	// Output:
	// Index of closest location: 0
}

func Example_findClosestLocation_TooFarAwayMeansNotClose() {
	ix := FindNearLocation(
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

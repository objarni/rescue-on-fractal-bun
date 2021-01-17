package scenes

import (
	"fmt"
	"github.com/faiface/pixel"
)

func ExampleFindClosestLocation_SinglePoint() {
	ix := FindClosestLocation(
		pixel.Vec{50, 50},
		[]Location{
			LocAt(50, 50),
		},
		1000,
	)
	fmt.Printf("Index of closest location: %v", ix)
	// Output:
	// Index of closest location: 0
}

func LocAt(x float64, y float64) Location {
	return Location{
		position:   pixel.Vec{x, y},
		discovered: false,
	}
}

func ExampleFindClosestLocation_TwoPoints() {
	ix := FindClosestLocation(
		pixel.Vec{0, 0},
		[]Location{
			LocAt(100, 100),
			LocAt(50, 50),
		},
		1000,
	)
	fmt.Printf("Index of closest location: %v", ix)
	// Output:
	// Index of closest location: 1
}

func ExampleFindClosestLocation_TwoPoints_first_is_closest() {
	ix := FindClosestLocation(
		pixel.Vec{99, 99},
		[]Location{
			LocAt(100, 100),
			LocAt(50, 50),
		},
		1000,
	)
	fmt.Printf("Index of closest location: %v", ix)
	// Output:
	// Index of closest location: 0
}

func ExampleFindClosestLocation_TooFarAwayMeansNotClose() {
	ix := FindClosestLocation(
		pixel.Vec{99, 99},
		[]Location{
			LocAt(1000, 1000),
		},
		30,
	)
	fmt.Println("Index of closest location:", ix)
	// Output:
	// Index of closest location: -1
}

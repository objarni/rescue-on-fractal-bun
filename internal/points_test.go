package internal

import (
	"fmt"
	"github.com/faiface/pixel"
)

func ExampleClosestPoint() {
	ix := ClosestPoint(
		pixel.Vec{X: 99, Y: 99},
		[]pixel.Vec{
			{X: 100, Y: 100},
			{X: 50, Y: 50},
		},
	)
	fmt.Printf("Index of closest point: %v", ix)
	// Output:
	// Index of closest point: 0
}

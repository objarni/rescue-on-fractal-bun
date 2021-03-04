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

func ExampleClosestPoint_single_vec() {
	ix := ClosestPoint(
		pixel.Vec{X: 50, Y: 50},
		[]pixel.Vec{
			{X: 50, Y: 50},
		},
	)
	fmt.Printf("Index of closest point: %v", ix)
	// Output:
	// Index of closest point: 0
}

func ExampleClosestPoint_twoPoints() {
	ix := ClosestPoint(
		pixel.Vec{},
		[]pixel.Vec{
			{X: 100, Y: 100},
			{X: 50, Y: 50},
		},
	)
	fmt.Printf("Index of closest location: %v", ix)
	// Output:
	// Index of closest location: 1
}

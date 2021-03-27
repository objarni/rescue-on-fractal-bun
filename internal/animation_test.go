package internal

import (
	"fmt"
)

func Example_anim3Frames() {
	var timeSeconds float64
	for timeSeconds = 0; timeSeconds < 1.0; timeSeconds += 0.34 {
		fmt.Println(Animation{Frames: 3, TargetFPS: 3}.FrameAtTime(timeSeconds))
	}
	// Output:
	// 0
	// 1
	// 2
}

func Example_anim10Frames() {
	var timeSeconds float64
	for timeSeconds = 0; timeSeconds < 0.45; timeSeconds += 0.1 {
		fmt.Println(Animation{Frames: 10, TargetFPS: 10}.FrameAtTime(timeSeconds))
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

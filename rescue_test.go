package main

import "fmt"

/*
We want to have be able to progress gametime with
the help of computer real time clock, measured
in milliseconds, but in a frame-based manner.

We want frame-rate-independent animation, so that
on a screen resolution with 60 Hz the game plays
at same speed as on a 75 Hz, or if the graphics
stutters at some point, and drops to 25-30 FPS,
the speed of the game is unaffected (except the
user experience which will be worse, including
reaction time, but that is outside of scope).

The logical ticks per second is 200 in the game,
so regardless of 60, 75, 80, 100 or 140 Hz screen,
it should play almost identically.

Each frame, we check the computer runtime clock.
We convert the raw format into milliseconds since
app start. By keeping track of the previous
frame time, we get a ms delta; e.g of 14 ms.

Since we want the game to be 200 logical frames
per second, this would mean 14/5 = 2 ticks.

At start, we initialize a 'accumulator' (ms) to 0.
Since we took 2 ticks, we have 4 ms 'in store' and
add to the accumulator, now value 4.

Next frame, another 14 ms has passed. 14+4 = 18.
We get 3 ticks this time around, and 3 left in accumulator.

The game time isn't interesting to animations, for now.

But the number of ticks is!

- stepping through 3 frames with 200 logical fps, and deltas of
   5.1 , 4.1, 1.1
  should give
   1 ticks 0.1 leftover
   0 ticks 4.2 leftover
   1 tick 0.3 leftover
- stepping through 3 frames with 200 logical fps, and deltas 11,5,5
  should give 2, 1, 1 and 1 left over ms.
- stepping through 3 frames with 200 logical fps, and deltas 11,5,5
  should give 2, 1, 1 and 1 left over ms.
*/

func Example_scenarioOne() {
	lengthOfTickInMS := 5.0
	deltas := []float64{5.1, 4.1, 1.1}
	accumulator := 0.0
	ticks := 0
	for _, delta := range deltas {
		accumulator, ticks = gameTimeSteps(accumulator, delta, lengthOfTickInMS)
		fmt.Printf("%v ticks %1.1f leftover\n", ticks, accumulator)
	}

	// Output:
	// 1 ticks 0.1 leftover
	// 0 ticks 4.2 leftover
	// 1 ticks 0.3 leftover
}

func Example_scenarioOne_lengthOfTickInMS100() {
	lengthOfTickInMS := 100.0
	deltas := []float64{5.1, 4.1, 1.1, 100.0}
	accumulator := 0.0
	ticks := 0
	for _, delta := range deltas {
		accumulator, ticks = gameTimeSteps(accumulator, delta, lengthOfTickInMS)
		fmt.Printf("%v ticks %1.1f leftover\n", ticks, accumulator)
	}

	// Output:
	// 0 ticks 5.1 leftover
	// 0 ticks 9.2 leftover
	// 0 ticks 10.3 leftover
	// 1 ticks 10.3 leftover
}

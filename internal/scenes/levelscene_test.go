package scenes

import (
	"fmt"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
)

type LevelSimulation struct {
	level               internal.Level
	playerStartMapPoint string
}

func (sim *LevelSimulation) HandleKeyDown(_ internal.ControlKey) {

}

func (sim *LevelSimulation) Tick() {

}

func ExampleLevelSimulation() {
	levelSim := LevelSimulation{
		level: internal.Level{
			Width:      3000,
			Height:     1000,
			ClearColor: colornames.White,
			SignPosts: []internal.SignPost{
				{
					Pos:        v(2700, 900),
					Discovered: true,
					Location:   "Hembyn",
				},
				{
					Pos:        v(100, 900),
					Discovered: false,
					Location:   "Korsningen",
				},
			},
		},
		playerStartMapPoint: "Hembyn",
	}
	printLevelSimulation(&levelSim)
	levelSim.HandleKeyDown(internal.Right)
	ticks := 100
	i := 0
	for i < ticks {
		levelSim.Tick()
		i++
	}
	// Output:
	// Level is 3000x1000 pixels big.
	// Background color is 255,255,255.
	// Player starts at 2700,900.
	// There are SignPosts here:
	//   2700,900 (Discovered)
	//   100,900 (not Discovered)
	// Simulation started.
	// Key Right Down
	// Tick x100
	// Simulation stopped. State:
	// Player is at 3000, 900.
}

func printLevelSimulation(sim *LevelSimulation) {
	fmt.Printf(`Level is %vx%v pixels big.
Background color is %v.
Player starts at 2700,900.
There are SignPosts here:
  2700,900 (Discovered)
  100,900 (not Discovered)
Simulation started.
Key Right Down
Tick x100
Simulation stopped. State:
Player is at 3000, 900.
`, sim.level.Width, sim.level.Height,
		projectColor(sim.level.ClearColor))
}

func projectColor(color color.RGBA) string {
	return fmt.Sprintf("%v,%v,%v",
		color.R, color.G, color.B)
}

/*
some scenarios
when entering level on "Hembyn", player starts
 at map point named "Hembyn"
when clicking jump over a map position (close enough),
 exits to map scene with hair cross at that location
when walking across a level, the camera pans
 (i.e player is centered on screen)
when walking to end of level, the camera stops panning,
 instead the player approaches side of screen
*/

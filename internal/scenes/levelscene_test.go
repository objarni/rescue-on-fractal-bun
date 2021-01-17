package scenes

import (
	"fmt"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
)

type LevelSimulation struct {
	level               Level
	playerStartMapPoint string
}

func (sim *LevelSimulation) HandleKeyDown(key internal.ControlKey) {

}

func (sim *LevelSimulation) Tick() {

}

func ExampleWalkingRight100Ticks() {
	levelSim := LevelSimulation{
		level: Level{
			width:      3000,
			height:     1000,
			clearColor: colornames.White,
			mapPoints: []MapPoint{
				{
					pos:        v(2700, 900),
					discovered: true,
					mapTarget:  "Hembyn",
				},
				{
					pos:        v(100, 900),
					discovered: false,
					mapTarget:  "Korsningen",
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
	// There are MapPoints here:
	//   2700,900 (discovered)
	//   100,900 (not discovered)
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
There are MapPoints here:
  2700,900 (discovered)
  100,900 (not discovered)
Simulation started.
Key Right Down
Tick x100
Simulation stopped. State:
Player is at 3000, 900.
`, sim.level.width, sim.level.height,
		projectColor(sim.level.clearColor))
}

func projectColor(color color.RGBA) string {
	return fmt.Sprintf("%v,%v,%v",
		color.R, color.G, color.B)
}

package scenes

import (
	"fmt"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
)

/*
- empty scenes gives elise and levelboundary
- scene with 1 spider gives - " -     and spider
*/

func Example_noEntitiesGivesError() {
	pos := px.V(0, 0)
	level := internal.Level{
		EntitySpawnPoints: nil,
	}
	_, err := SpawnEntities(pos, level)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
	// boom
}

func Example_emptyLevelHiFi() {
	pos := px.V(0, 0)
	level := internal.Level{
		EntitySpawnPoints: []internal.EntitySpawnPoint{},
	}
	res, _ := SpawnEntities(pos, level)
	for _, esp := range res {
		fmt.Println(esp.String())
	}
	// Output:
	// Elise Standing
	// Vel: <0.0,0.0>
	// Facing right
	// Gfx:
	// Moved 0 pixels right 50 pixels up:
	//   Image "IEliseWalk2"
	//
	// Level Boundary
	//
}

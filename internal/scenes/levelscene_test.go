package scenes

import (
	"fmt"
	px "github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
	"reflect"
)

/*
- empty scenes gives elise and levelboundary
- scene with 1 spider gives - " -     and spider
*/

func Example_emptyLevel() {
	pos := px.V(0, 0)
	level := internal.Level{
		EntitySpawnPoints: nil,
	}
	res := SpawnEntities(pos, level)
	for _, esp := range res {
		fmt.Println(reflect.TypeOf(esp))
	}
	// Output:
	// entities.Elise
	// entities.LevelBoundary
}

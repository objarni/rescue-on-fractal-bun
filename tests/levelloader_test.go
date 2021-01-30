package tests

import (
	"fmt"
	"objarni/rescue-on-fractal-bun/internal"
	"strings"
)

func templateThis(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

func printLevel(level internal.Level) {
	mapPoints := printMapPoints(level.MapPoints)
	fmt.Println(templateThis("Width: {w}   Height: {h}  (tiles)\n"+
		"There are {countMapPoints} MapPoint(s):\n"+
		"{mapPoints}\n"+
		"Walls:\n"+
		"...#\n"+
		"...#\n"+
		"....\n"+
		"Platforms:\n"+
		"....\n"+
		"....\n"+
		"####",
		"{w}", toString(level.Width),
		"{h}", toString(level.Height),
		"{countMapPoints}", toString(len(level.MapPoints)),
		"{mapPoints}", mapPoints,
	))
}

func printMapPoints(points []internal.MapPoint) string {
	s := ""
	for _, mp := range points {
		s += fmt.Sprintf("'%v' at %1.0f, %1.0f", mp.Location, mp.Pos.X, mp.Pos.Y)
	}
	return s
}

func toString(v interface{}) string {
	return fmt.Sprint(v)
}

func ExampleLoadingMiniLevel() {
	level := internal.LoadLevel("../testdata/MiniLevel.tmx")
	printLevel(level)
	// Output:
	// Width: 4   Height: 3  (tiles)
	// There are 1 MapPoint(s):
	// 'Korsningen' at 11, 56
	// Walls:
	// ...#
	// ...#
	// ....
	// Platforms:
	// ....
	// ....
	// ####
}

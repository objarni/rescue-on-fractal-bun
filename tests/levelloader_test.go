package tests

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
	"strings"
)

type Level struct {
	Width, Height int
	MapPoints     []scenes.MapPoint
}

func LoadLevel(path string) Level {
	level, err := tilepix.ReadFile(path)
	internal.PanicIfError(err)
	return Level{
		Width:     level.Width,
		Height:    level.Height,
		MapPoints: nil,
	}
}

func templateThis(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

func printLevel(level Level) {
	fmt.Println(templateThis("Width: {w}   Height: {h}  (tiles)\n"+
		"There are {countMapPoints} MapPoint(s):\n"+
		"'Korsningen' at 40, 40\n"+
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
	))
}

func toString(v interface{}) string {
	return fmt.Sprint(v)
}

func ExampleLoadingMiniLevel() {
	level := LoadLevel("../testdata/MiniLevel.tmx")
	printLevel(level)
	// Output:
	// Width: 4   Height: 3  (tiles)
	// There is 1 MapPoint:
	// 'Korsningen' at 40, 40
	// Walls:
	// ...#
	// ...#
	// ....
	// Platforms:
	// ....
	// ....
	// ####
}

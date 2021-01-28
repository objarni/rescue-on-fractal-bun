package tests

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
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
	points := []scenes.MapPoint{}
	for _, object := range level.ObjectGroups[0].Objects {
		x := object.X
		y := object.Y
		var mp = scenes.MapPoint{
			Pos:        pixel.Vec{x, y},
			Discovered: false,
			Location:   object.Name,
		}
		points = append(points, mp)
	}
	return Level{
		Width:     level.Width,
		Height:    level.Height,
		MapPoints: points,
	}
}

func templateThis(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

func printLevel(level Level) {
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

func printMapPoints(points []scenes.MapPoint) string {
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
	level := LoadLevel("../testdata/MiniLevel.tmx")
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

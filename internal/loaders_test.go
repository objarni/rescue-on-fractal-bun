package internal

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"strings"
)

func templateThis(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

func printLevel(level Level) {
	mapPoints := printMapPoints(level.SignPosts)
	fmt.Println(templateThis(
		"GetWidth: {Width}   GetHeight: {Height}  (tiles)\n"+
			"Background color: RGB={red},{green},{blue}\n"+
			"There are {countMapPoints} SignPost(s):\n"+
			"{mapPoints}\n"+
			"Walls:\n"+
			"...#\n"+
			"...#\n"+
			"....\n"+
			"Platforms:\n"+
			"....\n"+
			"....\n"+
			"####",
		"{Width}", toString(level.Width),
		"{Height}", toString(level.Height),
		"{countMapPoints}", toString(len(level.SignPosts)),
		"{mapPoints}", mapPoints,
		"{red}", toString(level.ClearColor.R),
		"{green}", toString(level.ClearColor.G),
		"{blue}", toString(level.ClearColor.B),
	))
}

func printMapPoints(points []SignPost) string {
	s := ""
	for _, mp := range points {
		s += fmt.Sprintf("'%v' at %1.0f, %1.0f\n", mp.Location, mp.Pos.X, mp.Pos.Y)
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
	// GetWidth: 4   GetHeight: 3  (tiles)
	// Background color: RGB=10,50,100
	// There are 2 SignPost(s):
	// 'Korsningen' at 12, 62
	// 'Hembyn' at 10, 41
	//
	// Walls:
	// ...#
	// ...#
	// ....
	// Platforms:
	// ....
	// ....
	// ####
}

func ExampleLoadingBrokenLevel() {
	brokenLevelPath := "../testdata/BrokenLevel.tmx"
	brokenLevel, _ := tilepix.ReadFile(brokenLevelPath)
	ValidateLevel(brokenLevelPath, brokenLevel)
	// Output:
	// ../testdata/BrokenLevel.tmx contains the following errors:
	// There is no Background layer
	// There is no Platforms layer
	// There is no Walls layer
	// There is no Foreground layer
	// There should be an object layer named "SignPosts", instead I found:
	// "Object Layer 1"
	// The BackgroundColor should be on web-color format #RRGGBB, instead I found:
	// ""
}

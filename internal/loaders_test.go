package internal

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/bcvery1/tilepix"
	"objarni/rescue-on-fractal-bun/tests"
	"strings"
	"testing"
)

func templateThis(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

func levelToString(level Level) string {
	mapPoints := mapPointsToString(level.SignPosts)
	entities := entitySpawnPointsToString(level.EntitySpawnPoints)
	levelString := templateThis(
		"GetWidth: {Width}   GetHeight: {Height}  (tiles)\n"+
			"Background color: RGB={red},{green},{blue}\n"+
			"There are {countMapPoints} SignPost(s):\n"+
			"{mapPoints}\n"+
			"Entities:\n"+
			"{entities}\n"+
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
		"{entities}", entities,
		"{red}", toString(level.ClearColor.R),
		"{green}", toString(level.ClearColor.G),
		"{blue}", toString(level.ClearColor.B),
	)
	return levelString
}

func entitySpawnPointsToString(spawnPoints []EntitySpawnPoint) string {
	result := ""
	for _, esp := range spawnPoints {
		result += fmt.Sprintf("'%v' at %v\n", esp.EntityType, esp.SpawnAt)
	}
	return result
}

func mapPointsToString(points []SignPost) string {
	s := ""
	for _, mp := range points {
		s += fmt.Sprintf("'%v' at %1.0f, %1.0f\n", mp.Text, mp.Pos.X, mp.Pos.Y)
	}
	return s
}

func toString(v interface{}) string {
	return fmt.Sprint(v)
}

func Test_loadingSimpleButCompleteLevel(t *testing.T) {
	level := LoadLevel("../testdata/MiniLevel.tmx")
	approvals.VerifyString(t, levelToString(level)+"\n")
}

func Example_loadingBrokenLevel() {
	brokenLevelPath := "../testdata/BrokenLevel.tmx"
	brokenLevel, _ := tilepix.ReadFile(brokenLevelPath)
	_ = ValidateLevel(brokenLevelPath, brokenLevel)
	// Output:
	// ../testdata/BrokenLevel.tmx contains the following errors:
	// There is no Background layer
	// There is no Platforms layer
	// There is no Walls layer
	// There is no Foreground layer
	// There should be an object layer named "SignPosts", instead I found:
	// "Object Layer 1"
	// There should be an object layer named "Entities", instead I found:
	// "Object Layer 1"
	// The BackgroundColor should be on web-color format #RRGGBB, instead I found:
	// ""
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

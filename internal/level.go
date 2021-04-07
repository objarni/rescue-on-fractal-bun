package internal

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"image/color"
)

type Level struct {
	Width, Height     int
	ClearColor        color.RGBA
	SignPosts         []SignPost
	TilepixMap        *tilepix.Map
	EntitySpawnPoints []EntitySpawnPoint
}

type SignPost struct {
	Pos  pixel.Vec
	Text string // Text on signpost
}

type EntitySpawnPoint struct {
	spawnAt    pixel.Vec
	entityType string
}

type MapSign struct {
	MapPos    pixel.Vec // X,Y coordinate on map
	LevelName string    // Name of level where MapSign stands
	LevelPos  pixel.Vec // X,Y coordinate on level
}

func BuildMapSignArray(levelMap map[string]Level) []MapSign {
	var mapSigns = []MapSign{}

	var positions = map[string]pixel.Vec{
		"Hembyn":     {X: 246, Y: 109},
		"Korsningen": {X: 355, Y: 235},
		"Skogen":     {X: 299, Y: 375},
	}
	for levelName, levelData := range levelMap {
		fmt.Println(levelName)
		for _, signPost := range levelData.SignPosts {
			mapSigns = append(mapSigns, MapSign{
				MapPos:    positions[signPost.Text],
				LevelPos:  signPost.Pos,
				LevelName: levelName,
			})
		}
	}
	return mapSigns
}

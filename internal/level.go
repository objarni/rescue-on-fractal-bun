package internal

import (
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
	SpawnAt    pixel.Rect
	EntityType string
}

type MapSign struct {
	MapPos    pixel.Vec // X,Y coordinate on map
	LevelName string    // Name of level where MapSign stands
	LevelPos  pixel.Vec // X,Y coordinate on level
	Text      string    // Text on sign, displayed to player
}

func BuildMapSignArray(levelMap map[string]Level) []MapSign {
	var mapSigns = make([]MapSign, 0)
	// TODO: move this data to .json file, or even Tiled 'somehow'
	var mapSignPositions = map[string]pixel.Vec{
		"Hembyn":     {X: 246, Y: 109},
		"Korsningen": {X: 355, Y: 235},
		"Skogen":     {X: 299, Y: 375},
	}
	for levelName, levelData := range levelMap {
		for _, signPost := range levelData.SignPosts {
			mapSigns = append(mapSigns, MapSign{
				MapPos:    mapSignPositions[signPost.Text],
				LevelName: levelName,
				LevelPos:  signPost.Pos,
				Text:      signPost.Text,
			})
		}
	}
	return mapSigns
}

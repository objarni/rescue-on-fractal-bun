package internal

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"image/color"
)

type Level struct {
	Width, Height int
	ClearColor    color.RGBA
	SignPosts     []SignPost
	TilepixMap    *tilepix.Map
}

// TODO: Discovered should probably be stored somewhere else
type SignPost struct {
	Pos        pixel.Vec
	Discovered bool
	Location   string
}

type MapSign struct {
	MapPos    pixel.Vec // X,Y coordinate on map image
	LevelPos  pixel.Vec // X,Y coordinate on tiled map
	LevelName string    // Name of level where MapSign stands
}

func BuildMapSignArray(levelMap map[string]Level) []MapSign {
	var mapSigns = []MapSign{}

	var positions = map[string]pixel.Vec{
		"Hembyn":     {X: 246, Y: 109},
		"Korsningen": {X: 355, Y: 235},
		"Skogen":     {X: 299, Y: 375},
	}
	for levelName, levelData := range levelMap {
		mapSigns = append(mapSigns, MapSign{
			MapPos:    positions[levelName],
			LevelPos:  levelData.SignPosts[0].Pos,
			LevelName: levelData.SignPosts[0].Location,
		})
	}
	return mapSigns
}

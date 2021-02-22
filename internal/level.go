package internal

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"image/color"
)

type Level struct {
	Width, Height int
	ClearColor    color.RGBA
	MapSigns      []MapPoint
	TilepixMap    *tilepix.Map
}

// TODO: Discovered should probably be stored somewhere else
type MapPoint struct {
	Pos        pixel.Vec
	Discovered bool
	Location   string
}

type MapSign struct {
	MapPos    pixel.Vec // X,Y coordinate on map image
	LevelPos  pixel.Vec // X,Y coordinate on tiled map
	LevelName string    // Name of level where MapSign stands
}

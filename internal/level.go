package internal

import (
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"image/color"
)

type Level struct {
	Width, Height int
	ClearColor    color.RGBA
	MapPoints     []MapPoint
	TilepixMap    *tilepix.Map
}

// TODO: Discovered should probably be stored somewhere else
type MapPoint struct {
	Pos        pixel.Vec
	Discovered bool
	Location   string
}

package internal

import (
	"github.com/faiface/beep"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type ImageMap map[Image]*pixel.Sprite

type Resources struct {
	Atlas    *text.Atlas
	Blip     *beep.Buffer
	FPS      float64
	MapSigns []MapSign
	Levels   map[string]Level
	ImageMap ImageMap
}

// This is so repetitive. Can we have a Python script to generate the enum+String func?
type Image int

const (
	IMap Image = iota
	IGhost
	IMapSymbol
	ISignPost
	ITemporaryPlayerImage
	IEliseWalk1
	IEliseWalk2
	IEliseWalk3
	IEliseWalk4
	IEliseWalk5
	IEliseWalk6
	IEliseCrouch
	AfterLastImage
)

func (image Image) String() string {
	return [...]string{
		"IMap",
		"IGhost",
		"IMapSymbol",
		"ISignPost",
		"ITemporaryPlayerImage",
		"IEliseWalk1",
		"IEliseWalk2",
		"IEliseWalk3",
		"IEliseWalk4",
		"IEliseWalk5",
		"IEliseWalk6",
		"IEliseCrouch",
		"AfterLastImage",
	}[image]
}

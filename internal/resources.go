package internal

import (
	"github.com/faiface/beep"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type Resources struct {
	Atlas          *text.Atlas
	Ghost          *pixel.Sprite
	MapPoint       *pixel.Sprite
	PlayerStanding *pixel.Sprite
	MapSymbol      *pixel.Sprite
	Blip           *beep.Buffer
	FPS            float64
	MapSigns       []MapSign
	Levels         map[string]Level
	ImageMap       map[Image]*pixel.Sprite
}

type Image int

const (
	Map Image = iota
	Ghost
	MapSymbol
	AfterLastImage
)

func (image Image) String() string {
	return [...]string{"Map", "Ghost", "MapSymbol"}[image]
}

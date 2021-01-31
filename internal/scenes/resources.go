package scenes

import (
	"github.com/faiface/beep"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/text"
)

type Resources struct {
	Atlas    *text.Atlas
	Ghost    *pixel.Sprite
	MapPoint *pixel.Sprite
	Blip     *beep.Buffer
}

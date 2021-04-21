package entities

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	d "objarni/rescue-on-fractal-bun/internal/draw"
)

type SpiderState int

const (
	AtBottom = iota
	GoingUp
	AtTop
	GoingDown
)
const spiderWidth = 30
const spiderHeight = 40

type Spider struct {
	pos      pixel.Vec
	min, max float64
	state    SpiderState
	timeout  float64
}

func (spider Spider) String() string {
	panic("implement me")
}

func (spider Spider) Handle(_ EventBox) Entity {
	return spider
}

func (spider Spider) HitBox() pixel.Rect {
	halfSize := pixel.V(spiderWidth/2, spiderHeight/2)
	min := spider.pos.Sub(halfSize)
	max := spider.pos.Add(halfSize)
	rect := pixel.Rect{Min: min, Max: max}
	return rect
}

func (spider Spider) Tick(gameTimeMs float64, _ EventBoxReceiver) Entity {
	movement := pixel.V(0, 0.1)
	pauseMs := 3000.0

	switch spider.state {
	case GoingUp:
		spider.pos = spider.pos.Add(movement)
		if spider.pos.Y >= spider.max {
			spider.timeout = gameTimeMs + pauseMs
			spider.state = AtTop
		}
	case AtTop:
		if gameTimeMs >= spider.timeout {
			spider.state = GoingDown
		}
	case GoingDown:
		spider.pos = spider.pos.Sub(movement)
		if spider.pos.Y <= spider.min {
			spider.timeout = gameTimeMs + pauseMs
			spider.state = AtBottom
		}
	case AtBottom:
		if gameTimeMs >= spider.timeout {
			spider.state = GoingUp
		}
	}
	return spider
}

func (spider Spider) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	image := d.Image(*imageMap, internal.ISpider)
	image = d.Color(colornames.Black, image)
	return d.Moved(spider.pos, image)
}

func MakeSpider(area pixel.Rect) Entity {
	return Spider{
		pos:   area.Center(),
		min:   area.Min.Y,
		max:   area.Max.Y,
		state: GoingUp,
	}
}

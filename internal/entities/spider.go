package entities

import (
	"github.com/faiface/pixel"
	d "github.com/objarni/pixelop"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
)

type SpiderState int

const (
	SpiderAtBottom = iota
	SpiderGoingUp
	SpiderAtTop
	SpiderGoingDown
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
	case SpiderGoingUp:
		spider.pos = spider.pos.Add(movement)
		if spider.pos.Y >= spider.max {
			spider.timeout = gameTimeMs + pauseMs
			spider.state = SpiderAtTop
		}
	case SpiderAtTop:
		if gameTimeMs >= spider.timeout {
			spider.state = SpiderGoingDown
		}
	case SpiderGoingDown:
		spider.pos = spider.pos.Sub(movement)
		if spider.pos.Y <= spider.min {
			spider.timeout = gameTimeMs + pauseMs
			spider.state = SpiderAtBottom
		}
	case SpiderAtBottom:
		if gameTimeMs >= spider.timeout {
			spider.state = SpiderGoingUp
		}
	}
	return spider
}

func (spider Spider) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	image := d.Image((*imageMap)[internal.ISpider], internal.ISpider.String())
	image = d.Color(colornames.White, image)
	return d.Moved(spider.pos, image)
}

func MakeSpider(area pixel.Rect) Entity {
	return Spider{
		pos:   area.Center(),
		min:   area.Min.Y,
		max:   area.Max.Y,
		state: SpiderGoingUp,
	}
}

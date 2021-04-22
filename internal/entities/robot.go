package entities

import (
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	d "objarni/rescue-on-fractal-bun/internal/draw"
	"objarni/rescue-on-fractal-bun/internal/events"
)

const robotWidth = 30
const robotHeight = 40

type Robot struct {
	pos      pixel.Vec
	min, max float64
	state    SpiderState
	timeout  float64
}

func (robot Robot) String() string {
	panic("implement me")
}

func (robot Robot) Handle(_ EventBox) Entity {
	return robot
}

func (robot Robot) HitBox() pixel.Rect {
	halfSize := pixel.V(robotWidth/2, robotHeight/2)
	min := robot.pos.Sub(halfSize)
	max := robot.pos.Add(halfSize)
	rect := pixel.Rect{Min: min, Max: max}
	return rect
}

func (robot Robot) Tick(gameTimeMs float64, ebr EventBoxReceiver) Entity {
	movement := pixel.V(0, 0.1)
	pauseMs := 3000.0

	switch robot.state {
	case GoingUp:
		robot.pos = robot.pos.Add(movement)
		if robot.pos.Y >= robot.max {
			robot.timeout = gameTimeMs + pauseMs
			robot.state = AtTop
		}
	case AtTop:
		if gameTimeMs >= robot.timeout {
			robot.state = GoingDown
			ebr.AddEventBox(EventBox{
				Event: events.RobotMove,
				Box:   pixel.Rect{},
			})
		}
	case GoingDown:
		robot.pos = robot.pos.Sub(movement)
		if robot.pos.Y <= robot.min {
			robot.timeout = gameTimeMs + pauseMs
			robot.state = AtBottom
		}
	case AtBottom:
		if gameTimeMs >= robot.timeout {
			robot.state = GoingUp
			ebr.AddEventBox(EventBox{
				Event: events.RobotMove,
				Box:   pixel.Rect{},
			})
		}
	}
	return robot
}

func (robot Robot) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	image := d.Image(*imageMap, internal.IRobot1)
	image = d.Color(colornames.White, image)
	return d.Moved(robot.pos, image)
}

func MakeRobot(area pixel.Rect) Entity {
	return Robot{
		pos:   area.Center(),
		min:   area.Min.Y,
		max:   area.Max.Y,
		state: GoingUp,
	}
}

package entities

import (
	"github.com/faiface/pixel"
	d "github.com/objarni/pixelop"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/events"
)

const robotWidth = 30
const robotHeight = 40

type RobotState int

const (
	RobotAtLeft = iota
	RobotGoingLeft
	RobotAtRight
	RobotGoingRight
)

type Robot struct {
	pos      pixel.Vec
	min, max float64
	state    RobotState
	timeout  float64
}

func (robot Robot) HitBoxes() []pixel.Rect {
	panic("implement me")
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
	movement := pixel.V(0.3, 0)
	pauseMs := 1000.0

	switch robot.state {
	case RobotGoingRight:
		robot.pos = robot.pos.Add(movement)
		if robot.pos.X >= robot.max {
			robot.timeout = gameTimeMs + pauseMs
			robot.state = RobotAtRight
		}
	case RobotAtRight:
		if gameTimeMs >= robot.timeout {
			robot.state = RobotGoingLeft
			ebr.AddEventBox(EventBox{
				Event: events.RobotMove,
				Box:   pixel.Rect{},
			})
		}
	case RobotGoingLeft:
		robot.pos = robot.pos.Sub(movement)
		if robot.pos.X <= robot.min {
			robot.timeout = gameTimeMs + pauseMs
			robot.state = RobotAtLeft
		}
	case RobotAtLeft:
		if gameTimeMs >= robot.timeout {
			robot.state = RobotGoingRight
			ebr.AddEventBox(EventBox{
				Event: events.RobotMove,
				Box:   pixel.Rect{},
			})
		}
	}
	return robot
}

func (robot Robot) GfxOp(imageMap *internal.ImageMap) d.WinOp {
	// TODO: switch imageMap type to Image -> d.ImageOp
	image := d.Image((*imageMap)[internal.IRobot1], internal.IRobot1.String())
	image = d.Color(colornames.White, image)
	if robot.state == RobotGoingLeft || robot.state == RobotAtLeft {
		image = d.Mirrored(image)
	}
	return d.Moved(robot.pos, image)
}

func MakeRobot(area pixel.Rect) Entity {
	return Robot{
		pos:   area.Center(),
		min:   area.Min.X,
		max:   area.Max.X,
		state: RobotGoingRight,
	}
}

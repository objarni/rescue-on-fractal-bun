package draw

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type ImdCircle struct {
	radius, x, y, thickness int
}

func (circle ImdCircle) String() string {
	return fmt.Sprintf("Circle radius %v center <%v, %v> thickness %v",
		circle.radius,
		circle.x, circle.y,
		circle.thickness)
}

func Circle(radius int, x int, y int, thickness int) ImdOp {
	return ImdCircle{
		radius:    radius,
		x:         x,
		y:         y,
		thickness: thickness,
	}
}

func (circle ImdCircle) Render(imd *imdraw.IMDraw) {
	imd.Push(pixel.Vec{X: float64(circle.x), Y: float64(circle.y)})
	imd.Circle(float64(circle.radius), float64(circle.thickness))
}

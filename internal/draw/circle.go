package draw

import (
	"fmt"
	"github.com/faiface/pixel/imdraw"
)

type ImdCircle struct {
	radius, thickness int
	center            Coordinate
}

func (circle ImdCircle) String() string {
	return fmt.Sprintf("Circle radius %v center %v thickness %v",
		circle.radius,
		circle.center.String(),
		circle.thickness)
}

func Circle(radius int, center Coordinate, thickness int) ImdOp {
	return ImdCircle{
		radius:    radius,
		thickness: thickness,
		center:    center,
	}
}

func (circle ImdCircle) Render(imd *imdraw.IMDraw) {
	imd.Push(circle.center.toVec())
	imd.Circle(float64(circle.radius), float64(circle.thickness))
}

func (circle ImdCircle) Lines() []string {
	return []string{circle.String()}
}

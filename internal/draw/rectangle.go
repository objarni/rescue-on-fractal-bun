package draw

import (
	"fmt"
	"github.com/faiface/pixel/imdraw"
)

type ImdRectangle struct {
	from, to  Coordinate
	thickness int
}

func (rectangle ImdRectangle) String() string {
	var rectangleStyle string
	rectangleStyle = "(filled)"
	if rectangle.thickness != 0 {
		rectangleStyle = fmt.Sprintf("thickness %v", rectangle.thickness)
	}
	return fmt.Sprintf("Rectangle from %v to %v %v",
		rectangle.from.String(),
		rectangle.to.String(),
		rectangleStyle)
}

func Rectangle(from Coordinate, to Coordinate, thickness int) ImdOp {
	return ImdRectangle{from: from, to: to, thickness: thickness}
}

func (rectangle ImdRectangle) Render(imd *imdraw.IMDraw) {
	imd.Push(rectangle.from.toVec())
	imd.Push(rectangle.to.toVec())
	imd.Rectangle(float64(rectangle.thickness))
}

func (rectangle ImdRectangle) Lines() []string {
	return []string{rectangle.String()}
}

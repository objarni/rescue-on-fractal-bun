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

func (line ImdRectangle) Render(imd *imdraw.IMDraw) {
	imd.Push(line.from.toVec())
	imd.Push(line.to.toVec())
	imd.Rectangle(float64(line.thickness))
}

func (line ImdRectangle) Lines() []string {
	return []string{line.String()}
}

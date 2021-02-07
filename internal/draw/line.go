package draw

import (
	"fmt"
	"github.com/faiface/pixel/imdraw"
)

type ImdLine struct {
	from, to  Coordinate
	thickness int
}

func (line ImdLine) String() string {
	return fmt.Sprintf("Line from %v to %v thickness %v",
		line.from.String(),
		line.to.String(),
		line.thickness)
}

func Line(from Coordinate, to Coordinate, thickness int) ImdOp {
	return ImdLine{from: from, to: to, thickness: thickness}
}

func (line ImdLine) Render(imd *imdraw.IMDraw) {
	imd.Push(line.from.toVec())
	imd.Push(line.to.toVec())
	imd.Line(float64(line.thickness))
}

func (line ImdLine) Lines() []string {
	return []string{line.String()}
}

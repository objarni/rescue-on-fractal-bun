package draw

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
)

type ImdOp interface {
	String() string
	Render(imd *imdraw.IMDraw)
}

// Types
type ImdCircle struct {
	radius, x, y, thickness int
}

type ImdColor struct {
	color     color.RGBA
	Operation ImdOp
}

type ImdSequence struct {
	imdOps []ImdOp
}

// Projectors
func (circle ImdCircle) String() string {
	return fmt.Sprintf("Circle radius %v center <%v, %v> thickness %v",
		circle.radius,
		circle.x, circle.y,
		circle.thickness)
}

func (color ImdColor) String() string {
	head := fmt.Sprintf("Color %v, %v, %v:\n  ",
		color.color.R, color.color.G, color.color.B)
	body := color.Operation.String()
	return head + body
}

func (sequence ImdSequence) String() string {
	result := "Sequence:\n"
	for _, imdOp := range sequence.imdOps {
		result += "  " + imdOp.String() + "\n"
	}
	return result
}

// Builders
func Colored(color color.RGBA, imdOp ImdOp) ImdOp {
	return ImdColor{
		color:     color,
		Operation: imdOp,
	}
}

func Circle(radius int, x int, y int, thickness int) ImdOp {
	return ImdCircle{
		radius:    radius,
		x:         x,
		y:         y,
		thickness: thickness,
	}
}

func Sequence(imdOps ...ImdOp) ImdOp {
	return ImdSequence{
		imdOps: imdOps,
	}
}

// Renderers
func (circle ImdCircle) Render(imd *imdraw.IMDraw) {
	imd.Push(pixel.Vec{X: float64(circle.x), Y: float64(circle.y)})
	imd.Circle(float64(circle.radius), float64(circle.thickness))
}

func (color ImdColor) Render(imd *imdraw.IMDraw) {
	// TODO: do we want to reset color to previous state?
	imd.Color = color.color
	color.Operation.Render(imd)
}

func (sequence ImdSequence) Render(imd *imdraw.IMDraw) {
	for _, imdOp := range sequence.imdOps {
		imdOp.Render(imd)
	}
}

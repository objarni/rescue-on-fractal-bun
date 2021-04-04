package draw

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"image/color"
	"strings"
)

type ImdCircle struct {
	radius, thickness int
	center            pixel.Vec
}

func C(x, y float64) pixel.Vec {
	return pixel.V(x, y)
}

func (circle ImdCircle) String() string {
	return fmt.Sprintf("Circle radius %v center %v thickness %v",
		circle.radius,
		circle.center.String(),
		circle.thickness)
}

func Circle(radius int, center pixel.Vec, thickness int) ImdOp {
	return ImdCircle{
		radius:    radius,
		thickness: thickness,
		center:    center,
	}
}

func (circle ImdCircle) Render(imd *imdraw.IMDraw) {
	imd.Push(circle.center)
	imd.Circle(float64(circle.radius), float64(circle.thickness))
}

func (circle ImdCircle) Lines() []string {
	return []string{circle.String()}
}

type ImdLine struct {
	from, to  pixel.Vec
	thickness int
}

func (line ImdLine) String() string {
	return fmt.Sprintf("Line from %v to %v thickness %v",
		line.from.String(),
		line.to.String(),
		line.thickness)
}

func Line(from pixel.Vec, to pixel.Vec, thickness int) ImdOp {
	return ImdLine{from: from, to: to, thickness: thickness}
}

func (line ImdLine) Render(imd *imdraw.IMDraw) {
	imd.Push(line.from)
	imd.Push(line.to)
	imd.Line(float64(line.thickness))
}

func (line ImdLine) Lines() []string {
	return []string{line.String()}
}

type ImdRectangle struct {
	from, to  pixel.Vec
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

func Rectangle(from pixel.Vec, to pixel.Vec, thickness int) ImdOp {
	return ImdRectangle{from: from, to: to, thickness: thickness}
}

func (rectangle ImdRectangle) Render(imd *imdraw.IMDraw) {
	imd.Push(rectangle.from)
	imd.Push(rectangle.to)
	imd.Rectangle(float64(rectangle.thickness))
}

func (rectangle ImdRectangle) Lines() []string {
	return []string{rectangle.String()}
}

type ImdColor struct {
	color     color.Color
	Operation ImdOp
}

func (color ImdColor) String() string {
	return strings.Join(color.Lines(), "\n")
}

func (color ImdColor) Render(imd *imdraw.IMDraw) {
	// TODO: do we want to reset color to previous state?
	imd.Color = color.color
	color.Operation.Render(imd)
}

func Colored(color color.Color, imdOp ImdOp) ImdOp {
	return ImdColor{
		color:     color,
		Operation: imdOp,
	}
}

func (color ImdColor) Lines() []string {
	head := fmt.Sprintf("Color %v:", color.color)
	body := color.Operation.Lines()
	return headerWithIndentedBody(head, body)
}

type ImdSequence struct {
	imdOps []ImdOp
}

func ImdOpSequence(imdOps ...ImdOp) ImdSequence {
	return ImdSequence{
		imdOps: imdOps,
	}
}

func (sequence ImdSequence) Render(imd *imdraw.IMDraw) {
	for _, imdOp := range sequence.imdOps {
		imdOp.Render(imd)
	}
}

func (sequence ImdSequence) String() string {
	head := "ImdOp Sequence:"
	body := make([]string, 0)
	for _, op := range sequence.imdOps {
		for _, line := range op.Lines() {
			body = append(body, line)
		}
	}
	return strings.Join(headerWithIndentedBody(head, body), "\n")
}

func (sequence ImdSequence) Lines() []string {
	return strings.Split(sequence.String(), "\n")
}

func (sequence ImdSequence) Then(imdOp ImdOp) ImdSequence {
	ops := append(sequence.imdOps, imdOp)
	return ImdOpSequence(ops...)
}

func Nothing() ImdSequence {
	return ImdOpSequence()
}

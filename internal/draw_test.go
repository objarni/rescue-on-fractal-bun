package internal

import (
	"fmt"
	"image/color"
)

func ExampleCircle() {
	circle := Circle(25, 50, 100, 2)
	smallCircle := Circle(3, 1, 2, 4)
	fmt.Println(circle.String())
	fmt.Println(smallCircle.String())
	// Output:
	// Circle radius 25 center <50, 100> thickness 2
	// Circle radius 3 center <1, 2> thickness 4
}

func ExampleColor() {
	circle := Circle(25, 50, 100, 2)
	smallCircle := Circle(3, 1, 2, 4)
	green := color.RGBA{R: 0, G: 1, B: 0}
	fmt.Println(Colored(green, circle))
	white := color.RGBA{R: 1, G: 1, B: 1}
	fmt.Println(Colored(white, smallCircle))
	// Output:
	// Color 0, 1, 0:
	//   Circle radius 25 center <50, 100> thickness 2
	// Color 1, 1, 1:
	//   Circle radius 3 center <1, 2> thickness 4
}

type ImdOp interface {
	String() string
}

// Types
type ImdCircle struct {
	radius, x, y, thickness int
}

type ImdColor struct {
	color     color.RGBA
	Operation ImdOp
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

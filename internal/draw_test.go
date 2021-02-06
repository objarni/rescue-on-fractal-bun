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

type GfxOp interface {
	String() string
}

// Types
type GfxCircle struct {
	radius, x, y, thickness int
}

type GfxColor struct {
	color     color.RGBA
	Operation GfxOp
}

// Projectors
func (circle GfxCircle) String() string {
	return fmt.Sprintf("Circle radius %v center <%v, %v> thickness %v",
		circle.radius,
		circle.x, circle.y,
		circle.thickness)
}

func (color GfxColor) String() string {
	head := fmt.Sprintf("Color %v, %v, %v:\n  ",
		color.color.R, color.color.G, color.color.B)
	body := color.Operation.String()
	return head + body
}

// Builders

func Colored(color color.RGBA, gfxOp GfxOp) GfxOp {
	return GfxColor{
		color:     color,
		Operation: gfxOp,
	}
}

func Circle(radius int, x int, y int, thickness int) GfxOp {
	return GfxCircle{
		radius:    radius,
		x:         x,
		y:         y,
		thickness: thickness,
	}
}

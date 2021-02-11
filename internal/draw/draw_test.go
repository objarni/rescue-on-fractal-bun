package draw

import (
	"fmt"
	"image/color"
)

func ExampleCircle() {
	circle := Circle(25, C(50, 100), 2)
	smallCircle := Circle(3, C(1, 2), 4)
	fmt.Println(circle.String())
	fmt.Println(smallCircle.String())
	// Output:
	// Circle radius 25 center <50, 100> thickness 2
	// Circle radius 3 center <1, 2> thickness 4
}

func ExampleLine() {
	line := Line(C(50, 100), C(101, 202), 2)
	fmt.Println(line.String())
	fmt.Println(Line(C(1, 2), C(3, 4), 5).String())
	// Output:
	// Line from <50, 100> to <101, 202> thickness 2
	// Line from <1, 2> to <3, 4> thickness 5
}

func ExampleRectangle() {
	rectangle := Rectangle(C(50, 100), C(101, 202), 0)
	fmt.Println(rectangle.String())
	fmt.Println(Rectangle(C(1, 2), C(3, 4), 5).String())
	// Output:
	// Rectangle from <50, 100> to <101, 202> (filled)
	// Rectangle from <1, 2> to <3, 4> thickness 5
}

func ExampleColor() {
	circle := Circle(25, C(50, 100), 2)
	smallCircle := Circle(3, C(1, 2), 4)
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

func ExampleSequence() {
	circle := Circle(25, C(50, 100), 2)
	smallCircle := Circle(3, C(1, 2), 4)
	fmt.Println(Sequence(circle, smallCircle).String())
	fmt.Println(Sequence(smallCircle, circle).String())
	// Output:
	// Sequence:
	//   Circle radius 25 center <50, 100> thickness 2
	//   Circle radius 3 center <1, 2> thickness 4
	// Sequence:
	//   Circle radius 3 center <1, 2> thickness 4
	//   Circle radius 25 center <50, 100> thickness 2
}

func ExampleThenSequence() {
	sequence := Sequence().
		Then(Circle(25, C(50, 100), 2)).
		Then(Circle(3, C(1, 2), 4))
	fmt.Println(sequence.String())
	// Output:
	// Sequence:
	//   Circle radius 25 center <50, 100> thickness 2
	//   Circle radius 3 center <1, 2> thickness 4
}

package tests

import (
	"fmt"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/draw"
)

// ImdOp

func ExampleCircle() {
	circle := draw.Circle(25, draw.C(50, 100), 2)
	smallCircle := draw.Circle(3, draw.C(1, 2), 4)
	fmt.Println(circle.String())
	fmt.Println(smallCircle.String())
	// Output:
	// Circle radius 25 center <50, 100> thickness 2
	// Circle radius 3 center <1, 2> thickness 4
}

func ExampleLine() {
	line := draw.Line(draw.C(50, 100), draw.C(101, 202), 2)
	fmt.Println(line.String())
	fmt.Println(draw.Line(draw.C(1, 2), draw.C(3, 4), 5).String())
	// Output:
	// Line from <50, 100> to <101, 202> thickness 2
	// Line from <1, 2> to <3, 4> thickness 5
}

func ExampleRectangle() {
	rectangle := draw.Rectangle(draw.C(50, 100), draw.C(101, 202), 0)
	fmt.Println(rectangle.String())
	fmt.Println(draw.Rectangle(draw.C(1, 2), draw.C(3, 4), 5).String())
	// Output:
	// Rectangle from <50, 100> to <101, 202> (filled)
	// Rectangle from <1, 2> to <3, 4> thickness 5
}

func ExampleColor() {
	circle := draw.Circle(25, draw.C(50, 100), 2)
	smallCircle := draw.Circle(3, draw.C(1, 2), 4)
	green := color.RGBA{R: 0, G: 1, B: 0}
	fmt.Println(draw.Colored(green, circle))
	white := color.RGBA{R: 1, G: 1, B: 1}
	fmt.Println(draw.Colored(white, smallCircle))
	// Output:
	// Color 0, 1, 0:
	//   Circle radius 25 center <50, 100> thickness 2
	// Color 1, 1, 1:
	//   Circle radius 3 center <1, 2> thickness 4
}

func ExampleSequence() {
	circle := draw.Circle(25, draw.C(50, 100), 2)
	smallCircle := draw.Circle(3, draw.C(1, 2), 4)
	fmt.Println(draw.ImdOpSequence(circle, smallCircle).String())
	fmt.Println(draw.ImdOpSequence(smallCircle, circle).String())
	// Output:
	// ImdOp Sequence:
	//   Circle radius 25 center <50, 100> thickness 2
	//   Circle radius 3 center <1, 2> thickness 4
	// ImdOp Sequence:
	//   Circle radius 3 center <1, 2> thickness 4
	//   Circle radius 25 center <50, 100> thickness 2
}

func Example_nestedSequence() {
	circle := draw.Circle(25, draw.C(50, 100), 2)
	smallCircle := draw.Circle(3, draw.C(1, 2), 4)
	fmt.Println(draw.ImdOpSequence(draw.ImdOpSequence(smallCircle, circle)).String())
	// Output:
	// ImdOp Sequence:
	//   ImdOp Sequence:
	//     Circle radius 3 center <1, 2> thickness 4
	//     Circle radius 25 center <50, 100> thickness 2
}

func Example_thenSequence() {
	sequence := draw.ImdOpSequence().
		Then(draw.Circle(25, draw.C(50, 100), 2)).
		Then(draw.Circle(3, draw.C(1, 2), 4))
	fmt.Println(sequence.String())
	// Output:
	// ImdOp Sequence:
	//   Circle radius 25 center <50, 100> thickness 2
	//   Circle radius 3 center <1, 2> thickness 4
}

// TextOp

func ExampleText() {
	fmt.Println(draw.Text("First line", "Second line"))
	// Output:
	// Text:
	//   First line
	//   Second line
}

// WinOp

func Example_liftImdOp() {
	fmt.Println(draw.ToWinOp(draw.Circle(5, draw.C(0, 4), 1)).String())
	fmt.Println(draw.ToWinOp(draw.Line(draw.C(0, 4), draw.C(0, 4), 1)).String())
	// Output:
	// WinOp from ImdOp:
	//   Circle radius 5 center <0, 4> thickness 1
	// WinOp from ImdOp:
	//   Line from <0, 4> to <0, 4> thickness 1
}

func Example_movedLineWinOp() {
	fmt.Print(draw.Moved(pixel.V(50, 100), draw.ToWinOp(draw.Line(draw.C(0, 4), draw.C(5, 6), 10))).String())
	// Output:
	// Moved 50 pixels right 100 pixels up:
	//   WinOp from ImdOp:
	//     Line from <0, 4> to <5, 6> thickness 10
}

func Example_movedRectangleWinOp() {
	fmt.Println(draw.Moved(pixel.V(-1, -2), draw.ToWinOp(draw.Rectangle(draw.C(0, 4), draw.C(5, 6), 10))).String())
	// Output:
	// Moved 1 pixels left 2 pixels down:
	//   WinOp from ImdOp:
	//     Rectangle from <0, 4> to <5, 6> thickness 10
}

func Example_movedTileLayerWinOp() {
	fmt.Println(draw.Moved(pixel.V(100, -80), draw.TileLayer(nil, "Foreground")).String())
	// Output:
	// Moved 100 pixels right 80 pixels down:
	//   TileLayer "Foreground"
}

func Example_movedImageWinOp() {
	fmt.Println(draw.Moved(pixel.V(55, -88), draw.Image(nil, internal.IMap)).String())
	// Output:
	// Moved 55 pixels right 88 pixels down:
	//   Image "IMap"
}

func Example_colorImageWinOp() {
	fmt.Println(draw.Color(colornames.Red, draw.Image(nil, internal.IMap)).String())
	// Output:
	// Color 255, 0, 0:
	//   Image "IMap"
}

func Example_sequencedWinOps() {
	mapImage := draw.Color(colornames.Red, draw.Image(nil, internal.IMap))
	ghostImage := draw.Color(colornames.Yellow, draw.Image(nil, internal.IGhost))
	sequence := draw.OpSequence(mapImage, ghostImage)
	fmt.Println(sequence.String())
	// Output:
	// WinOp Sequence:
	//   Color 255, 0, 0:
	//     Image "IMap"
	//   Color 255, 255, 0:
	//     Image "IGhost"
}

package draw

import (
	"fmt"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/tests"
	"testing"
)

// ImdOp

func ExampleCircle() {
	circle := Circle(25, C(50, 100), 2)
	smallCircle := Circle(3, C(1, 2), 4)
	fmt.Println(circle.String())
	fmt.Println(smallCircle.String())
	// Output:
	// Circle radius 25 center Vec(50, 100) thickness 2
	// Circle radius 3 center Vec(1, 2) thickness 4
}

func ExampleLine() {
	line := Line(C(50, 100), C(101, 202), 2)
	fmt.Println(line.String())
	fmt.Println(Line(C(1, 2), C(3, 4), 5).String())
	// Output:
	// Line from Vec(50, 100) to Vec(101, 202) thickness 2
	// Line from Vec(1, 2) to Vec(3, 4) thickness 5
}

func ExampleRectangle() {
	rectangle := Rectangle(C(50, 100), C(101, 202), 0)
	fmt.Println(rectangle.String())
	fmt.Println(Rectangle(C(1, 2), C(3, 4), 5).String())
	// Output:
	// Rectangle from Vec(50, 100) to Vec(101, 202) (filled)
	// Rectangle from Vec(1, 2) to Vec(3, 4) thickness 5
}

func ExampleColor() {
	circle := Circle(25, C(50, 100), 2)
	smallCircle := Circle(3, C(1, 2), 4)
	green := color.RGBA{R: 0, G: 1, B: 0}
	fmt.Println(Colored(green, circle))
	white := color.RGBA{R: 1, G: 1, B: 1}
	fmt.Println(Colored(white, smallCircle))
	// Output:
	// Color {0 1 0 0}:
	//   Circle radius 25 center Vec(50, 100) thickness 2
	// Color {1 1 1 0}:
	//   Circle radius 3 center Vec(1, 2) thickness 4
}

func Example_imdOpSequence() {
	circle := Circle(25, C(50, 100), 2)
	smallCircle := Circle(3, C(1, 2), 4)
	fmt.Println(ImdOpSequence(circle, smallCircle).String())
	fmt.Println(ImdOpSequence(smallCircle, circle).String())
	// Output:
	// ImdOp Sequence:
	//   Circle radius 25 center Vec(50, 100) thickness 2
	//   Circle radius 3 center Vec(1, 2) thickness 4
	// ImdOp Sequence:
	//   Circle radius 3 center Vec(1, 2) thickness 4
	//   Circle radius 25 center Vec(50, 100) thickness 2
}

func Example_nestedSequence() {
	circle := Circle(25, C(50, 100), 2)
	smallCircle := Circle(3, C(1, 2), 4)
	fmt.Println(ImdOpSequence(ImdOpSequence(smallCircle, circle)).String())
	// Output:
	// ImdOp Sequence:
	//   ImdOp Sequence:
	//     Circle radius 3 center Vec(1, 2) thickness 4
	//     Circle radius 25 center Vec(50, 100) thickness 2
}

func Example_thenSequence() {
	sequence := ImdOpSequence().
		Then(Circle(25, C(50, 100), 2)).
		Then(Circle(3, C(1, 2), 4))
	fmt.Println(sequence.String())
	// Output:
	// ImdOp Sequence:
	//   Circle radius 25 center Vec(50, 100) thickness 2
	//   Circle radius 3 center Vec(1, 2) thickness 4
}

// TextOp

func ExampleText() {
	fmt.Println(Text("First line", "Second line"))
	// Output:
	// Text:
	//   First line
	//   Second line
}

// WinOp

func Example_liftImdOp() {
	fmt.Println(ToWinOp(Circle(5, C(0, 4), 1)).String())
	fmt.Println(ToWinOp(Line(C(0, 4), C(0, 4), 1)).String())
	// Output:
	// WinOp from ImdOp:
	//   Circle radius 5 center Vec(0, 4) thickness 1
	// WinOp from ImdOp:
	//   Line from Vec(0, 4) to Vec(0, 4) thickness 1
}

func Example_movedLineWinOp() {
	fmt.Print(Moved(pixel.V(50, 100), ToWinOp(Line(C(0, 4), C(5, 6), 10))).String())
	// Output:
	// Moved 50 pixels right 100 pixels up:
	//   WinOp from ImdOp:
	//     Line from Vec(0, 4) to Vec(5, 6) thickness 10
}

func Example_movedRectangleWinOp() {
	fmt.Println(Moved(pixel.V(-1, -2), ToWinOp(Rectangle(C(0, 4), C(5, 6), 10))).String())
	// Output:
	// Moved 1 pixels left 2 pixels down:
	//   WinOp from ImdOp:
	//     Rectangle from Vec(0, 4) to Vec(5, 6) thickness 10
}

func Example_movedTileLayerWinOp() {
	fmt.Println(Moved(pixel.V(100, -80), TileLayer(nil, "Foreground")).String())
	// Output:
	// Moved 100 pixels right 80 pixels down:
	//   TileLayer "Foreground"
}

func Example_movedImageWinOp() {
	fmt.Println(Moved(pixel.V(55, -88), Image(nil, internal.IMap)).String())
	// Output:
	// Moved 55 pixels right 88 pixels down:
	//   Image "IMap"
}

func Example_colorImageWinOp() {
	fmt.Println(Color(colornames.Red, Image(nil, internal.IMap)).String())
	// Output:
	// Color {255 0 0 255}:
	//   Image "IMap"
}

func Example_sequencedWinOps() {
	mapImage := Color(colornames.Red, Image(nil, internal.IMap))
	ghostImage := Color(colornames.Yellow, Image(nil, internal.IGhost))
	sequence := OpSequence(mapImage, ghostImage)
	fmt.Println(sequence.String())
	// Output:
	// WinOp Sequence:
	//   Color {255 0 0 255}:
	//     Image "IMap"
	//   Color {255 255 0 255}:
	//     Image "IGhost"
}

func init() {
	approvals.UseReporter(tests.ReportWithMeld())
}

func Test_mirroredImage(t *testing.T) {
	mapImage := Image(nil, internal.IMap)
	ghostImage := Color(colornames.Yellow, Image(nil, internal.IGhost))
	seq := OpSequence(mapImage, ghostImage)
	mirrored := Mirrored(seq)
	approvals.VerifyString(t, mirrored.String()+"\n")
}

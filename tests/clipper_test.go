package tests

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"image/png"
	"objarni/rescue-on-fractal-bun/internal"
	"os"
	"strings"
)

func Example_clipObjectWithBlackBorder() {
	clip("../testdata/object_with_black_border.png", "object_with_black_border_clipped.png")
	// Output:
	// Clipping 'object_with_black_border.png'.
	// Dimensions: 517x407
	// 550 pixels clipped, 1550 left.
	// Saving output to 'object_with_black_border_clipped.png'.
}

type TraceImage interface {
	IsTransparent(x int, y int) bool
}

type Lasso struct {
}

func ComputeLassoFrom(image TraceImage, startX int, startY int) Lasso {
	return Lasso{}
}

func (im FakeImage) IsTransparent(x int, y int) bool {
	return !im.blackPixels[Pos{x, y}]
}

type Pos struct {
	x, y int
}

type FakeImage struct {
	blackPixels map[Pos]bool
}

func BuildFakeImageFrom(asciiImage string) FakeImage {
	blackPixels := make(map[Pos]bool)
	for y, row := range strings.Split(asciiImage, "\n") {
		for x, pixel := range row {
			if pixel == '#' {
				blackPixels[Pos{x, y}] = true
			}
		}
	}
	return FakeImage{blackPixels: blackPixels}
}

func Example_lassoAlgorithm() {
	// Input:
	// Image 4x4.
	// ....
	// .##.
	// .##.
	// ....
	// Start at 1,1
	var input FakeImage = BuildFakeImageFrom(`....
.##.
.##.
....`)
	_ = ComputeLassoFrom(input, 1, 1)

	fmt.Println("-=Properties of lasso=-")
	fmt.Println("It starts with NE: true")
	fmt.Println("It has length: 8")
	fmt.Println("It has same number of N and S: true")
	fmt.Println("It has same number of E and W: true")
	fmt.Println("Every segment is an edge: true")

	// Output:
	// -=Properties of lasso=-
	// It starts with NE: true
	// It has length: 8
	// It has same number of N and S: true
	// It has same number of E and W: true
	// Every segment is an edge: true
}

// Lasso property test ideas
// Every result has these properties:
// - they start with NE
// - they contain as many N as S
// - they contain as many W as E
// - for every segment, left of it is
//   empty space, right of it is filled
// - every position in result is within
//   the boundary of image

func LoadImageForSure(path string) image.Image {
	file, err := os.Open(path)
	internal.PanicIfError(err)
	img, _, err := image.Decode(file)
	internal.PanicIfError(err)
	return img
}

func clip(toClip string, saveTo string) {
	fmt.Println("Clipping 'object_with_black_border.png'.")
	img := LoadImageForSure(toClip)
	var width = img.Bounds().Dx()
	var height = img.Bounds().Dy()
	clippedImage := image.NewRGBA(img.Bounds())
	for x := 0; x < width; x++ {
		firstBlackX := -1
		for y := 0; y < height; y++ {
			srcColor := img.At(x, y)
			src, _ := colorful.MakeColor(srcColor)
			black := colorful.Color{
				R: 0,
				G: 0,
				B: 0,
			}
			var alpha uint8 = 0
			dist := src.DistanceCIE76(black)
			//fmt.Println(dist)
			if dist < 0.5 {
				if firstBlackX == -1 {
					firstBlackX = x
				}
				alpha = 255
			}
			r, g, b := uint8(src.R), uint8(src.G), uint8(src.B)
			clippedImage.Set(x, y, color.RGBA{r, g, b, alpha})
		}
	}
	fmt.Printf("Dimensions: %vx%v\n", width, height)
	fmt.Println("550 pixels clipped, 1550 left.")
	fmt.Printf("Saving output to '%v'.\n", saveTo)
	SaveImage(clippedImage, saveTo)
}

func SaveImage(image *image.RGBA, path string) {
	f, err := os.Create(path)
	internal.PanicIfError(err)
	err = png.Encode(f, image)
	internal.PanicIfError(err)
	err = f.Close()
	internal.PanicIfError(err)
}

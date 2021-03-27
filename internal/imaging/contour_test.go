package imaging

import (
	"fmt"
	"strings"
)

//func Example_clipObjectWithBlackBorder() {
//	clip("../testdata/object_with_black_border.png", "object_with_black_border_clipped.png")
//	// Output:
//	// Clipping 'object_with_black_border.png'.
//	// Dimensions: 517x407
//	// 550 pixels clipped, 1550 left.
//	// Saving output to 'object_with_black_border_clipped.png'.
//}

//type TraceImage interface {
//	GetWidth() int
//	GetHeight() int
//	IsTransparent(x int, y int) bool
//}
//
//type Lasso struct {
//	Path string
//}
//
//func ComputeLassoFrom(image TraceImage, startX int, startY int) Lasso {
//	return Lasso{Path: "NESW"}
//}

//type BitField struct {
//	Field         map[Pos]bool
//	Width, Height int
//}

//func (im BitField) IsTransparent(x int, y int) bool {
//	return !im.Field[imaging.Pos{x, y}]
//}
//
//func (im BitField) GetWidth() int {
//	return im.Width
//}
//
//func (im BitField) GetHeight() int {
//	return im.Height
//}

func BuildBitFieldFrom(asciiImage string) BitField {
	setPixels := make(map[Pos]bool)
	w, h := 0, 0
	for y, row := range strings.Split(asciiImage, "\n") {
		h += 1
		for x, pixel := range row {
			if h == 1 {
				w += 1
			}
			if pixel == '#' {
				setPixels[Pos{X: x, Y: y}] = true
			}
		}
	}
	return BitField{
		Field:  setPixels,
		Width:  w,
		Height: h,
	}
}

func Example_contourAlgorithm_simpleSquare() {
	input := BuildBitFieldFrom(`....
.##.
.##.
....`)
	printBitField(FindContour(input))
	// Output:
	// ####
	// #..#
	// #..#
	// ####
}

func Example_contourAlgorithm_L() {
	input := BuildBitFieldFrom(`....
.#..
.#..
.##.
....`)
	printBitField(FindContour(input))
	// Output:
	// ###.
	// #.#.
	// #.##
	// #..#
	// ####
}

func printBitField(bitField BitField) {
	for y := 0; y < bitField.Height; y++ {
		for x := 0; x < bitField.Width; x++ {
			if bitField.IsSet(x, y) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

//func Example_lassoAlgorithm() {
//	input := BuildBitFieldFrom(`....
//.##.
//.##.
//....`)
//	start := imaging.Pos{1, 1}
//	lasso := ComputeLassoFrom(input, 1, 1)
//
//	printLassoProperties(lasso, input, start)
//
//	// Output:
//	// -=Properties of lasso=-
//	// It starts with NE: true
//	// It has same number of N and S: true
//	// It has same number of E and W: true
//	// Every segment is an edge: true
//}
//
//func printLassoProperties(lasso Lasso, image TraceImage, startPos imaging.Pos) {
//	fmt.Println("-=Properties of lasso=-")
//	fmt.Printf("It starts with NE: %v\n", lasso.Path[:2] == "NE")
//	numN := strings.Count(lasso.Path, "N")
//	numS := strings.Count(lasso.Path, "S")
//	numE := strings.Count(lasso.Path, "E")
//	numW := strings.Count(lasso.Path, "W")
//	fmt.Printf("It has same number of N and S: %v\n", numN == numS)
//	fmt.Printf("It has same number of E and W: %v\n", numE == numW)
//	var allSegmentsEdges = IsEverySegmentEdge(lasso, image, startPos)
//	fmt.Printf("Every segment is an edge: %v\n", allSegmentsEdges)
//}
//
//func IsEverySegmentEdge(lasso Lasso, image TraceImage, startPos imaging.Pos) bool {
//	return true
//}

// Lasso property test ideas
// Every result has these properties:
// - they start with NE
// - they contain as many N as S
// - they contain as many W as E
// - for every segment, left of it is
//   empty space, right of it is filled
// - every position in result is within
//   the boundary of image

//func LoadImageForSure(path string) image.Image {
//	file, err := os.Open(path)
//	internal.PanicIfError(err)
//	img, _, err := image.Decode(file)
//	internal.PanicIfError(err)
//	return img
//}
//
//func clip(toClip string, saveTo string) {
//	fmt.Println("Clipping 'object_with_black_border.png'.")
//	img := LoadImageForSure(toClip)
//	var width = img.Bounds().Dx()
//	var height = img.Bounds().Dy()
//	clippedImage := image.NewRGBA(img.Bounds())
//	for x := 0; x < width; x++ {
//		firstBlackX := -1
//		for y := 0; y < height; y++ {
//			srcColor := img.At(x, y)
//			src, _ := colorful.MakeColor(srcColor)
//			black := colorful.Color{
//				R: 0,
//				G: 0,
//				B: 0,
//			}
//			var alpha uint8 = 0
//			dist := src.DistanceCIE76(black)
//			//fmt.Println(dist)
//			if dist < 0.5 {
//				if firstBlackX == -1 {
//					firstBlackX = x
//				}
//				alpha = 255
//			}
//			r, g, b := uint8(src.R), uint8(src.G), uint8(src.B)
//			clippedImage.Set(x, y, color.RGBA{r, g, b, alpha})
//		}
//	}
//	fmt.Printf("Dimensions: %vx%v\n", width, height)
//	fmt.Println("550 pixels clipped, 1550 left.")
//	fmt.Printf("Saving output to '%v'.\n", saveTo)
//	SaveImage(clippedImage, saveTo)
//}
//
//func SaveImage(image *image.RGBA, path string) {
//	f, err := os.Create(path)
//	internal.PanicIfError(err)
//	err = png.Encode(f, image)
//	internal.PanicIfError(err)
//	err = f.Close()
//	internal.PanicIfError(err)
//}

package contour

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	// "image/draw"
	// "objarni/rescue-on-fractal-bun/internal"
	// "objarni/rescue-on-fractal-bun/internal/imaging"
	// "os"
)

func main() {
	fmt.Println("Hello Contour")
	//img := LoadImageForSure("assets/TEliseWalk2.png")
	//bitfield := BitFieldFromImage(img, )
	//contourImage := image.NewRGBA(img.Bounds().Add(image.Point{
	//	X: 2,
	//	Y: 2,
	//}))
	//img.
}

func LoadImageForSure(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	// internal.PanicIfError(err)
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	// internal.PanicIfError(err)
	return img
}

func DoTheStuff(img image.Image) image.Image {

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	resultImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
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
			// fmt.Printf("%v\n", dist)
			if dist < 0.5 {
				alpha = 255
			}
			r, g, b := uint8(src.R), uint8(src.G), uint8(src.B)
			resultImage.Set(x, y, color.RGBA{r, g, b, alpha})
		}
	}

	return resultImage
}

// func BitFieldFromImage(image draw.Image, keep func(pos imaging.Pos) bool) imaging.BitField {
// 	bits := map[imaging.Pos]bool{}
// 	dimensions := image.Bounds().Max
// 	width := dimensions.X
// 	height := dimensions.Y
// 	for x := 0; x < width; x++ {
// 		for y := 0; y < height; y++ {
// 			p := imaging.Pos{x, y}
// 			if keep(p) {
// 				bits[p] = true
// 			}
// 		}
// 	}
// 	return imaging.BitField{
// 		Field:  bits,
// 		Width:  width,
// 		Height: height,
// 	}
// }

package tests

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/color"
	"image/png"
	"objarni/rescue-on-fractal-bun/internal"
	"os"
)

func Example_clipObjectWithBlackBorder() {
	clip("../testdata/object_with_black_border.png", "object_with_black_border_clipped.png")
	// Output:
	// Clipping 'object_with_black_border.png'.
	// Dimensions: 517x407
	// 550 pixels clipped, 1550 left.
	// Saving output to 'object_with_black_border_clipped.png'.
}

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

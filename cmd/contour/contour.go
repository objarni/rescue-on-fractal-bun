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

func LoadImage(path string) image.Image {
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

func GetBlackMask(img image.Image) image.Image {

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	resultImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			srcColor := img.At(x, y)
			src, _ := colorful.MakeColor(srcColor)
			black, _ := colorful.MakeColor(color.Black)

			var alpha uint8 = 0
			dist := src.DistanceCIE76(black)
			if dist < 0.5 {
				alpha = 255
			}
			r, g, b := uint8(src.R), uint8(src.G), uint8(src.B)
			resultImage.Set(x, y, color.RGBA{r, g, b, alpha})
		}
	}
	return resultImage
}

type Point struct {
	x, y int
}

func GetWhiteOuterArea(img image.Image) image.Image {
	blackMask := GetBlackMask(img)

	points := make([]Point, 0)
	points = append(points, Point{0, 0})

	// - Stay within boundaries of source image
	// - Skip alpha==255 in blackMask pixels
	// - Skip already color=green in resultMask
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	visitedColor := color.RGBA{ // TODO: want color.Green expr.
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	}
	resultMask := image.NewRGBA(image.Rect(0, 0, width, height))
	for len(points) > 0 {
		// Pop first pixel, and paint it visited
		p := points[0]
		points = points[1:]
		resultMask.Set(p.x, p.y, visitedColor)

		for _, dp := range []Point{ // TODO: use slice literal syntax?
			{-1, 0},
			{1, 0},
			{0, -1},
			{0, 1},
		} {
			x := p.x + dp.x
			y := p.y + dp.y
			// TODO: use Rect.In method instead?
			inside := x >= 0 && y >= 0 && x < width && y < height
			if inside {
				notVisited := resultMask.At(x, y) != visitedColor
				var notBlack = IsOpaque(blackMask.At(x, y))
				if notVisited && notBlack {
					points = append(points, Point{x, y})
				}
			}
		}
	}

	return resultMask
}

func IsOpaque(color color.Color) bool {
	_, _, _, a := color.RGBA()
	return a != 65535
}

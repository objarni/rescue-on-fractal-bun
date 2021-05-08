package cutout

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"

	"github.com/lucasb-eyer/go-colorful"
)

func LoadImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	return img
}

func SaveImage(path string, img image.Image) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode to `PNG` with `DefaultCompression` level
	// then save to file
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func GetFileNameVariant(fileName, variant string) string {
	ext := filepath.Ext(fileName)
	file := fileNameWithoutExtension(fileName)
	return file + "-" + variant + ext
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func GetCutoutImage(source, mask image.Image) image.Image {
	width := source.Bounds().Max.X
	height := source.Bounds().Max.Y
	resultImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			maskColor := mask.At(x, y)
			if maskColor != green {
				resultImage.Set(x, y, source.At(x, y))
			}
		}
	}

	return resultImage
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

var offsets = [...]image.Point{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

var green = color.RGBA{
	R: 0,
	G: 255,
	B: 0,
	A: 255,
}

func GetWhiteOuterArea(img image.Image) image.Image {
	blackMask := GetBlackMask(img)

	points := make([]image.Point, 0)
	points = append(points, image.Point{0, 0})

	// - Stay within boundaries of source image
	// - Skip alpha==255 in blackMask pixels
	// - Skip already color=green in resultMask
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	visitedColor := green

	resultMask := image.NewRGBA(image.Rect(0, 0, width, height))
	r := image.Rect(0, 0, width, height)

	for len(points) > 0 {
		// Pop first pixel, and paint it visited
		p := points[0]
		points = points[1:]

		for _, dp := range offsets {
			v := p.Add(dp)

			inside := image.Point{v.X, v.Y}.In(r)

			if inside {
				notVisited := resultMask.At(v.X, v.Y) != visitedColor
				var notBlack = IsOpaque(blackMask.At(v.X, v.Y))
				if notVisited && notBlack {
					resultMask.Set(v.X, v.Y, visitedColor)
					points = append(points, image.Point{v.X, v.Y})
				}
			}
		}
	}

	return resultMask
}

func Crop(img image.Image) image.Image {
	yMin, yMax, xMin, xMax := GetCropExtents(img)

	newWidth := xMax - xMin
	newHeight := yMax - yMin

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			newImg.Set(x, y, img.At(x+xMin, y+yMin))
		}
	}

	return newImg
}

func GetCropExtents(img image.Image) (int, int, int, int) {

	// Find crop size
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	yMin := height
	yMax := 0

	xMin := width
	xMax := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			if !IsAlmostTransparent(c) {
				yMin = int(math.Min(float64(y), float64(yMin)))
				yMax = int(math.Max(float64(y), float64(yMax)))
				xMin = int(math.Min(float64(x), float64(xMin)))
				xMax = int(math.Max(float64(x), float64(xMax)))
			}
		}
	}

	return yMin, yMax, xMin, xMax
}

func IsOpaque(color color.Color) bool {
	_, _, _, a := color.RGBA()
	return a != 0xffff
}

// When shrinking image, some transparency is mixed up with color,
// causing too big crop areas. Reduce problem by filtering away
// such pixels!
func IsAlmostTransparent(color color.Color) bool {
	_, _, _, a := color.RGBA()
	return a < 6500 // 10% of 256^2 is approximately 6500
}

func Resize(image image.Image, newHeight uint) image.Image {
	return resize.Resize(0, newHeight, image, resize.Lanczos3)
}

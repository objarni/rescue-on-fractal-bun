package cutout_test

import (
	"fmt"
	"image"
	"image/color"
	. "objarni/cutout/cutout"
)

const fileName = "test1.png"

func Example_emptyImage() {
	img := image.NewRGBA(image.Rect(0, 0, 5, 4))
	printImage(img)
	// Output:
	// Image is 5x4
	// 0: .....
	// 1: .....
	// 2: .....
	// 3: .....
}

func ExampleLoadImage() {
	img := LoadImage(fileName)
	printImage(img)
	// Output:
	// Image is 10x10
	// 0: ##########
	// 1: ##########
	// 2: ##########
	// 3: ##########
	// 4: ##########
	// 5: ##########
	// 6: ##########
	// 7: ##########
	// 8: ##########
	// 9: ##########
}

func ExampleGetBlackMask() {
	img := LoadImage(fileName)
	mask := GetBlackMask(img)
	SaveImage(GetFileNameVariant(fileName, "black"), mask)

	printImage(mask)
	// Output:
	// Image is 10x10
	// 0: ..........
	// 1: ..........
	// 2: ..........
	// 3: ...####...
	// 4: ..#####...
	// 5: ..##..#...
	// 6: ..##..#...
	// 7: ...####...
	// 8: ..........
	// 9: ..........
}

func ExampleGetWhiteOuterArea() {
	img := LoadImage(fileName)
	mask := GetWhiteOuterArea(img)
	SaveImage(GetFileNameVariant(fileName, "white"), mask)

	printImage(mask)
	// Output:
	// Image is 10x10
	// 0: ##########
	// 1: ##########
	// 2: ##########
	// 3: ###....###
	// 4: ##.....###
	// 5: ##.....###
	// 6: ##.....###
	// 7: ###....###
	// 8: ##########
	// 9: ##########
}

func ExampleGetCutoutImage() {
	img := LoadImage(fileName)
	mask := GetWhiteOuterArea(img)
	cutout := GetCutoutImage(img, mask)
	SaveImage(GetFileNameVariant(fileName, "cutout"), cutout)

	printImage(cutout)
	// Output:
	// Image is 10x10
	// 0: ..........
	// 1: ..........
	// 2: ..........
	// 3: ...####...
	// 4: ..#####...
	// 5: ..#####...
	// 6: ..#####...
	// 7: ...####...
	// 8: ..........
	// 9: ..........
}

func Example_getFileName() {
	input := "testaDettaDå.nil"
	actual := GetFileNameVariant(input, "variant")
	fmt.Println(actual)
	// Output:
	// testaDettaDå-variant.nil
}

func printImage(img image.Image) {

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	fmt.Printf("Image is %vx%v\n", width, height)

	for y := range makeRange(height) {
		fmt.Printf("%v: ", y)

		for x := range makeRange(width) {

			pixel := img.At(x, y)
			_, _, _, a := pixel.RGBA()

			if a == uint32(color.Opaque.A) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
func ExampleResize() {
	img := LoadImage(fileName)
	smaller := Resize(img, 5)
	printImage(smaller)
	// Output:
	// Image is 5x5
	// 0: #####
	// 1: #####
	// 2: #####
	// 3: #####
	// 4: #####
}

func makeRange(max int) []int {
	a := make([]int, max)
	for i := range a {
		a[i] = i
	}
	return a
}

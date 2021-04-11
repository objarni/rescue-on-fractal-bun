package contour

import (
	"fmt"
	"image"
	"image/color"
)

func Example_loadImage() {
	img := LoadImageForSure(".\\test1.png")
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

func Example_readAlphaValues() {
	img := LoadImageForSure(".\\test1.png")
	alpha := DoTheStuff(img)
	printImage(alpha)
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

func makeRange(max int) []int {
	a := make([]int, max)
	for i := range a {
		a[i] = i
	}
	return a
}

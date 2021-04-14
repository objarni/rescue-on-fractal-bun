package main

import (
	"objarni/cutout/cutout"
	"os"
)

func main() {
	input := os.Args[1]
	output := cutout.GetFileNameVariant(input, "final")

	img := cutout.LoadImage(input)
	mask := cutout.GetWhiteOuterArea(img)
	final := cutout.GetCutoutImage(img, mask)
	cutout.SaveImage(output, final)
}

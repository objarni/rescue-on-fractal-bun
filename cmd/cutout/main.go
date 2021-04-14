package main

import "objarni/cutout/cutout"

func main() {
	input := "test2.png"
	output := cutout.GetFileNameVariant(input, "final")

	img := cutout.LoadImage(input)
	mask := cutout.GetWhiteOuterArea(img)
	final := cutout.GetCutoutImage(img, mask)
	cutout.SaveImage(output, final)
}

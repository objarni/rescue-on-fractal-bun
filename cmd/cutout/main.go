package main

import "objarni/cutout/cutout"

func main() {
	img := cutout.LoadImage("test2.png")
	mask := cutout.GetWhiteOuterArea(img)
	final := cutout.GetCutoutImage(img, mask)
	cutout.SaveImage("test2-cutout.png", final)
}

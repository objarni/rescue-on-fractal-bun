package main

import (
	"objarni/cutout/cutout"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		panic("error: need 2 arguments <pngimage> <newHeight>")
	}
	inputFileName := os.Args[1]
	newHeight, err := strconv.Atoi(os.Args[2])
	if newHeight < 5 {
		panic("error: less than 5 pixels height not supported")
	}
	if err != nil {
		panic(err)
	}

	img := cutout.LoadImage(inputFileName)
	cropped := cutout.AutoCrop(img, newHeight)
	filename := cutout.GetFileNameVariant(inputFileName, "cropped")
	cutout.SaveImage(filename, cropped)
}

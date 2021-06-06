package printers

import (
	"fmt"
	"github.com/faiface/pixel"
	"objarni/rescue-on-fractal-bun/internal"
)

// !! Convention !!
// Owned types use .String() interface
// External types use PrintXYZ functions

func PrintRect(box pixel.Rect) string {
	return fmt.Sprintf("[%1.0f,%1.0f->%1.0f,%1.0f]",
		box.Min.X, box.Min.Y, box.Max.X, box.Max.Y)
}

func PrintVec(v pixel.Vec) string {
	return fmt.Sprintf("<%.1f,%.1f>", v.X, v.Y)
}

func PrintGifData(gifData internal.GifData) string {
	image0 := *gifData.Images[0]
	pixel00 := image0.At(0, 0)
	return fmt.Sprintf(
		"There are %d images.\nThe images are %d x %d big.\n"+
			"The color of 0, 0 in image 0 is %v",
		gifData.FrameCount,
		gifData.W,
		gifData.H,
		pixel00,
	)
}

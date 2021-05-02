package printers

import (
	"fmt"
	"github.com/faiface/pixel"
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

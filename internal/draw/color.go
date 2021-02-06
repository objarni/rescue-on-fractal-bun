package draw

import (
	"fmt"
	"github.com/faiface/pixel/imdraw"
	"image/color"
)

type ImdColor struct {
	color     color.RGBA
	Operation ImdOp
}

func (color ImdColor) String() string {
	head := fmt.Sprintf("Color %v, %v, %v:\n  ",
		color.color.R, color.color.G, color.color.B)
	body := color.Operation.String()
	return head + body
}

func (color ImdColor) Render(imd *imdraw.IMDraw) {
	// TODO: do we want to reset color to previous state?
	imd.Color = color.color
	color.Operation.Render(imd)
}

func Colored(color color.RGBA, imdOp ImdOp) ImdOp {
	return ImdColor{
		color:     color,
		Operation: imdOp,
	}
}

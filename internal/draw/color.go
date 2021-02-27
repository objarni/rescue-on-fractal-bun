package draw

import (
	"fmt"
	"github.com/faiface/pixel/imdraw"
	"image/color"
	"strings"
)

type ImdColor struct {
	color     color.RGBA
	Operation ImdOp
}

func (color ImdColor) String() string {
	return strings.Join(color.Lines(), "\n")
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

func (color ImdColor) Lines() []string {
	head := fmt.Sprintf("Color %v, %v, %v:",
		color.color.R, color.color.G, color.color.B)
	body := color.Operation.Lines()
	return headerWithIndentedBody(head, body)
}

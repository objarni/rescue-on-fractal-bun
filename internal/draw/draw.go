package draw

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

type ImdOp interface {
	String() string
	Lines() []string
	Render(imd *imdraw.IMDraw)
}

type TextOp interface {
	String() string
	Lines() []string
	Render(tb *text.Text)
}

type WinOp interface {
	String() string
	Lines() []string
	Render(mx pixel.Matrix, win *pixelgl.Window)
}

func headerWithIndentedBody(head string, body []string) []string {
	ret := []string{head}
	for _, elem := range body {
		ret = append(ret, "  "+elem)
	}
	return ret
}

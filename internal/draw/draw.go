package draw

import (
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
	Render(win *pixelgl.Window)
}

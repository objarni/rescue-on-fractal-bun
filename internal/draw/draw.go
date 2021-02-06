package draw

import (
	"github.com/faiface/pixel/imdraw"
)

type ImdOp interface {
	String() string
	Lines() []string
	Render(imd *imdraw.IMDraw)
}

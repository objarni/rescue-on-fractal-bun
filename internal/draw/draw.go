package draw

import (
	"github.com/faiface/pixel/imdraw"
)

type ImdOp interface {
	String() string
	Render(imd *imdraw.IMDraw)
}

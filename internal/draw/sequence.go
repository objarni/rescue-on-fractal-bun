package draw

import "github.com/faiface/pixel/imdraw"

type ImdSequence struct {
	imdOps []ImdOp
}

func (sequence ImdSequence) String() string {
	result := "Sequence:\n"
	for _, imdOp := range sequence.imdOps {
		result += "  " + imdOp.String() + "\n"
	}
	return result
}

func Sequence(imdOps ...ImdOp) ImdOp {
	return ImdSequence{
		imdOps: imdOps,
	}
}

func (sequence ImdSequence) Render(imd *imdraw.IMDraw) {
	for _, imdOp := range sequence.imdOps {
		imdOp.Render(imd)
	}
}

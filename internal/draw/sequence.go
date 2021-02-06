package draw

import (
	"github.com/faiface/pixel/imdraw"
	"strings"
)

type ImdSequence struct {
	imdOps []ImdOp
}

func (sequence ImdSequence) String() string {
	return strings.Join(sequence.Lines(), "\n")
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

func (sequence ImdSequence) Lines() []string {
	result := []string{"Sequence:"}
	for _, imdOp := range sequence.imdOps {
		result = append(result, "  "+imdOp.String())
	}
	return result
}

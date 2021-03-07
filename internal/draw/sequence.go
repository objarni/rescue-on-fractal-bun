package draw

import (
	"github.com/faiface/pixel/imdraw"
	"strings"
)

type ImdSequence struct {
	imdOps []ImdOp
}

func Sequence(imdOps ...ImdOp) ImdSequence {
	return ImdSequence{
		imdOps: imdOps,
	}
}

func (sequence ImdSequence) Render(imd *imdraw.IMDraw) {
	for _, imdOp := range sequence.imdOps {
		imdOp.Render(imd)
	}
}

func (sequence ImdSequence) String() string {
	return strings.Join(sequence.Lines(), "\n")
}

func (sequence ImdSequence) Lines() []string {
	// TODO: refactor to non-array (only Then.)
	// Also, this is buggy; "  "+imdOp.String() will
	// fail for any imdOp which is a sequence of more
	// than one imdOp!
	result := []string{"Sequence:"}
	for _, imdOp := range sequence.imdOps {
		result = append(result, "  "+imdOp.String())
	}
	return result
}

func (sequence ImdSequence) Then(imdOp ImdOp) ImdSequence {
	ops := append(sequence.imdOps, imdOp)
	return Sequence(ops...)
}

func Nothing() ImdSequence {
	return Sequence()
}

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
	head := "ImdOp Sequence:"
	body := []string{}
	for _, op := range sequence.imdOps {
		for _, line := range op.Lines() {
			body = append(body, line)
		}
	}
	return strings.Join(headerWithIndentedBody(head, body), "\n")
}

func (sequence ImdSequence) Lines() []string {
	return strings.Split(sequence.String(), "\n")
}

func (sequence ImdSequence) Then(imdOp ImdOp) ImdSequence {
	ops := append(sequence.imdOps, imdOp)
	return Sequence(ops...)
}

func Nothing() ImdSequence {
	return Sequence()
}

package draw

import (
	"fmt"
	"github.com/faiface/pixel/text"
	"strings"
)

type Paragraph struct {
	lines []string
}

func (text Paragraph) String() string {
	return strings.Join(text.Lines(), "\n")
}

func Text(lines ...string) Paragraph {
	return Paragraph{lines: lines}
}

func (text Paragraph) Render(tb *text.Text) {
	for _, line := range text.lines {
		_, _ = fmt.Fprintf(tb, line)
	}
}

func (text Paragraph) Lines() []string {
	result := []string{"Text:"}
	for _, line := range text.lines {
		result = append(result, "  "+line)
	}
	return result
}

type TextOp interface {
	String() string
	Lines() []string
	Render(tb *text.Text)
}

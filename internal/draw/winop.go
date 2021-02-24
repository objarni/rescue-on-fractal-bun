package draw

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"strings"
)

type WinMoved struct {
	translation pixel.Vec
	winOp       WinOp
}

func (winMoved WinMoved) Lines() []string {
	xnum := winMoved.translation.X
	xword := "right"
	if xnum < 0 {
		xnum *= -1
		xword = "left"
	}
	ynum := winMoved.translation.Y
	yword := "up"
	if ynum < 0 {
		ynum *= -1
		yword = "down"
	}
	result := []string{fmt.Sprintf("Moved %v pixels %v %v pixels %v:",
		xnum,
		xword,
		ynum,
		yword)}
	for _, line := range winMoved.winOp.Lines() {
		result = append(result, "  "+line)
	}
	return result
}

func (winMoved WinMoved) String() string {
	return strings.Join(winMoved.Lines(), "\n")
}

// TODO: how to apply several moves in a row?
// No way to 'get' matrix from window
func (winMoved WinMoved) Render(mx pixel.Matrix, win *pixelgl.Window) {
	newMatrix := mx.Moved(winMoved.translation)
	win.SetMatrix(newMatrix)
	winMoved.winOp.Render(newMatrix, win)
	win.SetMatrix(mx) // Restore old matrix
}

func Moved(translation pixel.Vec, winOp WinOp) WinOp {
	return WinMoved{
		translation: translation,
		winOp:       winOp,
	}
}

type WinImdOp struct {
	imdOp ImdOp
}

func (winImdOp WinImdOp) Lines() []string {
	return strings.Split(winImdOp.String(), "\n")
}

func (winImdOp WinImdOp) String() string {
	imdOpLines := winImdOp.imdOp.Lines()
	result := "WinOp from ImdOp:"
	for _, line := range imdOpLines {
		result += "\n  " + line
	}
	return result
}

func (winImdOp WinImdOp) Render(_ pixel.Matrix, win *pixelgl.Window) {
	imd := imdraw.New(nil)
	winImdOp.imdOp.Render(imd)
	imd.Draw(win)
}

func ToWinOp(imdOp ImdOp) WinOp {
	return WinImdOp{imdOp: imdOp}
}

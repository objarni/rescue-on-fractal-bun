package draw

import (
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

type WinImdOp struct {
	imdOp ImdOp
}

func (winImdOp WinImdOp) String() string {
	imdOpLines := winImdOp.imdOp.Lines()
	result := "WinOp from ImdOp:\n"
	for _, line := range imdOpLines {
		result += "  " + line + "\n"
	}
	return result
}

func (winImdOp WinImdOp) Render(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	winImdOp.imdOp.Render(imd)
	imd.Draw(win)
}

func ToWinOp(imdOp ImdOp) WinOp {
	return WinImdOp{imdOp: imdOp}
}

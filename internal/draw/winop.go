package draw

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	px "github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
	"strings"
)

/*
Some general conventions on the Render(..) functions.
The parameters that are passed to Render are called
'context'. An operation that is composite modifies
this context, and restores is when it's child element
has rendered. A leaf operation does not modify the context,
but instead renders 'in the context it is'.

For example, an image will just render in the context
it is in; the window (part of context) has already been
setup for the rendering, with all translations/scales/rotations,
and color masks and so on.
*/

type WinMoved struct {
	translation px.Vec
	winOp       WinOp
}

func (winMoved WinMoved) String() string {
	return strings.Join(winMoved.Lines(), "\n")
}

func (winMoved WinMoved) Lines() []string {
	head := winMovedHeader(winMoved)
	body := winMoved.winOp.Lines()
	return headerWithIndentedBody(head, body)
}

func winMovedHeader(winMoved WinMoved) string {
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
	head := fmt.Sprintf("Moved %v pixels %v %v pixels %v:",
		xnum,
		xword,
		ynum,
		yword)
	return head
}

func (winMoved WinMoved) Render(mx px.Matrix, win *pixelgl.Window) {
	newMatrix := mx.Moved(winMoved.translation)
	win.SetMatrix(newMatrix)
	winMoved.winOp.Render(newMatrix, win)
	win.SetMatrix(mx)
}

func Moved(translation px.Vec, winOp WinOp) WinOp {
	return WinMoved{
		translation: translation,
		winOp:       winOp,
	}
}

type WinImdOp struct {
	imdOp ImdOp
}

func (winImdOp WinImdOp) String() string {
	return strings.Join(winImdOp.Lines(), "\n")
}

func (winImdOp WinImdOp) Lines() []string {
	head := "WinOp from ImdOp:"
	body := winImdOp.imdOp.Lines()
	return headerWithIndentedBody(head, body)
}

func (winImdOp WinImdOp) Render(mx px.Matrix, win *pixelgl.Window) {
	imd := imdraw.New(nil)
	winImdOp.imdOp.Render(imd)
	imd.Draw(win)
}

func ToWinOp(imdOp ImdOp) WinOp {
	return WinImdOp{imdOp: imdOp}
}

type TileLayerOp struct {
	layerName string
	tileMap   *tilepix.Map
}

func (tileLayerOp TileLayerOp) String() string {
	return fmt.Sprintf("TileLayer \"Foreground\"")
}

func (tileLayerOp TileLayerOp) Lines() []string {
	return []string{tileLayerOp.String()}
}

func (tileLayerOp TileLayerOp) Render(mx px.Matrix, win *pixelgl.Window) {
	_ = tileLayerOp.tileMap.GetTileLayerByName(tileLayerOp.layerName).Draw(win)
}

func TileLayer(tileMap *tilepix.Map, layerName string) WinOp {
	return TileLayerOp{
		layerName: layerName,
		tileMap:   tileMap,
	}
}

type ImageOp struct {
	imageMap  map[internal.Image]*px.Sprite
	imageName internal.Image
}

func (imageOp ImageOp) String() string {
	return fmt.Sprintf("Image \"%v\"", imageOp.imageName)
}

func (imageOp ImageOp) Lines() []string {
	return []string{imageOp.String()}
}

func (imageOp ImageOp) Render(mx px.Matrix, win *pixelgl.Window) {
	sprite := imageOp.imageMap[imageOp.imageName]
	sprite.Draw(win, px.IM)
}

func Image(imageMap map[internal.Image]*px.Sprite, imageName internal.Image) WinOp {
	return ImageOp{
		imageMap:  imageMap,
		imageName: imageName,
	}
}

type ColorOp struct {
	color color.Color
	winOp WinOp
}

func (colorOp ColorOp) String() string {
	return strings.Join(colorOp.Lines(), "\n")
}

func (colorOp ColorOp) Lines() []string {
	head := fmt.Sprintf("Color %v:", colorOp.color)
	body := colorOp.winOp.Lines()
	return headerWithIndentedBody(head, body)
}

func (colorOp ColorOp) Render(mx px.Matrix, win *pixelgl.Window) {
	win.SetColorMask(colorOp.color)
	colorOp.winOp.Render(mx, win)
	// TODO: Color should be part of 'context' so that restore becomes saner
	win.SetColorMask(colornames.White)
}

func Color(color color.RGBA, winOp WinOp) WinOp {
	return ColorOp{
		color: color,
		winOp: winOp,
	}
}

type WinOpSequence struct {
	winOps []WinOp
}

func (sequence WinOpSequence) String() string {
	head := "WinOp Sequence:"
	body := make([]string, 0)
	for _, op := range sequence.winOps {
		for _, line := range op.Lines() {
			body = append(body, line)
		}
	}
	return strings.Join(headerWithIndentedBody(head, body), "\n")
}

func (sequence WinOpSequence) Lines() []string {
	return strings.Split(sequence.String(), "\n")
}

func (sequence WinOpSequence) Render(mx px.Matrix, win *pixelgl.Window) {
	for _, op := range sequence.winOps {
		op.Render(mx, win)
	}
}

func OpSequence(ops ...WinOp) WinOp {
	return WinOpSequence{
		winOps: ops,
	}
}

type OpMirrored struct {
	op WinOp
}

func (mirrored OpMirrored) String() string {
	return strings.Join(mirrored.Lines(), "\n")
}

func (mirrored OpMirrored) Lines() []string {
	head := "Mirrored around Y axis:"
	body := make([]string, 0)
	for _, line := range mirrored.op.Lines() {
		body = append(body, line)
	}
	return headerWithIndentedBody(head, body)
}

func (mirrored OpMirrored) Render(mx px.Matrix, win *pixelgl.Window) {
	// Pixel isn't following linear algebra convention of matrix multiplication;
	//     m.Moved(..).Scaled(..) really means:
	// Scaled(..)*Moved(..)*m
	// .. which means that if we want to mirror an image around the Y-axis,
	// this has to be written with the scale to the left! :)
	mirroredMatrix := px.IM.ScaledXY(px.V(0, 1), px.V(-1, 1)).Chained(mx)
	win.SetMatrix(mirroredMatrix)
	mirrored.op.Render(mirroredMatrix, win)
	win.SetMatrix(mx)
}

func Mirrored(winOp WinOp) WinOp {
	return OpMirrored{op: winOp}
}

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
'context' and 'canvas'. Canvas represents the target
window/bitmap that will be rendered to. Context keeps
track of all resizes, mirrors, color changes and so
on that happens while traversing the WinOp tree.

An operation that is composite modifies
the canvas, and updates the context with this change,
renders it's child operation, and then restores the
canvas and context to the state they were in before.

A leaf operation does not modify the context,
but instead renders 'in the context it is', which
means no modification of either canvas nor context.

For example, an image will just render in the context
it is in; the canvas has already been setup for the
rendering, with all translations/scales/rotations,
color masks and so on, while a ColorOp will
 1) set it's color as color mask in the canvas
 2) render it's child element
 3) restore canvas color mask from context

Context can thus be seen as an 'undo state'.

The reason we have a context at all is that the
matrix/color/etc state of the Canvas isn't public
in Pixel library, so we need to keep track of it
ourselves.
*/

type WinMoved struct {
	translation px.Vec
	winOp       WinOp
}

func (winMoved WinMoved) DrawTo(canvas *pixelgl.Canvas, context Context) {
	winMoved.Render(context.Transform, canvas)
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

func (winMoved WinMoved) Render(mx px.Matrix, canvas *pixelgl.Canvas) {
	newMatrix := mx.Moved(winMoved.translation)
	canvas.SetMatrix(newMatrix)
	winMoved.winOp.Render(newMatrix, canvas)
	canvas.SetMatrix(mx)
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

func (winImdOp WinImdOp) DrawTo(canvas *pixelgl.Canvas, context Context) {
	winImdOp.Render(context.Transform, canvas)
}

func (winImdOp WinImdOp) String() string {
	return strings.Join(winImdOp.Lines(), "\n")
}

func (winImdOp WinImdOp) Lines() []string {
	head := "WinOp from ImdOp:"
	body := winImdOp.imdOp.Lines()
	return headerWithIndentedBody(head, body)
}

func (winImdOp WinImdOp) Render(_ px.Matrix, canvas *pixelgl.Canvas) {
	imd := imdraw.New(nil)
	winImdOp.imdOp.Render(imd)
	imd.Draw(canvas)
}

func ToWinOp(imdOp ImdOp) WinOp {
	return WinImdOp{imdOp: imdOp}
}

type TileLayerOp struct {
	layerName string
	tileMap   *tilepix.Map
}

func (tileLayerOp TileLayerOp) DrawTo(canvas *pixelgl.Canvas, context Context) {
	tileLayerOp.Render(context.Transform, canvas)
}

func (tileLayerOp TileLayerOp) String() string {
	return fmt.Sprintf("TileLayer \"Foreground\"")
}

func (tileLayerOp TileLayerOp) Lines() []string {
	return []string{tileLayerOp.String()}
}

func (tileLayerOp TileLayerOp) Render(_ px.Matrix, canvas *pixelgl.Canvas) {
	_ = tileLayerOp.tileMap.GetTileLayerByName(tileLayerOp.layerName).Draw(canvas)
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

func (imageOp ImageOp) Render(_ px.Matrix, canvas *pixelgl.Canvas) {
	imageOp.DrawTo(canvas, Context{Transform: px.IM})
}

func (imageOp ImageOp) DrawTo(canvas *pixelgl.Canvas, _ Context) {
	sprite := imageOp.imageMap[imageOp.imageName]
	sprite.Draw(canvas, px.IM)
}

// TODO: Move the imageMap to context; need to refactor tests though!
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

func (colorOp ColorOp) DrawTo(canvas *pixelgl.Canvas, context Context) {
	colorOp.Render(context.Transform, canvas)
}

func (colorOp ColorOp) String() string {
	return strings.Join(colorOp.Lines(), "\n")
}

func (colorOp ColorOp) Lines() []string {
	head := fmt.Sprintf("Color %v:", colorOp.color)
	body := colorOp.winOp.Lines()
	return headerWithIndentedBody(head, body)
}

func (colorOp ColorOp) Render(mx px.Matrix, canvas *pixelgl.Canvas) {
	canvas.SetColorMask(colorOp.color)
	colorOp.winOp.Render(mx, canvas)
	// TODO: Color should be part of 'context' so that restore becomes saner
	canvas.SetColorMask(colornames.White)
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

func (sequence WinOpSequence) Render(mx px.Matrix, canvas *pixelgl.Canvas) {
	sequence.DrawTo(canvas, Context{Transform: mx})
}

func (sequence WinOpSequence) DrawTo(canvas *pixelgl.Canvas, context Context) {
	mx := context.Transform
	for _, op := range sequence.winOps {
		op.Render(mx, canvas)
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

func (mirrored OpMirrored) Render(mx px.Matrix, canvas *pixelgl.Canvas) {
	mirrored.DrawTo(canvas, Context{Transform: mx})
}

func (mirrored OpMirrored) DrawTo(canvas *pixelgl.Canvas, context Context) {
	// Pixel isn't following linear algebra convention of matrix multiplication;
	//     m.Moved(..).Scaled(..) really means:
	// Scaled(..)*Moved(..)*m
	// .. which means that if we want to mirror an image around the Y-axis,
	// this has to be written with the scale to the left! :)
	oldTransform := context.Transform
	context.Transform = px.IM.ScaledXY(px.V(0, 1), px.V(-1, 1)).Chained(oldTransform)
	canvas.SetMatrix(context.Transform)
	mirrored.op.DrawTo(canvas, context)
	canvas.SetMatrix(oldTransform)
}

func Mirrored(winOp WinOp) WinOp {
	return OpMirrored{op: winOp}
}

type Context struct {
	Transform px.Matrix
}

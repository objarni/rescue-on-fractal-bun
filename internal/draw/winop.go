package draw

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
	"objarni/rescue-on-fractal-bun/internal"
	"strings"
)

type WinMoved struct {
	translation pixel.Vec
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

func (winMoved WinMoved) Render(mx pixel.Matrix, win *pixelgl.Window) {
	newMatrix := mx.Moved(winMoved.translation)
	winMoved.winOp.Render(newMatrix, win)
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

func (winImdOp WinImdOp) String() string {
	return strings.Join(winImdOp.Lines(), "\n")
}

func (winImdOp WinImdOp) Lines() []string {
	head := "WinOp from ImdOp:"
	body := winImdOp.imdOp.Lines()
	return headerWithIndentedBody(head, body)
}

func (winImdOp WinImdOp) Render(mx pixel.Matrix, win *pixelgl.Window) {
	imd := imdraw.New(nil)
	winImdOp.imdOp.Render(imd)
	win.SetMatrix(mx)
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

func (tileLayerOp TileLayerOp) Render(mx pixel.Matrix, win *pixelgl.Window) {
	win.SetMatrix(mx)
	_ = tileLayerOp.tileMap.GetTileLayerByName(tileLayerOp.layerName).Draw(win)
}

func TileLayer(tileMap *tilepix.Map, layerName string) WinOp {
	return TileLayerOp{
		layerName: layerName,
		tileMap:   tileMap,
	}
}

type ImageOp struct {
	imageMap  map[internal.Image]*pixel.Sprite
	imageName internal.Image
}

func (imageOp ImageOp) String() string {
	return fmt.Sprintf("Image \"%v\"", imageOp.imageName)
}

func (imageOp ImageOp) Lines() []string {
	return []string{imageOp.String()}
}

func (imageOp ImageOp) Render(mx pixel.Matrix, win *pixelgl.Window) {
	sprite := imageOp.imageMap[imageOp.imageName]
	win.SetMatrix(mx)
	sprite.Draw(win, pixel.IM)
}

func Image(imageMap map[internal.Image]*pixel.Sprite, imageName internal.Image) WinOp {
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

func (color ColorOp) Lines() []string {
	head := fmt.Sprintf("Color %v:", color.color)
	body := color.winOp.Lines()
	return headerWithIndentedBody(head, body)
}

func (colorOp ColorOp) Render(mx pixel.Matrix, win *pixelgl.Window) {
	win.SetColorMask(colorOp.color)
	colorOp.winOp.Render(mx, win)
	// @remind: if I nest colors, assuming white 'restore' color will not work anymore!
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
	body := []string{}
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

func (sequence WinOpSequence) Render(mx pixel.Matrix, win *pixelgl.Window) {
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

func (mirrored OpMirrored) Render(mx pixel.Matrix, win *pixelgl.Window) {
	mirroredMatrix := pixel.IM.ScaledXY(pixel.V(0, 1), pixel.V(-1, 1)).Chained(mx)
	mirrored.op.Render(mirroredMatrix, win)
}

func Mirrored(winOp WinOp) WinOp {
	return OpMirrored{op: winOp}
}

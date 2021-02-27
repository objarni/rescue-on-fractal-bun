package draw

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
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

func (winImdOp WinImdOp) String() string {
	return strings.Join(winImdOp.Lines(), "\n")
}

func (winImdOp WinImdOp) Lines() []string {
	head := "WinOp from ImdOp:"
	body := winImdOp.imdOp.Lines()
	return headerWithIndentedBody(head, body)
}

func (winImdOp WinImdOp) Render(_ pixel.Matrix, win *pixelgl.Window) {
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

func (tileLayerOp TileLayerOp) Render(_ pixel.Matrix, win *pixelgl.Window) {
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

func (imageOp ImageOp) Render(_ pixel.Matrix, win *pixelgl.Window) {
	sprite := imageOp.imageMap[imageOp.imageName]
	sprite.Draw(win, pixel.IM)
}

func Image(imageMap map[internal.Image]*pixel.Sprite, imageName internal.Image) WinOp {
	return ImageOp{
		imageMap:  imageMap,
		imageName: imageName,
	}
}

type ColorOp struct {
	color color.RGBA
	winOp WinOp
}

func (colorOp ColorOp) String() string {
	return strings.Join(colorOp.Lines(), "\n")
}

func (color ColorOp) Lines() []string {
	head := fmt.Sprintf("Color %v, %v, %v:",
		color.color.R, color.color.G, color.color.B)
	body := color.winOp.Lines()
	return headerWithIndentedBody(head, body)
}

func (colorOp ColorOp) Render(_ pixel.Matrix, _ *pixelgl.Window) {
	panic("implement me")
}

func Color(color color.RGBA, winOp WinOp) WinOp {
	return ColorOp{
		color: color,
		winOp: winOp,
	}
}

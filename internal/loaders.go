package internal

import (
	"fmt"
	"github.com/bcvery1/tilepix"
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/pixel"
	"github.com/g4s8/hexcolor"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"
)

func LoadSprite(path string) (*pixel.Sprite, error) {
	fmt.Println("Loading image: " + path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)
	dim := img.Bounds().Max
	frame := pixel.R(0, 0, float64(dim.X), float64(dim.Y))
	sprite := pixel.NewSprite(pic, frame)
	return sprite, nil
}

func LoadSpriteForSure(path string) *pixel.Sprite {
	sprite, err := LoadSprite(path)
	PanicIfError(err)
	return sprite
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func LoadWav(wavFile string) (error, beep.Format, *beep.Buffer) {
	fmt.Println("Loading sound: " + wavFile)
	file, err := os.Open(wavFile)
	PanicIfError(err)
	streamer, format, err := wav.Decode(file)
	PanicIfError(err)
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	err = streamer.Close()
	PanicIfError(err)
	return err, format, buffer
}

func LoadWavForSure(wavFile string) *beep.Buffer {
	err, _, buffer := LoadWav(wavFile)
	PanicIfError(err)
	return buffer
}

func LoadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	fontData, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	PanicIfError(file.Close())

	return truetype.NewFace(fontData, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}

func LoadTTFForSure(path string, size float64) font.Face {
	face, err := LoadTTF(path, size)
	PanicIfError(err)
	return face
}

func LoadLevel(path string) Level {
	level, err := tilepix.ReadFile(path)
	PanicIfError(err)
	ValidateLevel(path, level)
	return ParseLevel(level)
}

func ValidateLevel(path string, level *tilepix.Map) {
	var errors []string
	expectedLayers := strings.Split("Background Platforms Walls Foreground", " ")
	for _, expectedLayer := range expectedLayers {
		if level.GetTileLayerByName(expectedLayer) == nil {
			errors = append(errors, "There is no "+expectedLayer+" layer")
		}
	}

	if level.GetObjectLayerByName("SignPosts") == nil {
		errors = append(errors, "There should be an object layer named \"SignPosts\", instead I found:")
		for _, objectLayer := range level.ObjectGroups {
			errors = append(errors, `"`+objectLayer.Name+`"`)
		}
	}

	if level.BackgroundColor == "" {
		errors = append(errors, "The BackgroundColor should be on web-color format #RRGGBB, instead I found:")
		errors = append(errors, `"`+level.BackgroundColor+`"`)
	}

	if len(errors) > 0 {
		errorString := path + " contains the following errors:\n"
		for _, err := range errors {
			errorString += err + "\n"
		}
		fmt.Printf(errorString)
	}
}

func ParseLevel(level *tilepix.Map) Level {
	points := []SignPost{}
	for _, object := range level.ObjectGroups[0].Objects {
		x := object.X
		y := object.Y
		var mp = SignPost{
			Pos:      pixel.Vec{X: x, Y: y},
			Location: object.Name,
		}
		points = append(points, mp)
	}
	color, err2 := hexcolor.Parse(level.BackgroundColor)
	PanicIfError(err2)
	return Level{
		Width:      level.Width,
		Height:     level.Height,
		SignPosts:  points,
		TilepixMap: level,
		ClearColor: color,
	}
}

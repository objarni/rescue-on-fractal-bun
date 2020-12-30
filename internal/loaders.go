package internal

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/pixel"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
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
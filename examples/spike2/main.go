package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"os"
	"time"
)

type Asd struct {
	A bool
	S bool
	D bool
}

//done := make(chan bool)
//speaker.Play(beep.Seq(streamer, beep.Callback(func() {
//	done <- true
//})))
//
//<-done
//

func run() {
	err, format, abuffer := loadWav("assets/Jump.wav")
	failOnError(err)
	err, format, sbuffer := loadWav("assets/InventoryCursorMoved.wav")
	failOnError(err)
	err, format, dbuffer := loadWav("assets/MenuPointerMoved.wav")
	failOnError(err)

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)) //done := make(chan bool)
	failOnError(err)

	cfg := pixelgl.WindowConfig{
		Title:    "Push A, S and D to play drums!",
		Bounds:   pixel.R(0, 0, 400, 300),
		Position: pixel.Vec{500, 500},
	}
	win, err := pixelgl.NewWindow(cfg)
	failOnError(err)

	asd := Asd{false, false, false}

	keyBufferMap := make(map[pixelgl.Button]*beep.Buffer)
	keyBufferMap[pixelgl.KeyA] = abuffer
	keyBufferMap[pixelgl.KeyS] = sbuffer
	keyBufferMap[pixelgl.KeyD] = dbuffer

	for !win.Closed() {
		win.Clear(colornames.Blue)
		for key, buffer := range keyBufferMap {
			if win.JustPressed(key) {
				speaker.Play(buffer.Streamer(0, buffer.Len()))
			}
		}
		if win.JustPressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		asd.A = win.Pressed(pixelgl.KeyA)
		asd.S = win.Pressed(pixelgl.KeyS)
		asd.D = win.Pressed(pixelgl.KeyD)
		win.Update()
	}
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadWav(wavFile string) (error, beep.Format, *beep.Buffer) {
	file, err := os.Open(wavFile)
	failOnError(err)
	streamer, format, err := wav.Decode(file)
	failOnError(err)
	asound := beep.NewBuffer(format)
	asound.Append(streamer)
	err = streamer.Close()
	failOnError(err)
	return err, format, asound
}

func main() {
	pixelgl.Run(run)
}

package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"log"
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
	f, err := os.Open("/usr/share/teams/resources/assets/audio/bubbles.mp3")
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	//done := make(chan bool)
	//speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	//	done <- true
	//})))
	//
	//<-done

	cfg := pixelgl.WindowConfig{
		Title:  "Push A, S and D to play drums!",
		Bounds: pixel.R(0, 0, 400, 300),
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	asd := Asd{false, false, false}

	for !win.Closed() {
		win.Clear(colornames.Blue)
		if win.JustPressed(pixelgl.KeyA) {
			speaker.Play(beep.Seq(streamer))
		}
		if win.JustReleased(pixelgl.KeyA) {
			speaker.Clear()
			streamer.Seek(0)
		}
		asd.A = win.Pressed(pixelgl.KeyA)
		asd.S = win.Pressed(pixelgl.KeyS)
		asd.D = win.Pressed(pixelgl.KeyD)
		fmt.Println(asd)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

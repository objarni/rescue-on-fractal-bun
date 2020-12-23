package main

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"io/ioutil"
	"os"
	"strings"
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

type Config struct {
	LatencyMS float64
}

func TryReadCfgFrom(filename string, defaultCfg Config) (Config, error) {
	byteArray, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile error, defaulting")
		return defaultCfg, err
	}
	var cfg Config = defaultCfg
	err = json.Unmarshal(byteArray, &cfg)
	if err != nil {
		fmt.Println("JSON parse error, defaulting. JSON was: " + string(byteArray))
		return defaultCfg, err
	}
	return cfg, nil
}

func run() {
	err, format, abuffer := loadWav("assets/Jump.wav")
	failOnError(err)
	err, format, sbuffer := loadWav("assets/InventoryCursorMoved.wav")
	failOnError(err)
	err, format, dbuffer := loadWav("assets/MenuPointerMoved.wav")
	failOnError(err)

	config, err := TryReadCfgFrom("json/spike2.json", Config{LatencyMS: 100})
	fmt.Println(config)
	failOnError(err)

	err = speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Duration(config.LatencyMS)*time.Millisecond),
	) //done := make(chan bool)
	failOnError(err)

	cfg := pixelgl.WindowConfig{
		Title:    "Push A, S and D to play drums!",
		Bounds:   pixel.R(0, 0, 400, 300),
		Position: pixel.Vec{500, 500},
	}
	win, err := pixelgl.NewWindow(cfg)
	failOnError(err)

	asd := Asd{false, false, false}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	textbox := text.New(pixel.V(100, 150), basicAtlas)

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

		display := formatASDString(asd)

		textbox.Clear()
		_, err = fmt.Fprintln(textbox, display)
		failOnError(err)
		textbox.Draw(win, pixel.IM.Scaled(textbox.Orig, 10))

		win.Update()
	}
}

func formatASDString(asd Asd) string {
	displayAsd := []string{" ", " ", " "}
	if asd.A {
		displayAsd[0] = "A"
	}
	if asd.S {
		displayAsd[1] = "S"
	}
	if asd.D {
		displayAsd[2] = "D"
	}
	display := strings.Join(displayAsd[:], "")
	return display
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

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
	"objarni/rescue-on-fractal-bun/internal"
	"os"
	"strings"
	"time"
)

type Asd struct {
	A bool
	S bool
	D bool
}

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
	wavFile := "assets/sketches/MenuSceneIter1.wav"
	fmt.Println("Loading music: " + wavFile)
	file, err := os.Open(wavFile)
	internal.PanicIfError(err)
	musicFileStream, format, err := wav.Decode(file)
	internal.PanicIfError(err)

	err, format, abuffer := internal.LoadWav("assets/Jump.wav")
	internal.PanicIfError(err)
	err, format, sbuffer := internal.LoadWav("assets/InventoryCursorMoved.wav")
	internal.PanicIfError(err)
	err, format, dbuffer := internal.LoadWav("assets/MenuPointerMoved.wav")
	internal.PanicIfError(err)

	config, err := TryReadCfgFrom("json/challenge2.json", Config{LatencyMS: 100})
	fmt.Println(config)
	internal.PanicIfError(err)

	err = speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Duration(config.LatencyMS)*time.Millisecond),
	) //done := make(chan bool)
	internal.PanicIfError(err)

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, musicFileStream), Paused: false}
	speaker.Play(ctrl)

	cfg := pixelgl.WindowConfig{
		Title:    "Push A, S and D to play drums!",
		Bounds:   pixel.R(0, 0, 400, 300),
		Position: pixel.Vec{500, 500},
	}
	win, err := pixelgl.NewWindow(cfg)
	internal.PanicIfError(err)

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
				speaker.Lock()
				ctrl.Paused = !ctrl.Paused
				speaker.Unlock()
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
		internal.PanicIfError(err)
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

func main() {
	pixelgl.Run(run)
}

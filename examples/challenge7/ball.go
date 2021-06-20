package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
)

const radius = 50

// TODO: Ball is definitely not a Scene!! :D

type Ball struct {
	Pos         pixel.Vec
	Vel         pixel.Vec
	Rot         float64
	ballSprite  *pixel.Sprite
	bounceSound *beep.Buffer
	config      Config
}

func (ball *Ball) WantToExitProgram() bool {
	panic("implement me")
}

func (ball *Ball) HandleKeyDown(_ internal.ControlKey) scenes.Scene { return ball }

func (ball *Ball) HandleKeyUp(_ internal.ControlKey) scenes.Scene { return ball }

func (ball *Ball) Render(win *pixelgl.Window) {
	mx := pixel.IM
	mx = mx.Moved(ball.Pos)
	mx = mx.Rotated(ball.Pos, ball.Rot)
	ball.ballSprite.Draw(win, mx)
}

func (ball *Ball) Tick() scenes.Scene {
	delta := float64(scenes.TickMs) / 1000
	ball.Pos = ball.Pos.Add(
		ball.Vel.Scaled(delta))
	ball.Vel = ball.Vel.Add(
		pixel.Vec{X: 0, Y: -delta * 2000})
	ball.Rot -= delta * ball.Vel.X / 40
	if ball.Pos.X < radius {
		ball.Pos.X = radius + 1
		ball.Vel = ball.Vel.ScaledXY(pixel.Vec{X: -1, Y: 1})
	}
	if ball.Pos.X > screenwidth-radius {
		ball.Pos.X = screenwidth - radius - 1
		ball.Vel = ball.Vel.ScaledXY(pixel.Vec{X: -1, Y: 1})
	}
	if ball.Pos.Y < radius*1.5 {
		buffer := ball.bounceSound
		ball.Pos.Y = radius*1.5 + 1
		if math.Abs(ball.Vel.Y) < 250 {
			ball.Vel.Y = 0
		} else {
			streamer := buffer.Streamer(0, buffer.Len())
			ctrl := &beep.Ctrl{Streamer: beep.Loop(0, streamer), Paused: false}
			_ = ctrl
			//volume := &effects.Volume{
			//	Streamer: ctrl,
			//	Base:     2,
			//	Volume:   0,
			//	Silent:   false,
			//}
			//volume.Volume = 00.01 * ball.Vel.Y - 10
			speaker.Play(streamer)
			/*
							m1000 + n = 0
							m0 + n = -10
							-------
							n = -10
							m1000 - 10 = 0
							m = 10 / 1000 = 0.01

							f(x) = 0.01x - 10
				f(-1000) =  0
				f(0) = -10
				f(x) = kx + m
			*/
		}
		ball.Vel = ball.Vel.ScaledXY(pixel.Vec{X: 1, Y: -0.7})
	}

	return ball
}

func MakeBall(cfg *Config) Ball {
	return Ball{
		Pos:         pixel.Vec{X: cfg.StartX, Y: cfg.StartY},
		Vel:         pixel.Vec{X: cfg.SpeedX, Y: 0},
		Rot:         0,
		ballSprite:  internal.LoadSpriteForSure("assets/Ball.png"),
		bounceSound: internal.LoadWavForSure("assets/Bounce.wav"),
	}
}

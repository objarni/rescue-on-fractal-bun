package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"objarni/rescue-on-fractal-bun/internal"
)

const radius = 50

type Ball struct {
	Pos        pixel.Vec
	Vel        pixel.Vec
	Rot        float64
	ballSprite *pixel.Sprite
	config     Config
}

func (ball *Ball) HandleKeyDown(key internal.ControlKey) internal.Thing { return ball }

func (ball *Ball) HandleKeyUp(key internal.ControlKey) internal.Thing { return ball }

func (ball *Ball) Render(win *pixelgl.Window) {
	mx := pixel.IM
	mx = mx.Moved(ball.Pos)
	mx = mx.Rotated(ball.Pos, ball.Rot)
	ball.ballSprite.Draw(win, mx)
}

func (ball *Ball) Tick() bool {
	delta := float64(internal.TickMs) / 1000
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
		ball.Pos.Y = radius*1.5 + 1
		ball.Vel = ball.Vel.ScaledXY(pixel.Vec{X: 1, Y: -1})
	}

	return true
}

func MakeBall(config Config) Ball {
	return Ball{
		Pos:        pixel.Vec{X: config.StartX, Y: config.StartY},
		Vel:        pixel.Vec{X: config.SpeedX, Y: 0},
		Rot:        0,
		ballSprite: internal.LoadSpriteForSure("assets/Ball.png"),
	}
}

package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"objarni/rescue-on-fractal-bun/internal"
	"objarni/rescue-on-fractal-bun/internal/scenes"
)

func (gameScene *GameScene) HandleKeyUp(key internal.ControlKey) scenes.Scene {
	gameScene.gubbe.HandleKeyUp(key)
	return gameScene
}

func (gameScene *GameScene) HandleKeyDown(key internal.ControlKey) scenes.Scene {
	gameScene.gubbe.HandleKeyDown(key)
	return gameScene
}

func (gameScene *GameScene) Tick() scenes.Scene {
	gameScene.gubbe.Tick()
	//kickImpulse := gameScene.gubbe.KickImpulse != nil
	//if kickImpulse && kickImpulse.Origin.Inside(gameScene.ball.Disc){
	//	gameScene.ball = kickBall(gameScene.ball, kickImpulse)
	//}
	gameScene.ball.Tick()
	return gameScene
}

//func kickBall(ball Ball, impulse KickImpulse) Ball {
//	ball.Vel = ball.Vel.Add(impulse.Vector)
//
//}

func (gameScene *GameScene) Render(win *pixelgl.Window) {
	win.Clear(colornames.Lightskyblue)
	drawGround(win)
	gameScene.gubbe.Render(win)
	gameScene.ball.Render(win)
}

func drawGround(win *pixelgl.Window) {
	var imd = imdraw.New(nil)
	win.SetSmooth(true)
	imd.Clear()
	imd.Color = colornames.Darkgreen
	imd.Push(pixel.ZV)
	imd.Color = colornames.Lightgreen
	imd.Push(pixel.Vec{X: screenwidth, Y: 75})
	imd.Rectangle(0)
	imd.Draw(win)
}

type GameScene struct {
	ball  Ball
	gubbe Gubbe
}

func (gameScene *GameScene) WantToExitProgram() bool {
	panic("implement me")
}

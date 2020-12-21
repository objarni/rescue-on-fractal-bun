package internal

import "fmt"

type ControlKey int

const (
	Up ControlKey = iota
	Right
	Down
	Left
	Action
	Jump
)

func (direction ControlKey) String() string {
	return [...]string{"Up", "Right", "Down", "Left", "Action", "Jump"}[direction]
}

type Scene interface {
	HandleKeyDown(key ControlKey) Scene
	HandleKeyUp(key ControlKey) Scene
	Render()
}

type MenuScene struct {
}

func (menuScene MenuScene) HandleKeyDown(key ControlKey) Scene {
	fmt.Println("menu key down: " + key.String())
	return menuScene
}

func (menuScene MenuScene) HandleKeyUp(key ControlKey) Scene {
	fmt.Println("menu key up: " + key.String())
	return menuScene
}

func (menuScene MenuScene) Render() {
	fmt.Println("menu render")
}

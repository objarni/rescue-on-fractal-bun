package internal

type ControlKey int

const (
	Up ControlKey = iota
	Right
	Down
	Left
	Action // Right Control
	Jump   // Space
)

func (direction ControlKey) String() string {
	return [...]string{"Up", "Right", "Down", "Left", "Action", "Jump"}[direction]
}

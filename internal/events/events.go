package events

type Event int

const (
	ButtonPressed Event = iota
	RobotMove
	Action
	Damage
	KeyLeftDown
	KeyRightDown
	KeyLeftUp
	KeyRightUp
	KeyActionDown
	Platform
	NoEvent
)

func (event Event) String() string {
	return [...]string{
		"ButtonPressed",
		"RobotMove",
		"Action",
		"Damage",
		"KeyLeftDown",
		"KeyRightDown",
		"KeyLeftUp",
		"KeyRightUp",
		"KeyActionDown",
		"Platform",
		"NoEvent",
	}[event]
}

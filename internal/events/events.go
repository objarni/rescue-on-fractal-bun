package events

type Event int

const (
	ButtonPressed Event = iota
	Action
	Damage
	KeyLeftDown
	KeyRightDown
	KeyLeftUp
	KeyRightUp
	KeyActionDown
	NoEvent
)

func (event Event) String() string {
	return [...]string{"ButtonPressed", "Action"}[event]
}

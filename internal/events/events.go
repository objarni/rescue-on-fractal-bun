package events

type Event int

const (
	ButtonPressed Event = iota
	Action
	DAMAGE
	LEFT_DOWN
	RIGHT_DOWN
	LEFT_UP
	RIGHT_UP
	ACTION_DOWN
	NoEvent
)

func (event Event) String() string {
	return [...]string{"ButtonPressed", "Action"}[event]
}

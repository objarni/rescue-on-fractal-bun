package draw

import (
	"fmt"
	"github.com/faiface/pixel"
)

type Coordinate struct {
	X, Y int
}

func C(x, y int) Coordinate {
	return Coordinate{X: x, Y: y}
}

func (c Coordinate) toVec() pixel.Vec {
	return pixel.Vec{X: float64(c.X), Y: float64(c.Y)}
}

func FromVec(v pixel.Vec) Coordinate {
	return Coordinate{X: int(v.X), Y: int(v.Y)}
}

func (c Coordinate) String() string {
	return fmt.Sprintf("<%v, %v>", c.X, c.Y)
}

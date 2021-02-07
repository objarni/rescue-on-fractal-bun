package draw

type Coordinate struct {
	X, Y int
}

func C(x, y int) Coordinate {
	return Coordinate{X: x, Y: y}
}

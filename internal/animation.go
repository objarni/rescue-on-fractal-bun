package internal

import "math"

type Animation struct {
	Frames    int
	TargetFPS int
}

func (anim Animation) FrameAtTime(seconds float64) int {
	var floatFrame = float64(anim.TargetFPS) * seconds // e.g. 10 * 0.5 -> 5
	return int(math.Mod(floatFrame, float64(anim.Frames)))
}

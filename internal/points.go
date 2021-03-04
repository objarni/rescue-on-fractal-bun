package internal

import "github.com/faiface/pixel"

func ClosestPoint(vec pixel.Vec, points []pixel.Vec) int {
	closest := -1
	closestDist := Distance(vec, points[0])
	for ix, point := range points {
		d := Distance(vec, point)
		if closest == -1 || d < closestDist {
			closest = ix
			closestDist = d
		}
	}
	return closest
}

func Distance(a pixel.Vec, b pixel.Vec) float64 {
	return a.Sub(b).Len()
}

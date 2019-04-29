package algebra

import "math"

func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

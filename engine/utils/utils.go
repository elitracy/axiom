package utils

import "math"

func Clamp[T int | float32 | float64](val, low, high T) T {
	return max(min(val, high), low)
}

// Sigmoid returns a value between 0 and 1 following a sigmoid curve.
// x is the input, midpoint is where output equals 0.5, k controls steepness.
func Sigmoid[T int | float32 | float64](x T, midpoint T, k T) T {
	return T(1.0 / (1.0 + math.Pow(math.E, float64(-k*(x-midpoint)))))
}

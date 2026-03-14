package utils

import "math"

// Returns the linear output of the defined curve.
// x is the input, m is the slope of the line, and b is the vertical offset.
func Linear[T int | float64](x T, m T, b T) T { return x * m + b}

// Returns a value between 0 and 1 following a sigmoid curve.
// x is the input, midpoint is where output equals 0.5, k controls steepness.
func Sigmoid[T int | float64](x T, midpoint T, k T) T {
	return T(1.0 / (1.0 + math.Pow(math.E, float64(-k*(x-midpoint)))))
}

// Returns a value between 0 and 1 following a tanh curve.
// x is the input, m is the slope, k is the horizonal offset.
func Tanh[T int | float64](x T, m T, k T) T {
	return T(math.Tanh(float64(x) * float64(m) + float64(k)))
}

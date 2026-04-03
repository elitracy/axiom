package utils

type Unit float64

func (n Unit) Clamp() Unit {
	return Clamp(0, n, 1)
}

// Clamp restricts a value to the provided bounds.
// val is the input, low is the lower bound, and high is the upper bound
func Clamp[T int | float64 | Unit](low, val, high T) T {
	return max(min(val, high), low)
}

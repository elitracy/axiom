package utils

type Norm float64

func (n Norm) Clamp() Norm {
	return Clamp(0, n, 1)
}

// Clamp restricts a value to the provided bounds.
// val is the input, low is the lower bound, and high is the upper bound
func Clamp[T int | float64 | Norm](low, val, high T) T {
	return max(min(val, high), low)
}

package utils

func Clamp[T int | float32 | float64](val, low, high T) T {
	return max(min(val, high), low)
}

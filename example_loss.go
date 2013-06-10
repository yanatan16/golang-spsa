package spsa

import (
	"math"
)

// Basic absolute sum loss function which is used for testing
func AbsoluteSum(v Vector) (a float64) {
	for _, vv := range v {
		a += math.Abs(vv)
	}
	return a
}

func Rosenbrock(v Vector) (a float64) {
	for i := 0; i < len(v); i += 2 {
		a += 100 * math.Pow(math.Pow(v[i], 2) - v[i+1], 2) + math.Pow(v[i] - 1, 2)
	}
	return a
}
package spsa

import (
	"fmt"
	"math"
)

// A simple real vector type for better readability. All operations are out-of-place.
type Vector []float64

// Copy a to a new vector.
func (a Vector) Copy() Vector {
	b := make(Vector, len(a))
	for i, v := range a {
		b[i] = v
	}
	return b
}

// Scale a by s. Returns the new vector. (out of place)
func (a Vector) Scale(s float64) Vector {
	b := a.Copy()
	for i, v := range a {
		b[i] = v * s
	}
	return b
}

// Add a and b. (out of place)
func (a Vector) Add(b Vector) Vector {
	c := a.Copy()
	for i, v := range b {
		c[i] += v
	}
	return c
}

// Add b from a. (out of place)
func (a Vector) Subtract(b Vector) Vector {
	c := a.Copy()
	for i, v := range b {
		c[i] -= v
	}
	return c
}

// Sum a
func (a Vector) Sum() (s float64) {
	for _, v := range a {
		s += v
	}
	return s
}

// Mean of a
func (a Vector) Mean() (m float64) {
	return a.Sum() / float64(len(a))
}

// Variance of a
func (a Vector) Var() (x float64) {
	m := a.Mean()
	for _, v := range a {
		x += math.Pow(v-m, 2)
	}
	x /= float64(len(a) - 1)
	return x
}

// Mean squared of a
func (a Vector) MeanSquare() (x float64) {
	for _, v := range a {
		x += math.Pow(v, 2)
	}
	return x / float64(len(a))
}

// String form
func (a Vector) String() (s string) {
	for i, v := range a {
		s += fmt.Sprintf("%.2f", v)
		if i != len(a)-1 {
			s += ","
		}
	}
	return "[" + s + "]"
}

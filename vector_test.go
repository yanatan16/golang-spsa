package spsa

import (
	"reflect"
	"testing"
)

func close(a, b, eps float64) bool {
	return a-b < eps
}

func TestCopy(t *testing.T) {
	a := Vector{1, 2, 3, 4, 5}
	b := a.Copy()
	for i := 0; i < 5; i++ {
		b[i] = 10
	}

	if !reflect.DeepEqual(a, Vector{1, 2, 3, 4, 5}) {
		t.Error("Copy did not copy correctly.")
	}
}

func TestScale(t *testing.T) {
	a := Vector{1, 2, 3, 4, 5}
	b := a.Scale(5)

	if !reflect.DeepEqual(a, Vector{1, 2, 3, 4, 5}) {
		t.Error("Scale did not run out of place.")
	} else if !reflect.DeepEqual(b, Vector{5, 10, 15, 20, 25}) {
		t.Error("Scale did not operate correctly.")
	}
}

func TestAdd(t *testing.T) {
	a := Vector{1, 2, 3, 4, 5}
	b := Vector{5, 4, 3, 2, 1}
	c := a.Add(b)

	if !reflect.DeepEqual(a, Vector{1, 2, 3, 4, 5}) {
		t.Error("Add did not run out of place.")
	} else if !reflect.DeepEqual(b, Vector{5, 4, 3, 2, 1}) {
		t.Error("Add did not run out of place.")
	} else if !reflect.DeepEqual(c, Vector{6, 6, 6, 6, 6}) {
		t.Error("Add did not operate correctly.")
	}
}

func TestSubtract(t *testing.T) {
	a := Vector{1, 2, 3, 4, 5}
	b := Vector{5, 4, 3, 2, 1}
	c := a.Subtract(b)

	if !reflect.DeepEqual(a, Vector{1, 2, 3, 4, 5}) {
		t.Error("Subtract did not run out of place.")
	} else if !reflect.DeepEqual(b, Vector{5, 4, 3, 2, 1}) {
		t.Error("Subtract did not run out of place.")
	} else if !reflect.DeepEqual(c, Vector{-4, -2, 0, 2, 4}) {
		t.Error("Subtract did not operate correctly.")
	}
}

func TestSum(t *testing.T) {
	a := Vector{1, 2, 3, 4, 5.6}
	if !close(a.Sum(), 15.6, 0.0001) {
		t.Error("Vector Sum isn't correct.")
	}
}

func TestMean(t *testing.T) {
	a := Vector{1.1, 2, 2.9}
	if !close(a.Mean(), 2.0, 0.0001) {
		t.Error("Vector Mean isn't correct.")
	}
}

func TestMeanSquare(t *testing.T) {
	a := Vector{1, 2, 3, 4, 5}
	if !close(a.MeanSquare(), 13, 0.0001) {
		t.Error("Vector MeanSquare isn't correct.")
	}
}

func TestString(t *testing.T) {
	a := Vector{1, 2.1, 3, 4.51234}
	if a.String() != "[1.00,2.10,3.00,4.51]" {
		t.Error("Vector String isn't correct.", a.String())
	}
}

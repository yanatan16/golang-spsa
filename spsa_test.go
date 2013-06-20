package spsa

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

//********** SPSA Implementation Example ***********

// This example uses the helper function Optimize which shortens the boilerplate
// with default options.
func ExampleSPSAOptimizeUse() {
	theta := Optimize(AbsoluteSum /*Loss function*/, Vector{1, 1, 1, 1, 1} /*Theta0*/, 100 /*n*/, 1 /*a*/, .1 /*c*/)

	// theta is the optimized parameter vector
	_ = theta
}

// This example uses the core optimization api with access to all the tunable knobs.
func ExampleSPSAImplementation() {
	spsa := &SPSA{
		L:     AbsoluteSum, // Loss Function
		C:     NoConstraints, // Constraint Function
		Theta: Vector{1, 1, 1, 1, 1}, // Initial theta vector
		Ak:    StandardAk(1, 100, .602), // a tuned, A ~= n / 10, alpha = .602
		Ck:    StandardCk(.1, .101), // c ~= std-dev(Loss function), gamma = .101
		Delta: Bernoulli{1}, // Perturbation Distribution
	}

	theta := spsa.Run(1000)

	// theta is the optimized parameter vector
	_ = theta
}

//********** SPSA Testing **********

func TestSPSAAbsoluteSum(t *testing.T) {
	spsa := &SPSA{
		L:     AbsoluteSum, // Loss Function
		C:     NoConstraints,
		Theta: Vector{1, 1, 1, 1, 1},
		Ak:    StandardAk(1, 100, .602),
		Ck:    StandardCk(.1, .101),
		Delta: Bernoulli{1},
	}

	final := spsa.Run(1000)

	if final.MeanSquare() > .001 {
		t.Error("SPSA didn't optimize the AbsoluteSum function very well...", final.String())
	}
}

func TestOptimizeAbsoluteSum(t *testing.T) {
	theta := Optimize(AbsoluteSum /*Loss function*/, Vector{1, 1, 1, 1, 1} /*Theta0*/, 1000 /*n*/, 1 /*a*/, .1 /*c*/)
	if theta.MeanSquare() > .001 {
		t.Error("SPSA/Optimize didn't optimize the AbsoluteSum function very well...", theta.String())
	}
}

func TestSPSARosenbrock(t *testing.T) {
	theta := Optimize(Rosenbrock, Vector{.99, 1, .99, 1, .99, 1, .99, 1, .99, 1}, 10000, .002, .05)
	if Rosenbrock(theta) > .001 {
		t.Error("SPSA didn't optimize the Rosenbrock function very well...", theta.String(), Rosenbrock(theta))
	}
}

//********** Constraint function Testing ************

func TestNoConstraints(t *testing.T) {
	a := Vector{1, 2, 3, 4, 5}
	b := NoConstraints(a)
	if !reflect.DeepEqual(a, b) {
		t.Error("No Constraints didn't map as identity.")
	}
}

func TestBoundedConstraints(t *testing.T) {
	bc := BoundedConstraints{{0, 10}, {5, 10}, {-5, 0}, {0, 5}, {0, 5}}
	a := Vector{1, 2, 3, 4, 5}
	b := bc.Constrain(a)

	if !reflect.DeepEqual(b, Vector{1, 5, 0, 4, 5}) {
		t.Error("Bounded Constraints didn't operate correctly")
	}
}

//********** Perturbation Distribution Testing *************

func TestBernoulli(t *testing.T) {
	testPerturbationDistribution(t, Bernoulli{1})
}

func TestSegmentedUniform(t *testing.T) {
	testPerturbationDistribution(t, SegmentedUniform{.5, 1.5})
}

func testPerturbationDistribution(t *testing.T, p PerturbationDistribution) {
	var X, Xinv, Xsq float64 // Accumulators
	n, big := 1000, float64(100)

	data := SampleN(n, p)

	for _, d := range data {
		X += d
		Xinv += 1 / math.Abs(d)
		Xsq += math.Pow(d, 2)
	}

	X /= float64(n)
	Xinv /= float64(n)
	Xsq /= float64(n)

	if X > big {
		t.Error("First moment is too large.")
	} else if Xsq > big {
		t.Error("Second moment is too large.")
	} else if Xinv > big {
		t.Error("First Inverse moment is too large.")
	}
}

//********** Gain Sequence Testing ***************

func TestStandardAk(t *testing.T) {
	testGainSequence(t, StandardAk(rand.Float64()*100, rand.Float64()*100, rand.Float64()))
}

func TestStandardCk(t *testing.T) {
	testGainSequence(t, StandardCk(rand.Float64()*100, rand.Float64()))
}

func testGainSequence(t *testing.T, g GainSequence) {
	last := <-g
	for i := 0; i < 100; i++ {
		cur := <-g
		if cur >= last {
			t.Error("GainSequence is not monotonically decreasing.")
		} else if cur <= 0 {
			t.Error("GainSequence is not positive.")
		}
	}
}

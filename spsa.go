// Package spsa provides the Simultaneous Perturbation Stochastic Approximation method.
//
// Much of the notation is taken from Introduction To Stochastic Search and Optimization
// (ISSO), James Spall's book on Stochastic Optimization. (published by Wiley 2003)
package spsa

import (
	"math/rand"
	"math"
)

//********** Type Definitions ************

// Gain sequences are infinite iterators of floats. The must follow the conditions
// specified in ISSO.
type GainSequence <-chan float64

// A perturbation distribution is used to simultaneously perturb the otimization
// criteria to approximate the loss function's gradient. It must have special
// properties, the most restrictive is E[1/X] is bounded. This rules out
// uniform and normal. The asymptotically optimal distribution is Bernoulli +/- 1.
type PerturbationDistribution interface {
	Sample() float64
}

// A loss function is a vector-valued to real function. It will be minimized in SPSA.
// (Negate maximization functions to act as Loss functions.)
type LossFunction func(Vector) float64

// Map the parameter vector to a constrained version of itself.
type ConstraintFunction func(Vector) Vector

// An instance of the SPSA optimization algorithm.
// Initialize with all the parameters as object instantiation.
type SPSA struct {
	// The parameter vector in question. Initialize with Theta0 starting point.
	Theta Vector

	L LossFunction
	Ak, Ck GainSequence
	Delta PerturbationDistribution
	C ConstraintFunction
}

//****************** SPSA Implementation ****************

// A helper function to optimize a loss function using SPSA using mostly default options.
// It uses standard ak and ck gain sequences, bernoulli +/- 1 perturbation distribution
// and n rounds. The constraint function is optional.
func Optimize(L LossFunction, theta0 Vector, n int, a, c float64, C ...ConstraintFunction) Vector {
	constraint := NoConstraints
	if len(C) > 0 {
		constraint = C[0]
	}

	spsa := &SPSA{
		Theta: theta0,
		L: L,
		Ak: StandardAk(a, float64(n / 10), .602),
		Ck: StandardCk(c, .101),
		Delta: Bernoulli{1},
		C: constraint,
	}

	return spsa.Run(n)
}

// Helper function to run many rounds of SPSA and return the current Theta value.
func (spsa *SPSA) Run(rounds int) Vector {
	for i := 0; i < rounds; i++ {
		spsa.round()
	}
	return spsa.Theta
}

// Run one round of SPSA.
func (spsa *SPSA) round() {
	// Estimate gradient and scale it by ak
	Gk := spsa.estimateGradient().Scale(<- spsa.Ak)

	// Adjust theta via SA
	spsa.Theta = spsa.Theta.Subtract(Gk)

	// Correct any constraints
	spsa.Theta = spsa.C(spsa.Theta)
}

// Estimate the gradient in one round of spsa
func (spsa *SPSA) estimateGradient() Vector {
	n := len(spsa.Theta)

	// Get delta vector
	delta := SampleN(n, spsa.Delta).Scale(<- spsa.Ck)

	// Evaluate theta + ck * delta
	tpos := spsa.Theta.Add(delta)
	fpos := spsa.L(tpos)

	// Evaluate theta - ck * delta
	tneg := spsa.Theta.Subtract(delta)
	fneg := spsa.L(tneg)

	// Calculate estimated gradient
	grad := make([]float64, n)
	for i, d := range delta {
		grad[i] = (fpos - fneg) / (2 * d)
	}

	return grad
}

//********** Constrain function helpers ***********

// A ConstraintFunction that is just the identity mapper
func NoConstraints(a Vector) Vector {
	return a
}

// A pair of bounds on a variable. The first is the lower bound and the
// second is the upper bound
type Bounds struct {
	Lower, Upper float64
}

// An array of bounds on an array of variables. This object's Constrain function
// can be used as a ConstraintFunction for SPSA.
type BoundedConstraints []Bounds

// Constrain theta by mapping each value into its bounded domain. (in place)
func (bc BoundedConstraints) Constrain(theta Vector) Vector {
	for i, t := range theta {
		theta[i] = math.Min(math.Max(t, bc[i].Lower), bc[i].Upper)
	}
	return theta
}

//********** Gain Sequences *************

// Create an infinite iterator of a_k gain values in standard form.
// Standard form is a_k = a / (k + 1 + A) ^ alpha
// Semiauomatic tuning says that A should be roughly 10% of the iterations,
// and alpha should be .602. However, if running for a very long time,
// alpha = 1.0 might be more advantageous.
func StandardAk(a, A, alpha float64) GainSequence {
	c := make(chan float64)
	go func() {
		for k := 1; true; k++ {
			c <- a / math.Pow(float64(k)+A, alpha)
		}
	}()
	return GainSequence(c)
}

// Create an infinite iterator of c_k gain values in standard form.
// Standard form is c_k = c / (k + 1) ^ gamma
// Semiautomatic tuning says that c = sqrt(Var(L(x))) where L is the loss function
// being optimized. Gamma has restrictions based on alpha (see ISSO),
// but for best results in finite samples, use gamma = .101.
func StandardCk(c, gamma float64) GainSequence {
	return StandardAk(c, 0, gamma)
}


//********** Perturbation Distribution *************

func SampleN(n int, d PerturbationDistribution) Vector {
	a := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = d.Sample()
	}
	return a
}

// The bernoulli +/- r distribution.
type Bernoulli struct {
	r float64
}

func (b Bernoulli) Sample() float64 {
	if rand.Float32() > .5 {
		return b.r
	} else {
		return -b.r
	}
}

// The segmented/mirrored uniform distribution. Samples with equal probability
// all real numbers in [a,b] U [-b,-a] where 0 < a < b.
type SegmentedUniform struct {
	a, b float64
}

func (su SegmentedUniform) Sample() float64 {
	r := rand.Float64() - .5
	return math.Copysign(r, math.Abs(r) * 2 * (su.b - su.a) + su.a)
}
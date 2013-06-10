# SPSA in Go [![Build Status][1]][2]

[Simultaneous Perturbation Stochastic Approximation](http://jhuapl.edu/SPSA) is a global optimizer for continuous loss functions of any dimension. It has a strong theoretical foundation with very few knobs that must be tuned. While it doesn't get the press of Genetic Algorithms or other Evolutionary Computing, it is on the same level with a better foundation.

## Documentation

See [godoc](http://godoc.org/github.com/yanatan16/spsa) or the [examples](https://github.com/blob/master/spsa_test.go)

## Using SPSA

SPSA minimizes a loss function, so maximization problems must negate the fitness function. It can optionally use a constraint mapper that runs each round of SPSA. SPSA is an iterative or recursive optimization algorithm based as a stochastic extension of gradient descent. Like the Finite Difference (FDSA) approximation to the gradient, SPSA uses Simultaneous Perturbation to estimate the gradient of the loss function being minimized. However, it only uses two loss measurements, regardless of the dimension of the parameter vector, whereas FDSA uses 2p where p is the dimension of the parameter vector. Surprisingly, this shortcut has no ill effect on convergence rate and only a small effect on convergence criteria.

Since SPSA uses a small amount of randomness in its gradient estimate, it also produces some noise. This noise is the good kind however, because it promotes SPSA to become a global optimizer with the same convergence rate as FDSA!

SPSA has fewer tuning knobs than many other global optimization methods, but it still requires some work. The most important parameter is `a` in the `a_k = a / (k + 1 + A) ^ &alpha;` gain sequence. It can wildly affect results. If you wish to limit function measurements, since SPSA uses two function measurements per iteration, pass in N / 2 as the number of rounds to run when using SPSA.

For more information on SPSA, please see Spall's papers from his [website](http://jhuapl.edu/SPSA).

## References

- [SPSA Website](http://jhuapl.edu/SPSA)
- Introduction to Stochastic Search and Optimization. James Spall. Wiley 2003.

## License

The MIT License found in the LICENSE file.

[1]: https://travis-ci.org/yanatan16/golang-spsa.png?branch=master
[2]: http://travis-ci.org/yanatan16/golang-spsa
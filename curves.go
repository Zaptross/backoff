package backoff

import "math"

// Default is the recommended default curve for backoff. It is a logistic curve
// which generates values in a sigmoid or S-curve shape based on the maximum
// number of attempts.
func Default(attempts int, limit float64) func(float64) float64 {
	incline := 1 / (5 / float64(attempts))
	return func(x float64) float64 {
		return Logistic(x, incline, limit, float64(attempts))
	}
}

// Logistic is a function that returns a value based on the logistic function.
// It generates values in a sigmoid or S-curve shape.
// https://en.wikipedia.org/wiki/Logistic_function
//
// `x` is the input value.
// `k` is the steepness of the curve.
// `L` is the curve's maximum value.
// `x0` is the x-value of the sigmoid's midpoint.
func Logistic(x, k, L, x0 float64) float64 {
	return L / (1 + math.E*math.Exp(-k*(x-x0)))
}

// Linear is a function that returns a value based on a linear function.
// It generates values in a straight line.
//
// `x` is the input value.
// `mul` is the multiplier.
func Linear(x float64, mul float64) float64 {
	return x * mul
}

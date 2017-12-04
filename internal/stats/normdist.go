// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// https://github.com/golang/perf

package stats

import "math"

// NormalDist is a normal (Gaussian) distribution with mean Mu and
// standard deviation Sigma.
type NormalDist struct {
	Mu, Sigma float64
}

// StdNormal is the standard normal distribution (Mu = 0, Sigma = 1)
var StdNormal = NormalDist{0, 1}

// 1/sqrt(2 * pi)
const invSqrt2Pi = 0.39894228040143267793994605993438186847585863116493465766592583

// PDF is the probability density of the normal distribution
func (n NormalDist) PDF(x float64) float64 {
	z := x - n.Mu
	return math.Exp(-z*z/(2*n.Sigma*n.Sigma)) * invSqrt2Pi / n.Sigma
}

// CDF is the cumulative distribution of the normal distribution
func (n NormalDist) CDF(x float64) float64 {
	return math.Erfc(-(x-n.Mu)/(n.Sigma*math.Sqrt2)) / 2
}

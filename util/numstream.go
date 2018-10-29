package util

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// NumStream accumulates stats over a numeric stream.
type NumStream struct{ vv *Vector }

// WrapNumStream inits a new stream observer from a vector.
func WrapNumStream(vv *Vector) NumStream {
	if vv == nil {
		panic("reason: received nil vector")
	}
	return NumStream{vv: vv}
}

// NewNumStream inits a new stream observer.
func NewNumStream() NumStream {
	return WrapNumStream(NewVector())
}

// Observe adds a new observation.
func (s NumStream) Observe(value float64) {
	s.ObserveWeight(value, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (s NumStream) ObserveWeight(value, weight float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) || weight <= 0 {
		return
	}

	wv := weight * value
	s.vv.Add(0, weight)
	s.vv.Add(1, wv)
	s.vv.Add(2, wv*value)
}

// TotalWeight returns the total weight observed.
func (s NumStream) TotalWeight() float64 {
	return s.vv.At(0)
}

// Sum returns the weighted sum of all observations.
func (s NumStream) Sum() float64 {
	return s.vv.At(1)
}

// Mean returns a mean average
func (s NumStream) Mean() float64 {
	return s.Sum() / s.TotalWeight()
}

// Variance is the sample variance of the series
func (s NumStream) Variance() float64 {
	w := s.TotalWeight()
	if w <= 1 {
		return math.NaN()
	}
	z := s.Sum()
	return (s.vv.At(2) - (z * z / w)) / (w - 1)
}

// StdDev is the sample standard deviation of the series
func (s NumStream) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

// Prob calculates the gaussian probability density of a value
func (s NumStream) Prob(value float64) float64 {
	if sig := s.StdDev(); !math.IsNaN(sig) {
		return distuv.Normal{Mu: s.Mean(), Sigma: sig}.Prob(value)
	}
	return math.NaN()
}

// Estimate estimates weight boundaries for a given value
func (s NumStream) Estimate(value float64) (lessThan float64, equalTo float64, greaterThan float64) {
	total := s.TotalWeight()
	equalTo = s.Prob(value) * total

	if sig := s.StdDev(); !math.IsNaN(sig) {
		lessThan = distuv.Normal{
			Mu:    s.Mean(),
			Sigma: sig,
		}.CDF(value)*total - equalTo
	} else {
		lessThan = math.NaN()
	}

	if greaterThan = total - equalTo - lessThan; greaterThan < 0 {
		greaterThan = 0
	}

	return
}

// --------------------------------------------------------------------

// NumStreams accumulates stats of multiple numeric streams.
type NumStreams struct{ mat *Matrix }

// WrapNumStreams init stats with a matrix.
func WrapNumStreams(mat *Matrix) NumStreams {
	if mat == nil {
		panic("reason: received nil matrix")
	}
	return NumStreams{mat: mat}
}

// NewNumStreams inits a new streams observer.
func NewNumStreams() NumStreams {
	return WrapNumStreams(NewMatrix())
}

// Observe adds a new observation.
func (s NumStreams) Observe(cat int, value float64) {
	s.ObserveWeight(cat, value, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (s NumStreams) ObserveWeight(cat int, value, weight float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) || weight <= 0 {
		return
	}

	wv := weight * value
	s.mat.Add(cat, 0, weight)
	s.mat.Add(cat, 1, wv)
	s.mat.Add(cat, 2, wv*value)
}

// ForEach iterates over each category.
func (s NumStreams) ForEach(iter func(int)) {
	rows := s.NumRows()
	for i := 0; i < rows; i++ {
		if s.mat.At(i, 0) != 0 {
			iter(i)
		}
	}
}

// NumRows returns the number of rows, including blanks.
func (s NumStreams) NumRows() int {
	return s.mat.NumRows()
}

// NumCategories returns the number of categories.
func (s NumStreams) NumCategories() int {
	n := 0
	s.ForEach(func(_ int) { n++ })
	return n
}

// TotalWeight returns the total weight observed for cat.
func (s NumStreams) TotalWeight(cat int) float64 {
	return s.mat.At(cat, 0)
}

// Sum returns the weighted sum of all observations of cat.
func (s NumStreams) Sum(cat int) float64 {
	return s.mat.At(cat, 1)
}

// Mean returns a mean average of ovserved cat values.
func (s NumStreams) Mean(cat int) float64 {
	return s.Sum(cat) / s.TotalWeight(cat)
}

// Variance is the sample variance of the cat series.
func (s NumStreams) Variance(cat int) float64 {
	w := s.TotalWeight(cat)
	if w <= 1 {
		return math.NaN()
	}
	z := s.Sum(cat)
	return (s.mat.At(cat, 2) - (z * z / w)) / (w - 1)
}

// StdDev is the sample standard deviation of the cat series.
func (s NumStreams) StdDev(cat int) float64 {
	return math.Sqrt(s.Variance(cat))
}

// Prob calculates the gaussian probability density of a value for cat.
func (s NumStreams) Prob(cat int, value float64) float64 {
	if sig := s.StdDev(cat); !math.IsNaN(sig) {
		return distuv.Normal{Mu: s.Mean(cat), Sigma: sig}.Prob(value)
	}
	return math.NaN()
}

// Estimate estimates weight boundaries for a given cat value.
func (s NumStreams) Estimate(cat int, value float64) (lessThan float64, equalTo float64, greaterThan float64) {
	total := s.TotalWeight(cat)
	equalTo = s.Prob(cat, value) * total

	if sig := s.StdDev(cat); !math.IsNaN(sig) {
		lessThan = distuv.Normal{
			Mu:    s.Mean(cat),
			Sigma: sig,
		}.CDF(value)*total - equalTo
	} else {
		lessThan = math.NaN()
	}

	if greaterThan = total - equalTo - lessThan; greaterThan < 0 {
		greaterThan = 0
	}

	return
}

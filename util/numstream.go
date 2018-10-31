package util

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// NewNumStream inits a new stream observer.
func NewNumStream() *NumStream {
	return new(NumStream)
}

// Observe adds a new observation.
func (s *NumStream) Observe(value float64) {
	s.ObserveWeight(value, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (s *NumStream) ObserveWeight(value, weight float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) || weight <= 0 {
		return
	}

	if s.Weight == 0 || value < s.Min {
		s.Min = value
	}
	if s.Weight == 0 || value > s.Max {
		s.Max = value
	}

	wv := weight * value
	s.Weight += weight
	s.Sum += wv
	s.SumSquares += wv * value

}

// Mean returns a mean average
func (s *NumStream) Mean() float64 {
	return s.Sum / s.Weight
}

// Variance is the sample variance of the series
func (s *NumStream) Variance() float64 {
	if s.Weight <= 1 {
		return math.NaN()
	}
	return (s.SumSquares - (s.Sum * s.Sum / s.Weight)) / (s.Weight - 1)
}

// StdDev is the sample standard deviation of the series
func (s *NumStream) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

// Prob calculates the gaussian probability density of a value
func (s *NumStream) Prob(value float64) float64 {
	if sig := s.StdDev(); !math.IsNaN(sig) {
		return distuv.Normal{Mu: s.Mean(), Sigma: sig}.Prob(value)
	}
	return math.NaN()
}

// Estimate estimates weight boundaries for a given value
func (s *NumStream) Estimate(value float64) (lessThan float64, equalTo float64, greaterThan float64) {
	equalTo = s.Prob(value) * s.Weight

	if sig := s.StdDev(); !math.IsNaN(sig) {
		lessThan = distuv.Normal{
			Mu:    s.Mean(),
			Sigma: sig,
		}.CDF(value)*s.Weight - equalTo
	} else {
		lessThan = math.NaN()
	}

	if greaterThan = s.Weight - equalTo - lessThan; greaterThan < 0 {
		greaterThan = 0
	}

	return
}

// --------------------------------------------------------------------

// NewNumStreams inits a new numeric streams distribution.
func NewNumStreams() *NumStreams {
	return new(NumStreams)
}

// Observe adds a new observation.
func (s *NumStreams) Observe(cat int, value float64) {
	s.ObserveWeight(cat, value, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (s *NumStreams) ObserveWeight(cat int, value, weight float64) {
	if cat < 0 || math.IsNaN(value) || math.IsInf(value, 0) || weight <= 0 {
		return
	}

	if n := cat + 1; n > cap(s.Data) {
		data := make([]NumStream, n, 2*n)
		copy(data, s.Data)
		s.Data = data
	} else if n > len(s.Data) {
		s.Data = s.Data[:n]
	}
	s.Data[cat].ObserveWeight(value, weight)
}

// NumRows returns the number of rows, including blanks.
func (s *NumStreams) NumRows() int {
	return len(s.Data)
}

// NumCategories returns the number of categories.
func (s *NumStreams) NumCategories() int {
	n := 0
	for _, t := range s.Data {
		if t.Weight > 0 {
			n++
		}
	}
	return n
}

// WeightSum returns the total weight observed.
func (s *NumStreams) WeightSum() float64 {
	sum := 0.0
	for _, t := range s.Data {
		sum += t.Weight
	}
	return sum
}

// At returns the stream at the given cat (or nil).
// Please note that a copy of the stream is returned. Mutating the
// value will not affect the distribution.
func (s *NumStreams) At(cat int) *NumStream {
	if cat > -1 && cat < len(s.Data) {
		if t := s.Data[cat]; t.Weight > 0 {
			return &t
		}
	}
	return nil
}

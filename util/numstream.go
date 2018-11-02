package util

import (
	"math"
	"sort"

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
func (s *NumStreams) Observe(target int, prediction float64) {
	s.ObserveWeight(target, prediction, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (s *NumStreams) ObserveWeight(target int, prediction, weight float64) {
	if target < 0 || math.IsNaN(prediction) || math.IsInf(prediction, 0) || weight <= 0 {
		return
	}

	if n := target + 1; n > cap(s.Data) {
		data := make([]NumStream, n, 2*n)
		copy(data, s.Data)
		s.Data = data
	} else if n > len(s.Data) {
		s.Data = s.Data[:n]
	}
	s.Data[target].ObserveWeight(prediction, weight)
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

// At returns the stream at the given target (or nil).
// Please note that a copy of the stream is returned. Mutating the
// value will not affect the distribution.
func (s *NumStreams) At(target int) *NumStream {
	if target > -1 && target < len(s.Data) {
		if t := s.Data[target]; t.Weight > 0 {
			return &t
		}
	}
	return nil
}

// --------------------------------------------------------------------

// NewNumStreamBuckets inits new bucketed stream distribution.
func NewNumStreamBuckets(maxBuckets uint32) *NumStreamBuckets {
	return &NumStreamBuckets{MaxBuckets: maxBuckets}
}

// Observe adds a new observation.
func (s *NumStreamBuckets) Observe(target, prediction float64) {
	s.ObserveWeight(target, prediction, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (s *NumStreamBuckets) ObserveWeight(target, prediction, weight float64) {
	if math.IsNaN(target) || math.IsInf(target, 0) || math.IsNaN(prediction) || math.IsInf(prediction, 0) || weight <= 0 {
		return
	}

	// ensure buckets are set
	if s.MaxBuckets == 0 {
		s.MaxBuckets = 12
	}

	// upsert bucket
	slot := s.findSlot(target)
	if slot < len(s.Buckets) {
		if s.Buckets[slot].Threshold != target {
			s.Buckets = append(s.Buckets, NumStreamBuckets_Bucket{})
			copy(s.Buckets[slot+1:], s.Buckets[slot:])
			s.Buckets[slot] = NumStreamBuckets_Bucket{Threshold: target}
		}
	} else {
		s.Buckets = append(s.Buckets, NumStreamBuckets_Bucket{Threshold: target})
	}
	s.Buckets[slot].ObserveWeight(prediction, weight)

	// prune buckets
	for uint32(len(s.Buckets)) > s.MaxBuckets {
		delta := math.MaxFloat64
		slot := 0
		for i := 0; i < len(s.Buckets)-1; i++ {
			if x := s.Buckets[i+1].Threshold - s.Buckets[i].Threshold; x < delta {
				slot, delta = i, x
			}
		}

		b1, b2 := s.Buckets[slot], s.Buckets[slot+1]
		weightSum := b1.Weight + b2.Weight
		threshold := (b1.Threshold*b1.Weight + b2.Threshold*b2.Weight) / weightSum
		s.Buckets[slot+1] = NumStreamBuckets_Bucket{
			Threshold: threshold,
			NumStream: NumStream{
				Weight:     weightSum,
				Sum:        b1.Sum + b2.Sum,
				SumSquares: b1.SumSquares + b2.SumSquares,
			},
		}
		s.Buckets = s.Buckets[:slot+copy(s.Buckets[slot:], s.Buckets[slot+1:])]
	}
}

// WeightSum returns the total weight observed.
func (s *NumStreamBuckets) WeightSum() float64 {
	sum := 0.0
	for _, b := range s.Buckets {
		sum += b.Weight
	}
	return sum
}

func (s *NumStreamBuckets) findSlot(v float64) int {
	return sort.Search(len(s.Buckets), func(i int) bool { return s.Buckets[i].Threshold >= v })
}

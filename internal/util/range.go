package util

import "math"

type NumRange struct {
	Min, Max float64
}

// NewNumRange creates a range with min set to +infinity and max set to -infinity
func NewNumRange() *NumRange {
	r := new(NumRange)
	r.Reset()
	return r
}

// Reset sets min to +infinity and max to -infinity
func (r *NumRange) Reset() {
	r.Min = math.Inf(1)
	r.Max = math.Inf(-1)
}

// SplitPoints splits range into n parts and returns the split points
func (r *NumRange) SplitPoints(n int) []float64 {
	delta := r.Max - r.Min
	if delta <= 0 {
		return nil
	}

	min := r.Min
	stp := delta / float64(n+1)
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = min + stp*float64(i+1)
	}
	return res
}

// Update updates range by includin the value
func (r *NumRange) Update(v float64) {
	if v < r.Min {
		r.Min = v
	}
	if v > r.Max {
		r.Max = v
	}
}

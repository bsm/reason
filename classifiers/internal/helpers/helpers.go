package helpers

import (
	"math"

	"github.com/bsm/reason/util"
)

type minMaxRange struct{ Min, Max float64 }

// newMinMaxRange creates a range with min set to +infinity and max set to -infinity
func newMinMaxRange() *minMaxRange {
	r := new(minMaxRange)
	r.Reset()
	return r
}

// Reset sets min to +infinity and max to -infinity
func (r *minMaxRange) Reset() {
	r.Min = math.Inf(1)
	r.Max = math.Inf(-1)
}

// SplitPoints splits range into n parts and returns the split points
func (r *minMaxRange) SplitPoints(n int) []float64 {
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
func (r *minMaxRange) Update(v float64) {
	if v < r.Min {
		r.Min = v
	}
	if v > r.Max {
		r.Max = v
	}
}

// --------------------------------------------------------------------

type minMaxRanges struct {
	min, max util.Vector
}

func newMinMaxRanges() *minMaxRanges {
	return &minMaxRanges{
		min: util.NewVector(),
		max: util.NewVector(),
	}
}

func (m *minMaxRanges) Len() int          { return m.min.Count() }
func (m *minMaxRanges) ByteSize() int     { return 16 + m.min.ByteSize() + m.max.ByteSize() }
func (m *minMaxRanges) Min(i int) float64 { return m.min.Get(i) }
func (m *minMaxRanges) Max(i int) float64 { return m.max.Get(i) }
func (m *minMaxRanges) Update(i int, val float64) {
	if min := m.Min(i); min == 0 {
		m.min = m.min.Set(i, val)
		m.max = m.max.Set(i, val)
	} else {
		if val < min {
			m.min = m.min.Set(i, val)
		}
		if val > m.Max(i) {
			m.max = m.max.Set(i, val)
		}
	}
}

func (m *minMaxRanges) SplitPoints(n int) []float64 {
	rng := newMinMaxRange()
	m.min.ForEach(func(i int, _ float64) {
		rng.Update(m.Min(i))
		rng.Update(m.Max(i))
	})
	return rng.SplitPoints(n)
}

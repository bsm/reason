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
	min, max util.SparseVector
}

func newMinMaxRanges() *minMaxRanges {
	return &minMaxRanges{
		min: util.NewSparseVector(),
		max: util.NewSparseVector(),
	}
}

func (m *minMaxRanges) Len() int            { return len(m.min) }
func (m *minMaxRanges) Min(pos int) float64 { return m.min.Get(pos) }
func (m *minMaxRanges) Max(pos int) float64 { return m.max.Get(pos) }
func (m *minMaxRanges) Update(pos int, val float64) {
	if _, ok := m.min[pos]; ok {
		if val < m.Min(pos) {
			m.min.Set(pos, val)
		}
		if val > m.Max(pos) {
			m.max.Set(pos, val)
		}
	} else {
		m.min.Set(pos, val)
		m.max.Set(pos, val)
	}
}

func (m *minMaxRanges) SplitPoints(n int) []float64 {
	rng := newMinMaxRange()
	for i := range m.min {
		rng.Update(m.Min(i))
		rng.Update(m.Max(i))
	}
	return rng.SplitPoints(n)
}

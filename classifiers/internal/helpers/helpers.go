package helpers

import (
	"math"

	"github.com/bsm/reason/internal/msgpack"
	"github.com/bsm/reason/util"
)

func init() {
	msgpack.Register(7735, (*MinMaxRange)(nil))
	msgpack.Register(7736, (*MinMaxRanges)(nil))
	msgpack.Register(7746, Observation{})
}

type Observation struct{ PVal, TVal, Weight float64 }

func (o Observation) EncodeTo(enc *msgpack.Encoder) error {
	return enc.Encode(o.PVal, o.TVal, o.Weight)
}

func (o *Observation) DecodeFrom(dec *msgpack.Decoder) error {
	return dec.Decode(&o.PVal, &o.TVal, &o.Weight)
}

// --------------------------------------------------------------------

type MinMaxRange struct{ Min, Max float64 }

// NewMinMaxRange creates a range with min set to +infinity and max set to -infinity
func NewMinMaxRange() *MinMaxRange {
	r := new(MinMaxRange)
	r.Reset()
	return r
}

// Reset sets min to +infinity and max to -infinity
func (r *MinMaxRange) Reset() {
	r.Min = math.Inf(1)
	r.Max = math.Inf(-1)
}

// SplitPoints splits range into n parts and returns the split points
func (r *MinMaxRange) SplitPoints(n int) []float64 {
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
func (r *MinMaxRange) Update(v float64) {
	if v < r.Min {
		r.Min = v
	}
	if v > r.Max {
		r.Max = v
	}
}

func (r *MinMaxRange) EncodeTo(enc *msgpack.Encoder) error   { return enc.Encode(r.Min, r.Max) }
func (r *MinMaxRange) DecodeFrom(dec *msgpack.Decoder) error { return dec.Decode(&r.Min, &r.Max) }

// --------------------------------------------------------------------

type MinMaxRanges struct {
	Min, Max util.Vector
}

func NewMinMaxRanges() *MinMaxRanges {
	return &MinMaxRanges{
		Min: util.NewVector(),
		Max: util.NewVector(),
	}
}

func (m *MinMaxRanges) Len() int             { return m.Min.Count() }
func (m *MinMaxRanges) ByteSize() int        { return 16 + m.Min.ByteSize() + m.Max.ByteSize() }
func (m *MinMaxRanges) GetMin(i int) float64 { return m.Min.Get(i) }
func (m *MinMaxRanges) GetMax(i int) float64 { return m.Max.Get(i) }
func (m *MinMaxRanges) Update(i int, val float64) {
	if min := m.GetMin(i); min == 0 {
		m.Min = m.Min.Set(i, val)
		m.Max = m.Max.Set(i, val)
	} else {
		if val < min {
			m.Min = m.Min.Set(i, val)
		}
		if val > m.GetMax(i) {
			m.Max = m.Max.Set(i, val)
		}
	}
}

func (m *MinMaxRanges) SplitPoints(n int) []float64 {
	rng := NewMinMaxRange()
	m.Min.ForEach(func(i int, _ float64) {
		rng.Update(m.GetMin(i))
		rng.Update(m.GetMax(i))
	})
	return rng.SplitPoints(n)
}

func (m *MinMaxRanges) EncodeTo(enc *msgpack.Encoder) error   { return enc.Encode(m.Min, m.Max) }
func (m *MinMaxRanges) DecodeFrom(dec *msgpack.Decoder) error { return dec.Decode(&m.Min, &m.Max) }

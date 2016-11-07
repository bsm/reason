package util

import "github.com/bsm/reason/internal/calc"

// NumVector is a slice of float values with helpers
type NumVector []float64

func (vv NumVector) Count() int       { return len(vv) }
func (vv NumVector) Sum() float64     { return calc.Sum(vv) }
func (vv NumVector) Min() float64     { return calc.Min(vv) }
func (vv NumVector) Max() float64     { return calc.Max(vv) }
func (vv NumVector) Mean() float64    { return calc.Mean(vv) }
func (vv NumVector) Entropy() float64 { return calc.Entropy(vv) }

func (vv NumVector) Variance() float64 { return calc.Variance(vv) }
func (vv NumVector) StdDev() float64   { return calc.StdDev(vv) }

func (vv NumVector) SampleVariance() float64 { return calc.SampleVariance(vv) }
func (vv NumVector) SampleStdDev() float64   { return calc.SampleStdDev(vv) }

func (vv NumVector) Normalize() NumVector {
	size := len(vv)
	if size == 0 {
		return nil
	}

	sum := vv.Sum()
	nv := make(NumVector, size)
	for i, v := range vv {
		nv[i] = v / sum
	}
	return nv
}

func (vv NumVector) Get(i int) float64 {
	if i > -1 && i < len(vv) {
		return vv[i]
	}
	return 0.0
}

func (vv NumVector) Set(i int, v float64) NumVector {
	vv = vv.grow(i + 1)
	vv[i] = v
	return vv
}

func (vv NumVector) Incr(i int, v float64) NumVector {
	vv = vv.grow(i + 1)
	vv[i] += v
	return vv
}

func (vv NumVector) grow(n int) NumVector {
	if n > len(vv) {
		nv := make(NumVector, n)
		copy(nv, vv)
		vv = nv
	}
	return vv
}

// BoolVector is a slice of bool values with helpers
type BoolVector []bool

func (vv BoolVector) Get(i int) bool {
	if i > -1 && i < len(vv) {
		return vv[i]
	}
	return false
}

func (vv BoolVector) Set(i int, v bool) BoolVector {
	vv = vv.grow(i + 1)
	vv[i] = v
	return vv
}

func (vv BoolVector) grow(n int) BoolVector {
	if n > len(vv) {
		nv := make(BoolVector, n)
		copy(nv, vv)
		vv = nv
	}
	return vv
}

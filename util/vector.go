package util

import (
	"math"

	"github.com/bsm/reason/internal/msgpack"
)

func init() {
	msgpack.Register(7733, (*DenseVector)(nil))
	msgpack.Register(7734, (SparseVector)(nil))
}

// Vector interface represents number vectors
type Vector interface {
	// Count returns the number of non-zero elements in the vector
	Count() int
	// Get returns a value
	Get(int) float64
	// Set sets a value
	Set(int, float64) Vector
	// Incr increments a value
	Incr(int, float64) Vector
	// ForEach iterates over each index/value
	ForEach(VectorIterator)
	// ForEachValue iterates over each value
	ForEachValue(VectorValueIterator)
	// Clear removes all values from the vector
	Clear()
	// Normalize normalizes all values to the 0..1 range
	Normalize()
	// Sum returns the sum of values
	Sum() float64
	// Min returns the minimum value
	Min() float64
	// Max returns the maximum value
	Max() float64
	// Mean returns the mean value
	Mean() float64
	// Variance calculates the variance
	Variance() float64
	// StdDev calculates the standard deviation
	StdDev() float64
	// Entropy calculates the entropy of the vector
	Entropy() float64
	// ByteSize estimates the required heap-size
	ByteSize() int
}

// NewVector creates a new default vector
func NewVector() Vector { return NewDenseVector() }

type VectorIterator func(int, float64)
type VectorValueIterator func(float64)

// --------------------------------------------------------------------

// Sum returns the sum
func vectorSum(vv Vector) (Σ float64) {
	vv.ForEachValue(func(v float64) {
		Σ += v
	})
	return
}

func vectorMin(vv Vector) float64 {
	min := math.MaxFloat64
	vv.ForEachValue(func(v float64) {
		if v < min {
			min = v
		}
	})
	return min
}

func vectorMax(vv Vector) float64 {
	max := -math.MaxFloat64
	vv.ForEachValue(func(v float64) {
		if v > max {
			max = v
		}
	})
	return max
}

func vectorMean(vv Vector) (µ float64) {
	n := vv.Count()
	if n == 0 {
		return
	}

	µ = vv.Sum() / float64(n)
	return
}

func vectorVariance(vv Vector) (V float64) {
	n := vv.Count()
	if n == 0 {
		return
	}

	µ := vv.Mean()
	vv.ForEachValue(func(v float64) {
		Δ := v - µ
		V += Δ * Δ
	})

	V /= float64(n)
	return
}

func vectorSampleVariance(vv Vector) (sV float64) {
	n := vv.Count()
	if n < 2 {
		return
	}

	µ := vv.Mean()
	vv.ForEachValue(func(v float64) {
		Δ := v - µ
		sV += Δ * Δ
	})

	sV /= float64(n - 1)
	return
}

func vectorStdDev(v float64) (σ float64) {
	if v == 0 {
		return
	}

	σ = math.Sqrt(v)
	return
}

func vectorEntropy(vv Vector) float64 {
	ent := 0.0
	sum := 0.0
	vv.ForEachValue(func(v float64) {
		ent -= v * math.Log2(v)
		sum += v
	})
	if sum > 0 {
		return (ent + sum*math.Log2(sum)) / sum
	}
	return 0.0
}

// --------------------------------------------------------------------

const sparseVectorBaseSize = 8 * sizeOfInt

// SparseVector is a (sparse) series of float64 numbers
type SparseVector map[int]float64

// NewSparseVector returns a blank vector
func NewSparseVector() SparseVector {
	return make(SparseVector)
}

// Count returns the number of non-zero elements in the vector
func (vv SparseVector) Count() int { return len(vv) }

// Clone creates a copy of the vector
func (vv SparseVector) Clone() Vector {
	sv := make(SparseVector, len(vv))
	for i, v := range vv {
		sv[i] = v
	}
	return sv
}

// Get returns a value at index i
func (vv SparseVector) Get(i int) float64 {
	if i > -1 {
		return vv[i]
	}
	return 0.0
}

// Set sets a value v at index i
func (vv SparseVector) Set(i int, v float64) Vector {
	if i > -1 {
		vv[i] = v
	}
	return vv
}

// ForEach iterates over each index/value
func (vv SparseVector) ForEach(iter VectorIterator) {
	for i, v := range vv {
		iter(i, v)
	}
}

// ForEachValue iterates over each value
func (vv SparseVector) ForEachValue(iter VectorValueIterator) {
	for _, v := range vv {
		iter(v)
	}
}

// Incr increments a value at index i by delta
func (vv SparseVector) Incr(i int, delta float64) Vector {
	if i > -1 {
		vv[i] += delta
	}
	return vv
}

// Clear removes all values from the vector
func (vv SparseVector) Clear() {
	for i := range vv {
		delete(vv, i)
	}
}

// Sum returns the sum
func (vv SparseVector) Sum() float64 { return vectorSum(vv) }

// Min returns the minimum value
func (vv SparseVector) Min() float64 { return vectorMin(vv) }

// Max returns the maximum value
func (vv SparseVector) Max() float64 { return vectorMax(vv) }

// Mean returns the mean value
func (vv SparseVector) Mean() float64 { return vectorMean(vv) }

// Variance calculates the variance
func (vv SparseVector) Variance() float64 { return vectorVariance(vv) }

// StdDev calculates the standard deviation
func (vv SparseVector) StdDev() float64 { return vectorStdDev(vv.Variance()) }

// SampleVariance calculates the sample variance
func (vv SparseVector) SampleVariance() float64 { return vectorSampleVariance(vv) }

// SampleStdDev calculates the sample standard deviation
func (vv SparseVector) SampleStdDev() float64 { return vectorStdDev(vv.SampleVariance()) }

// Normalize normalizes all values to the 0..1 range
func (vv SparseVector) Normalize() {
	if vv.Count() == 0 {
		return
	}

	sum := vv.Sum()
	for i, v := range vv {
		vv[i] = v / sum
	}
}

// Entropy calculates the entropy of the vector
func (vv SparseVector) Entropy() float64 { return vectorEntropy(vv) }

// ByteSize estimates the required heap-size
func (vv SparseVector) ByteSize() int {
	return 24 + len(vv)*sparseVectorBaseSize
}

// --------------------------------------------------------------------

// DenseVector is a dense series of float64 numbers
type DenseVector struct {
	vv []float64
}

// NewDenseVector returns a blank vector
func NewDenseVector() *DenseVector {
	return new(DenseVector)
}

// NewDenseVectorFromSlice returns a vector using slice values
func NewDenseVectorFromSlice(vv ...float64) *DenseVector {
	return &DenseVector{vv: vv}
}

// Count returns the number of non-zero elements in the vector
func (x *DenseVector) Count() int {
	n := 0
	for _, v := range x.vv {
		if v != 0 {
			n++
		}
	}
	return n
}

// Clone creates a copy of the vector
func (x *DenseVector) Clone() Vector {
	nv := &DenseVector{
		vv: make([]float64, len(x.vv)),
	}
	copy(nv.vv, x.vv)
	return nv
}

// Get returns a value at index i
func (x *DenseVector) Get(i int) float64 {
	if i > -1 && i < len(x.vv) {
		return x.vv[i]
	}
	return 0.0
}

// Set sets a value v at index i
func (x *DenseVector) Set(i int, v float64) Vector {
	if i < 0 {
		return x
	}

	if n := (i + 1); n > len(x.vv) {
		if n > 200 {
			if cnt := x.Count(); cnt > 10 && cnt*10 < n {
				return x.convertToSparse()
			}
		}

		vv := make([]float64, n)
		copy(vv, x.vv)
		x.vv = vv
	}
	x.vv[i] = v
	return x
}

// Incr increments a value at index i by delta
func (x *DenseVector) Incr(i int, delta float64) Vector {
	if i < 0 {
		return x
	}

	if i < len(x.vv) {
		x.vv[i] += delta
		return x
	}
	return x.Set(i, delta)
}

// ForEach iterates over each index/value
func (x DenseVector) ForEach(iter VectorIterator) {
	for i, v := range x.vv {
		if v != 0 {
			iter(i, v)
		}
	}
}

// ForEachValue iterates over each value
func (x DenseVector) ForEachValue(iter VectorValueIterator) {
	for _, v := range x.vv {
		if v != 0 {
			iter(v)
		}
	}
}

// Clear removes all values from the vector
func (x *DenseVector) Clear() {
	x.vv = x.vv[:0]
}

// Sum returns the sum
func (x *DenseVector) Sum() float64 { return vectorSum(x) }

// Min returns the minimum value
func (x *DenseVector) Min() float64 { return vectorMin(x) }

// Max returns the maximum value
func (x *DenseVector) Max() float64 { return vectorMax(x) }

// Mean returns the mean value
func (x *DenseVector) Mean() float64 { return vectorMean(x) }

// Variance calculates the variance
func (x *DenseVector) Variance() float64 { return vectorVariance(x) }

// StdDev calculates the standard deviation
func (x *DenseVector) StdDev() float64 { return vectorStdDev(x.Variance()) }

// SampleVariance calculates the sample variance
func (x *DenseVector) SampleVariance() float64 { return vectorSampleVariance(x) }

// SampleStdDev calculates the sample standard deviation
func (x *DenseVector) SampleStdDev() float64 { return vectorStdDev(x.SampleVariance()) }

// Normalize normalizes all values to the 0..1 range
func (x *DenseVector) Normalize() {
	sum := x.Sum()
	if sum == 0 {
		return
	}

	for i, v := range x.vv {
		x.vv[i] = v / sum
	}
}

// Entropy calculates the entropy of the vector
func (x *DenseVector) Entropy() float64 { return vectorEntropy(x) }

// ByteSize estimates the required heap-size
func (x *DenseVector) ByteSize() int {
	return 24 + cap(x.vv)*8
}

func (x *DenseVector) EncodeTo(enc *msgpack.Encoder) error {
	if x.vv == nil {
		x.vv = make([]float64, 0)
	}
	return enc.Encode(x.vv)
}
func (x *DenseVector) DecodeFrom(dec *msgpack.Decoder) error { return dec.Decode(&x.vv) }

func (x *DenseVector) convertToSparse() SparseVector {
	nv := make(SparseVector, x.Count())
	x.ForEach(func(i int, v float64) { nv.Set(i, v) })
	return nv
}

// --------------------------------------------------------------------

// VectorDistribution is a distribution of vectors by predicate
type VectorDistribution map[int]Vector

// NewVectorDistribution creates a new distribution
func NewVectorDistribution() VectorDistribution {
	return make(VectorDistribution)
}

// NumPredicates returns the number of predicates
func (m VectorDistribution) NumPredicates() int { return len(m) }

// NumTargets returns the number of targets
func (m VectorDistribution) NumTargets() int {
	var n int
	for _, row := range m {
		if l := row.Count(); l > n {
			n = l
		}
	}
	return n
}

// Weights returns the weight distribution by predicate
func (m VectorDistribution) Weights() map[int]float64 {
	vv := make(map[int]float64, len(m))
	for i, row := range m {
		vv[i] = row.Sum()
	}
	return vv
}

// Get increments row at pos or nil
func (m VectorDistribution) Get(pos int) Vector {
	if pos > -1 {
		return m[pos]
	}
	return nil
}

// Incr increments value at row/col index
func (m VectorDistribution) Incr(row, col int, delta float64) {
	vv, ok := m[row]
	if !ok {
		vv = NewVector()
	}
	m[row] = vv.Incr(col, delta)
}

// ByteSize estimates the required heap-size
func (m VectorDistribution) ByteSize() int {
	size := 24
	for _, vv := range m {
		size += 16 + sizeOfInt + vv.ByteSize()
	}
	return size
}

package util

import "math"

// Vector interface represents number vectors
type Vector interface {
	// Count returns the number of non-zero elements in the vector
	Count() int
	// Get returns a value
	Get(int) float64
	// Set sets a value
	Set(int, float64)
	// Incr increments a value
	Incr(int, float64)
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
	// HeapSize estimates the required heap-size
	HeapSize() int
}

type VectorIterator func(int, float64)
type VectorValueIterator func(float64)

// --------------------------------------------------------------------

// SparseVector is a (sparse) series of float64 numbers
type SparseVector map[int]float64

// NewSparseVector returns a blank vector
func NewSparseVector() SparseVector {
	return make(SparseVector)
}

// NewSparseVectorFromSlice returns a vector using slice values
func NewSparseVectorFromSlice(vv ...float64) SparseVector {
	sv := make(SparseVector, len(vv))
	for i, v := range vv {
		sv.Set(i, v)
	}
	return sv
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
func (vv SparseVector) Set(i int, v float64) {
	if i > -1 {
		vv[i] = v
	}
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
func (vv SparseVector) Incr(i int, delta float64) {
	if i > -1 {
		vv[i] += delta
	}
}

// Clear removes all values from the vector
func (vv SparseVector) Clear() {
	for i := range vv {
		delete(vv, i)
	}
}

// Sum returns the sum
func (vv SparseVector) Sum() (Σ float64) {
	for _, v := range vv {
		Σ += v
	}
	return
}

// Min returns the minimum value
func (vv SparseVector) Min() float64 {
	if vv.Count() == 0 {
		return math.NaN()
	}

	min := math.MaxFloat64
	for _, v := range vv {
		if v < min {
			min = v
		}
	}
	return min
}

// Max returns the maximum value
func (vv SparseVector) Max() float64 {
	if vv.Count() == 0 {
		return math.NaN()
	}

	max := -math.MaxFloat64
	for _, v := range vv {
		if v > max {
			max = v
		}
	}
	return max
}

// Mean returns the mean value
func (vv SparseVector) Mean() (µ float64) {
	n := vv.Count()
	if n == 0 {
		return
	}

	µ = vv.Sum() / float64(n)
	return
}

// Variance calculates the variance
func (vv SparseVector) Variance() (V float64) {
	n := vv.Count()
	if n == 0 {
		return
	}

	µ := vv.Mean()
	for _, v := range vv {
		Δ := v - µ
		V += Δ * Δ
	}

	V /= float64(n)
	return
}

// StdDev calculates the standard deviation
func (vv SparseVector) StdDev() (σ float64) {
	V := vv.Variance()
	if V == 0 {
		return
	}

	σ = math.Sqrt(V)
	return
}

// SampleVariance calculates the sample variance
func (vv SparseVector) SampleVariance() (sV float64) {
	n := vv.Count()
	if n < 2 {
		return
	}

	µ := vv.Mean()
	for _, v := range vv {
		Δ := v - µ
		sV += Δ * Δ
	}

	sV /= float64(n - 1)
	return
}

// SampleStdDev calculates the sample standard deviation
func (vv SparseVector) SampleStdDev() (sσ float64) {
	sV := vv.SampleVariance()
	if sV == 0 {
		return
	}

	sσ = math.Sqrt(sV)
	return
}

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
func (vv SparseVector) Entropy() float64 {
	ent := 0.0
	sum := 0.0
	for _, v := range vv {
		if v > 0 {
			ent -= v * math.Log2(v)
			sum += v
		}
	}
	if sum > 0 {
		return (ent + sum*math.Log2(sum)) / sum
	}
	return 0.0
}

// HeapSize estimates the required heap-size
func (vv SparseVector) HeapSize() int {
	return 24 + len(vv)*60
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
func (x *DenseVector) Set(i int, v float64) {
	if i < 0 {
		return
	}

	if n := (i + 1); n > len(x.vv) {
		vv := make([]float64, n)
		copy(vv, x.vv)
		x.vv = vv
	}
	x.vv[i] = v
}

// Incr increments a value at index i by delta
func (x *DenseVector) Incr(i int, delta float64) {
	if i < 0 {
		return
	}

	if i < len(x.vv) {
		x.vv[i] += delta
	} else {
		x.Set(i, delta)
	}
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
func (x *DenseVector) Sum() (Σ float64) {
	for _, v := range x.vv {
		Σ += v
	}
	return
}

// Min returns the minimum value
func (x *DenseVector) Min() float64 {
	if x.Count() == 0 {
		return math.NaN()
	}

	min := math.MaxFloat64
	for _, v := range x.vv {
		if v != 0.0 && v < min {
			min = v
		}
	}
	return min
}

// Max returns the maximum value
func (x *DenseVector) Max() float64 {
	if x.Count() == 0 {
		return math.NaN()
	}

	max := -math.MaxFloat64
	for _, v := range x.vv {
		if v != 0.0 && v > max {
			max = v
		}
	}
	return max
}

// Mean returns the mean value
func (x *DenseVector) Mean() (µ float64) {
	n := x.Count()
	if n == 0 {
		return
	}

	µ = x.Sum() / float64(n)
	return
}

// Variance calculates the variance
func (x *DenseVector) Variance() (V float64) {
	n := 0
	µ := x.Mean()
	for _, v := range x.vv {
		if v != 0 {
			Δ := v - µ
			V += Δ * Δ
			n++
		}
	}
	if n < 1 {
		return 0.0
	}

	V /= float64(n)
	return
}

// StdDev calculates the standard deviation
func (x *DenseVector) StdDev() (σ float64) {
	V := x.Variance()
	if V == 0 {
		return
	}

	σ = math.Sqrt(V)
	return
}

// SampleVariance calculates the sample variance
func (x *DenseVector) SampleVariance() (sV float64) {
	n := 0
	µ := x.Mean()
	for _, v := range x.vv {
		if v != 0 {
			Δ := v - µ
			sV += Δ * Δ
			n++
		}
	}
	if n < 2 {
		return 0.0
	}

	sV /= float64(n - 1)
	return
}

// SampleStdDev calculates the sample standard deviation
func (x *DenseVector) SampleStdDev() (sσ float64) {
	sV := x.SampleVariance()
	if sV == 0 {
		return
	}

	sσ = math.Sqrt(sV)
	return
}

// Normalize normalizes all values to the 0..1 range
func (x *DenseVector) Normalize() {
	if x.Count() == 0 {
		return
	}

	sum := x.Sum()
	for i, v := range x.vv {
		x.vv[i] = v / sum
	}
}

// Entropy calculates the entropy of the vector
func (x *DenseVector) Entropy() float64 {
	ent := 0.0
	sum := 0.0
	for _, v := range x.vv {
		if v > 0 {
			ent -= v * math.Log2(v)
			sum += v
		}
	}
	if sum > 0 {
		return (ent + sum*math.Log2(sum)) / sum
	}
	return 0.0
}

// HeapSize estimates the required heap-size
func (x *DenseVector) HeapSize() int {
	return 24 + len(x.vv)*8
}

// converts the vector to a sparse one
func (x *DenseVector) convertToSparse() SparseVector {
	nv := make(SparseVector, x.Count())
	x.ForEach(func(i int, v float64) { nv.Set(i, v) })
	return nv
}

// --------------------------------------------------------------------

// VectorDistribution is a slice of sparse vectors
type VectorDistribution map[int]Vector

// NewVectorDistribution creates a new sparse matrix
func NewVectorDistribution() VectorDistribution {
	return make(VectorDistribution)
}

// NumRows returns the number of rows
func (m VectorDistribution) NumRows() int { return len(m) }

// NumCols returns the number of cols
func (m VectorDistribution) NumCols() int {
	var n int
	for _, row := range m {
		if l := row.Count(); l > n {
			n = l
		}
	}
	return n
}

// Weights returns the weight distribution
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
	m.initRow(row)
	m[row].Incr(col, delta)
}

// HeapSize estimates the required heap-size
func (m VectorDistribution) HeapSize() int {
	size := 24 + len(m)*16
	for _, vv := range m {
		size += vv.HeapSize()
	}
	return size
}

func (m VectorDistribution) initRow(pos int) {
	if _, ok := m[pos]; !ok {
		m[pos] = NewSparseVector()
	}
}

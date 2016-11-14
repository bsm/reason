package util

import "math"

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
func (vv SparseVector) Clone() SparseVector {
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
	if len(vv) == 0 {
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
	if len(vv) == 0 {
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
	n := len(vv)
	if n == 0 {
		return
	}

	µ = vv.Sum() / float64(n)
	return
}

// Variance calculates the variance
func (vv SparseVector) Variance() (V float64) {
	n := len(vv)
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
	n := len(vv)
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
	if len(vv) == 0 {
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

/*

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
}

*/

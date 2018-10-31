package util

import (
	"math"
)

// NewVector inits a new (sparse) vector.
func NewVector() *Vector {
	return new(Vector)
}

// NewVectorFromSlice initalizes a vector using a slice of weights.
func NewVectorFromSlice(weights ...float64) *Vector {
	return &Vector{Data: weights}
}

// Dims returns the number of rows and cols.
func (vv *Vector) Dims() (rows, cols int) {
	return 1, vv.NumCols()
}

// NumRows returns the number of rows.
func (vv *Vector) NumRows() int {
	return 1
}

// NumCols returns the number of cols.
func (vv *Vector) NumCols() int {
	return len(vv.Data)
}

// NNZ returns the number of non-zero weights in the vector.
func (vv *Vector) NNZ() int {
	n := 0
	for _, v := range vv.Data {
		if v != 0 {
			n++
		}
	}
	return n
}

// WeightSum returns the sum of all weights of the vector.
func (vv *Vector) WeightSum() float64 {
	sum := 0.0
	for _, v := range vv.Data {
		sum += v
	}
	return sum
}

// Min returns the position of the minimum weight in the vector (with value).
func (vv *Vector) Min() (pos int, weight float64) {
	pos = -1
	for i, w := range vv.Data {
		if w > 0 && (weight == 0 || w < weight) {
			pos = i
			weight = w
		}
	}
	return pos, weight
}

// Max returns the position od the maximum weight in the vector (with value).
func (vv *Vector) Max() (pos int, weight float64) {
	pos = -1
	for i, w := range vv.Data {
		if w > weight {
			pos = i
			weight = w
		}
	}
	return
}

// Clone creates a copy of the vector
func (vv *Vector) Clone() *Vector {
	nv := new(Vector)
	if vv.Data != nil {
		nv.Data = make([]float64, len(vv.Data))
		copy(nv.Data, vv.Data)
	}
	return nv
}

// At returns weight at index.
func (vv *Vector) At(index int) float64 {
	if index > -1 && index < len(vv.Data) {
		return vv.Data[index]
	}
	return 0.0
}

// Set sets a weight at index (auto-expanding).
// Returns true if the vector has expanded.
func (vv *Vector) Set(index int, weight float64) (expanded bool) {
	if index < 0 {
		return
	}

	if n := index + 1; n > cap(vv.Data) {
		data := make([]float64, n, 2*n)
		copy(data, vv.Data)
		vv.Data = data
		expanded = true
	} else if n > len(vv.Data) {
		vv.Data = vv.Data[:n]
	}
	vv.Data[index] = weight
	return
}

// Add increments a weight at index by delta.
// Returns true if the vector has expanded.
func (vv *Vector) Add(index int, delta float64) (expanded bool) {
	if index < 0 {
		return
	}

	if index < len(vv.Data) {
		vv.Data[index] += delta
	} else {
		expanded = vv.Set(index, delta)
	}
	return
}

// Mean returns the mean value of non-zero weights.
func (vv *Vector) Mean() float64 {
	_, _, mu := vv.csm()
	return mu
}

// Variance calculates the sample variance of non-zero weights.
func (vv *Vector) Variance() float64 {
	if n, _, mu := vv.csm(); n > 1 {
		var ss, cs float64
		for _, w := range vv.Data {
			if w > 0 {
				dt := w - mu
				ss += dt * dt
				cs += dt
			}
		}
		return (ss - cs*cs/float64(n)) / float64(n-1)
	}
	return math.NaN()
}

// Entropy calculates the entropy of non-zero weights.
func (vv *Vector) Entropy() float64 {
	ent := 0.0
	sum := 0.0
	for _, w := range vv.Data {
		if w > 0 {
			ent -= w * math.Log2(w)
			sum += w
		}
	}
	if sum > 0 {
		return (ent + sum*math.Log2(sum)) / sum
	}
	return 0.0
}

// StdDev calculates the standard deviation of non-zero weights.
func (vv *Vector) StdDev() float64 {
	return math.Sqrt(vv.Variance())
}

// Normalize normalizes all values to the 0..1 range.
func (vv *Vector) Normalize() {
	sum := vv.WeightSum()
	if sum == 0 {
		return
	}

	for i, w := range vv.Data {
		if w > 0 {
			vv.Set(i, w/sum)
		}
	}
}

func (vv *Vector) csm() (n int, sum, mu float64) {
	for _, w := range vv.Data {
		if w > 0 {
			n++
			sum += w
		}
	}
	if n > 0 {
		mu = sum / float64(n)
	}
	return
}

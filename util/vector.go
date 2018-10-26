package util

import (
	"math"
)

// NewVectorFromSlice initalizes a vector using a slice of weights.
func NewVectorFromSlice(weights ...float64) *Vector {
	return &Vector{Data: weights}
}

// NewVector inits a new (sparse) vector.
func NewVector() *Vector {
	return new(Vector)
}

// IsSparse returns true if the vector uses sparse representation.
func (vv *Vector) IsSparse() bool {
	return vv.SparseDict != nil && len(vv.SparseDict) == len(vv.Data)
}

// MakeSparse converts a vector to sparse storage.
func (vv *Vector) MakeSparse() {
	if vv.IsSparse() {
		return
	}

	nv := Vector{
		Data:       make([]float64, 0, len(vv.Data)*5),
		SparseDict: make([]uint32, 0, len(vv.Data)*5),
	}
	vv.ForEach(func(i int, w float64) bool {
		nv.Set(i, w)
		return true
	})
	*vv = nv
}

// MakeDense converts a vector to dense storage.
func (vv *Vector) MakeDense() {
	if !vv.IsSparse() {
		return
	}

	nv := Vector{
		Data: make([]float64, 0, len(vv.SparseDict)*5),
	}
	vv.ForEach(func(i int, w float64) bool {
		nv.Set(i, w)
		return true
	})
	*vv = nv
}

// Len returns the number of non-zero weights in the vector.
func (vv *Vector) Len() int {
	n := 0
	for _, v := range vv.Data {
		if v != 0 {
			n++
		}
	}
	return n
}

// Weight returns the total of all weights of the vector.
func (vv *Vector) Weight() float64 {
	sum := 0.0
	for _, v := range vv.Data {
		sum += v
	}
	return sum
}

// Min returns the position of the minimum weight in the vector (with value).
func (vv *Vector) Min() (pos int, weight float64) {
	pos = -1
	vv.ForEach(func(i int, w float64) bool {
		if i == 0 || w < weight {
			pos = i
			weight = w
		}
		return true
	})
	return pos, weight
}

// Max returns the position od the maximum weight in the vector (with value).
func (vv *Vector) Max() (pos int, weight float64) {
	pos = -1
	vv.ForEach(func(i int, w float64) bool {
		if w > weight {
			pos = i
			weight = w
		}
		return true
	})
	return
}

// Clone creates a copy of the vector
func (vv *Vector) Clone() *Vector {
	nv := new(Vector)
	if vv.Data != nil {
		nv.Data = make([]float64, len(vv.Data))
		copy(nv.Data, vv.Data)
	}
	if vv.SparseDict != nil {
		nv.SparseDict = make([]uint32, len(vv.SparseDict))
		copy(nv.SparseDict, vv.SparseDict)
	}
	return nv
}

// At returns weight at index.
func (vv *Vector) At(index int) float64 {
	if index < 0 {
		return 0.0
	}

	if vv.IsSparse() {
		index = vv.sparsePos(index)
	}
	if index < len(vv.Data) {
		return vv.Data[index]
	}
	return 0.0
}

// Set sets a weight at index.
func (vv *Vector) Set(index int, weight float64) {
	if index < 0 {
		return
	}

	if vv.IsSparse() {
		if pos := vv.sparsePos(index); pos < len(vv.Data) {
			vv.Data[pos] = weight
		} else {
			vv.Data = append(vv.Data, weight)
			vv.SparseDict = append(vv.SparseDict, uint32(index))
		}
		return
	}

	if n := index + 1; n > cap(vv.Data) {
		data := make([]float64, n, 2*n)
		copy(data, vv.Data)
		vv.Data = data
	} else if n > len(vv.Data) {
		vv.Data = vv.Data[:n]
	}
	vv.Data[index] = weight
}

// Add increments a weight at index by delta.
func (vv *Vector) Add(index int, delta float64) {
	if index < 0 {
		return
	}

	if vv.IsSparse() {
		if pos := vv.sparsePos(index); pos < len(vv.Data) {
			vv.Data[pos] += delta
		} else {
			vv.Set(index, delta)
		}
		return
	}

	if index < len(vv.Data) {
		vv.Data[index] += delta
	} else {
		vv.Set(index, delta)
	}
}

// ForEach iterates over each index/value
func (vv *Vector) ForEach(iter func(int, float64) bool) {
	if vv.IsSparse() {
		for i, u := range vv.SparseDict {
			if w := vv.Data[i]; w > 0 && !iter(int(u), w) {
				break
			}
		}
		return
	}

	for i, w := range vv.Data {
		if w > 0 && !iter(i, w) {
			break
		}
	}
}

// ForEachValue iterates over each value
func (vv *Vector) ForEachValue(iter func(float64) bool) {
	if vv.IsSparse() {
		for i := range vv.SparseDict {
			if w := vv.Data[i]; w > 0 && !iter(w) {
				break
			}
		}
		return
	}

	for _, w := range vv.Data {
		if w > 0 && !iter(w) {
			break
		}
	}
}

// Mean returns the mean value of non-zero weights.
func (vv *Vector) Mean() float64 {
	_, _, mu := vv.csm()
	return mu
}

// Variance calculates the variance of non-zero weights.
func (vv *Vector) Variance() float64 {
	if n, _, mu := vv.csm(); n > 1 {
		var ss, cs float64
		vv.ForEachValue(func(w float64) bool {
			dt := w - mu
			ss += dt * dt
			cs += dt
			return true
		})
		return (ss - cs*cs/float64(n)) / float64(n-1)
	}
	return math.NaN()
}

// Entropy calculates the entropy of non-zero weights.
func (vv *Vector) Entropy() float64 {
	ent := 0.0
	sum := 0.0
	vv.ForEachValue(func(w float64) bool {
		ent -= w * math.Log2(w)
		sum += w
		return true
	})
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
	sum := vv.Weight()
	if sum == 0 {
		return
	}

	vv.ForEach(func(i int, w float64) bool {
		vv.Set(i, w/sum)
		return true
	})
}

func (vv *Vector) csm() (n int, sum, mu float64) {
	vv.ForEachValue(func(w float64) bool {
		n++
		sum += w
		return true
	})
	if n > 0 {
		mu = sum / float64(n)
	}
	return
}

func (vv *Vector) sparsePos(index int) int {
	x := uint32(index)
	for i, u := range vv.SparseDict {
		if u == x {
			return i
		}
	}
	return len(vv.Data)
}

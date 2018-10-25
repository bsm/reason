package util

import (
	"math"

	"github.com/bsm/reason/internal/sparsedense"
)

// NewVectorFromSlice initalizes a vector using a slice of weights.
func NewVectorFromSlice(weights ...float64) *Vector {
	return &Vector{Data: weights}
}

// Len returns the number of non-zero weights in the vector.
func (vv *Vector) Len() int {
	if vv.Sparse != nil {
		return len(vv.Sparse)
	}

	n := 0
	for _, v := range vv.Data {
		if v != 0 {
			n++
		}
	}
	return n
}

// ToDense converts the vector to dense.
func (vv *Vector) ToDense() {
	if vv.Sparse == nil {
		return
	}

	data := make([]float64, 0, len(vv.Sparse))
	vv.ForEach(func(i int, w float64) bool {
		if n := i + 1; n > cap(data) {
			newdata := make([]float64, n, 2*n)
			copy(newdata, data)
			data = newdata
		} else if n > len(data) {
			data = data[:n]
		}
		data[i] = w
		return true
	})
	vv.Data = data
	vv.Sparse = nil
}

// ToSparse converts the vector to sparse.
func (vv *Vector) ToSparse() {
	if vv.Sparse != nil {
		return
	}

	data := make([]float64, 0, len(vv.Data)/2)
	sparse := make([]int64, 0, len(vv.Data)/2)
	vv.ForEach(func(i int, w float64) bool {
		sparse = append(sparse, int64(i))
		data = append(data, w)
		return true
	})
	vv.Data = data
	vv.Sparse = sparse
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
	if vv.Sparse != nil {
		nv.Sparse = make([]int64, len(vv.Sparse))
		copy(nv.Sparse, vv.Sparse)
	}
	if vv.Data != nil {
		nv.Data = make([]float64, len(vv.Data))
		copy(nv.Data, vv.Data)
	}
	return nv
}

// Get returns a value at index
func (vv *Vector) Get(index int) float64 {
	if index < 0 {
		return 0.0
	}

	if vv.Sparse != nil {
		index = vv.findSparse(index)
	}

	if index < len(vv.Data) {
		return vv.Data[index]
	}
	return 0.0
}

// Set sets a weight at index
func (vv *Vector) Set(index int, weight float64) {
	if index < 0 {
		return
	}

	if vv.Sparse != nil {
		if stored := vv.findSparse(index); stored < len(vv.Data) {
			vv.Data[stored] = weight
		} else {
			vv.Data = append(vv.Data, weight)
			vv.Sparse = append(vv.Sparse, int64(index))
		}
		vv.tryConvertToDense()
		return
	}

	if n := index + 1; n > cap(vv.Data) {
		dense := make([]float64, n, 2*n)
		copy(dense, vv.Data)
		vv.Data = dense
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

	if vv.Sparse != nil {
		if stored := vv.findSparse(index); stored < len(vv.Data) {
			vv.Data[stored] += delta
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
	if vv.Sparse != nil {
		size := len(vv.Data)
		for pos, i := range vv.Sparse {
			if pos >= size {
				break
			}
			if !iter(int(i), vv.Data[pos]) {
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
	if vv.Sparse != nil {
		size := len(vv.Data)
		for pos := range vv.Sparse {
			if pos >= size {
				break
			}
			if !iter(vv.Data[pos]) {
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

// Clear removes all weights from the vector.
func (vv *Vector) Clear() {
	for i := range vv.Data {
		vv.Data[i] = 0
	}
	vv.Data = vv.Data[:0]
	vv.Sparse = vv.Sparse[:0]
}

// Mean returns the mean value of the vector.
func (vv *Vector) Mean() float64 {
	_, _, mu := vv.csm()
	return mu
}

// Variance calculates the variance of the vector weights.
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

// Entropy calculates the entropy of the vector weights.
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

// StdDev calculates the standard deviation
func (vv *Vector) StdDev() float64 { return math.Sqrt(vv.Variance()) }

// Normalize normalizes all values to the 0..1 range
func (vv *Vector) Normalize() {
	sum := vv.Weight()
	if sum == 0 {
		return
	}

	for i, w := range vv.Data {
		vv.Data[i] = w / sum
	}
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

func (vv *Vector) tryConvertToDense() {
	capacity := 0
	for _, i := range vv.Sparse {
		if n := int(i) + 1; n > capacity {
			capacity = n
		}
	}

	if sparsedense.BetterOffDense(len(vv.Sparse), capacity) {
		vv.ToDense()
	}
}

func (vv *Vector) findSparse(index int) int {
	i64 := int64(index)
	for pos, i := range vv.Sparse {
		if i == i64 {
			return pos
		}
	}
	return len(vv.Sparse)
}

// --------------------------------------------------------------------

// Get returns the vector at index
func (x *VectorDistribution) Get(index int) *Vector {
	if index < 0 {
		return nil
	}

	if x.Dense != nil && index < len(x.Dense) {
		return x.Dense[index].Vector
	} else if x.Sparse != nil {
		return x.Sparse[int64(index)]
	}
	return nil
}

// Add increments a weight for index at value by delta.
func (x *VectorDistribution) Add(index, value int, delta float64) {
	if index < 0 {
		return
	}

	if x.Dense != nil {
		x.fetchDense(index).Add(value, delta)
		return
	}

	if x.Sparse == nil {
		x.Sparse = make(map[int64]*Vector)
	}
	x.fetchSparse(int64(index)).Add(value, delta)
}

// Len returns the number of elements in the distribution.
func (x *VectorDistribution) Len() int {
	if x.Dense != nil {
		n := 0
		x.ForEach(func(_ int, _ *Vector) bool { n++; return true })
		return n
	}

	return len(x.Sparse)
}

// ForEach iterates over each stats eleemnt in the distribution.
func (x *VectorDistribution) ForEach(iter func(index int, vv *Vector) bool) {
	if x.Dense != nil {
		for i, s := range x.Dense {
			if s.Vector != nil && !iter(i, s.Vector) {
				break
			}
		}
	} else if x.Sparse != nil {
		for i, s := range x.Sparse {
			if !iter(int(i), s) {
				break
			}
		}
	}
}

func (x *VectorDistribution) fetchDense(index int) *Vector {
	if n := index + 1; n > len(x.Dense) {
		dense := make([]VectorDistribution_Dense, n, n*2)
		copy(dense, x.Dense)
		x.Dense = dense
	} else if n > cap(x.Dense) {
		x.Dense = x.Dense[:n]
	}

	s := x.Dense[index].Vector
	if s == nil {
		s = new(Vector)
		s.ToSparse()
		x.Dense[index].Vector = s
	}
	return s
}

func (x *VectorDistribution) fetchSparse(index int64) *Vector {
	if n := index + 1; n > x.SparseCap {
		x.SparseCap = n
	}

	s, ok := x.Sparse[index]
	if !ok {
		if sparsedense.BetterOffDense(len(x.Sparse), int(x.SparseCap)) {
			x.convertToDense()
			return x.fetchDense(int(index))
		}

		s = new(Vector)
		s.ToSparse()
		x.Sparse[index] = s
	}
	return s
}

func (x *VectorDistribution) convertToDense() {
	dense := make([]VectorDistribution_Dense, int(x.SparseCap))
	x.ForEach(func(i int, vv *Vector) bool {
		dense[i].Vector = vv
		return true
	})
	x.Dense = dense
	x.Sparse = nil
	x.SparseCap = 0
}

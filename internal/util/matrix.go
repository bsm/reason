package util

import "github.com/bsm/reason/internal/calc"

// NumMatrix is a multi-dimensional float vector
type NumMatrix [][]float64

// NumRows returns the number of rows
func (m NumMatrix) NumRows() int { return len(m) }

// NumCols returns the number of cols
func (m NumMatrix) NumCols() int {
	var c int
	for _, row := range m {
		if l := len(row); l > c {
			c = l
		}
	}
	return c
}

// SumRows returns a vector of matrix row sums
func (m NumMatrix) SumRows() []float64 {
	sums := make([]float64, len(m))
	for i, row := range m {
		sums[i] = calc.Sum(row)
	}
	return sums
}

// SumCols returns a vector of matrix col sums
func (m NumMatrix) SumCols() []float64 {
	sums := make([]float64, m.NumCols())
	for _, row := range m {
		for i, v := range row {
			sums[i] += v
		}
	}
	return sums
}

// SumRowsPlusTotal returns a vector of matrix row sums and the overall matrix sum
func (m NumMatrix) SumRowsPlusTotal() ([]float64, float64) { return calc.MatrixRowSumsPlusTotal(m) }

// GetRow returns the row at index i. Returns nil if row doesn't exist
func (m NumMatrix) GetRow(i int) []float64 {
	if i > -1 && i < len(m) {
		return m[i]
	}
	return nil
}

// SetRow sets row at index i resizing the matrix if necessary.
// Returns the updated and resized matrix
func (m NumMatrix) SetRow(i int, vv []float64) NumMatrix {
	m = m.grow(i + 1)
	m[i] = vv
	return m
}

// Grow increases the number of rows and returns the resized matrix
func (m NumMatrix) grow(n int) NumMatrix {
	if d := n - len(m); d > 0 {
		nv := make(NumMatrix, n)
		copy(nv, m)
		m = nv
	}
	return m
}

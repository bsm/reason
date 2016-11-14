package util

// SparseMatrix is a slice of sparse vectors
type SparseMatrix map[int]SparseVector

// NewSparseMatrix creates a new sparse matrix
func NewSparseMatrix() SparseMatrix {
	return make(SparseMatrix)
}

// NumRows returns the number of rows
func (m SparseMatrix) NumRows() int { return len(m) }

// NumCols returns the number of cols
func (m SparseMatrix) NumCols() int {
	var n int
	for _, row := range m {
		if l := row.Count(); l > n {
			n = l
		}
	}
	return n
}

// Weights returns the weight distribution
func (m SparseMatrix) Weights() map[int]float64 {
	vv := make(map[int]float64, len(m))
	for i, row := range m {
		vv[i] = row.Sum()
	}
	return vv
}

// Get increments row at pos or nil
func (m SparseMatrix) Get(pos int) SparseVector {
	if pos > -1 {
		return m[pos]
	}
	return nil
}

// Incr increments value at row/col index
func (m SparseMatrix) Incr(row, col int, delta float64) {
	m.initRow(row)
	m[row].Incr(col, delta)
}

func (m SparseMatrix) initRow(pos int) {
	if _, ok := m[pos]; !ok {
		m[pos] = NewSparseVector()
	}
}

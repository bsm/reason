package util

// NewMatrix inits a new (dense) matrix.
func NewMatrix() *Matrix {
	return new(Matrix)
}

// Dims returns the number of rows and cols.
func (m *Matrix) Dims() (rows, cols int) {
	cols = m.NumCols()
	if cols != 0 {
		rows = len(m.Data) / cols
	}
	return
}

// NumRows returns the number of rows.
func (m *Matrix) NumRows() int {
	rows, _ := m.Dims()
	return rows
}

// NumCols returns the number of cols.
func (m *Matrix) NumCols() int {
	return int(m.Stride)
}

// At gets the weight of cell (i, j).
func (m *Matrix) At(i, j int) float64 {
	if i < 0 || j < 0 {
		return 0.0
	}

	rows, cols := m.Dims()
	if i < rows && j < cols {
		return m.Data[i*cols+j]
	}
	return 0.0
}

// Set sets the weight of cell (i, j) to w.
func (m *Matrix) Set(i, j int, w float64) {
	if i < 0 || j < 0 {
		return
	}

	m.expand(i, j)
	m.Data[i*int(m.Stride)+j] = w
}

// Incr increments the weight of cell (i, j) by delta.
func (m *Matrix) Incr(i, j int, delta float64) {
	if i < 0 || j < 0 {
		return
	}

	m.expand(i, j)
	m.Data[i*int(m.Stride)+j] += delta
}

// Row returns the weight at row i.
func (m *Matrix) Row(i int) []float64 {
	if i < 0 {
		return nil
	}

	stride := int(m.Stride)
	min := i * stride
	max := min + stride
	if max > len(m.Data) {
		return nil
	}

	return m.Data[min:max]
}

// IsRowZero returns true if all fields in row i are 0.
func (m *Matrix) IsRowZero(i int) bool {
	if i < 0 {
		return true
	}

	stride := int(m.Stride)
	min := i * stride
	max := min + stride
	if max > len(m.Data) {
		return true
	}

	for i := min; i < max; i++ {
		if m.Data[i] != 0 {
			return false
		}
	}
	return true
}

// RowSum returns the sum of all weights in row i.
func (m *Matrix) RowSum(i int) float64 {
	if i < 0 {
		return 0.0
	}

	min := i * int(m.Stride)
	max := min + int(m.Stride)
	if max > len(m.Data) {
		return 0.0
	}

	sum := 0.0
	for i := min; i < max; i++ {
		sum += m.Data[i]
	}
	return sum
}

// RowNNZ counts the number of non-zero weights in row i.
func (m *Matrix) RowNNZ(i int) int {
	if i < 0 {
		return 0
	}

	min := i * int(m.Stride)
	max := min + int(m.Stride)
	if max > len(m.Data) {
		return 0
	}

	n := 0
	for i := min; i < max; i++ {
		if m.Data[i] != 0 {
			n++
		}
	}
	return n
}

// ColSum returns the sum of all weights in col i.
func (m *Matrix) ColSum(i int) float64 {
	s := int(m.Stride)
	if i < 0 || i >= s {
		return 0.0
	}

	sum := 0.0
	for x := i; x < len(m.Data); x += s {
		sum += m.Data[x]
	}
	return sum
}

// ColNNZ counts the number of non-zero weights in col i.
func (m *Matrix) ColNNZ(i int) int {
	s := int(m.Stride)
	if i < 0 || i >= s {
		return 0
	}

	n := 0
	for x, s := i, int(m.Stride); x < len(m.Data); x += s {
		if m.Data[x] != 0 {
			n++
		}
	}
	return n
}

// WeightSum returns the sum of all weights.
func (m *Matrix) WeightSum() float64 {
	sum := 0.0
	for _, v := range m.Data {
		sum += v
	}
	return sum
}

func (m *Matrix) expand(i, j int) {
	rows, cols := m.Dims()
	if i < rows && j < cols {
		return
	}

	newrows := maxInt(rows, i+1)
	newcols := maxInt(cols, j+1)

	data := make([]float64, newrows*newcols)
	for row := 0; row < rows; row++ {
		copy(data[row*newcols:], m.Data[row*cols:row*cols+cols])
	}
	m.Stride = uint32(newcols)
	m.Data = data
}

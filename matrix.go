package etumag1matrix

type Matrix struct {
	// Number of rows
	rows int
	// Number of columns
	cols int
	// Matrix stored as a flat array: Aij = elements[i*cols + j]
	elements []float64
}

func NewEmptyMatrix(rows, cols int) *Matrix {
	return &Matrix{rows, cols, make([]float64, 0, rows*cols)}
}

func NewMatrixFromSlice(rows, cols int, elements []float64) *Matrix {
	if rows*cols <= len(elements) {
		elements = elements[:rows*cols]
	}

	els := make([]float64, 0, rows*cols)
	els = append(els, elements...)

	return &Matrix{rows, cols, els}
}

func (m *Matrix) CountRows() int {
	return m.rows
}

func (m *Matrix) CountCols() int {
	return m.cols
}

func (m *Matrix) GetElement(i int, j int) float64 {
	if offset := i*m.cols + j; len(m.elements) > offset {
		return m.elements[offset]
	}

	return 0
}

func (m *Matrix) GetRow(i int) []float64 {
	return m.elements[i*m.cols : i*m.cols+m.rows]
}

func (m *Matrix) GetCol(j int) []float64 {
	col := make([]float64, m.rows)
	for i := 0; i < m.rows; i++ {
		col[i] = m.elements[i*m.cols+j]
	}
	return col
}

func (m *Matrix) SetElement(i int, j int, v float64) {
	if offset := i*m.cols + j; cap(m.elements) > offset {
		if len(m.elements) <= offset {
			m.elements = append(m.elements, make([]float64, offset-len(m.elements)+1)...)
		}

		m.elements[offset] = v
	}
}

func (m *Matrix) CalcCoord(n int) (int, int) {
	return n % m.cols, (n - n%m.rows) / m.rows
}

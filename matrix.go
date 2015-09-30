package etumag1matrix

import "fmt"

type Matrix struct {
	rows int // Number of rows
	cols int // Number of columns
	// Matrix stored as a float array: Aij = elements[i*cols + j]
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
	if i >= m.rows {
		return nil
	}
	var row = make([]float64, m.cols)
	for j := 0; j < m.cols; j++ {
		row[j] = m.GetElement(i, j)
	}
	return row
}

func (m *Matrix) GetCol(j int) []float64 {
	if j >= m.cols {
		return nil
	}
	var col = make([]float64, m.rows)
	for i := 0; i < m.rows; i++ {
		col[i] = m.GetElement(i, j)
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

func (m *Matrix) CalcIJ(n int) (int, int) {
	return (n - n%m.cols) / m.rows, n % m.cols
}

func (m *Matrix) String() (output string) {
	for n := 0; n < m.cols*m.rows; n++ {
		if n%m.cols == 0 {
			output += "\n"
		}
		output += fmt.Sprintf("%10.2f", m.GetElement(m.CalcIJ(n)))
	}
	return
}

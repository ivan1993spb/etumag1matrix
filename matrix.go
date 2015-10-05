package etumag1matrix

import (
	// "encoding/xml"
	"fmt"
)

type Matrix struct {
	// XMLName  xml.Name  `xml:"matrix"`
	Rows     int       `xml:"rows>i4"`      // Number of rows
	Cols     int       `xml:"cols>i4"`      // Number of columns
	Elements []float64 `xml:"array>double"` // Matrix stored as a float array: Aij = elements[i*cols + j]
}

func NewEmptyMatrix(rows, cols int) *Matrix {
	return &Matrix{Rows: rows, Cols: cols, Elements: make([]float64, 0, rows*cols)}
}

func NewMatrixFromSlice(rows, cols int, elements []float64) *Matrix {
	if rows*cols <= len(elements) {
		elements = elements[:rows*cols]
	}

	els := make([]float64, 0, rows*cols)
	els = append(els, elements...)

	return &Matrix{Rows: rows, Cols: cols, Elements: els}
}

func (m *Matrix) CountRows() int {
	return m.Rows
}

func (m *Matrix) CountCols() int {
	return m.Cols
}

func (m *Matrix) GetElement(i int, j int) float64 {
	if offset := m.CalcIndex(i, j); len(m.Elements) > offset {
		return m.Elements[offset]
	}

	return 0
}

func (m *Matrix) GetRow(i int) []float64 {
	if i >= m.Rows {
		return nil
	}
	var row = make([]float64, m.Cols)
	for j := 0; j < m.Cols; j++ {
		row[j] = m.GetElement(i, j)
	}
	return row
}

func (m *Matrix) GetCol(j int) []float64 {
	if j >= m.Cols {
		return nil
	}
	var col = make([]float64, m.Rows)
	for i := 0; i < m.Rows; i++ {
		col[i] = m.GetElement(i, j)
	}
	return col
}

func (m *Matrix) SetElement(i int, j int, v float64) {
	if offset := m.CalcIndex(i, j); cap(m.Elements) > offset {
		if len(m.Elements) <= offset {
			m.Elements = append(m.Elements, make([]float64, offset-len(m.Elements)+1)...)
		}

		m.Elements[offset] = v
	}
}

func (m *Matrix) CalcIJ(n int) (int, int) {
	return (n - n%m.Cols) / m.Rows, n % m.Cols
}

func (m *Matrix) CalcIndex(i, j int) int {
	return i*m.Cols + j
}

func (m *Matrix) String() (output string) {
	for n := 0; n < m.Cols*m.Rows; n++ {
		if n%m.Cols == 0 {
			output += "\n"
		}
		output += fmt.Sprintf("%10.2f", m.GetElement(m.CalcIJ(n)))
	}
	return
}

func (m1 *Matrix) Equals(m2 *Matrix) bool {
	if m1.Rows != m2.Rows {
		return false
	}
	if m1.Cols != m2.Cols {
		return false
	}
	if m1.Elements == nil && m2.Elements == nil {
		return true
	}
	if m1.Elements == nil || m2.Elements == nil {
		return false
	}
	if len(m1.Elements) != len(m2.Elements) {
		return false
	}
	for i := range m1.Elements {
		if m1.Elements[i] != m2.Elements[i] {
			return false
		}
	}
	return true
}

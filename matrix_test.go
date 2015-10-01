package etumag1matrix

import "testing"

var emptyMatrix10_12 = NewEmptyMatrix(10, 12)

var matrixFromSlice5_6 = NewMatrixFromSlice(5, 6, []float64{
	1, 2, 3, 4, 5, 6,
	7, 8, 9, 10, 11, 12,
	13, 14, 15, 16, 17, 18,
	19, 20, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30,
})

var col_6_5_2 = []float64{3, 9, 15, 21, 27}

var row_6_5_3 = []float64{19, 20, 21, 22, 23, 24}

func TestNewEmptyMatrix(t *testing.T) {
	if emptyMatrix10_12.CountRows() != 10 {
		t.Error("Matrix.CountRows() fail")
	}
	if emptyMatrix10_12.CountCols() != 12 {
		t.Error("Matrix.CountCols() fail")
	}
	if cap(emptyMatrix10_12.elements) != 10*12 {
		t.Error("cap(matrix.elements)")
	}
	if len(emptyMatrix10_12.elements) != 0 {
		t.Error("len(matrix.elements)")
	}
}

func TestNewMatrixFromSlice(t *testing.T) {
	if matrixFromSlice5_6.CountRows() != 5 {
		t.Error("Matrix.CountRows() fail")
	}
	if matrixFromSlice5_6.CountCols() != 6 {
		t.Error("Matrix.CountCols() fail")
	}
	if cap(matrixFromSlice5_6.elements) != 6*5 {
		t.Error("cap(matrix.elements)")
	}
	if len(matrixFromSlice5_6.elements) != 6*5 {
		t.Error("len(matrix.elements)")
	}
}

func TestSetGetElement(t *testing.T) {
	emptyMatrix10_12.SetElement(10, 12, 22)
	if emptyMatrix10_12.GetElement(10, 12) != 0 {
		t.Error("Matrix.GetElement() fail")
	}
	if cap(emptyMatrix10_12.elements) != 10*12 {
		t.Error("cap(matrix.elements)")
	}
	if len(emptyMatrix10_12.elements) != 0 {
		t.Error("len(matrix.elements)")
	}

	emptyMatrix10_12.SetElement(3, 3, 22)
	if cap(emptyMatrix10_12.elements) != 10*12 {
		t.Error("cap(matrix.elements)")
	}
	if len(emptyMatrix10_12.elements) == 3*emptyMatrix10_12.cols+3 {
		t.Error("len(matrix.elements)")
	}
	if emptyMatrix10_12.elements[3*emptyMatrix10_12.cols+3] != 22 {
		t.Error("matrix.elements[i*cols+j] != val")
	}

	emptyMatrix10_12.SetElement(4, 4, 22)
	if cap(emptyMatrix10_12.elements) != 10*12 {
		t.Error("cap(matrix.elements)")
	}
	if len(emptyMatrix10_12.elements) == 4*emptyMatrix10_12.cols+4 {
		t.Error("len(matrix.elements)")
	}
	if emptyMatrix10_12.elements[4*emptyMatrix10_12.cols+4] != 22 {
		t.Error("matrix.elements[i*cols+j] != val")
	}
}

func TestGetColRow(t *testing.T) {
	if matrixFromSlice5_6.GetCol(22) != nil {
		t.Error("matrix.GetCol() error")
	}
	if matrixFromSlice5_6.GetRow(22) != nil {
		t.Error("matrix.GetRow() error")
	}
	if !compareFloat64Slices(col_6_5_2, matrixFromSlice5_6.GetCol(2)) {
		t.Error("matrix.GetCol() error")
	}
	if !compareFloat64Slices(row_6_5_3, matrixFromSlice5_6.GetRow(3)) {
		t.Error("matrix.GetRow() error")
	}

}

func compareFloat64Slices(a, b []float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCalcIJ(t *testing.T) {
	for n := 0; n < len(matrixFromSlice5_6.elements); n++ {
		if matrixFromSlice5_6.elements[n] != matrixFromSlice5_6.GetElement(matrixFromSlice5_6.CalcIJ(n)) {
			t.Error("matrix.CalcIJ(n)")
		}
	}
}

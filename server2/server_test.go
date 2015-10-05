package main

import (
	"testing"

	"github.com/ivan1993spb/etumag1matrix"
)

var matrixFromSlice1 = etumag1matrix.NewMatrixFromSlice(3, 3, []float64{
	1, 2, 3,
	4, 5, 6,
	7, 8, 9,
})

var matrixFromSlice2 = etumag1matrix.NewMatrixFromSlice(3, 3, []float64{
	7, 8, 9,
	4, 5, 6,
	1, 2, 3,
})

var matrixFromSlice3 = etumag1matrix.NewMatrixFromSlice(3, 4, []float64{
	1, 2, 3,
	4, 5, 6,
	7, 8, 9,
	10, 11, 12,
})

var matrixFromSlice4 = etumag1matrix.NewMatrixFromSlice(4, 3, []float64{
	1, 4, 7, 10,
	2, 5, 8, 11,
	3, 6, 9, 12,
})

var matrixFromSlice5 = etumag1matrix.NewMatrixFromSlice(3, 6, []float64{
	1, 2, 3, 4, 5, 6,
	0, 0, 0, 0, 0, 0,
	7, 8, 9, 10, 11, 12,
})

var matrixFromSlice6 = etumag1matrix.NewMatrixFromSlice(6, 2, []float64{
	1, 1,
	2, 2,
	3, 3,
	3, 3,
	2, 2,
	1, 1,
})

var res1 = etumag1matrix.NewMatrixFromSlice(3, 3, []float64{
	18, 24, 30,
	54, 69, 84,
	90, 114, 138,
})

var res2 = etumag1matrix.NewMatrixFromSlice(3, 3, []float64{
	69, 77, 74,
	169, 181, 182,
	269, 285, 290,
})

var res3 = etumag1matrix.NewMatrixFromSlice(3, 2, []float64{
	42, 42,
	0, 0,
	114, 114,
})

func TestMultiply(t *testing.T) {
	res_, err := multiply(matrixFromSlice1, matrixFromSlice2)
	if err != nil {
		t.Error("TestMultiply error:", err)
	}

	if !res1.Equals(res_) {
		t.Error("TestMultiply error")
	}

	res_, err = multiplyFast(matrixFromSlice3, matrixFromSlice4)
	if err != nil {
		t.Error("TestMultiply error:", err)
	}

	if !res2.Equals(res_) {
		t.Error("TestMultiply error")
	}

	res_, err = multiply(matrixFromSlice5, matrixFromSlice6)
	if err != nil {
		t.Error("TestMultiply error:", err)
	}

	if !res3.Equals(res_) {
		t.Error("TestMultiply error", err)
	}
}

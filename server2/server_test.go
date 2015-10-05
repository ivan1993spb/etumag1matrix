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

var res1 = etumag1matrix.NewMatrixFromSlice(3, 3, []float64{
	18, 24, 30,
	54, 69, 84,
	90, 114, 138,
})

func TestMultiply(t *testing.T) {
	res2, err := multiply(matrixFromSlice1, matrixFromSlice2)
	if err != nil {
		t.Error("TestMultiply error:", err)
	}

	if !res1.Equals(res2) {
		t.Error("TestMultiply error")
	}
}

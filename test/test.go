package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ivan1993spb/etumag1matrix"
)

const (
	COLS = 20
	ROWS = 20
)

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

var res3 = etumag1matrix.NewMatrixFromSlice(3, 2, []float64{
	42, 42,
	0, 0,
	114, 114,
})

func main() {
	if len(os.Args) < 2 {
		fmt.Println("enter http addresses")
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	elements := make([]float64, COLS*ROWS)
	for i := 0; i < COLS*ROWS; i++ {
		elements[i] = rand.Float64()
	}

	matrix := etumag1matrix.NewMatrixFromSlice(COLS, ROWS, elements)
	client, err := etumag1matrix.NewClient(os.Args[1:]...)

	fmt.Println("try to multiply matrix")

	start := time.Now()
	matrix2, err := client.MultiplyMatrix(matrix, matrix)
	fmt.Println("time:", time.Since(start))

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	if matrix2.CountRows()*matrix2.CountCols() <= 100 {
		fmt.Println(matrix2)
	}

	matrix, err = client.MultiplyMatrix(matrixFromSlice5, matrixFromSlice6)
	if !matrix.Equals(res3) {
		fmt.Println("ERROR")
	} else {
		fmt.Println("SUCCESS ^_^")
	}
}

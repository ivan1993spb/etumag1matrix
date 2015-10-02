package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ivan1993spb/etumag1matrix"
)

const (
	COLS = 100
	ROWS = 100
)

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
	fmt.Println("SUCCESS")
}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/ivan1993spb/etumag1matrix"
)

const (
	COLS = 50
	ROWS = 50
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

	_, err = client.MultiplyMatrix(matrix, matrix)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(0)
	}

	fmt.Println("SUCCESS")
}

package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ivan1993spb/etumag1matrix"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("starting server")

	var addr string
	flag.StringVar(&addr, "http", ":8080", "server address")
	flag.Parse()
	if addr == "" {
		log.Fatalln("empty server address")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("accepted connection")

		var pr *packreq
		if err := xml.NewDecoder(r.Body).Decode(&pr); err == nil {
			t := time.Now()
			request := pr.MultMatrix
			result, err := multiplyFast(request.A, request.B)
			log.Printf("[%dx%d] X [%dx%d] %s\n", request.A.Cols, request.A.Rows,
				request.B.Cols, request.B.Rows, time.Since(t))
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
			} else {
				fmt.Fprintln(w, xml.Header)
				fmt.Fprintln(w, `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body>`)
				xml.NewEncoder(w).Encode(&multiplyResult{result})
				fmt.Fprintln(w, `</soap:Body></soap:Envelope>`)
			}
		} else {
			log.Println(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	})

	log.Fatalln(http.ListenAndServe(addr, nil))
}

type packreq struct {
	XMLName    xml.Name        `xml:"Envelope"`
	MultMatrix *multiplyMatrix `xml:"Body>multiplyMatrix"`
}

type multiplyMatrix struct {
	A, B *etumag1matrix.Matrix
}

type multiplyResult struct {
	Result *etumag1matrix.Matrix
}

func multiply(A, B *etumag1matrix.Matrix) (*etumag1matrix.Matrix, error) {
	if A.CountCols() != B.CountRows() {
		return nil, errors.New("number of columns in A is not equal to the number of rows in B")
	}

	cin := make(chan *task)
	chans := make([]<-chan *result, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		chans[i] = runprocess(cin)
	}
	cout := merge(chans...)

	go func() {
		for i := 0; i < A.CountRows(); i++ {
			for j := 0; j < B.CountCols(); j++ {
				cin <- &task{i*B.CountCols() + j, A.GetRow(i), B.GetCol(j)}
			}
		}
		close(cin)
	}()

	elements := make([]float64, A.CountRows()*B.CountCols())
	for r := range cout {
		if r.index < len(elements) {
			elements[r.index] = r.value
		}
	}

	return etumag1matrix.NewMatrixFromSlice(A.CountRows(), B.CountCols(), elements), nil
}

func runprocess(cin <-chan *task) <-chan *result {
	cout := make(chan *result)

	go func() {
		for t := range cin {
			if len(t.col) != len(t.row) {
				continue
			}
			r := &result{t.index, 0}
			for i := 0; i < len(t.col); i++ {
				r.value += t.col[i] * t.row[i]
			}
			cout <- r
		}
		close(cout)
	}()

	return cout
}

func merge(cs ...<-chan *result) <-chan *result {
	if len(cs) == 1 {
		return cs[0]
	}

	var wg sync.WaitGroup
	cout := make(chan *result)

	output := func(c <-chan *result) {
		for r := range c {
			cout <- r
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(cout)
	}()

	return cout
}

type task struct {
	index int
	row   []float64
	col   []float64
}

type result struct {
	index int
	value float64
}

func multiplyFast(A, B *etumag1matrix.Matrix) (*etumag1matrix.Matrix, error) {
	if A.CountCols() != B.CountRows() {
		return nil, errors.New("number of columns in A is not equal to the number of rows in B")
	}

	elements := make([]float64, A.Rows*B.Cols)
	stopch := make(chan struct{})

	proc := func(n int) {
		for m := n; m < len(elements); m += runtime.NumCPU() {
			i, j := (m-m%B.Cols)/B.Cols, m%B.Cols
			row := A.GetRow(i)
			col := B.GetCol(j)
			if len(row) == len(col) {
				for k := 0; k < len(row); k++ {
					elements[m] += row[k] * col[k]
				}
			}
		}
		stopch <- struct{}{}
	}

	for n := 0; n < runtime.NumCPU(); n++ {
		go proc(n)
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		<-stopch
	}
	close(stopch)

	return etumag1matrix.NewMatrixFromSlice(A.CountRows(), B.CountCols(), elements), nil
}

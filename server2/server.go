package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/ivan1993spb/etumag1matrix"
)

var procs int
var goroutines int

func main() {

	var addr string
	flag.StringVar(&addr, "http", ":8080", "server address")
	flag.IntVar(&procs, "procs", runtime.NumCPU(), "proc count")
	flag.IntVar(&goroutines, "goroutines", runtime.NumCPU(), "proc count")
	flag.Parse()

	runtime.GOMAXPROCS(procs)

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

	log.Println("starting server:", addr)
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

func multiplyFast(A, B *etumag1matrix.Matrix) (*etumag1matrix.Matrix, error) {
	if A.CountCols() != B.CountRows() {
		return nil, errors.New("number of columns in A is not equal to the number of rows in B")
	}

	elements := make([]float64, A.Rows*B.Cols)
	stopch := make(chan struct{})

	proc := func(n int) {
		for m := n; m < len(elements); m += goroutines {
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

	for n := 0; n < goroutines; n++ {
		go proc(n)
	}

	for i := 0; i < goroutines; i++ {
		<-stopch
	}
	close(stopch)

	return etumag1matrix.NewMatrixFromSlice(A.CountRows(), B.CountCols(), elements), nil
}

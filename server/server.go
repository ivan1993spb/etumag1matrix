package main

import (
	"encoding/xml"
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr string
	flag.StringVar(&addr, "http", ":8080", "server address")
	flag.Parse()
	if addr == "" {
		log.Fatalln("empty server address")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var requests []*reqfield
		if err := xml.NewDecoder(r.Body).Decode(&requests); err == nil {
			responses := make([]*respfield, 0)

			for _, req := range requests {
				responses = append(responses, handle(req))
			}

			xml.NewEncoder(w).Encode(responses)

		} else {
			log.Println(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	})

	log.Fatalln(http.ListenAndServe(addr, nil))
}

func handle(req *reqfield) *respfield {
	resp := &respfield{Index: req.Index}
	if len(req.Col)*len(req.Row) == 0 || len(req.Col) != len(req.Row) {
		resp.Status = 1
		return resp
	}

	for i := 0; i < len(req.Col); i++ {
		resp.Value += req.Col[i] * req.Row[i]
	}

	return resp
}

type reqfield struct {
	Index int       `xml:"index,attr"`    // Field index
	Col   []float64 `xml:"src>col>value"` // Column of matrix A
	Row   []float64 `xml:"src>row>value"` // Row of matrix B
}

type respfield struct {
	Index  int     `xml:"index,attr"`   // Field index
	Value  float64 `xml:"result>value"` // Result
	Status int     `xml:"status"`       // Response status
}

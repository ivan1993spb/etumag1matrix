package main

import (
	"encoding/xml"
	"flag"
	"log"
	"net/http"
)

const STATUS_INVALID_INPUT = 1

func main() {
	log.Println("starting server")

	var addr string
	flag.StringVar(&addr, "http", ":8080", "server address")
	flag.Parse()
	if addr == "" {
		log.Fatalln("empty server address")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("accepted connection")

		var reqfields *packreqfield
		if err := xml.NewDecoder(r.Body).Decode(&reqfields); err == nil {
			respfields := &packrespfield{Respfields: make([]*respfield, 0)}
			log.Println("fields:", len(reqfields.Reqfields))

			for _, req := range reqfields.Reqfields {
				respfields.Respfields = append(respfields.Respfields, handle(req))
			}

			xml.NewEncoder(w).Encode(respfields)
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
		resp.Status = STATUS_INVALID_INPUT
		return resp
	}

	for i := 0; i < len(req.Col); i++ {
		resp.Value += req.Col[i] * req.Row[i]
	}

	return resp
}

type packreqfield struct {
	XMLName   xml.Name    `xml:"reqfields"`
	Reqfields []*reqfield `xml:"reqfield"`
}

type reqfield struct {
	Index int       `xml:"index,attr"`    // Field index
	Col   []float64 `xml:"src>col>value"` // Column of matrix A
	Row   []float64 `xml:"src>row>value"` // Row of matrix B
}

type packrespfield struct {
	XMLName    xml.Name     `xml:"respfields"`
	Respfields []*respfield `xml:"respfield"`
}

type respfield struct {
	Index  int     `xml:"index,attr"`   // Field index
	Value  float64 `xml:"result>value"` // Result
	Status int     `xml:"status"`       // Response status
}

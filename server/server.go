package main

import (
	"encoding/xml"
	"flag"
	"log"
	"net/http"
)

type errServer struct {
	str string
}

func (e *errServer) Error() string {
	return "server error: " + e.str
}

func main() {
	var addr string
	flag.StringVar(&addr, "http", ":8080", "server address")
	flag.Parse()
	if addr == "" {
		log.Fatalln(&errServer{"empty server address"})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var requests []*request
		if err := xml.NewDecoder(r.Body).Decode(&requests); err == nil {
			responses := make([]*response, 0)

			for _, req := range requests {
				responses = append(responses, handle(req))
			}

			xml.NewEncoder(w).Encode(responses)

		} else {
			log.Println(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	})

	log.Fatal(http.ListenAndServe(addr, nil))

}

func handle(req *request) (resp *response) {
	resp = &response{Id: req.Id}
	if len(req.Col)*len(req.Row) == 0 || len(req.Col) != len(req.Row) {
		resp.Status = 1
		return
	}

	for i := 0; i < len(req.Col); i++ {
		resp.Value += req.Col[i] * req.Row[i]
	}

	return
}

type request struct {
	Id  string    `xml:"id"`        // Request UUID
	Col []float64 `xml:"col>value"` // Column of matrix A
	Row []float64 `xml:"row>value"` // Row of matrix B
}

type response struct {
	Id     string  `xml:"id"`
	Value  float64 `xml:"value"`
	Status int     `xml:"status"`
}

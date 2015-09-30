package etumag1matrix

import (
	"bytes"
	"encoding/xml"
	"errors"
	"math"
	"net"
	"net/http"
	"sync"
)

type ErrMultiply struct {
	str string
}

func (err *ErrMultiply) Error() string {
	return "multiplying error: " + err.str
}

var ErrMultiplyColsRowsNumber = &ErrMultiply{"number of columns in A is not equal to the number of rows in B"}

type Client struct {
	servers []*server
	cursor  int
}

func NewClient(servers ...string) (*Client, error) {
	client := &Client{make([]*server, len(servers)), 0}

	for i, server := range servers {
		addr, err := net.ResolveTCPAddr("tcp", server)
		if err != nil {
			return nil, err
		}
		client.servers[i].addr = addr
	}

	return client, nil
}

func (c *Client) MultiplyMatrix(A, B *Matrix) (*Matrix, error) {
	if A.CountCols() != B.CountRows() {
		return nil, ErrMultiplyColsRowsNumber
	}

	var (
		cin   = make(chan *reqfield)
		count = A.CountRows() * B.CountCols()
		cout  = c.initSession(cin, count)
	)

	for i := 0; i < A.CountRows(); i++ {
		for j := 0; j < B.CountCols(); j++ {
			cin <- &reqfield{i*B.CountCols() + j, A.GetRow(i), B.GetCol(j)}
		}
	}

	close(cin)

	var (
		err      error
		elements = make([]float64, count)
	)

	for field := range cout {
		if field.Err() == nil {
			if field.Index < count {
				elements[field.Index] = field.Value
			}
		} else {
			err = field.Err()
		}
	}

	return NewMatrixFromSlice(A.CountRows(), B.CountCols(), elements), err
}

func (c *Client) initSession(cin <-chan *reqfield, count int) <-chan *respfield {
	serverList := c.getServers(c.calculateServerCount(count))

	var wg sync.WaitGroup
	сout := make(chan *respfield)

	output := func(ch <-chan *respfield) {
		for resp := range ch {
			сout <- resp
		}
		wg.Done()
	}
	wg.Add(len(serverList))

	for _, server := range serverList {
		scout := server.startSession(cin)
		go output(scout)
	}

	go func() {
		wg.Wait()
		close(сout)
	}()

	return сout
}

func (c *Client) getServers(count int) []*server {
	if count < 1 {
		return nil
	}
	if count >= len(c.servers) {
		return c.servers
	}
	var ncurs int
	defer func() {
		c.cursor = ncurs - 1
	}()
	if c.cursor+count < len(c.servers) {
		ncurs = c.cursor + count
		return c.servers[c.cursor:ncurs]
	}
	ncurs = c.cursor + count - len(c.servers)
	return append(c.servers[c.cursor:], c.servers[:ncurs]...)
}

func (c *Client) calculateServerCount(calcCount int) int {
	if calcCount < 1 {
		return 0
	}
	count := int(math.Floor(math.Log10(float64(calcCount))))
	if count > len(c.servers) {
		return len(c.servers)
	}
	return count
}

func (c *Client) MultiplyMatrixCallback(A, B *Matrix, callback func(res *Matrix, err error)) {
	go func() {
		callback(c.MultiplyMatrix(A, B))
	}()
}

type server struct {
	addr net.Addr
}

func (s *server) startSession(cin <-chan *reqfield) <-chan *respfield {
	cout := make(chan *respfield)
	go func() {
		var buff bytes.Buffer

		buff.WriteString(xml.Header)

		enc := xml.NewEncoder(&buff)
		for reqf := range cin {
			enc.Encode(reqf)
		}

		httpreq, _ := http.NewRequest("GET", s.addr.String(), &buff)
		client := &http.Client{}
		httpresp, _ := client.Do(httpreq)

		var respfields []*respfield
		xml.NewDecoder(httpresp.Body).Decode(&respfields)

		for _, respf := range respfields {
			cout <- respf
		}

		close(cout)
	}()

	return cout
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

func (r *respfield) Err() error {
	switch r.Status {
	case 0:
		return nil
	case 1:
		return errors.New("invalid number of source values")
	}
	return errors.New("unknown status")
}

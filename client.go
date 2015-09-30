package etumag1matrix

import (
	"bytes"
	"encoding/xml"
	"errors"
	"math"
	"net"
	"net/http"
	"sync"

	"github.com/satori/go.uuid"
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

	count := A.CountRows() * B.CountCols()
	cin := make(chan *request)
	cout := c.initSession(cin, count)
	positions := map[string]int{}

	for i := 0; i < A.CountRows(); i++ {
		for j := 0; j < B.CountCols(); j++ {
			rid := uuid.NewV4().String()
			positions[rid] = i*B.CountCols() + j
			cin <- &request{rid, A.GetRow(i), B.GetCol(j)}
		}
	}

	close(cin)

	res := NewEmptyMatrix(A.CountRows(), B.CountCols())
	var err error

	for resp := range cout {
		if resp.Err() == nil {
			i, j := res.CalcCoord(positions[resp.Id])
			res.SetElement(i, j, resp.Value)
		} else {
			err = resp.Err()
		}
	}

	return res, err
}

func (c *Client) initSession(cin <-chan *request, count int) <-chan *response {
	serverList := c.getServers(c.calculateServerCount(count))

	var wg sync.WaitGroup
	сout := make(chan *response)

	output := func(c <-chan *response) {
		for n := range c {
			сout <- n
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

func (s *server) startSession(cin <-chan *request) <-chan *response {
	cout := make(chan *response)
	go func() {
		var buff bytes.Buffer

		buff.WriteString(xml.Header)

		enc := xml.NewEncoder(&buff)
		for req := range cin {
			enc.Encode(req)
		}

		hreq, _ := http.NewRequest("GET", s.addr.String(), &buff)
		client := &http.Client{}
		hresp, _ := client.Do(hreq)

		var responses []*response
		xml.NewDecoder(hresp.Body).Decode(&responses)

		for _, resp := range responses {
			cout <- resp
		}

		close(cout)
	}()

	return cout
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

func (r *response) Err() error {
	switch r.Status {
	case 1:
		return errors.New("invalid number of values")
	default:
		return errors.New("unknown status")
	}

	return nil
}

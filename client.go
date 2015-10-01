package etumag1matrix

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
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

type Client struct {
	addrs  []net.Addr
	cursor int
}

func NewClient(addrs ...string) (*Client, error) {
	client := &Client{make([]net.Addr, len(addrs)), 0}

	for i := range addrs {
		addr, err := net.ResolveTCPAddr("tcp", addrs[i])
		if err != nil {
			return nil, fmt.Errorf("cannot create client: %s", err)
		}
		client.addrs[i] = addr
	}

	return client, nil
}

func (c *Client) MultiplyMatrix(A, B *Matrix) (*Matrix, error) {
	if A.CountCols() != B.CountRows() {
		return nil, &ErrMultiply{"number of columns in A is not equal to the number of rows in B"}
	}

	var (
		cin   = make(chan *reqfield)
		count = A.CountRows() * B.CountCols()
	)

	cout := c.initSession(cin, count)
	for i := 0; i < A.CountRows(); i++ {
		for j := 0; j < B.CountCols(); j++ {
			cin <- &reqfield{i*B.CountCols() + j, A.GetRow(i), B.GetCol(j)}
		}
	}

	close(cin)

	var (
		err      error // First happened error
		elements = make([]float64, count)
	)

	for field := range cout {
		if field.Err() == nil {
			if field.Index < count {
				elements[field.Index] = field.Value
			} else if err == nil {
				err = &ErrMultiply{"server return invalid field index"}
			}
		} else if err == nil {
			err = &ErrMultiply{field.Err().Error()}
		}
	}

	if len(elements) == 0 && err == nil {
		return nil, &ErrMultiply{"data was not received"}
	}

	return NewMatrixFromSlice(A.CountRows(), B.CountCols(), elements), err
}

func (c *Client) initSession(cin <-chan *reqfield, count int) <-chan *respfield {
	var (
		scount = c.calculateServerCount(count)
		wg     sync.WaitGroup
		сout   = make(chan *respfield)
	)
	if scount == 0 {
		return nil
	}

	wg.Add(scount)
	addrs := c.getAddrs(scount)

	sendReceive := func(addr net.Addr) {
		defer wg.Done()
		var buff bytes.Buffer

		buff.WriteString(xml.Header)
		buff.WriteString("<reqfields>")
		enc := xml.NewEncoder(&buff)

		for reqfld := range cin {
			enc.Encode(reqfld)
		}

		buff.WriteString("</reqfields>")

		httpreq, err := http.NewRequest("GET", "http://"+addr.String(), &buff)
		if err != nil {
			return
		}
		httpresp, err := http.DefaultClient.Do(httpreq)
		if err != nil {
			return
		}

		var respfields *packrespfield
		err = xml.NewDecoder(httpresp.Body).Decode(&respfields)

		if err == nil {
			for _, respfld := range respfields.Respfields {
				сout <- respfld
			}
		}
	}

	for _, addr := range addrs {
		go sendReceive(addr)
	}

	go func() {
		wg.Wait()
		close(сout)
	}()

	return сout
}

func (c *Client) getAddrs(count int) []net.Addr {
	if count < 1 {
		return nil
	}
	if count >= len(c.addrs) {
		return c.addrs
	}

	var cursor = c.cursor

	if cursor+count <= len(c.addrs) {
		c.cursor = cursor + count
		if c.cursor == len(c.addrs) {
			c.cursor = 0
		}
		return c.addrs[cursor : cursor+count]
	}

	c.cursor = cursor + count - len(c.addrs)
	return append(c.addrs[cursor:], c.addrs[:c.cursor]...)
}

func (c *Client) calculateServerCount(calcCount int) int {
	if calcCount < 1 {
		return 0
	}

	count := int(math.Ceil(math.Log10(float64(calcCount))))

	if count > len(c.addrs) {
		return len(c.addrs)
	}
	return count
}

func (c *Client) MultiplyMatrixCallback(A, B *Matrix, callback func(res *Matrix, err error)) {
	go func() {
		callback(c.MultiplyMatrix(A, B))
	}()
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

const (
	STATUS_OK            = 0
	STATUS_INVALID_INPUT = 1
)

func (r *respfield) Err() error {
	switch r.Status {
	case STATUS_OK:
		return nil
	case STATUS_INVALID_INPUT:
		return errors.New("invalid input")
	}
	return errors.New("unknown status")
}

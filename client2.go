package etumag1matrix

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
)

type Client2 struct {
	addr net.Addr
}

func NewClient2(address string) (*Client2, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("cannot create client2: %s", err)
	}

	return &Client2{addr}, nil
}

func (c *Client2) MultiplyMatrix(A, B *Matrix) (*Matrix, error) {
	var buff bytes.Buffer

	buff.WriteString(xml.Header)

	if err := xml.NewEncoder(&buff).Encode(&multiplyMatrix{A, B}); err != nil {
		return nil, err
	}

	httpreq, err := http.NewRequest("GET", "http://"+c.addr.String(), &buff)
	if err != nil {
		return nil, err
	}
	httpresp, err := http.DefaultClient.Do(httpreq)
	if err != nil {
		return nil, err
	}

	var res *multiplyResult
	if err = xml.NewDecoder(httpresp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client2) MultiplyMatrixCallback(A, B *Matrix, callback func(res *Matrix, err error)) {
	go func() {
		callback(c.MultiplyMatrix(A, B))
	}()
}

type multiplyMatrix struct {
	A, B *Matrix
}

type multiplyResult struct {
	Result *Matrix
}

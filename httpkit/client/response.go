package client

import (
	"io"
	"net/http"

	"github.com/go-swagger/go-swagger/client"
)

var _ client.Response = response{}

type response struct {
	resp *http.Response
}

func (r response) Code() int {
	return r.resp.StatusCode
}

func (r response) Message() string {
	return r.resp.Status
}

func (r response) GetHeader(name string) string {
	return r.resp.Header.Get(name)
}

func (r response) Body() io.ReadCloser {
	return r.resp.Body
}

package client

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/casualjim/go-swagger/client"
)

// NewRequest creates a new swagger http client request
func NewRequest(method, pathPattern string, writer client.RequestWriter) *Request {
	return &Request{
		pathPattern: pathPattern,
		method:      method,
		writer:      writer,
	}
}

// Request represents a swagger client request.
// This is backed by a http request and implements the client.Request interface
// so it can be passed into a RequestWriter
type Request struct {
	pathPattern string
	method      string
	writer      client.RequestWriter

	queryParams url.Values
	formParams  url.Values
	headers     http.Header
	isMultipart bool
	body        io.Writer
	// HTTPRequest is the underlying http request that is being built
	HTTPRequest *http.Request
}

func (r *Request) AddHeaderParam(string, ...string) {}

func (r *Request) AddQueryParam(string, ...string) {}

func (r *Request) AddFormParam(string, ...string) {}

func (r *Request) AddPathParam(string, string) {}

func (r *Request) AddFileParam(string, *os.File) {}

func (r *Request) SetBodyParam(io.Writer) {}

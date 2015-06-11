package client

import (
	"io"
	"os"

	"github.com/casualjim/go-swagger/strfmt"
)

// RequestWriter is an interface for things that know how to write to a request
type RequestWriter interface {
	WriteToRequest(Request, strfmt.Registry) error
}

// Request is an interface for things that know how to
// add information to a swagger client request
type Request interface {
	AddHeaderParam(string, ...string)

	AddQueryParam(string, ...string)

	AddFormParam(string, ...string)

	AddPathParam(string, string)

	AddFileParam(string, *os.File)

	SetBodyParam(io.Writer)
}

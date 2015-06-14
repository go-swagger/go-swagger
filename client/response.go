package client

import (
	"io"

	"github.com/go-swagger/go-swagger/httpkit"
)

// A Response represents a client response
// This bridges between responses obtained from different transports
type Response interface {
	Code() int
	Message() string
	GetHeader(string) string
	Body() io.ReadCloser
}

// A ResponseReaderFunc turns a function into a ResponseReader interface implementation
type ResponseReaderFunc func(Response, httpkit.Consumer) (interface{}, error)

// ReadResponse reads the response
func (read ResponseReaderFunc) ReadResponse(resp Response, consumer httpkit.Consumer) (interface{}, error) {
	return read(resp, consumer)
}

// A ResponseReader is an interface for things want to read a response.
// An application of this is to create structs from response values
type ResponseReader interface {
	ReadResponse(Response, httpkit.Consumer) (interface{}, error)
}

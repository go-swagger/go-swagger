package client

import "github.com/go-swagger/go-swagger/strfmt"

// RequestWriterFunc converts a function to a request writer interface
type RequestWriterFunc func(Request, strfmt.Registry) error

// WriteToRequest adds data to the request
func (fn RequestWriterFunc) WriteToRequest(req Request, reg strfmt.Registry) error {
	return fn(req, reg)
}

// RequestWriter is an interface for things that know how to write to a request
type RequestWriter interface {
	WriteToRequest(Request, strfmt.Registry) error
}

// Request is an interface for things that know how to
// add information to a swagger client request
type Request interface {
	SetHeaderParam(string, ...string) error

	SetQueryParam(string, ...string) error

	SetFormParam(string, ...string) error

	SetPathParam(string, string) error

	SetFileParam(string, string) error

	SetBodyParam(interface{}) error
}

package client

import "github.com/go-swagger/go-swagger/strfmt"

// Operation represents the context for a swagger operation to be submitted to the transport
type Operation struct {
	ID       string
	AuthInfo AuthInfoWriter
	Params   RequestWriter
	Reader   ResponseReader
}

// A Transport implementor knows how to submit Request objects to some destination
type Transport interface {
	//Submit(string, RequestWriter, ResponseReader, AuthInfoWriter) (interface{}, error)
	Submit(*Operation) (interface{}, error)
}

// AuthInfoWriterFunc converts a function to a request writer interface
type AuthInfoWriterFunc func(Request, strfmt.Registry) error

// AuthenticateRequest adds authentication data to the request
func (fn AuthInfoWriterFunc) AuthenticateRequest(req Request, reg strfmt.Registry) error {
	return fn(req, reg)
}

// An AuthInfoWriter implementor knows how to write authentication info to a request
type AuthInfoWriter interface {
	AuthenticateRequest(Request, strfmt.Registry) error
}

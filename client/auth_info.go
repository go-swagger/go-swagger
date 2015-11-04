package client

import "github.com/go-swagger/go-swagger/strfmt"

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

// Package client contains a client to send http requests
// to a swagger API. This implementation is untyped
package client

import "fmt"

type methodAndPath struct {
	Method      string
	PathPattern string
}

// APIError wraps an error model and captures the status code
type APIError struct {
	OperationName string
	Payload       []byte
	Code          int
}

func (a *APIError) Error() string {
	return fmt.Sprintf("%s (status %d): %+v ", a.OperationName, a.Code, string(a.Payload))
}

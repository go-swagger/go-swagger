// Package client contains a client to send http requests
// to a swagger API. This implementation is untyped
package client

import (
	"encoding/base64"
	"fmt"

	"github.com/go-swagger/go-swagger/client"
	"github.com/go-swagger/go-swagger/strfmt"
)

type methodAndPath struct {
	Method      string
	PathPattern string
	Schemes     []string
}

// NewAPIError creates a new API error
func NewAPIError(opName string, payload []byte, code int) *APIError {
	return &APIError{
		OperationName: opName,
		Payload:       payload,
		Code:          code,
	}
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

// BasicAuth provides a basic auth info writer
func BasicAuth(username, password string) client.AuthInfoWriter {
	return client.AuthInfoWriterFunc(func(r client.Request, fmts strfmt.Registry) error {
		encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		r.SetHeaderParam("Authorization", "Basic "+encoded)
		return nil
	})
}

// APIKeyAuth provides an API key auth info writer
func APIKeyAuth(name, in, value string) client.AuthInfoWriter {
	if in == "query" {
		return client.AuthInfoWriterFunc(func(r client.Request, _ strfmt.Registry) error {
			r.SetQueryParam(name, value)
			return nil
		})
	} else if in == "header" {
		return client.AuthInfoWriterFunc(func(r client.Request, _ strfmt.Registry) error {
			r.SetHeaderParam(name, value)
			return nil
		})
	}
	return nil
}

// BearerToken provides a header based oauth2 bearer access token auth info writer
func BearerToken(token string) client.AuthInfoWriter {
	return client.AuthInfoWriterFunc(func(r client.Request, _ strfmt.Registry) error {
		r.SetHeaderParam("Authorization", "Bearer "+token)
		return nil
	})
}

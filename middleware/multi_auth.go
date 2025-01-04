package middleware

import (
	"github.com/go-openapi/runtime"
)

// MultiAuthenticator combines multiple authentication methods
type MultiAuthenticator struct {
	authenticators []runtime.Authenticator
}

// NewMultiAuthenticator creates a new multi-authenticator
func NewMultiAuthenticator(authenticators ...runtime.Authenticator) *MultiAuthenticator {
	return &MultiAuthenticator{authenticators: authenticators}
}

// Authenticate implements the runtime.Authenticator interface
func (m *MultiAuthenticator) Authenticate(params interface{}) (bool, interface{}, error) {
	for _, auth := range m.authenticators {
		if ok, principal, err := auth.Authenticate(params); ok {
			return true, principal, err
		}
	}
	return false, nil, nil
}
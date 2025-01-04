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
	if len(authenticators) == 0 {
		panic("at least one authenticator is required")
	}
	return &MultiAuthenticator{authenticators: authenticators}
}

// Authenticate implements the runtime.Authenticator interface
func (m *MultiAuthenticator) Authenticate(params interface{}) (bool, interface{}, error) {
	for _, auth := range m.authenticators {
		if auth == nil {
			continue
		}

		ok, principal, err := auth.Authenticate(params)
		if err != nil {
			return false, nil, err
		}
		if ok {
			return true, principal, nil
		}
	}

	return false, nil, nil
}
package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIVerificationFailed(t *testing.T) {
	err := &APIVerificationFailed{
		Section:              "consumer",
		MissingSpecification: []string{"application/json", "application/x-yaml"},
		MissingRegistration:  []string{"text/html", "application/xml"},
	}

	expected := `missing [text/html, application/xml] consumer registrations
missing from spec file [application/json, application/x-yaml] consumer`
	assert.Equal(t, expected, err.Error())
}

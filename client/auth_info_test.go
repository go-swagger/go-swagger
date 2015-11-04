package client

import (
	"testing"

	"github.com/go-swagger/go-swagger/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestAuthInfoWriter(t *testing.T) {
	hand := AuthInfoWriterFunc(func(r Request, _ strfmt.Registry) error {
		r.SetHeaderParam("authorization", "Bearer the-token-goes-here")
		return nil
	})

	tr := new(trw)
	hand.AuthenticateRequest(tr, nil)
	assert.Equal(t, "Bearer the-token-goes-here", tr.Headers.Get("Authorization"))
}

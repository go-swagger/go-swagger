package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuth(t *testing.T) {
	r, _ := newRequest("GET", "/", nil)

	writer := BasicAuth("someone", "with a password")
	writer.AuthenticateRequest(r, nil)

	req := new(http.Request)
	req.Header = make(http.Header)
	req.Header.Set("Authorization", r.header.Get("Authorization"))
	usr, pw, ok := req.BasicAuth()
	if assert.True(t, ok) {
		assert.Equal(t, "someone", usr)
		assert.Equal(t, "with a password", pw)
	}
}

func TestAPIKeyAuth_Query(t *testing.T) {
	r, _ := newRequest("GET", "/", nil)

	writer := APIKeyAuth("api_key", "query", "the-shared-key")
	writer.AuthenticateRequest(r, nil)

	assert.Equal(t, "the-shared-key", r.query.Get("api_key"))
}

func TestAPIKeyAuth_Header(t *testing.T) {
	r, _ := newRequest("GET", "/", nil)

	writer := APIKeyAuth("x-api-token", "header", "the-shared-key")
	writer.AuthenticateRequest(r, nil)

	assert.Equal(t, "the-shared-key", r.header.Get("x-api-token"))
}

func TestBearerTokenAuth(t *testing.T) {
	r, _ := newRequest("GET", "/", nil)

	writer := BearerToken("the-shared-token")
	writer.AuthenticateRequest(r, nil)

	assert.Equal(t, "Bearer the-shared-token", r.header.Get("Authorization"))
}

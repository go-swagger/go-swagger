package security

import (
	"net/http"
	"testing"

	"github.com/casualjim/go-swagger/errors"
	"github.com/stretchr/testify/assert"
)

var tokenAuth = TokenAuthentication(func(token string) (interface{}, error) {
	if token == "token123" {
		return "admin", nil
	}
	return nil, errors.Unauthenticated("token")
})

func TestInvalidApiKeyAuthInitialization(t *testing.T) {
	assert.Panics(t, func() { APIKeyAuth("api_key", "qery", tokenAuth) })
}

func TestValidApiKeyAuth(t *testing.T) {
	ta := APIKeyAuth("api_key", "query", tokenAuth)
	ta2 := APIKeyAuth("X-API-KEY", "header", tokenAuth)

	req1, _ := http.NewRequest("GET", "/blah?api_key=token123", nil)

	ok, usr, err := ta.Authenticate(req1)
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)

	req2, _ := http.NewRequest("GET", "/blah", nil)
	req2.Header.Set("X-API-KEY", "token123")

	ok, usr, err = ta2.Authenticate(req2)
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)
}

func TestInvalidApiKeyAuth(t *testing.T) {
	ta := APIKeyAuth("api_key", "query", tokenAuth)
	ta2 := APIKeyAuth("X-API-KEY", "header", tokenAuth)

	req1, _ := http.NewRequest("GET", "/blah?api_key=token124", nil)

	ok, usr, err := ta.Authenticate(req1)
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)

	req2, _ := http.NewRequest("GET", "/blah", nil)
	req2.Header.Set("X-API-KEY", "token124")

	ok, usr, err = ta2.Authenticate(req2)
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)
}

func TestMissingApiKeyAuth(t *testing.T) {
	ta := APIKeyAuth("api_key", "query", tokenAuth)
	ta2 := APIKeyAuth("X-API-KEY", "header", tokenAuth)

	req1, _ := http.NewRequest("GET", "/blah", nil)
	req1.Header.Set("X-API-KEY", "token123")

	ok, usr, err := ta.Authenticate(req1)
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)

	req2, _ := http.NewRequest("GET", "/blah?api_key=token123", nil)

	ok, usr, err = ta2.Authenticate(req2)
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)
}

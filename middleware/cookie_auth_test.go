package middleware_test

import (
	"net/http"
	"testing"

	"github.com/go-swagger/go-swagger/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCookieAuthenticator(t *testing.T) {
	auth := middleware.NewCookieAuthenticator("session")

	t.Run("valid cookie", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: "session", Value: "valid-token"})

		ok, principal, err := auth.Authenticate(req)
		assert.True(t, ok)
		assert.Equal(t, "valid-token", principal)
		assert.NoError(t, err)
	})

	t.Run("missing cookie", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)

		ok, principal, err := auth.Authenticate(req)
		assert.False(t, ok)
		assert.Nil(t, principal)
		assert.NoError(t, err)
	})

	t.Run("invalid params", func(t *testing.T) {
		ok, principal, err := auth.Authenticate("not a request")
		assert.False(t, ok)
		assert.Nil(t, principal)
		assert.NoError(t, err)
	})
}
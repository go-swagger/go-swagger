package middleware_test

import (
	"net/http"
	"testing"

	"github.com/go-swagger/go-swagger/middleware"
	"github.com/stretchr/testify/require"
)

func TestCookieAuthenticator(t *testing.T) {
	auth := middleware.NewCookieAuthenticator("session")

	t.Run("valid cookie", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "session", Value: "valid-token"})

		ok, principal, err := auth.Authenticate(req)
		require.True(t, ok)
		require.Equal(t, "valid-token", principal)
		require.NoError(t, err)
	})

	t.Run("missing cookie", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		ok, principal, err := auth.Authenticate(req)
		require.False(t, ok)
		require.Nil(t, principal)
		require.NoError(t, err)
	})

	t.Run("invalid params", func(t *testing.T) {
		ok, principal, err := auth.Authenticate("not a request")
		require.False(t, ok)
		require.Nil(t, principal)
		require.NoError(t, err)
	})
}
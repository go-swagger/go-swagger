package middleware_test

import (
	"net/http"
	"testing"

	"github.com/go-swagger/go-swagger/middleware"
	"github.com/stretchr/testify/require"
)

func TestCookieAuthenticator(t *testing.T) {
	t.Parallel()
	auth := middleware.NewCookieAuthenticator("session")

	t.Run("valid cookie", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		req.AddCookie(&http.Cookie{Name: "session", Value: "valid-token"})

		ok, principal, err := auth.Authenticate(req)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "valid-token", principal)
	})

	t.Run("missing cookie", func(t *testing.T) {
		t.Parallel()
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)

		ok, principal, err := auth.Authenticate(req)
		require.NoError(t, err)
		require.False(t, ok)
		require.Nil(t, principal)
	})

	t.Run("invalid params", func(t *testing.T) {
		t.Parallel()
		ok, principal, err := auth.Authenticate("not a request")
		require.NoError(t, err)
		require.False(t, ok)
		require.Nil(t, principal)
	})
}
package middleware_test

import (
	"net/http"
	"testing"

	"github.com/go-swagger/go-swagger/middleware"
	"github.com/stretchr/testify/require"
)

func TestCookieAuthenticator(t *testing.T) {
	t.Parallel()

	const testCookieName = "session"
	const testCookieValue = "valid-token"

	t.Run("valid cookie", func(t *testing.T) {
		t.Parallel()
		auth := middleware.NewCookieAuthenticator(testCookieName)
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{Name: testCookieName, Value: testCookieValue})

		ok, principal, err := auth.Authenticate(req)
		require.NoError(t, err)
		require.True(t, ok, "authentication should succeed with valid cookie")
		require.Equal(t, testCookieValue, principal)
	})

	t.Run("empty cookie value", func(t *testing.T) {
		t.Parallel()
		auth := middleware.NewCookieAuthenticator(testCookieName)
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		req.AddCookie(&http.Cookie{Name: testCookieName, Value: ""})

		ok, principal, err := auth.Authenticate(req)
		require.NoError(t, err, "empty cookie should not cause error")
		require.False(t, ok, "authentication should fail with empty cookie")
		require.Nil(t, principal, "principal should be nil for empty cookie")
	})

	t.Run("missing cookie", func(t *testing.T) {
		t.Parallel()
		auth := middleware.NewCookieAuthenticator(testCookieName)
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		ok, principal, err := auth.Authenticate(req)
		require.NoError(t, err, "missing cookie should not cause error")
		require.False(t, ok, "authentication should fail with missing cookie")
		require.Nil(t, principal, "principal should be nil for missing cookie")
	})

	t.Run("invalid params", func(t *testing.T) {
		t.Parallel()
		auth := middleware.NewCookieAuthenticator(testCookieName)
		
		ok, principal, err := auth.Authenticate("not a request")
		require.NoError(t, err, "invalid params should not cause error")
		require.False(t, ok, "authentication should fail with invalid params")
		require.Nil(t, principal, "principal should be nil for invalid params")
	})

	t.Run("default cookie name", func(t *testing.T) {
		t.Parallel()
		auth := middleware.NewCookieAuthenticator("")
		require.Equal(t, testCookieName, auth.Name, "should use default cookie name")
	})
}
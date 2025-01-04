package middleware_test

import (
	"errors"
	"testing"

	"github.com/go-swagger/go-swagger/middleware"
	"github.com/stretchr/testify/require"
)

type mockAuthenticator struct {
	shouldAuthenticate bool
	principal          interface{}
	err                error
}

func (m *mockAuthenticator) Authenticate(interface{}) (bool, interface{}, error) {
	return m.shouldAuthenticate, m.principal, m.err
}

func TestMultiAuthenticator(t *testing.T) {
	t.Parallel()

	t.Run("first authenticator succeeds", func(t *testing.T) {
		t.Parallel()
		auth1 := &mockAuthenticator{shouldAuthenticate: true, principal: "user1"}
		auth2 := &mockAuthenticator{shouldAuthenticate: false}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "user1", principal)
	})

	t.Run("second authenticator succeeds", func(t *testing.T) {
		t.Parallel()
		auth1 := &mockAuthenticator{shouldAuthenticate: false}
		auth2 := &mockAuthenticator{shouldAuthenticate: true, principal: "user2"}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "user2", principal)
	})

	t.Run("all authenticators fail", func(t *testing.T) {
		t.Parallel()
		auth1 := &mockAuthenticator{shouldAuthenticate: false}
		auth2 := &mockAuthenticator{shouldAuthenticate: false}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.NoError(t, err)
		require.False(t, ok)
		require.Nil(t, principal)
	})

	t.Run("authenticator returns error", func(t *testing.T) {
		t.Parallel()
		auth1 := &mockAuthenticator{err: errors.New("auth failed")}
		auth2 := &mockAuthenticator{shouldAuthenticate: true}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.Error(t, err)
		require.False(t, ok)
		require.Nil(t, principal)
		require.EqualError(t, err, "auth failed")
	})

	t.Run("nil authenticator is skipped", func(t *testing.T) {
		t.Parallel()
		auth1 := (*mockAuthenticator)(nil)
		auth2 := &mockAuthenticator{shouldAuthenticate: true, principal: "user2"}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "user2", principal)
	})

	t.Run("no authenticators", func(t *testing.T) {
		t.Parallel()
		require.Panics(t, func() {
			_ = middleware.NewMultiAuthenticator()
		})
	})
}
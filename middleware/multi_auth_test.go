package middleware_test

import (
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
}
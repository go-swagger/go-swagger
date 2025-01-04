package middleware_test

import (
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-swagger/go-swagger/middleware"
	"github.com/stretchr/testify/require"
)

type mockAuthenticator struct {
	shouldAuthenticate bool
	principal          interface{}
	err                error
}

func (m *mockAuthenticator) Authenticate(params interface{}) (bool, interface{}, error) {
	return m.shouldAuthenticate, m.principal, m.err
}

func TestMultiAuthenticator(t *testing.T) {
	t.Run("first authenticator succeeds", func(t *testing.T) {
		auth1 := &mockAuthenticator{shouldAuthenticate: true, principal: "user1"}
		auth2 := &mockAuthenticator{shouldAuthenticate: false}

		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.True(t, ok)
		require.Equal(t, "user1", principal)
		require.NoError(t, err)
	})

	t.Run("second authenticator succeeds", func(t *testing.T) {
		auth1 := &mockAuthenticator{shouldAuthenticate: false}
		auth2 := &mockAuthenticator{shouldAuthenticate: true, principal: "user2"}

		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.True(t, ok)
		require.Equal(t, "user2", principal)
		require.NoError(t, err)
	})

	t.Run("all authenticators fail", func(t *testing.T) {
		auth1 := &mockAuthenticator{shouldAuthenticate: false}
		auth2 := &mockAuthenticator{shouldAuthenticate: false}

		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.False(t, ok)
		require.Nil(t, principal)
		require.NoError(t, err)
	})
}
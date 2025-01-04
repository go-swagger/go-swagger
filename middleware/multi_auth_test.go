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
		require.NoError(t, err, "first authenticator should succeed without error")
		require.True(t, ok, "authentication should succeed when first authenticator succeeds")
		require.Equal(t, "user1", principal, "should return principal from first authenticator")
	})

	t.Run("second authenticator succeeds", func(t *testing.T) {
		t.Parallel()
		auth1 := &mockAuthenticator{shouldAuthenticate: false}
		auth2 := &mockAuthenticator{shouldAuthenticate: true, principal: "user2"}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.NoError(t, err, "second authenticator should succeed without error")
		require.True(t, ok, "authentication should succeed when second authenticator succeeds")
		require.Equal(t, "user2", principal, "should return principal from second authenticator")
	})

	t.Run("all authenticators fail", func(t *testing.T) {
		t.Parallel()
		auth1 := &mockAuthenticator{shouldAuthenticate: false}
		auth2 := &mockAuthenticator{shouldAuthenticate: false}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.NoError(t, err, "should not return error when all authenticators fail")
		require.False(t, ok, "authentication should fail when all authenticators fail")
		require.Nil(t, principal, "principal should be nil when all authenticators fail")
	})

	t.Run("authenticator returns error", func(t *testing.T) {
		t.Parallel()
		expectedErr := errors.New("auth failed")
		auth1 := &mockAuthenticator{err: expectedErr}
		auth2 := &mockAuthenticator{shouldAuthenticate: true}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.ErrorIs(t, err, expectedErr, "should return error from authenticator")
		require.False(t, ok, "authentication should fail when authenticator returns error")
		require.Nil(t, principal, "principal should be nil when authenticator returns error")
	})

	t.Run("nil authenticator is skipped", func(t *testing.T) {
		t.Parallel()
		auth1 := (*mockAuthenticator)(nil)
		auth2 := &mockAuthenticator{shouldAuthenticate: true, principal: "user2"}
		multi := middleware.NewMultiAuthenticator(auth1, auth2)

		ok, principal, err := multi.Authenticate(nil)
		require.NoError(t, err, "should not return error when skipping nil authenticator")
		require.True(t, ok, "authentication should succeed when valid authenticator succeeds")
		require.Equal(t, "user2", principal, "should return principal from valid authenticator")
	})

	t.Run("no authenticators panic", func(t *testing.T) {
		t.Parallel()
		require.PanicsWithValue(t, "at least one authenticator is required", func() {
			_ = middleware.NewMultiAuthenticator()
		}, "should panic when no authenticators are provided")
	})
}
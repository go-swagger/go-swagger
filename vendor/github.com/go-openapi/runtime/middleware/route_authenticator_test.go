package middleware

import (
	"errors"
	"net/http"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/require"
)

type countAuthenticator struct {
	count     int
	applies   bool
	principal interface{}
	err       error
}

func (c *countAuthenticator) Authenticate(params interface{}) (bool, interface{}, error) {
	c.count++
	return c.applies, c.principal, c.err
}

func newCountAuthenticator(applies bool, principal interface{}, err error) *countAuthenticator {
	return &countAuthenticator{applies: applies, principal: principal, err: err}
}

var (
	successAuth = runtime.AuthenticatorFunc(func(_ interface{}) (bool, interface{}, error) {
		return true, "the user", nil
	})
	failAuth = runtime.AuthenticatorFunc(func(_ interface{}) (bool, interface{}, error) {
		return true, nil, errors.New("unauthenticated")
	})
	noApplyAuth = runtime.AuthenticatorFunc(func(_ interface{}) (bool, interface{}, error) {
		return false, nil, nil
	})
)

func TestAuthenticateSingle(t *testing.T) {
	ra := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth1": successAuth,
		},
		Schemes: []string{"auth1"},
		Scopes:  map[string][]string{"auth1": nil},
	}
	ras := RouteAuthenticators([]RouteAuthenticator{ra})

	require.False(t, ras.AllowsAnonymous())

	req, _ := http.NewRequest("GET", "/", nil)
	route := &MatchedRoute{}
	ok, prin, err := ras.Authenticate(req, route)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "the user", prin)

	require.Equal(t, ra, *route.Authenticator)
}

func TestAuthenticateLogicalOr(t *testing.T) {
	ra1 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth1": noApplyAuth,
		},
		Schemes: []string{"auth1"},
		Scopes:  map[string][]string{"auth1": nil},
	}
	ra2 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth2": successAuth,
		},
		Schemes: []string{"auth2"},
		Scopes:  map[string][]string{"auth2": nil},
	}
	// right side matches
	ras := RouteAuthenticators([]RouteAuthenticator{ra1, ra2})

	require.False(t, ras.AllowsAnonymous())

	req, _ := http.NewRequest("GET", "/", nil)
	route := &MatchedRoute{}
	ok, prin, err := ras.Authenticate(req, route)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "the user", prin)

	require.Equal(t, ra2, *route.Authenticator)

	// left side matches
	ras = RouteAuthenticators([]RouteAuthenticator{ra2, ra1})

	require.False(t, ras.AllowsAnonymous())

	req, _ = http.NewRequest("GET", "/", nil)
	route = &MatchedRoute{}
	ok, prin, err = ras.Authenticate(req, route)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "the user", prin)

	require.Equal(t, ra2, *route.Authenticator)
}

func TestAuthenticateLogicalAnd(t *testing.T) {
	ra1 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth1": noApplyAuth,
		},
		Schemes: []string{"auth1"},
		Scopes:  map[string][]string{"auth1": nil},
	}
	auther := newCountAuthenticator(true, "the user", nil)
	ra2 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth2": auther,
			"auth3": auther,
		},
		Schemes: []string{"auth2", "auth3"},
		Scopes:  map[string][]string{"auth2": nil},
	}
	ras := RouteAuthenticators([]RouteAuthenticator{ra1, ra2})

	require.False(t, ras.AllowsAnonymous())

	req, _ := http.NewRequest("GET", "/", nil)
	route := &MatchedRoute{}
	ok, prin, err := ras.Authenticate(req, route)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "the user", prin)

	require.Equal(t, ra2, *route.Authenticator)
	require.Equal(t, 2, auther.count)

	var count int
	successA := runtime.AuthenticatorFunc(func(_ interface{}) (bool, interface{}, error) {
		count++
		return true, "the user", nil
	})
	failA := runtime.AuthenticatorFunc(func(_ interface{}) (bool, interface{}, error) {
		count++
		return true, nil, errors.New("unauthenticated")
	})

	ra3 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth2": successA,
			"auth3": failA,
			"auth4": successA,
		},
		Schemes: []string{"auth2", "auth3", "auth4"},
		Scopes:  map[string][]string{"auth2": nil},
	}
	ras = RouteAuthenticators([]RouteAuthenticator{ra1, ra3})

	require.False(t, ras.AllowsAnonymous())

	req, _ = http.NewRequest("GET", "/", nil)
	route = &MatchedRoute{}
	ok, prin, err = ras.Authenticate(req, route)
	require.Error(t, err)
	require.True(t, ok)
	require.Nil(t, prin)

	require.Equal(t, ra3, *route.Authenticator)
	require.Equal(t, 2, count)

	ra4 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth2": successA,
			"auth3": successA,
			"auth4": failA,
		},
		Schemes: []string{"auth2", "auth3", "auth4"},
		Scopes:  map[string][]string{"auth2": nil},
	}
	ras = RouteAuthenticators([]RouteAuthenticator{ra1, ra4})

	require.False(t, ras.AllowsAnonymous())

	req, _ = http.NewRequest("GET", "/", nil)
	route = &MatchedRoute{}
	ok, prin, err = ras.Authenticate(req, route)
	require.Error(t, err)
	require.True(t, ok)
	require.Nil(t, prin)

	require.Equal(t, ra4, *route.Authenticator)
	require.Equal(t, 5, count)
}

func TestAuthenticateOptional(t *testing.T) {
	ra1 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth1": noApplyAuth,
		},
		Schemes: []string{"auth1"},
		Scopes:  map[string][]string{"auth1": nil},
	}
	ra2 := RouteAuthenticator{
		allowAnonymous: true,
	}

	ra3 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth2": noApplyAuth,
		},
		Schemes: []string{"auth2"},
		Scopes:  map[string][]string{"auth2": nil},
	}

	ras := RouteAuthenticators([]RouteAuthenticator{ra1, ra2, ra3})
	require.True(t, ras.AllowsAnonymous())

	req, _ := http.NewRequest("GET", "/", nil)
	route := &MatchedRoute{}
	ok, prin, err := ras.Authenticate(req, route)
	require.NoError(t, err)
	require.True(t, ok)
	require.Nil(t, prin)

	require.Equal(t, ra2, *route.Authenticator)

	ra4 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth1": noApplyAuth,
		},
		Schemes: []string{"auth1"},
		Scopes:  map[string][]string{"auth1": nil},
	}
	ra5 := RouteAuthenticator{
		allowAnonymous: true,
	}

	ra6 := RouteAuthenticator{
		Authenticator: map[string]runtime.Authenticator{
			"auth2": failAuth,
		},
		Schemes: []string{"auth2"},
		Scopes:  map[string][]string{"auth2": nil},
	}

	ras = RouteAuthenticators([]RouteAuthenticator{ra4, ra5, ra6})
	require.True(t, ras.AllowsAnonymous())

	req, _ = http.NewRequest("GET", "/", nil)
	route = &MatchedRoute{}
	ok, prin, err = ras.Authenticate(req, route)
	require.Error(t, err)
	require.True(t, ok)
	require.Nil(t, prin)

	require.Equal(t, ra6, *route.Authenticator)
}

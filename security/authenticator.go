package security

import (
	"net/http"
	"strings"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
)

// httpAuthenticator is a function that authenticates a HTTP request
func httpAuthenticator(handler func(*http.Request) (bool, interface{}, error)) swagger.Authenticator {
	return swagger.AuthenticatorFunc(func(params interface{}) (bool, interface{}, error) {
		if request, ok := params.(*http.Request); ok {
			return handler(request)
		}
		return false, nil, nil
	})
}

// UserPassAuthentication authentication function
type UserPassAuthentication func(string, string) (interface{}, error)

// TokenAuthentication authentication function
type TokenAuthentication func(string) (interface{}, error)

// BasicAuth creates a basic auth authenticator with the provided authentication function
func BasicAuth(authenticate UserPassAuthentication) swagger.Authenticator {
	return httpAuthenticator(func(r *http.Request) (bool, interface{}, error) {
		if usr, pass, ok := r.BasicAuth(); ok {
			p, err := authenticate(usr, pass)
			return true, p, err
		}

		return false, nil, nil
	})
}

// APIKeyAuth creates an authenticator that uses a token for authorization.
// This token can be obtained from either a header or a query string
func APIKeyAuth(name, in string, authenticate TokenAuthentication) swagger.Authenticator {
	inl := strings.ToLower(in)
	if inl != "query" && inl != "header" {
		// panic because this is most likely a typo
		panic(errors.New(500, "api key auth: in value needs to be either \"query\" or \"header\"."))
	}

	var getToken func(*http.Request) string
	switch inl {
	case "header":
		getToken = func(r *http.Request) string { return r.Header.Get(name) }
	case "query":
		getToken = func(r *http.Request) string { return r.URL.Query().Get(name) }
	}

	return httpAuthenticator(func(r *http.Request) (bool, interface{}, error) {
		token := getToken(r)
		if token == "" {
			return false, nil, nil
		}

		p, err := authenticate(token)
		return true, p, err
	})
}

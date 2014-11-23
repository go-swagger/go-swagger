package schema

import (
	"encoding/json"

	"github.com/fatih/structs"
)

type SecurityRequirement map[string][]string
type SecurityRequirements []SecurityRequirement

func BasicAuth() *SecurityScheme {
	return &SecurityScheme{Type: "basic"}
}

func ApiKeyAuth(fieldName, valueSource string) *SecurityScheme {
	return &SecurityScheme{Type: "apiKey", Name: fieldName, In: valueSource}
}

func OAuth2Implicit(authorizationURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:             "oauth2",
		Flow:             "implicit",
		AuthorizationURL: authorizationURL,
	}
}

func OAuth2Password(tokenURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:     "oauth2",
		Flow:     "password",
		TokenURL: tokenURL,
	}
}

func OAuth2Application(tokenURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:     "oauth2",
		Flow:     "application",
		TokenURL: tokenURL,
	}
}

func OAuth2AccessToken(authorizationURL, tokenURL string) *SecurityScheme {
	return &SecurityScheme{
		Type:             "oauth2",
		Flow:             "accessCode",
		AuthorizationURL: authorizationURL,
		TokenURL:         tokenURL,
	}
}

type SecurityScheme struct {
	Description      string                 `structs:"description,omitempty"`
	Extensions       map[string]interface{} `structs:"-"` // custom extensions, omitted when empty
	Type             string                 `structs:"type"`
	Name             string                 `structs:"name,omitempty"`             // api key
	In               string                 `structs:"in,omitempty"`               // api key
	Flow             string                 `structs:"flow,omitempty"`             // oauth2
	AuthorizationURL string                 `structs:"authorizationUrl,omitempty"` // oauth2
	TokenURL         string                 `structs:"tokenUrl,omitempty"`         // oauth2
	Scopes           map[string]string      `structs:"scopes,omitempty"`           // oauth2
}

func (s *SecurityScheme) AddScope(scope, description string) {
	if s.Scopes == nil {
		s.Scopes = make(map[string]string)
	}
	s.Scopes[scope] = description
}

func (s SecurityScheme) Map() map[string]interface{} {
	res := structs.Map(s)
	addExtensions(res, s.Extensions)
	return res
}

func (s SecurityScheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Map())
}

func (s SecurityScheme) MarshalYAML() (interface{}, error) {
	return s.Map(), nil
}

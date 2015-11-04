package client

import (
	"encoding/base64"

	"github.com/go-swagger/go-swagger/client"
	"github.com/go-swagger/go-swagger/strfmt"
)

// BasicAuth provides a basic auth info writer
func BasicAuth(username, password string) client.AuthInfoWriter {
	return client.AuthInfoWriterFunc(func(r client.Request, _ strfmt.Registry) error {
		encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		r.SetHeaderParam("Authorization", "Basic "+encoded)
		return nil
	})
}

// APIKeyAuth provides an API key auth info writer
func APIKeyAuth(name, in, value string) client.AuthInfoWriter {
	if in == "query" {
		return client.AuthInfoWriterFunc(func(r client.Request, _ strfmt.Registry) error {
			r.SetQueryParam(name, value)
			return nil
		})
	}

	if in == "header" {
		return client.AuthInfoWriterFunc(func(r client.Request, _ strfmt.Registry) error {
			r.SetHeaderParam(name, value)
			return nil
		})
	}
	return nil
}

// BearerToken provides a header based oauth2 bearer access token auth info writer
func BearerToken(token string) client.AuthInfoWriter {
	return client.AuthInfoWriterFunc(func(r client.Request, _ strfmt.Registry) error {
		r.SetHeaderParam("Authorization", "Bearer "+token)
		return nil
	})
}

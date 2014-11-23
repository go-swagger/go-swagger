package schema

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthSerialization(t *testing.T) {
	Convey("Auth should", t, func() {
		Convey("serialize a basic auth security scheme", func() {
			auth := BasicAuth()
			So(auth, validateJSON, `{"type":"basic"}`)
		})

		Convey("serialize a header key model", func() {
			auth := ApiKeyAuth("api-key", "header")
			So(auth, validateJSON, `{"in":"header","name":"api-key","type":"apiKey"}`)
		})

		Convey("serialize an oauth2 implicit flow model", func() {
			auth := OAuth2Implicit("http://foo.com/authorization")
			So(auth, validateJSON, `{"authorizationUrl":"http://foo.com/authorization","flow":"implicit","type":"oauth2"}`)
		})

		Convey("serialize an oauth2 password flow model", func() {
			auth := OAuth2Password("http://foo.com/token")
			So(auth, validateJSON, `{"flow":"password","tokenUrl":"http://foo.com/token","type":"oauth2"}`)
		})

		Convey("serialize an oauth2 application flow model", func() {
			auth := OAuth2Application("http://foo.com/token")
			So(auth, validateJSON, `{"flow":"application","tokenUrl":"http://foo.com/token","type":"oauth2"}`)
		})

		Convey("serialize an oauth2 access code flow model", func() {
			auth := OAuth2AccessToken("http://foo.com/authorization", "http://foo.com/token")
			So(auth, validateJSON, `{"authorizationUrl":"http://foo.com/authorization","flow":"accessCode","tokenUrl":"http://foo.com/token","type":"oauth2"}`)
		})

		Convey("serialize an oauth2 implicit flow model with scopes", func() {
			auth := OAuth2Implicit("http://foo.com/authorization")
			auth.AddScope("email", "read your email")
			So(auth, validateJSON, `{"authorizationUrl":"http://foo.com/authorization","flow":"implicit","scopes":{"email":"read your email"},"type":"oauth2"}`)
		})

		Convey("serialize an oauth2 password flow model with scopes", func() {
			auth := OAuth2Password("http://foo.com/token")
			auth.AddScope("email", "read your email")
			So(auth, validateJSON, `{"flow":"password","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`)
		})

		Convey("serialize an oauth2 application flow model with scopes", func() {
			auth := OAuth2Application("http://foo.com/token")
			auth.AddScope("email", "read your email")
			So(auth, validateJSON, `{"flow":"application","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`)
		})

		Convey("serialize an oauth2 access code flow model with scopes", func() {
			auth := OAuth2AccessToken("http://foo.com/authorization", "http://foo.com/token")
			auth.AddScope("email", "read your email")
			So(auth, validateJSON, `{"authorizationUrl":"http://foo.com/authorization","flow":"accessCode","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`)
		})
	})
}

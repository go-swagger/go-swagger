package swagger

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthSerialization(t *testing.T) {
	Convey("Auth should", t, func() {
		Convey("serialize", func() {
			Convey("basic auth security scheme", func() {
				auth := BasicAuth()
				So(auth, ShouldSerializeJSON, `{"type":"basic"}`)
			})

			Convey("header key model", func() {
				auth := ApiKeyAuth("api-key", "header")
				So(auth, ShouldSerializeJSON, `{"in":"header","name":"api-key","type":"apiKey"}`)
			})

			Convey("oauth2 implicit flow model", func() {
				auth := OAuth2Implicit("http://foo.com/authorization")
				So(auth, ShouldSerializeJSON, `{"authorizationUrl":"http://foo.com/authorization","flow":"implicit","type":"oauth2"}`)
			})

			Convey("oauth2 password flow model", func() {
				auth := OAuth2Password("http://foo.com/token")
				So(auth, ShouldSerializeJSON, `{"flow":"password","tokenUrl":"http://foo.com/token","type":"oauth2"}`)
			})

			Convey("oauth2 application flow model", func() {
				auth := OAuth2Application("http://foo.com/token")
				So(auth, ShouldSerializeJSON, `{"flow":"application","tokenUrl":"http://foo.com/token","type":"oauth2"}`)
			})

			Convey("oauth2 access code flow model", func() {
				auth := OAuth2AccessToken("http://foo.com/authorization", "http://foo.com/token")
				So(auth, ShouldSerializeJSON, `{"authorizationUrl":"http://foo.com/authorization","flow":"accessCode","tokenUrl":"http://foo.com/token","type":"oauth2"}`)
			})

			Convey("oauth2 implicit flow model with scopes", func() {
				auth := OAuth2Implicit("http://foo.com/authorization")
				auth.AddScope("email", "read your email")
				So(auth, ShouldSerializeJSON, `{"authorizationUrl":"http://foo.com/authorization","flow":"implicit","scopes":{"email":"read your email"},"type":"oauth2"}`)
			})

			Convey("oauth2 password flow model with scopes", func() {
				auth := OAuth2Password("http://foo.com/token")
				auth.AddScope("email", "read your email")
				So(auth, ShouldSerializeJSON, `{"flow":"password","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`)
			})

			Convey("oauth2 application flow model with scopes", func() {
				auth := OAuth2Application("http://foo.com/token")
				auth.AddScope("email", "read your email")
				So(auth, ShouldSerializeJSON, `{"flow":"application","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`)
			})

			Convey("oauth2 access code flow model with scopes", func() {
				auth := OAuth2AccessToken("http://foo.com/authorization", "http://foo.com/token")
				auth.AddScope("email", "read your email")
				So(auth, ShouldSerializeJSON, `{"authorizationUrl":"http://foo.com/authorization","flow":"accessCode","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`)
			})
		})
		Convey("deserialize", func() {
			Convey("basic auth security scheme", func() {
				auth := BasicAuth()
				So(`{"type":"basic"}`, ShouldParseJSON, auth)
			})

			Convey("header key model", func() {
				auth := ApiKeyAuth("api-key", "header")
				So(`{"in":"header","name":"api-key","type":"apiKey"}`, ShouldParseJSON, auth)
			})

			Convey("oauth2 implicit flow model", func() {
				auth := OAuth2Implicit("http://foo.com/authorization")
				So(`{"authorizationUrl":"http://foo.com/authorization","flow":"implicit","type":"oauth2"}`, ShouldParseJSON, auth)
			})

			Convey("oauth2 password flow model", func() {
				auth := OAuth2Password("http://foo.com/token")
				So(`{"flow":"password","tokenUrl":"http://foo.com/token","type":"oauth2"}`, ShouldParseJSON, auth)
			})

			Convey("oauth2 application flow model", func() {
				auth := OAuth2Application("http://foo.com/token")
				So(`{"flow":"application","tokenUrl":"http://foo.com/token","type":"oauth2"}`, ShouldParseJSON, auth)
			})

			Convey("oauth2 access code flow model", func() {
				auth := OAuth2AccessToken("http://foo.com/authorization", "http://foo.com/token")
				So(`{"authorizationUrl":"http://foo.com/authorization","flow":"accessCode","tokenUrl":"http://foo.com/token","type":"oauth2"}`, ShouldParseJSON, auth)
			})

			Convey("oauth2 implicit flow model with scopes", func() {
				auth := OAuth2Implicit("http://foo.com/authorization")
				auth.AddScope("email", "read your email")
				So(`{"authorizationUrl":"http://foo.com/authorization","flow":"implicit","scopes":{"email":"read your email"},"type":"oauth2"}`, ShouldParseJSON, auth)
			})

			Convey("oauth2 password flow model with scopes", func() {
				auth := OAuth2Password("http://foo.com/token")
				auth.AddScope("email", "read your email")
				So(`{"flow":"password","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`, ShouldParseJSON, auth)
			})

			Convey("oauth2 application flow model with scopes", func() {
				auth := OAuth2Application("http://foo.com/token")
				auth.AddScope("email", "read your email")
				So(`{"flow":"application","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`, ShouldParseJSON, auth)
			})

			Convey("oauth2 access code flow model with scopes", func() {
				auth := OAuth2AccessToken("http://foo.com/authorization", "http://foo.com/token")
				auth.AddScope("email", "read your email")
				So(`{"authorizationUrl":"http://foo.com/authorization","flow":"accessCode","scopes":{"email":"read your email"},"tokenUrl":"http://foo.com/token","type":"oauth2"}`, ShouldParseJSON, auth)
			})
		})
	})
}

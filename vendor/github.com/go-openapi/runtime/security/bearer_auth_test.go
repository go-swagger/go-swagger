package security

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/go-openapi/errors"
	"github.com/stretchr/testify/assert"
)

var bearerAuth = ScopedTokenAuthentication(func(token string, requiredScopes []string) (interface{}, error) {
	if token == "token123" {
		return "admin", nil
	}
	return nil, errors.Unauthenticated("bearer")
})

func TestValidBearerAuth(t *testing.T) {
	ba := BearerAuth("owners_auth", bearerAuth)

	req1, _ := http.NewRequest("GET", "/blah?access_token=token123", nil)

	ok, usr, err := ba.Authenticate(&ScopedAuthRequest{Request: req1})
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)

	req2, _ := http.NewRequest("GET", "/blah", nil)
	req2.Header.Set("Authorization", "Bearer token123")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req2})
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)

	body := url.Values(map[string][]string{})
	body.Set("access_token", "token123")
	req3, _ := http.NewRequest("POST", "/blah", strings.NewReader(body.Encode()))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req3})
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)

	mpbody := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(mpbody)
	_ = writer.WriteField("access_token", "token123")
	writer.Close()
	req4, _ := http.NewRequest("POST", "/blah", mpbody)
	req4.Header.Set("Content-Type", writer.FormDataContentType())

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req4})
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)
}

func TestInvalidBearerAuth(t *testing.T) {
	ba := BearerAuth("owners_auth", bearerAuth)

	req1, _ := http.NewRequest("GET", "/blah?access_token=token124", nil)

	ok, usr, err := ba.Authenticate(&ScopedAuthRequest{Request: req1})
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)

	req2, _ := http.NewRequest("GET", "/blah", nil)
	req2.Header.Set("Authorization", "Bearer token124")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req2})
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)

	body := url.Values(map[string][]string{})
	body.Set("access_token", "token124")
	req3, _ := http.NewRequest("POST", "/blah", strings.NewReader(body.Encode()))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req3})
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)

	mpbody := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(mpbody)
	_ = writer.WriteField("access_token", "token124")
	writer.Close()
	req4, _ := http.NewRequest("POST", "/blah", mpbody)
	req4.Header.Set("Content-Type", writer.FormDataContentType())

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req4})
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)
}

func TestMissingBearerAuth(t *testing.T) {
	ba := BearerAuth("owners_auth", bearerAuth)

	req1, _ := http.NewRequest("GET", "/blah?access_toke=token123", nil)

	ok, usr, err := ba.Authenticate(&ScopedAuthRequest{Request: req1})
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)

	req2, _ := http.NewRequest("GET", "/blah", nil)
	req2.Header.Set("Authorization", "Beare token123")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req2})
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)

	body := url.Values(map[string][]string{})
	body.Set("access_toke", "token123")
	req3, _ := http.NewRequest("POST", "/blah", strings.NewReader(body.Encode()))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req3})
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)

	mpbody := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(mpbody)
	_ = writer.WriteField("access_toke", "token123")
	writer.Close()
	req4, _ := http.NewRequest("POST", "/blah", mpbody)
	req4.Header.Set("Content-Type", writer.FormDataContentType())

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req4})
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)
}

var bearerAuthCtx = ScopedTokenAuthenticationCtx(func(ctx context.Context, token string, requiredScopes []string) (context.Context, interface{}, error) {
	if token == "token123" {
		return context.WithValue(ctx, extra, extraWisdom), "admin", nil
	}
	return context.WithValue(ctx, reason, expReason), nil, errors.Unauthenticated("bearer")
})

func TestValidBearerAuthCtx(t *testing.T) {
	ba := BearerAuthCtx("owners_auth", bearerAuthCtx)

	req1, _ := http.NewRequest("GET", "/blah?access_token=token123", nil)
	req1 = req1.WithContext(context.WithValue(req1.Context(), original, wisdom))
	ok, usr, err := ba.Authenticate(&ScopedAuthRequest{Request: req1})
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)
	assert.Equal(t, wisdom, req1.Context().Value(original))
	assert.Equal(t, extraWisdom, req1.Context().Value(extra))
	assert.Nil(t, req1.Context().Value(reason))

	req2, _ := http.NewRequest("GET", "/blah", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), original, wisdom))
	req2.Header.Set("Authorization", "Bearer token123")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req2})
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)
	assert.Equal(t, wisdom, req2.Context().Value(original))
	assert.Equal(t, extraWisdom, req2.Context().Value(extra))
	assert.Nil(t, req2.Context().Value(reason))

	body := url.Values(map[string][]string{})
	body.Set("access_token", "token123")
	req3, _ := http.NewRequest("POST", "/blah", strings.NewReader(body.Encode()))
	req3 = req3.WithContext(context.WithValue(req3.Context(), original, wisdom))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req3})
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)
	assert.Equal(t, wisdom, req3.Context().Value(original))
	assert.Equal(t, extraWisdom, req3.Context().Value(extra))
	assert.Nil(t, req3.Context().Value(reason))

	mpbody := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(mpbody)
	_ = writer.WriteField("access_token", "token123")
	writer.Close()
	req4, _ := http.NewRequest("POST", "/blah", mpbody)
	req4 = req4.WithContext(context.WithValue(req4.Context(), original, wisdom))
	req4.Header.Set("Content-Type", writer.FormDataContentType())

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req4})
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.NoError(t, err)
	assert.Equal(t, wisdom, req4.Context().Value(original))
	assert.Equal(t, extraWisdom, req4.Context().Value(extra))
	assert.Nil(t, req4.Context().Value(reason))
}

func TestInvalidBearerAuthCtx(t *testing.T) {
	ba := BearerAuthCtx("owners_auth", bearerAuthCtx)

	req1, _ := http.NewRequest("GET", "/blah?access_token=token124", nil)
	req1 = req1.WithContext(context.WithValue(req1.Context(), original, wisdom))
	ok, usr, err := ba.Authenticate(&ScopedAuthRequest{Request: req1})
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)
	assert.Equal(t, wisdom, req1.Context().Value(original))
	assert.Equal(t, expReason, req1.Context().Value(reason))
	assert.Nil(t, req1.Context().Value(extra))

	req2, _ := http.NewRequest("GET", "/blah", nil)
	req2 = req2.WithContext(context.WithValue(req2.Context(), original, wisdom))
	req2.Header.Set("Authorization", "Bearer token124")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req2})
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)
	assert.Equal(t, wisdom, req2.Context().Value(original))
	assert.Equal(t, expReason, req2.Context().Value(reason))
	assert.Nil(t, req2.Context().Value(extra))

	body := url.Values(map[string][]string{})
	body.Set("access_token", "token124")
	req3, _ := http.NewRequest("POST", "/blah", strings.NewReader(body.Encode()))
	req3 = req3.WithContext(context.WithValue(req3.Context(), original, wisdom))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req3})
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)
	assert.Equal(t, wisdom, req3.Context().Value(original))
	assert.Equal(t, expReason, req3.Context().Value(reason))
	assert.Nil(t, req3.Context().Value(extra))

	mpbody := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(mpbody)
	_ = writer.WriteField("access_token", "token124")
	writer.Close()
	req4, _ := http.NewRequest("POST", "/blah", mpbody)
	req4 = req4.WithContext(context.WithValue(req4.Context(), original, wisdom))
	req4.Header.Set("Content-Type", writer.FormDataContentType())

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req4})
	assert.True(t, ok)
	assert.Equal(t, nil, usr)
	assert.Error(t, err)
	assert.Equal(t, wisdom, req4.Context().Value(original))
	assert.Equal(t, expReason, req4.Context().Value(reason))
	assert.Nil(t, req4.Context().Value(extra))
}

func TestMissingBearerAuthCtx(t *testing.T) {
	ba := BearerAuthCtx("owners_auth", bearerAuthCtx)

	req1, _ := http.NewRequest("GET", "/blah?access_toke=token123", nil)
	req1 = req1.WithContext(context.WithValue(req1.Context(), original, wisdom))
	ok, usr, err := ba.Authenticate(&ScopedAuthRequest{Request: req1})
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)
	assert.Equal(t, wisdom, req1.Context().Value(original))
	assert.Nil(t, req1.Context().Value(reason))
	assert.Nil(t, req1.Context().Value(extra))

	req2, _ := http.NewRequest("GET", "/blah", nil)
	req2.Header.Set("Authorization", "Beare token123")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req2})
	req2 = req2.WithContext(context.WithValue(req2.Context(), original, wisdom))
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)
	assert.Equal(t, wisdom, req2.Context().Value(original))
	assert.Nil(t, req2.Context().Value(reason))
	assert.Nil(t, req2.Context().Value(extra))

	body := url.Values(map[string][]string{})
	body.Set("access_toke", "token123")
	req3, _ := http.NewRequest("POST", "/blah", strings.NewReader(body.Encode()))
	req3 = req3.WithContext(context.WithValue(req3.Context(), original, wisdom))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req3})
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)
	assert.Equal(t, wisdom, req3.Context().Value(original))
	assert.Nil(t, req3.Context().Value(reason))
	assert.Nil(t, req3.Context().Value(extra))

	mpbody := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(mpbody)
	_ = writer.WriteField("access_toke", "token123")
	writer.Close()
	req4, _ := http.NewRequest("POST", "/blah", mpbody)
	req4 = req4.WithContext(context.WithValue(req4.Context(), original, wisdom))
	req4.Header.Set("Content-Type", writer.FormDataContentType())

	ok, usr, err = ba.Authenticate(&ScopedAuthRequest{Request: req4})
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
	assert.NoError(t, err)
	assert.Equal(t, wisdom, req4.Context().Value(original))
	assert.Nil(t, req4.Context().Value(reason))
	assert.Nil(t, req4.Context().Value(extra))
}

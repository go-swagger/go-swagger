// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package security

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-openapi/errors"
	"github.com/stretchr/testify/assert"
)

var basicAuthHandler = UserPassAuthentication(func(user, pass string) (interface{}, error) {
	if user == "admin" && pass == "123456" {
		return "admin", nil
	}
	return "", errors.Unauthenticated("basic")
})

func TestValidBasicAuth(t *testing.T) {
	ba := BasicAuth(basicAuthHandler)

	req, _ := http.NewRequest("GET", "/blah", nil)
	req.SetBasicAuth("admin", "123456")
	ok, usr, err := ba.Authenticate(req)

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
}

func TestInvalidBasicAuth(t *testing.T) {
	ba := BasicAuth(basicAuthHandler)

	req, _ := http.NewRequest("GET", "/blah", nil)
	req.SetBasicAuth("admin", "admin")
	ok, usr, err := ba.Authenticate(req)

	assert.Error(t, err)
	assert.True(t, ok)
	assert.Equal(t, "", usr)
}

func TestMissingbasicAuth(t *testing.T) {
	ba := BasicAuth(basicAuthHandler)

	req, _ := http.NewRequest("GET", "/blah", nil)

	ok, usr, err := ba.Authenticate(req)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, nil, usr)
}

func TestNoRequestBasicAuth(t *testing.T) {
	ba := BasicAuth(basicAuthHandler)

	ok, usr, err := ba.Authenticate("token")

	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, usr)
}

type secTestKey uint8

const (
	original secTestKey = iota
	extra
	reason
)

const (
	wisdom      = "The man who is swimming against the stream knows the strength of it."
	extraWisdom = "Our greatest glory is not in never falling, but in rising every time we fall."
	expReason   = "I like the dreams of the future better than the history of the past."
)

var basicAuthHandlerCtx = UserPassAuthenticationCtx(func(ctx context.Context, user, pass string) (context.Context, interface{}, error) {
	if user == "admin" && pass == "123456" {
		return context.WithValue(ctx, extra, extraWisdom), "admin", nil
	}
	return context.WithValue(ctx, reason, expReason), "", errors.Unauthenticated("basic")
})

func TestValidBasicAuthCtx(t *testing.T) {
	ba := BasicAuthCtx(basicAuthHandlerCtx)

	req, _ := http.NewRequest("GET", "/blah", nil)
	req = req.WithContext(context.WithValue(req.Context(), original, wisdom))
	req.SetBasicAuth("admin", "123456")
	ok, usr, err := ba.Authenticate(req)

	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "admin", usr)
	assert.Equal(t, wisdom, req.Context().Value(original))
	assert.Equal(t, extraWisdom, req.Context().Value(extra))
	assert.Nil(t, req.Context().Value(reason))
}

func TestInvalidBasicAuthCtx(t *testing.T) {
	ba := BasicAuthCtx(basicAuthHandlerCtx)

	req, _ := http.NewRequest("GET", "/blah", nil)
	req = req.WithContext(context.WithValue(req.Context(), original, wisdom))
	req.SetBasicAuth("admin", "admin")
	ok, usr, err := ba.Authenticate(req)

	assert.Error(t, err)
	assert.True(t, ok)
	assert.Equal(t, "", usr)
	assert.Equal(t, wisdom, req.Context().Value(original))
	assert.Nil(t, req.Context().Value(extra))
	assert.Equal(t, expReason, req.Context().Value(reason))
}

func TestMissingbasicAuthCtx(t *testing.T) {
	ba := BasicAuthCtx(basicAuthHandlerCtx)

	req, _ := http.NewRequest("GET", "/blah", nil)
	req = req.WithContext(context.WithValue(req.Context(), original, wisdom))
	ok, usr, err := ba.Authenticate(req)
	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Equal(t, nil, usr)

	assert.Equal(t, wisdom, req.Context().Value(original))
	assert.Nil(t, req.Context().Value(extra))
	assert.Nil(t, req.Context().Value(reason))
}

func TestNoRequestBasicAuthCtx(t *testing.T) {
	ba := BasicAuthCtx(basicAuthHandlerCtx)

	ok, usr, err := ba.Authenticate("token")

	assert.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, usr)
}

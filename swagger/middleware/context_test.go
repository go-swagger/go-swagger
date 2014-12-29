package middleware

import (
	"net/http"
	"testing"

	"github.com/casualjim/go-swagger/swagger/httputils"
	"github.com/casualjim/go-swagger/swagger/testing/petstore"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func TestContextValidRoute(t *testing.T) {
	ctx := NewContext(petstore.NewAPI(t))
	request, _ := http.NewRequest("GET", "http://localhost:8080/api/pets", nil)

	// check there's nothing there
	_, ok := context.GetOk(request, ctxMatchedRoute)
	assert.False(t, ok)

	matched, ok := ctx.RouteInfo(request)
	assert.True(t, ok)
	assert.NotNil(t, matched)

	// check it was cached
	_, ok = context.GetOk(request, ctxMatchedRoute)
	assert.True(t, ok)

	matched, ok = ctx.RouteInfo(request)
	assert.True(t, ok)
	assert.NotNil(t, matched)
}

func TestContextInvalidRoute(t *testing.T) {
	ctx := NewContext(petstore.NewAPI(t))
	request, _ := http.NewRequest("DELETE", "http://localhost:8080/api/pets", nil)

	// check there's nothing there
	_, ok := context.GetOk(request, ctxMatchedRoute)
	assert.False(t, ok)

	matched, ok := ctx.RouteInfo(request)
	assert.False(t, ok)
	assert.Nil(t, matched)

	// check it was cached
	_, ok = context.GetOk(request, ctxMatchedRoute)
	assert.False(t, ok)

	matched, ok = ctx.RouteInfo(request)
	assert.False(t, ok)
	assert.Nil(t, matched)
}

func TestContextValidContentType(t *testing.T) {
	ct := "application/json"
	ctx := NewContext(nil, nil)

	request, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	request.Header.Set(httputils.HeaderContentType, ct)

	// check there's nothing there
	_, ok := context.GetOk(request, ctxContentType)
	assert.False(t, ok)

	// trigger the parse
	mt, _, err := ctx.ContentType(request)
	assert.NoError(t, err)
	assert.Equal(t, ct, mt)

	// check it was cached
	_, ok = context.GetOk(request, ctxContentType)
	assert.True(t, ok)

	// check if the cast works and fetch from cache too
	mt, _, err = ctx.ContentType(request)
	assert.NoError(t, err)
	assert.Equal(t, ct, mt)
}

func TestContextInvalidContentType(t *testing.T) {
	ct := "application("
	ctx := NewContext(nil, nil)

	request, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	request.Header.Set(httputils.HeaderContentType, ct)

	// check there's nothing there
	_, ok := context.GetOk(request, ctxContentType)
	assert.False(t, ok)

	// trigger the parse
	mt, _, err := ctx.ContentType(request)
	assert.Error(t, err)
	assert.Empty(t, mt)

	// check it was not cached
	_, ok = context.GetOk(request, ctxContentType)
	assert.False(t, ok)

	// check if the failure continues
	_, _, err = ctx.ContentType(request)
	assert.Error(t, err)
}

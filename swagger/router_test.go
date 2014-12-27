package swagger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naoina/denco"
	"github.com/stretchr/testify/assert"
)

func TestRouterPathConversions(t *testing.T) {
	router := DefaultRouter().(*defaultRouter)

	data := []struct{ Original, Expected string }{
		{"/", "/"},
		{"", ""},
		{"/hello", "/hello"},
		{"/hello/", "/hello"},
		{"/users/{id}", "/users/:id"},
		{"/users/{id}/", "/users/:id"},
		{"/users/{id}/friends", "/users/:id/friends"},
		{"/users/{id}/friends/", "/users/:id/friends"},
		{"/users/{id}/friends/{friend_id}", "/users/:id/friends/:friend_id"},
		{"/users/{id}/friends/{friend_id}/", "/users/:id/friends/:friend_id"},
		{"/users/{id}/friends/{friend_id}/comments", "/users/:id/friends/:friend_id/comments"},
		{"/users/{id}/friends/{friend_id}/comments/", "/users/:id/friends/:friend_id/comments"},
		{"/users/{id}/friends/{friend_id}/comments/{comment_id}", "/users/:id/friends/:friend_id/comments/:comment_id"},
		{"/users/{id}/friends/{friend_id}/{comment_id}/", "/users/:id/friends/:friend_id/:comment_id"},
	}

	for _, d := range data {
		assert.Equal(t, d.Expected, router.convertPathPattern(d.Original))
	}
}

var emptyHandler HandlerFunc = func(_ http.ResponseWriter, _ *http.Request, _ RouteParams) {}

func TestAddingRoutes(t *testing.T) {
	router := DefaultRouter().(*defaultRouter)

	data := []struct {
		Method, Path string
		Handler      denco.Handler
	}{
		{"GET", "/hello", denco.Handler{Method: "GET", Path: "/hello"}},
		{"POST", "/users", denco.Handler{Method: "POST", Path: "/users"}},
		{"PUT", "/users/{id}", denco.Handler{Method: "PUT", Path: "/users/:id"}},
		{"DELETE", "/users/{id}", denco.Handler{Method: "DELETE", Path: "/users/:id"}},
		{"PATCH", "/users/{id}", denco.Handler{Method: "PATCH", Path: "/users/:id"}},
		{"HEAD", "/users/{id}", denco.Handler{Method: "HEAD", Path: "/users/:id"}},
		{"OPTIONS", "/users/{id}", denco.Handler{Method: "OPTIONS", Path: "/users/:id"}},
	}

	for _, d := range data {
		router.AddRoute(d.Method, d.Path, emptyHandler)
	}

	for i, d := range data {
		assert.Equal(t, router.handlers[i].Method, d.Handler.Method)
		assert.Equal(t, router.handlers[i].Path, d.Handler.Path)
	}
}

func TestRouterIntegration(t *testing.T) {
	var recvParams RouteParams
	router := DefaultRouter().(*defaultRouter)

	router.AddRoute("GET", "/todo/:name", func(rw http.ResponseWriter, r *http.Request, p RouteParams) {
		recvParams = p
		rw.WriteHeader(http.StatusNoContent)
	})

	handler, err := router.Build()
	assert.NoError(t, err)

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/todo/the-thing", nil)
	assert.NoError(t, err)
	handler.ServeHTTP(rw, req)
	assert.Len(t, recvParams, 1)
	rp := recvParams.Get("name")
	assert.Equal(t, "the-thing", rp)
}

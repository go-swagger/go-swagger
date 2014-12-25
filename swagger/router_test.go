package swagger

import (
	"testing"

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

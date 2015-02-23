package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteParams(t *testing.T) {
	coll1 := RouteParams([]RouteParam{
		{"blah", "foo"},
		{"abc", "bar"},
		{"ccc", "efg"},
	})

	v := coll1.Get("blah")
	assert.Equal(t, v, "foo")
	v2 := coll1.Get("abc")
	assert.Equal(t, v2, "bar")
	v3 := coll1.Get("ccc")
	assert.Equal(t, v3, "efg")
	v4 := coll1.Get("ydkdk")
	assert.Empty(t, v4)
}

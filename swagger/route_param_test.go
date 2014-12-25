package swagger

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

	v, ok := coll1.Get("blah")
	assert.True(t, ok)
	assert.Equal(t, v, "foo")
	v2, ok := coll1.Get("abc")
	assert.True(t, ok)
	assert.Equal(t, v2, "bar")
	v3, ok := coll1.Get("ccc")
	assert.True(t, ok)
	assert.Equal(t, v3, "efg")
}

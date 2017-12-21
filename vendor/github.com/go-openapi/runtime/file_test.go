package runtime

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileImplementsIOReader(t *testing.T) {
	var file interface{} = &File{}
	expected := "that File implements io.Reader"
	assert.Implements(t, new(io.Reader), file, expected)
}

func TestFileImplementsIOReadCloser(t *testing.T) {
	var file interface{} = &File{}
	expected := "that File implements io.ReadCloser"
	assert.Implements(t, new(io.ReadCloser), file, expected)
}

package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-swagger/go-swagger/generator"
)

func TestDeprecatedFlag(t *testing.T) {
	t.Run("should detect deprecated flag and force it to false", func(t *testing.T) {
		s := Server{
			WithContext: true,
		}
		s.apply(new(generator.GenOpts))
		assert.False(t, s.WithContext)
	})
}

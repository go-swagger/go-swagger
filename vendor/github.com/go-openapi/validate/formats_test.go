package validate

import (
	"reflect"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

// Validator for string formats
func TestFormatValidator_EdgeCases(t *testing.T) {
	// Apply
	v := formatValidator{
		KnownFormats: strfmt.Default,
	}

	// formatValidator applies to: Items, Parameter,Schema

	p := spec.Parameter{}
	p.Typed("string", "email")
	s := spec.Schema{}
	s.Typed("string", "uuid")
	i := spec.Items{}
	i.Typed("string", "datetime")

	sources := []interface{}{&p, &s, &i}

	for _, source := range sources {
		// Default formats for strings
		assert.True(t, v.Applies(source, reflect.String))
		// Do not apply for number formats
		assert.False(t, v.Applies(source, reflect.Int))
	}

	assert.False(t, v.Applies("A string", reflect.String))
	assert.False(t, v.Applies(nil, reflect.String))

}

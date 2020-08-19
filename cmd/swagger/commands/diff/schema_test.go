package diff

import (
	"testing"
	"time"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestGetTypeFromSimpleSchema(t *testing.T) {
	s := spec.SimpleSchema{Type: "string"}
	ty, a := getTypeFromSimpleSchema(&s)
	assert.Equal(t, "string", ty)
	assert.False(t, a)

	arr := spec.SimpleSchema{Type: "array", Items: spec.NewItems().Typed("integer", "int32")}
	ty, a = getTypeFromSimpleSchema(&arr)
	assert.Equal(t, "integer.int32", ty)
	assert.True(t, a)

}

func TestIsArray(t *testing.T) {
	arr := spec.SimpleSchema{Type: "array", Items: spec.NewItems().Typed("integer", "int32")}
	assert.True(t, isArray(&arr))
	assert.False(t, isArray(&time.Time{}))
}

func TestIsPrimitive(t *testing.T) {
	sa := spec.StringOrArray{"string"}
	assert.True(t, isPrimitive(sa))

	s := spec.Schema{SchemaProps: spec.SchemaProps{Type: sa}}
	assert.True(t, isPrimitive(&s))
	assert.False(t, isPrimitive(&time.Time{}))

	sc := spec.Schema{}
	assert.False(t, isPrimitive(&sc))
}

func TestGetSchemaType(t *testing.T) {
	tt, a := getSchemaType(time.Time{})
	assert.False(t, a)
	assert.Equal(t, "unknown", tt)

	s := spec.SimpleSchema{Type: "string"}
	tt, a = getSchemaType(s)
	assert.False(t, a)
	assert.Equal(t, "string", tt)

}

func TestDefinitionFromRef(t *testing.T) {
	refStr := definitionFromRef(spec.MustCreateRef(""))
	assert.True(t, len(refStr) == 0)
}

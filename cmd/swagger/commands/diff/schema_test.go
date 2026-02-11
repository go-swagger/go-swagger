// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package diff

import (
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"

	"github.com/go-openapi/spec"
)

func TestGetTypeFromSimpleSchema(t *testing.T) {
	s := spec.SimpleSchema{Type: "string"}
	ty, a := getTypeFromSimpleSchema(&s)
	assert.EqualT(t, "string", ty)
	assert.FalseT(t, a)

	arr := spec.SimpleSchema{Type: "array", Items: spec.NewItems().Typed("integer", "int32")}
	ty, a = getTypeFromSimpleSchema(&arr)
	assert.EqualT(t, "integer.int32", ty)
	assert.TrueT(t, a)
}

func TestIsArray(t *testing.T) {
	arr := spec.SimpleSchema{Type: "array", Items: spec.NewItems().Typed("integer", "int32")}
	assert.TrueT(t, isArray(&arr))
	assert.FalseT(t, isArray(&time.Time{}))
}

func TestIsPrimitive(t *testing.T) {
	sa := spec.StringOrArray{"string"}
	assert.TrueT(t, isPrimitive(sa))

	s := spec.Schema{SchemaProps: spec.SchemaProps{Type: sa}}
	assert.TrueT(t, isPrimitive(&s))
	assert.FalseT(t, isPrimitive(&time.Time{}))

	sc := spec.Schema{}
	assert.FalseT(t, isPrimitive(&sc))
}

func TestGetSchemaType(t *testing.T) {
	tt, a := getSchemaType(time.Time{})
	assert.FalseT(t, a)
	assert.EqualT(t, "unknown", tt)

	s := spec.SimpleSchema{Type: "string"}
	tt, a = getSchemaType(s)
	assert.FalseT(t, a)
	assert.EqualT(t, "string", tt)
}

func TestDefinitionFromRef(t *testing.T) {
	refStr := definitionFromRef(spec.MustCreateRef(""))
	assert.Empty(t, refStr)
}

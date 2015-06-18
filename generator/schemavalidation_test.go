package generator

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func reqm(str string) *regexp.Regexp {
	return regexp.MustCompile(regexp.QuoteMeta(str))
}

func TestSchemaValidation_Strings(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NamedString"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			if assertValidation(t, "m", "m", gm.GenSchema) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("named_string.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "err := validate.MinLength", res)
						assertInCode(t, "err := validate.MaxLength", res)
						assertInCode(t, "err := validate.Pattern", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func assertInCode(t testing.TB, expr, code string) bool {
	return assert.Regexp(t, reqm(expr), code)
}

func assertValidation(t testing.TB, pth, expr string, gm GenSchema) bool {
	if !assert.True(t, gm.HasValidations) {
		return false
	}
	if !assert.Equal(t, pth, gm.Path) {
		return false
	}
	if !assert.Equal(t, expr, gm.ValueExpression) {
		return false
	}
	return true
}

//func TestSchemaValidation_StringProps(t *testing.T) {
//specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
//if assert.NoError(t, err) {
//k := "StringValidations"
//schema := specDoc.Spec().Definitions[k]

//gm, err := makeGenDefinition(k, "models", schema, specDoc)
//if assert.NoError(t, err) {
//assert.True(t, gm.HasValidations)
//}
//}
//}

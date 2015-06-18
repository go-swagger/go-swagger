package generator

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
	"github.com/stretchr/testify/assert"
)

func reqm(str string) *regexp.Regexp {
	return regexp.MustCompile(regexp.QuoteMeta(str))
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

func TestSchemaValidation_RequiredProps(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "RequiredProps"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			assert.Len(t, gm.Properties, 6)
			for _, p := range gm.Properties {
				if assert.True(t, p.Required) {
					buf := bytes.NewBuffer(nil)
					err := modelTemplate.Execute(buf, gm)
					if assert.NoError(t, err) {
						formatted, err := formatGoFile("required_props.go", buf.Bytes())
						if assert.NoError(t, err) {
							res := string(formatted)
							assertInCode(t, k+") Validate(formats", res)
							assertInCode(t, "validate"+swag.ToGoName(p.Name), res)
							assertInCode(t, "err := validate.Required", res)
							assertInCode(t, "errors.CompositeValidationError(res...)", res)
						}
					}
				}
			}
		}
	}
}

func TestSchemaValidation_Strings(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NamedString"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			if assertValidation(t, "", "m", gm.GenSchema) {
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

func TestSchemaValidation_StringProps(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "StringValidations"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			prop := gm.Properties[0]
			if assertValidation(t, "\"name\"", "m.Name", prop) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("string_validations.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "m.validateName(formats", res)
						assertInCode(t, "err := validate.MinLength(\"name", res)
						assertInCode(t, "err := validate.MaxLength(\"name", res)
						assertInCode(t, "err := validate.Pattern(\"name", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func TestSchemaValidation_NamedNumber(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NamedNumber"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			if assertValidation(t, "", "m", gm.GenSchema) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("named_number.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						//fmt.Println(res)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "err := validate.Minimum", res)
						assertInCode(t, "err := validate.Maximum", res)
						assertInCode(t, "err := validate.MultipleOf", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func TestSchemaValidation_NumberProps(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NumberValidations"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			prop := gm.Properties[0]
			if assertValidation(t, "\"age\"", "m.Age", prop) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("number_validations.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "m.validateAge(formats", res)
						assertInCode(t, "err := validate.Minimum(\"age", res)
						assertInCode(t, "err := validate.Maximum(\"age", res)
						assertInCode(t, "err := validate.MultipleOf(\"age", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func TestSchemaValidation_NamedArray(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NamedArray"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			if assertValidation(t, "", "m", gm.GenSchema) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("named_array.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "err := validate.MinItems(\"\"", res)
						assertInCode(t, "err := validate.MaxItems(\"\"", res)
						assertInCode(t, "err := validate.MinLength(strconv.Itoa(i)", res)
						assertInCode(t, "err := validate.MaxLength(strconv.Itoa(i)", res)
						assertInCode(t, "err := validate.Pattern(strconv.Itoa(i)", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func TestSchemaValidation_ArrayProps(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "ArrayValidations"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			prop := gm.Properties[0]
			if assertValidation(t, "\"tags\"", "m.Tags", prop) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("array_validations.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "m.validateTags(formats", res)
						assertInCode(t, "err := validate.MinItems(\"tags\"", res)
						assertInCode(t, "err := validate.MaxItems(\"tags\"", res)
						assertInCode(t, "err := validate.MinLength(\"tags\"+\".\"+strconv.Itoa(i)", res)
						assertInCode(t, "err := validate.MaxLength(\"tags\"+\".\"+strconv.Itoa(i)", res)
						assertInCode(t, "err := validate.Pattern(\"tags\"+\".\"+strconv.Itoa(i)", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

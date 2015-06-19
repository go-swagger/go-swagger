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
						assertInCode(t, "err := validate.MinLength(\"name\",", res)
						assertInCode(t, "err := validate.MaxLength(\"name\",", res)
						assertInCode(t, "err := validate.Pattern(\"name\",", res)
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
						assertInCode(t, "err := validate.Minimum(\"age\",", res)
						assertInCode(t, "err := validate.Maximum(\"age\",", res)
						assertInCode(t, "err := validate.MultipleOf(\"age\",", res)
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
						assertInCode(t, "err := validate.MinLength(strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MaxLength(strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.Pattern(strconv.Itoa(i),", res)
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
						assertInCode(t, "err := validate.MinLength(\"tags\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MaxLength(\"tags\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.Pattern(\"tags\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func TestSchemaValidation_NamedNestedArray(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NamedNestedArray"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			if assertValidation(t, "", "m", gm.GenSchema) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("named_nested_array.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "iNamedNestedArraySize := int64(len(m))", res)
						assertInCode(t, "iiNamedNestedArraySize := int64(len(m[i]))", res)
						assertInCode(t, "iiiNamedNestedArraySize := int64(len(m[i][ii]))", res)
						assertInCode(t, "err := validate.MinItems(\"\"", res)
						assertInCode(t, "err := validate.MaxItems(\"\"", res)
						assertInCode(t, "err := validate.MinItems(strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MaxItems(strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MinItems(strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
						assertInCode(t, "err := validate.MaxItems(strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
						assertInCode(t, "err := validate.MinLength(strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "err := validate.MaxLength(strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "err := validate.Pattern(strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func TestSchemaValidation_NestedArrayProps(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NestedArrayValidations"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			prop := gm.Properties[0]
			if assertValidation(t, "\"tags\"", "m.Tags", prop) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("nested_array_validations.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "m.validateTags(formats", res)
						assertInCode(t, "iTagsSize := int64(len(m.Tags))", res)
						assertInCode(t, "iiTagsSize := int64(len(m.Tags[i]))", res)
						assertInCode(t, "iiiTagsSize := int64(len(m.Tags[i][ii]))", res)
						assertInCode(t, "err := validate.MinItems(\"tags\"", res)
						assertInCode(t, "err := validate.MaxItems(\"tags\"", res)
						assertInCode(t, "err := validate.MinItems(\"tags\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MaxItems(\"tags\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MinItems(\"tags\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
						assertInCode(t, "err := validate.MaxItems(\"tags\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
						assertInCode(t, "err := validate.MinLength(\"tags\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "err := validate.MaxLength(\"tags\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "err := validate.Pattern(\"tags\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func TestSchemaValidation_NamedNestedObject(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NamedNestedObject"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			if assertValidation(t, "", "m", gm.GenSchema) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("named_nested_object.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, k+") validateMeta(formats", res)
						assertInCode(t, "err := validate.MinLength(\"meta\"+\".\"+\"first\",", res)
						assertInCode(t, "err := validate.MaxLength(\"meta\"+\".\"+\"first\",", res)
						assertInCode(t, "err := validate.Pattern(\"meta\"+\".\"+\"first\",", res)
						assertInCode(t, "err := validate.Minimum(\"meta\"+\".\"+\"second\",", res)
						assertInCode(t, "err := validate.Maximum(\"meta\"+\".\"+\"second\",", res)
						assertInCode(t, "err := validate.MultipleOf(\"meta\"+\".\"+\"second\",", res)
						assertInCode(t, "iThirdSize := int64(len(m.Meta.Third))", res)
						assertInCode(t, "err := validate.MinItems(\"meta\"+\".\"+\"third\",", res)
						assertInCode(t, "err := validate.MaxItems(\"meta\"+\".\"+\"third\",", res)
						assertInCode(t, "err := validate.Minimum(\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.Maximum(\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MultipleOf(\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "iFourthSize := int64(len(m.Meta.Fourth))", res)
						assertInCode(t, "iiFourthSize := int64(len(m.Meta.Fourth[i]))", res)
						assertInCode(t, "iiiFourthSize := int64(len(m.Meta.Fourth[i][ii]))", res)
						assertInCode(t, "err := validate.MinItems(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MaxItems(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MinItems(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
						assertInCode(t, "err := validate.MaxItems(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
						assertInCode(t, "err := validate.Minimum(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "err := validate.Maximum(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "err := validate.MultipleOf(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

func TestSchemaValidation_NestedObjectProps(t *testing.T) {
	specDoc, err := spec.Load("../fixtures/codegen/todolist.schemavalidation.yml")
	if assert.NoError(t, err) {
		k := "NestedObjectValidations"
		schema := specDoc.Spec().Definitions[k]

		gm, err := makeGenDefinition(k, "models", schema, specDoc)
		if assert.NoError(t, err) {
			prop := gm.Properties[0]
			if assertValidation(t, "\"args\"", "m.Args", prop) {
				buf := bytes.NewBuffer(nil)
				err := modelTemplate.Execute(buf, gm)
				if assert.NoError(t, err) {
					formatted, err := formatGoFile("nested_object_validations.go", buf.Bytes())
					if assert.NoError(t, err) {
						res := string(formatted)
						assertInCode(t, k+") Validate(formats", res)
						assertInCode(t, "m.validateArgs(formats", res)
						assertInCode(t, "err := validate.MinLength(\"args\"+\".\"+\"meta\"+\".\"+\"first\",", res)
						assertInCode(t, "err := validate.MaxLength(\"args\"+\".\"+\"meta\"+\".\"+\"first\",", res)
						assertInCode(t, "err := validate.Pattern(\"args\"+\".\"+\"meta\"+\".\"+\"first\",", res)
						assertInCode(t, "err := validate.Minimum(\"args\"+\".\"+\"meta\"+\".\"+\"second\",", res)
						assertInCode(t, "err := validate.Maximum(\"args\"+\".\"+\"meta\"+\".\"+\"second\",", res)
						assertInCode(t, "err := validate.MultipleOf(\"args\"+\".\"+\"meta\"+\".\"+\"second\",", res)
						assertInCode(t, "iThirdSize := int64(len(m.Args.Meta.Third))", res)
						assertInCode(t, "err := validate.MinItems(\"args\"+\".\"+\"meta\"+\".\"+\"third\",", res)
						assertInCode(t, "err := validate.MaxItems(\"args\"+\".\"+\"meta\"+\".\"+\"third\",", res)
						assertInCode(t, "err := validate.Minimum(\"args\"+\".\"+\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.Maximum(\"args\"+\".\"+\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MultipleOf(\"args\"+\".\"+\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "iFourthSize := int64(len(m.Args.Meta.Fourth))", res)
						assertInCode(t, "iiFourthSize := int64(len(m.Args.Meta.Fourth[i]))", res)
						assertInCode(t, "iiiFourthSize := int64(len(m.Args.Meta.Fourth[i][ii]))", res)
						assertInCode(t, "err := validate.MinItems(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MaxItems(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i),", res)
						assertInCode(t, "err := validate.MinItems(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
						assertInCode(t, "err := validate.MaxItems(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
						assertInCode(t, "err := validate.Minimum(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "err := validate.Maximum(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "err := validate.MultipleOf(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
						assertInCode(t, "errors.CompositeValidationError(res...)", res)
					}
				}
			}
		}
	}
}

// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"bytes"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaValidation_RequiredProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "RequiredProps"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)
	assert.Len(t, gm.Properties, 6)

	for _, p := range gm.Properties {
		require.True(t, p.Required)

		buf := bytes.NewBuffer(nil)
		err := templates.MustGet("model").Execute(buf, gm)
		require.NoError(t, err)

		formatted, err := opts.LanguageOpts.FormatContent("required_props.go", buf.Bytes())
		require.NoErrorf(t, err, buf.String())

		res := string(formatted)
		assertInCode(t, k+") Validate(formats", res)
		assertInCode(t, "validate"+swag.ToGoName(p.Name), res)
		assertInCode(t, "err := validate.Required", res)
		assertInCode(t, "errors.CompositeValidationError(res...)", res)
	}
}

func TestSchemaValidation_Strings(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedString"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_string.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "err := validate.MinLength", res)
	assertInCode(t, "err := validate.MaxLength", res)
	assertInCode(t, "err := validate.Pattern", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_StringProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "StringValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"name\"", "m.Name", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("string_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.validateName(formats", res)
	assertInCode(t, "err := validate.MinLength(\"name\",", res)
	assertInCode(t, "err := validate.MaxLength(\"name\",", res)
	assertInCode(t, "err := validate.Pattern(\"name\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NamedNumber(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)
	k := "NamedNumber"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_number.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "err := validate.MinimumInt", res)
	assertInCode(t, "err := validate.MaximumInt", res)
	assertInCode(t, "err := validate.MultipleOf(", res) // factor is not an integer
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NumberProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NumberValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"age\"", "m.Age", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("number_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.validateAge(formats", res)
	assertInCode(t, "err := validate.MinimumInt(\"age\",", res)
	assertInCode(t, "err := validate.MaximumInt(\"age\",", res)
	assertInCode(t, "err := validate.MultipleOf(\"age\",", res) // factor is not an integer
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NamedArray(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedArray"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_array.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "err := validate.MinItems(\"\"", res)
	assertInCode(t, "err := validate.MaxItems(\"\"", res)
	assertInCode(t, "err := validate.MinLength(strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.MaxLength(strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.Pattern(strconv.Itoa(i),", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_ArrayProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "ArrayValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"tags\"", "m.Tags", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("array_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

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

func TestSchemaValidation_NamedNestedArray(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedNestedArray"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_nested_array.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

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

func TestSchemaValidation_NestedArrayProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NestedArrayValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"tags\"", "m.Tags", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("nested_array_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

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

func TestSchemaValidation_NamedNestedObject(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedNestedObject"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_nested_object.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, k+"Meta) Validate(formats", res)
	assertInCode(t, k+") validateMeta(formats", res)
	assertInCode(t, "m.Meta.Validate(formats)", res)
	assertInCode(t, "err := validate.MinLength(\"meta\"+\".\"+\"first\",", res)
	assertInCode(t, "err := validate.MaxLength(\"meta\"+\".\"+\"first\",", res)
	assertInCode(t, "err := validate.Pattern(\"meta\"+\".\"+\"first\",", res)
	assertInCode(t, "err := validate.Minimum(\"meta\"+\".\"+\"second\",", res)
	assertInCode(t, "err := validate.Maximum(\"meta\"+\".\"+\"second\",", res)
	assertInCode(t, "err := validate.MultipleOf(\"meta\"+\".\"+\"second\",", res)
	assertInCode(t, "iThirdSize := int64(len(m.Third))", res)
	assertInCode(t, "err := validate.MinItems(\"meta\"+\".\"+\"third\",", res)
	assertInCode(t, "err := validate.MaxItems(\"meta\"+\".\"+\"third\",", res)
	assertInCode(t, "err := validate.Minimum(\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.Maximum(\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.MultipleOf(\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "iFourthSize := int64(len(m.Fourth))", res)
	assertInCode(t, "iiFourthSize := int64(len(m.Fourth[i]))", res)
	assertInCode(t, "iiiFourthSize := int64(len(m.Fourth[i][ii]))", res)
	assertInCode(t, "err := validate.MinItems(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.MaxItems(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.MinItems(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
	assertInCode(t, "err := validate.MaxItems(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
	assertInCode(t, "err := validate.Minimum(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
	assertInCode(t, "err := validate.Maximum(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
	assertInCode(t, "err := validate.MultipleOf(\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NestedObjectProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NestedObjectValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"args\"", "m.Args", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("nested_object_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, k+"Args) Validate(formats", res)
	assertInCode(t, k+"ArgsMeta) Validate(formats", res)
	assertInCode(t, "m.validateArgs(formats", res)
	assertInCode(t, "err := validate.MinLength(\"args\"+\".\"+\"meta\"+\".\"+\"first\",", res)
	assertInCode(t, "err := validate.MaxLength(\"args\"+\".\"+\"meta\"+\".\"+\"first\",", res)
	assertInCode(t, "err := validate.Pattern(\"args\"+\".\"+\"meta\"+\".\"+\"first\",", res)
	assertInCode(t, "err := validate.Minimum(\"args\"+\".\"+\"meta\"+\".\"+\"second\",", res)
	assertInCode(t, "err := validate.Maximum(\"args\"+\".\"+\"meta\"+\".\"+\"second\",", res)
	assertInCode(t, "err := validate.MultipleOf(\"args\"+\".\"+\"meta\"+\".\"+\"second\",", res)
	assertInCode(t, "iThirdSize := int64(len(m.Third))", res)
	assertInCode(t, "err := validate.MinItems(\"args\"+\".\"+\"meta\"+\".\"+\"third\",", res)
	assertInCode(t, "err := validate.MaxItems(\"args\"+\".\"+\"meta\"+\".\"+\"third\",", res)
	assertInCode(t, "err := validate.Minimum(\"args\"+\".\"+\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.Maximum(\"args\"+\".\"+\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.MultipleOf(\"args\"+\".\"+\"meta\"+\".\"+\"third\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "iFourthSize := int64(len(m.Fourth))", res)
	assertInCode(t, "iiFourthSize := int64(len(m.Fourth[i]))", res)
	assertInCode(t, "iiiFourthSize := int64(len(m.Fourth[i][ii]))", res)
	assertInCode(t, "err := validate.MinItems(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.MaxItems(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "err := validate.MinItems(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
	assertInCode(t, "err := validate.MaxItems(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
	assertInCode(t, "err := validate.Minimum(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
	assertInCode(t, "err := validate.Maximum(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
	assertInCode(t, "err := validate.MultipleOf(\"args\"+\".\"+\"meta\"+\".\"+\"fourth\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NamedArrayMulti(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedArrayMulti"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_array_multi.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, k+") validateP0(formats", res)
	assertInCode(t, k+") validateP1(formats", res)
	assertInCode(t, "err := validate.Required(\"0\",", res)
	assertInCode(t, "err := validate.MinLength(\"0\",", res)
	assertInCode(t, "err := validate.MaxLength(\"0\",", res)
	assertInCode(t, "err := validate.Pattern(\"0\",", res)
	assertInCode(t, "err := validate.Required(\"1\",", res)
	assertInCode(t, "err := validate.Minimum(\"1\",", res)
	assertInCode(t, "err := validate.Maximum(\"1\",", res)
	assertInCode(t, "err := validate.MultipleOf(\"1\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_ArrayMultiProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "ArrayMultiValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"args\"", "m.Args", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("array_multi_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.validateArgs(formats", res)
	assertInCode(t, "err := validate.Required(\"P0\",", res)
	assertInCode(t, "err := validate.MinLength(\"P0\",", res)
	assertInCode(t, "err := validate.MaxLength(\"P0\",", res)
	assertInCode(t, "err := validate.Pattern(\"P0\",", res)
	assertInCode(t, "err := validate.Required(\"P1\",", res)
	assertInCode(t, "err := validate.Minimum(\"P1\",", res)
	assertInCode(t, "err := validate.Maximum(\"P1\",", res)
	assertInCode(t, "err := validate.MultipleOf(\"P1\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NamedArrayAdditional(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedArrayAdditional"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_array_additional.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, k+") validateP0(formats", res)
	assertInCode(t, k+") validateP1(formats", res)
	assertInCode(t, "err := validate.Required(\"0\",", res)
	assertInCode(t, "err := validate.MinLength(\"0\",", res)
	assertInCode(t, "err := validate.MaxLength(\"0\",", res)
	assertInCode(t, "err := validate.Pattern(\"0\",", res)
	assertInCode(t, "err := validate.Required(\"1\",", res)
	assertInCode(t, "err := validate.Minimum(\"1\",", res)
	assertInCode(t, "err := validate.Maximum(\"1\",", res)
	assertInCode(t, "err := validate.MultipleOf(\"1\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
	assertInCode(t, "m.NamedArrayAdditionalItems[i]", res)
}

func TestSchemaValidation_ArrayAdditionalProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "ArrayAdditionalValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"args\"", "m.Args", prop)
	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("array_additional_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.validateArgs(formats", res)
	assertInCode(t, "err := validate.Required(\"P0\",", res)
	assertInCode(t, "err := validate.MinLength(\"P0\",", res)
	assertInCode(t, "err := validate.MaxLength(\"P0\",", res)
	assertInCode(t, "err := validate.Pattern(\"P0\",", res)
	assertInCode(t, "err := validate.Required(\"P1\",", res)
	assertInCode(t, "err := validate.Minimum(\"P1\",", res)
	assertInCode(t, "err := validate.Maximum(\"P1\",", res)
	assertInCode(t, "err := validate.MultipleOf(\"P1\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
	assertInCode(t, "m.ArrayAdditionalValidationsArgsTuple0Items[i]", res)
}

func TestSchemaValidation_NamedMap(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedMap"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_map.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "for k := range m {", res)
	assertInCode(t, "err := validate.MinimumInt(k,", res)
	assertInCode(t, "err := validate.MaximumInt(k,", res)
	assertInCode(t, "err := validate.MultipleOfInt(k,", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_MapProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "MapValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"meta\"", "m.Meta", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("map_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.validateMeta(formats", res)
	assertInCode(t, "for k := range m.Meta {", res)
	assertInCode(t, "err := validate.MinimumInt(\"meta\"+\".\"+k,", res)
	assertInCode(t, "err := validate.MaximumInt(\"meta\"+\".\"+k,", res)
	assertInCode(t, "err := validate.MultipleOfInt(\"meta\"+\".\"+k,", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NamedMapComplex(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedMapComplex"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_map_complex.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "for k := range m {", res)
	assertInCode(t, "val.Validate(formats)", res)
	assertInCode(t, "err := validate.MinLength(\"name\",", res)
	assertInCode(t, "err := validate.MaxLength(\"name\",", res)
	assertInCode(t, "err := validate.Pattern(\"name\",", res)
	assertInCode(t, "err := validate.MinimumInt(\"age\",", res)
	assertInCode(t, "err := validate.MaximumInt(\"age\",", res)
	assertInCode(t, "err := validate.MultipleOfInt(\"age\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_MapComplexProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "MapComplexValidations"
	schema := specDoc.Spec().Definitions[k]
	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"meta\"", "m.Meta", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("map_complex_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "for k := range m.Meta {", res)
	assertInCode(t, "val.Validate(formats)", res)
	assertInCode(t, "err := validate.MinLength(\"name\",", res)
	assertInCode(t, "err := validate.MaxLength(\"name\",", res)
	assertInCode(t, "err := validate.Pattern(\"name\",", res)
	assertInCode(t, "err := validate.MinimumInt(\"age\",", res)
	assertInCode(t, "err := validate.MaximumInt(\"age\",", res)
	assertInCode(t, "err := validate.MultipleOfInt(\"age\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NamedNestedMap(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedNestedMap"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_nested_map.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "for k := range m {", res)
	assertInCode(t, "for kk := range m[k] {", res)
	assertInCode(t, "for kkk := range m[k][kk] {", res)
	assertInCode(t, "err := validate.MinimumInt(k+\".\"+kk+\".\"+kkk,", res)
	assertInCode(t, "err := validate.MaximumInt(k+\".\"+kk+\".\"+kkk,", res)
	assertInCode(t, "err := validate.MultipleOfInt(k+\".\"+kk+\".\"+kkk,", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NestedMapProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NestedMapValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"meta\"", "m.Meta", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("nested_map_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.validateMeta(formats", res)
	assertInCode(t, "for k := range m.Meta {", res)
	assertInCode(t, "for kk := range m.Meta[k] {", res)
	assertInCode(t, "for kkk := range m.Meta[k][kk] {", res)
	assertInCode(t, "err := validate.MinimumInt(\"meta\"+\".\"+k+\".\"+kk+\".\"+kkk,", res)
	assertInCode(t, "err := validate.MaximumInt(\"meta\"+\".\"+k+\".\"+kk+\".\"+kkk,", res)
	assertInCode(t, "err := validate.MultipleOfInt(\"meta\"+\".\"+k+\".\"+kk+\".\"+kkk,", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestAdditionalProperties_Simple(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedMapComplex"
	schema := specDoc.Spec().Definitions[k]
	tr := &typeResolver{
		ModelsPackage: "",
		ModelName:     k,
		Doc:           specDoc,
	}

	sg := schemaGenContext{
		Path:         "",
		Name:         k,
		Receiver:     "m",
		IndexVar:     "i",
		ValueExpr:    "m",
		Schema:       schema,
		Required:     false,
		TypeResolver: tr,
		Named:        true,
		ExtraSchemas: make(map[string]GenSchema),
	}

	fsm, lsm, err := newMapStack(&sg)
	require.NoError(t, err)

	assert.NotNil(t, fsm.Type)
	assert.Equal(t, &schema, fsm.Type)
	assert.Equal(t, fsm, lsm)
	assert.NotNil(t, fsm.Type.AdditionalProperties)
	assert.NotNil(t, fsm.Type.AdditionalProperties.Schema)
	assert.NotNil(t, fsm.NewObj)
	assert.Nil(t, fsm.Next)
	assert.Equal(t, "#/definitions/NamedMapComplexAnon", fsm.Type.AdditionalProperties.Schema.Ref.GetURL().String())

	require.NoError(t, lsm.Build())
}

func TestAdditionalProperties_Nested(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedNestedMapComplex"
	schema := specDoc.Spec().Definitions[k]
	tr := &typeResolver{
		ModelsPackage: "",
		ModelName:     k,
		Doc:           specDoc,
	}

	sg := schemaGenContext{
		Path:         "",
		Name:         k,
		Receiver:     "m",
		IndexVar:     "i",
		ValueExpr:    "m",
		Schema:       schema,
		Required:     false,
		TypeResolver: tr,
		Named:        true,
		ExtraSchemas: make(map[string]GenSchema),
	}

	fsm, lsm, err := newMapStack(&sg)
	require.NoError(t, err)

	assert.NotNil(t, fsm.Type)
	assert.Equal(t, &schema, fsm.Type)
	assert.Empty(t, fsm.Context.Path)
	assert.NotNil(t, schema.AdditionalProperties.Schema)

	require.NotNil(t, fsm.Next)
	require.Nil(t, fsm.Previous)

	assert.NotNil(t, fsm.Type)
	assert.Equal(t, &schema, fsm.Type)
	assert.NotEqual(t, fsm, lsm)
	assert.NotNil(t, fsm.Type.AdditionalProperties)
	assert.NotNil(t, fsm.Type.AdditionalProperties.Schema)
	assert.Nil(t, fsm.NewObj)
	assert.Nil(t, fsm.Next.NewObj)
	assert.NotNil(t, fsm.Next.Previous)
	assert.NotNil(t, fsm.Next.Next)
	assert.NotNil(t, fsm.Next.Next.NewObj)
	assert.NotNil(t, fsm.Next.Next.ValueRef)
	assert.Nil(t, fsm.Next.Next.Next)
	assert.Equal(t, fsm.Next.Next, lsm)
	require.NoError(t, lsm.Build())
}

func TestSchemaValidation_NamedNestedMapComplex(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedNestedMapComplex"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)
	require.True(t, gm.AdditionalProperties.HasValidations)
	require.True(t, gm.AdditionalProperties.AdditionalProperties.HasValidations)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_nested_map_complex.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "for k := range m {", res)
	assertInCode(t, "for kk := range m[k] {", res)
	assertInCode(t, "for kkk := range m[k][kk] {", res)
	assertInCode(t, "val.Validate(formats)", res)
	assertInCode(t, "err := validate.MinLength(\"name\",", res)
	assertInCode(t, "err := validate.MaxLength(\"name\",", res)
	assertInCode(t, "err := validate.Pattern(\"name\",", res)
	assertInCode(t, "err := validate.MinimumInt(\"age\",", res)
	assertInCode(t, "err := validate.MaximumInt(\"age\",", res)
	assertInCode(t, "err := validate.MultipleOfInt(\"age\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NestedMapPropsComplex(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NestedMapComplexValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"meta\"", "m.Meta", prop)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("nested_map_complex_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.validateMeta(formats", res)
	assertInCode(t, "for k := range m.Meta {", res)
	assertInCode(t, "for kk := range m.Meta[k] {", res)
	assertInCode(t, "for kkk := range m.Meta[k][kk] {", res)
	assertInCode(t, "val.Validate(formats)", res)
	assertInCode(t, "err := validate.MinLength(\"name\",", res)
	assertInCode(t, "err := validate.MaxLength(\"name\",", res)
	assertInCode(t, "err := validate.Pattern(\"name\",", res)
	assertInCode(t, "err := validate.MinimumInt(\"age\",", res)
	assertInCode(t, "err := validate.MaximumInt(\"age\",", res)
	assertInCode(t, "err := validate.MultipleOfInt(\"age\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_NamedAllOf(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "NamedAllOf"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	requireValidation(t, "", "m", gm.GenSchema)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("named_all_of.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, k+") validateName(formats", res)
	assertInCode(t, k+") validateAge(formats", res)
	assertInCode(t, k+") validateArgs(formats", res)
	assertInCode(t, k+") validateAssoc(formats", res)
	assertInCode(t, k+") validateOpts(formats", res)
	assertInCode(t, k+") validateExtOpts(formats", res)
	assertInCode(t, k+") validateCoords(formats", res)
	assertInCode(t, "validate.MinLength(\"name\",", res)
	assertInCode(t, "validate.MinimumInt(\"age\",", res)
	assertInCode(t, "validate.MinItems(\"args\",", res)
	assertInCode(t, "validate.MinItems(\"assoc\",", res)
	assertInCode(t, "validate.MinItems(\"assoc\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "validate.MinItems(\"assoc\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
	assertInCode(t, "validate.MinLength(\"assoc\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
	assertInCode(t, "validate.Minimum(\"opts\"+\".\"+k,", res)
	assertInCode(t, "validate.MinimumInt(\"extOpts\"+\".\"+k+\".\"+kk+\".\"+kkk,", res)
	assertInCode(t, "validate.MinLength(\"coords\"+\".\"+\"name\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_AllOfProps(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "AllOfValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	prop := gm.Properties[0]
	requireValidation(t, "\"meta\"", "m.Meta", prop)
	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("all_of_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.validateMeta(formats", res)
	assertInCode(t, "validate.MinLength(\"meta\"+\".\"+\"name\",", res)
	assertInCode(t, "validate.MinimumInt(\"meta\"+\".\"+\"age\",", res)
	assertInCode(t, "validate.MinItems(\"meta\"+\".\"+\"args\",", res)
	assertInCode(t, "validate.MinItems(\"meta\"+\".\"+\"assoc\",", res)
	assertInCode(t, "validate.MinItems(\"meta\"+\".\"+\"assoc\"+\".\"+strconv.Itoa(i),", res)
	assertInCode(t, "validate.MinItems(\"meta\"+\".\"+\"assoc\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii),", res)
	assertInCode(t, "validate.MinLength(\"meta\"+\".\"+\"assoc\"+\".\"+strconv.Itoa(i)+\".\"+strconv.Itoa(ii)+\".\"+strconv.Itoa(iii),", res)
	assertInCode(t, "validate.MinimumInt(\"meta\"+\".\"+\"opts\"+\".\"+k,", res)
	assertInCode(t, "validate.MinimumInt(\"meta\"+\".\"+\"extOpts\"+\".\"+k+\".\"+kk+\".\"+kkk,", res)
	assertInCode(t, "validate.MinLength(\"meta\"+\".\"+\"coords\"+\".\"+\"name\",", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_RefedAllOf(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "RefedAllOfValidations"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)
	require.Len(t, gm.AllOf, 2)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("all_of_validations.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "m.NamedString.Validate(formats)", res)
	assertInCode(t, "m.NamedNumber.Validate(formats)", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_SimpleZeroAllowed(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "SimpleZeroAllowed"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("simple_zero_allowed.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "swag.IsZero(m.ID)", res)
	assertInCode(t, "validate.Required(\"name\", \"body\", m.Name)", res)
	assertInCode(t, "validate.MinLength(\"id\", \"body\", m.ID, 2)", res)
	assertInCode(t, "validate.Required(\"urls\", \"body\", m.Urls)", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_Pet(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "Pet"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("pet.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "swag.IsZero(m.Status)", res)
	assertInCode(t, "swag.IsZero(m.Tags)", res)
	assertInCode(t, "validate.Required(\"name\", \"body\", m.Name)", res)
	assertInCode(t, "validate.Required(\"photoUrls\", \"body\", m.PhotoUrls)", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

func TestSchemaValidation_UpdateOrg(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.schemavalidation.yml")
	require.NoError(t, err)

	k := "UpdateOrg"
	schema := specDoc.Spec().Definitions[k]

	opts := opts()
	gm, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, gm)
	require.NoError(t, err)

	formatted, err := opts.LanguageOpts.FormatContent("pet.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(formatted)
	assertInCode(t, k+") Validate(formats", res)
	assertInCode(t, "swag.IsZero(m.TagExpiration)", res)
	assertInCode(t, "validate.MinimumInt(\"tag_expiration\", \"body\", *m.TagExpiration", res)
	assertInCode(t, "validate.MaximumInt(\"tag_expiration\", \"body\", *m.TagExpiration", res)
	assertInCode(t, "errors.CompositeValidationError(res...)", res)
}

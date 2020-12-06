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
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnum_StringThing(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "StringThing"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("string_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var stringThingEnum []interface{}", res)
	assertInCode(t, k+") validateStringThingEnum(path, location string, value StringThing)", res)
	assertInCode(t, "m.validateStringThingEnum(\"\", \"body\", m)", res)
}

func TestEnum_ComposedThing(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ComposedThing"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("composed_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "m.StringThing.Validate(formats)", res)
	assertInCode(t, "var composedThingTypeNamePropEnum []interface{}", res)
	assertInCode(t, "m.validateNameEnum(\"name\", \"body\", *m.Name)", res)
	assertInCode(t, k+") validateNameEnum(path, location string, value string)", res)
}

func TestEnum_IntThing(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "IntThing"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("int_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var intThingEnum []interface{}", res)
	assertInCode(t, k+") validateIntThingEnum(path, location string, value IntThing)", res)
	assertInCode(t, "m.validateIntThingEnum(\"\", \"body\", m)", res)
}

func TestEnum_FloatThing(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "FloatThing"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("float_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var floatThingEnum []interface{}", res)
	assertInCode(t, k+") validateFloatThingEnum(path, location string, value FloatThing)", res)
	assertInCode(t, "m.validateFloatThingEnum(\"\", \"body\", m)", res)
}

func TestEnum_SliceThing(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "SliceThing"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("slice_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var sliceThingEnum []interface{}", res)
	assertInCode(t, k+") validateSliceThingEnum(path, location string, value []string)", res)
	assertInCode(t, "m.validateSliceThingEnum(\"\", \"body\", m)", res)
}

func TestEnum_SliceAndItemsThing(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "SliceAndItemsThing"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("slice_and_items_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var sliceAndItemsThingEnum []interface{}", res)
	assertInCode(t, k+") validateSliceAndItemsThingEnum(path, location string, value []string)", res)
	assertInCode(t, "m.validateSliceAndItemsThingEnum(\"\", \"body\", m)", res)
	assertInCode(t, "var sliceAndItemsThingItemsEnum []interface{}", res)
	assertInCode(t, k+") validateSliceAndItemsThingItemsEnum(path, location string, value string)", res)
	assertInCode(t, "m.validateSliceAndItemsThingItemsEnum(strconv.Itoa(i), \"body\", m[i])", res)
}

func TestEnum_SliceAndAdditionalItemsThing(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "SliceAndAdditionalItemsThing"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("slice_and_additional_items_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var sliceAndAdditionalItemsThingEnum []interface{}", res)
	assertInCode(t, k+") validateSliceAndAdditionalItemsThingEnum(path, location string, value *SliceAndAdditionalItemsThing)", res)
	assertInCode(t, "var sliceAndAdditionalItemsThingTypeP0PropEnum []interface{}", res)
	assertInCode(t, k+") validateP0Enum(path, location string, value string)", res)
	assertInCode(t, "m.validateP0Enum(\"0\", \"body\", *m.P0)", res)
	assertInCode(t, "var sliceAndAdditionalItemsThingItemsEnum []interface{}", res)
	assertInCode(t, k+") validateSliceAndAdditionalItemsThingItemsEnum(path, location string, value float32)", res)
	assertInCode(t, "m.validateSliceAndAdditionalItemsThingItemsEnum(strconv.Itoa(i+1), \"body\", m.SliceAndAdditionalItemsThingItems[i])", res)
}

func TestEnum_MapThing(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "MapThing"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("map_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var mapThingEnum []interface{}", res)
	assertInCode(t, k+") validateMapThingEnum(path, location string, value MapThing)", res)
	assertInCode(t, "m.validateMapThingEnum(\"\", \"body\", m)", res)
	assertInCode(t, "var mapThingValueEnum []interface{}", res)
	assertInCode(t, k+") validateMapThingValueEnum(path, location string, value string)", res)
	assertInCode(t, "m.validateMapThingValueEnum(k, \"body\", m[k])", res)
}

func TestEnum_ObjectThing(t *testing.T) {
	// verify that additionalItems render the same from an expanded and a flattened spec
	// known issue: there are some slight differences in generated code and variables for enum,
	// depending on how the spec has been preprocessed
	specs := []string{
		"../fixtures/codegen/todolist.enums.yml",
		"../fixtures/codegen/todolist.enums.flattened.json", // this one is the first one, after "swagger flatten"
	}
	k := "ObjectThing"
	for _, fixture := range specs {
		t.Logf("%s from spec: %s", k, fixture)
		specDoc, err := loads.Spec(fixture)
		require.NoError(t, err)

		definitions := specDoc.Spec().Definitions
		schema := definitions[k]
		opts := opts()
		genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
		require.NoError(t, err)

		buf := bytes.NewBuffer(nil)
		err = templates.MustGet("model").Execute(buf, genModel)
		require.NoError(t, err)

		ff, err := opts.LanguageOpts.FormatContent("object_thing.go", buf.Bytes())
		require.NoErrorf(t, err, buf.String())

		res := string(ff)
		// all these remain unaffected
		assertInCode(t, "var objectThingTypeNamePropEnum []interface{}", res)
		assertInCode(t, "var objectThingTypeFlowerPropEnum []interface{}", res)
		assertInCode(t, "var objectThingTypeFlourPropEnum []interface{}", res)
		assertInCode(t, "var objectThingTypeWolvesPropEnum []interface{}", res)
		assertInCode(t, "var objectThingWolvesValueEnum []interface{}", res)
		assertInCode(t, "var objectThingCatsItemsEnum []interface{}", res)
		assertInCode(t, k+") validateNameEnum(path, location string, value string)", res)
		assertInCode(t, k+") validateFlowerEnum(path, location string, value int32)", res)
		assertInCode(t, k+") validateFlourEnum(path, location string, value float32)", res)
		assertInCode(t, k+") validateWolvesEnum(path, location string, value map[string]string)", res)
		assertInCode(t, k+") validateWolvesValueEnum(path, location string, value string)", res)
		assertInCode(t, k+") validateCatsItemsEnum(path, location string, value string)", res)
		assertInCode(t, k+") validateCats(", res)
		assertInCode(t, "m.validateNameEnum(\"name\", \"body\", *m.Name)", res)
		assertInCode(t, "m.validateFlowerEnum(\"flower\", \"body\", m.Flower)", res)
		assertInCode(t, "m.validateFlourEnum(\"flour\", \"body\", m.Flour)", res)
		assertInCode(t, "m.validateWolvesEnum(\"wolves\", \"body\", m.Wolves)", res)
		assertInCode(t, "m.validateWolvesValueEnum(\"wolves\"+\".\"+k, \"body\", m.Wolves[k])", res)
		assertInCode(t, "m.validateCatsItemsEnum(\"cats\"+\".\"+strconv.Itoa(i), \"body\", m.Cats[i])", res)

		// small naming differences may be found between the expand and the flatten version of spec
		namingDifference := "Tuple0"
		pathDifference := "P"
		if strings.Contains(fixture, "flattened") {
			// when expanded, all defs are in the same template for AdditionalItems
			schema := definitions["objectThingLions"]
			genModel, err = makeGenDefinition("ObjectThingLions", "models", schema, specDoc, opts)
			require.NoError(t, err)

			buf = bytes.NewBuffer(nil)
			err := templates.MustGet("model").Execute(buf, genModel)
			require.NoError(t, err)

			ff, err := opts.LanguageOpts.FormatContent("object_thing_lions.go", buf.Bytes())
			require.NoErrorf(t, err, buf.String())

			res = string(ff)
			namingDifference = ""
			pathDifference = ""
		}
		// now common check resumes
		assertInCode(t, "var objectThingLions"+namingDifference+"TypeP0PropEnum []interface{}", res)
		assertInCode(t, "var objectThingLions"+namingDifference+"TypeP1PropEnum []interface{}", res)
		assertInCode(t, "var objectThingLions"+namingDifference+"ItemsEnum []interface{}", res)
		assertInCode(t, "m.validateP1Enum(\""+pathDifference+"1\", \"body\", *m.P1)", res)
		assertInCode(t, "m.validateP0Enum(\""+pathDifference+"0\", \"body\", *m.P0)", res)
		assertInCode(t, k+"Lions"+namingDifference+") validateObjectThingLions"+namingDifference+"ItemsEnum(path, location string, value float64)", res)

		if namingDifference != "" {
			assertInCode(t, "m.validateObjectThingLions"+namingDifference+"ItemsEnum(strconv.Itoa(i), \"body\", m.ObjectThingLions"+namingDifference+"Items[i])", res)
			assertInCode(t, "var objectThingTypeLionsPropEnum []interface{}", res)
			assertInCode(t, k+") validateLionsEnum(path, location string, value float64)", res)
		} else {
			assertInCode(t, "m.validateObjectThingLions"+namingDifference+"ItemsEnum(strconv.Itoa(i+2), \"body\", m.ObjectThingLions"+namingDifference+"Items[i])", res)
			assertInCode(t, "var objectThingLionsItemsEnum []interface{}", res)
			assertInCode(t, k+"Lions) validateObjectThingLionsItemsEnum(path, location string, value float64)", res)
		}
	}
}

func TestEnum_ComputeInstance(t *testing.T) {
	// ensure that the enum validation for the anonymous object under the delegate property
	// is rendered.

	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "ComputeInstance"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("object_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "Region *string `json:\"region\"`", res)
	assertInCode(t, "var computeInstanceTypeRegionPropEnum []interface{}", res)
	assertInCode(t, "m.validateRegionEnum(\"region\", \"body\", *m.Region)", res)
}

func TestEnum_Cluster(t *testing.T) {
	// ensure that the enum validation for the anonymous object under the delegate property
	// is rendered.
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "Cluster"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("object_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "Data *ClusterData `json:\"data\"`", res)
	assertInCode(t, `ClusterDataStatusScheduled string = "scheduled"`, res)
	assertInCode(t, `ClusterDataStatusBuilding string = "building"`, res)
	assertInCode(t, `ClusterDataStatusUp string = "up"`, res)
	assertInCode(t, `ClusterDataStatusDeleting string = "deleting"`, res)
	assertInCode(t, `ClusterDataStatusExited string = "exited"`, res)
	assertInCode(t, `ClusterDataStatusError string = "error"`, res)
}

func TestEnum_NewPrototype(t *testing.T) {
	// ensure that the enum validation for the anonymous object under the delegate property
	// is rendered.
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "NewPrototype"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("object_thing.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "ActivatingUser *NewPrototypeActivatingUser `json:\"activating_user,omitempty\"`", res)
	assertInCode(t, "Delegate *NewPrototypeDelegate `json:\"delegate\"`", res)
	assertInCode(t, "Role *string `json:\"role\"`", res)
	assertInCode(t, "var newPrototypeTypeRolePropEnum []interface{}", res)
	assertInCode(t, "var newPrototypeDelegateTypeKindPropEnum []interface{}", res)
	assertInCode(t, "m.validateDelegate(formats)", res)
	assertInCode(t, "m.validateRole(formats)", res)
	assertInCode(t, "m.validateActivatingUser(formats)", res)
	assertInCode(t, "m.Delegate.Validate(formats)", res)
	assertInCode(t, "m.ActivatingUser.Validate(formats)", res)
}

func TestEnum_Issue265(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/sodabooth.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "SodaBrand"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("soda_brand.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assert.Equal(t, 1, strings.Count(res, "m.validateSodaBrandEnum"))
}

func TestGenerateModel_Issue303(t *testing.T) {
	specDoc, e := loads.Spec("../fixtures/enhancements/303/swagger.yml")
	require.NoError(t, e)

	opts := opts()
	tt := templateTest{t, templates.MustGet("model").Lookup("schema")}
	definitions := specDoc.Spec().Definitions
	for name, schema := range definitions {
		genModel, err := makeGenDefinition(name, "models", schema, specDoc, opts)
		require.NoError(t, err)

		assert.Equal(t, name, genModel.Name)
		assert.Equal(t, name, genModel.GoType)
		assert.True(t, genModel.IsEnumCI)

		extension := genModel.Extensions["x-go-enum-ci"]
		require.NotNil(t, extension)

		xGoEnumCI, ok := extension.(bool)
		assert.True(t, ok)
		assert.True(t, xGoEnumCI)

		buf := bytes.NewBuffer(nil)
		require.NoError(t, tt.template.Execute(buf, genModel))

		ff, err := opts.LanguageOpts.FormatContent("case_insensitive_enum_definition.go", buf.Bytes())
		require.NoErrorf(t, err, buf.String())

		res := string(ff)
		assertInCode(t, `if err := validate.EnumCase(path, location, value, vegetableEnum, false); err != nil {`, res)
	}
}

func TestEnum_Issue325(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/sodabooths.json")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "SodaBrand"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("soda_brand.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, "var sodaBrandEnum []interface{}", res)
	assertInCode(t, "err := validate.EnumCase(path, location, value, sodaBrandEnum, true)", res)
	assert.Equal(t, 1, strings.Count(res, "m.validateSodaBrandEnum"))

	k = "Soda"
	schema = definitions[k]
	genModel, err = makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf = bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err = opts.LanguageOpts.FormatContent("soda.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res = string(ff)
	assertInCode(t, "var sodaTypeBrandPropEnum []interface{}", res)
	assertInCode(t, "err := validate.EnumCase(path, location, value, sodaTypeBrandPropEnum, true)", res)
	assert.Equal(t, 1, strings.Count(res, "m.validateBrandEnum"))
}

func TestEnum_Issue352(t *testing.T) {
	specDoc, err := loads.Spec("../fixtures/codegen/todolist.enums.yml")
	require.NoError(t, err)

	definitions := specDoc.Spec().Definitions
	k := "slp_action_enum"
	schema := definitions[k]
	opts := opts()
	genModel, err := makeGenDefinition(k, "models", schema, specDoc, opts)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	err = templates.MustGet("model").Execute(buf, genModel)
	require.NoError(t, err)

	ff, err := opts.LanguageOpts.FormatContent("slp_action_enum.go", buf.Bytes())
	require.NoErrorf(t, err, buf.String())

	res := string(ff)
	assertInCode(t, ", value SlpActionEnum", res)
}

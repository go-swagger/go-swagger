package parser

import (
	goparser "go/parser"
	"log"
	"testing"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestParamsParser(t *testing.T) {
	docFile := "../fixtures/goparsing/classification/operations/noparams.go"
	fileTree, err := goparser.ParseFile(classificationProg.Fset, docFile, nil, goparser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	sp := newParameterParser(classificationProg)
	noParamOps := make(map[string]spec.Operation)
	err = sp.Parse(fileTree, noParamOps)
	if err != nil {
		log.Fatal(err)
	}
	assert.Len(t, noParamOps, 2)

	op, ok := noParamOps["someOperation"]
	assert.True(t, ok)
	assert.Len(t, op.Parameters, 6)

	for _, params := range op.Parameters {
		switch params.Name {
		case "id":
			assert.Equal(t, "ID of this no model instance.\nids in this application start at 11 and are smaller than 1000", params.Description)
			assert.Equal(t, "path", params.In)
			assert.Equal(t, "number", params.Type)
			assert.Equal(t, "int64", params.Format)
			assert.True(t, params.Required)
			assert.Equal(t, "ID", params.Extensions["x-go-name"])
			assert.EqualValues(t, 1000, *params.Maximum)
			assert.True(t, params.ExclusiveMaximum)
			assert.EqualValues(t, 10, *params.Minimum)
			assert.True(t, params.ExclusiveMinimum)

		case "score":
			assert.Equal(t, "The Score of this model", params.Description)
			assert.Equal(t, "query", params.In)
			assert.Equal(t, "number", params.Type)
			assert.Equal(t, "int32", params.Format)
			assert.True(t, params.Required)
			assert.Equal(t, "Score", params.Extensions["x-go-name"])
			assert.EqualValues(t, 45, *params.Maximum)
			assert.False(t, params.ExclusiveMaximum)
			assert.EqualValues(t, 3, *params.Minimum)
			assert.False(t, params.ExclusiveMinimum)

		case "x-hdr-name":
			assert.Equal(t, "Name of this no model instance", params.Description)
			assert.Equal(t, "header", params.In)
			assert.Equal(t, "string", params.Type)
			assert.True(t, params.Required)
			assert.Equal(t, "Name", params.Extensions["x-go-name"])
			assert.EqualValues(t, 4, *params.MinLength)
			assert.EqualValues(t, 50, *params.MaxLength)
			assert.Equal(t, "[A-Za-z0-9-.]*", params.Pattern)

		case "created":
			assert.Equal(t, "Created holds the time when this entry was created", params.Description)
			assert.Equal(t, "query", params.In)
			assert.Equal(t, "string", params.Type)
			assert.Equal(t, "date-time", params.Format)
			assert.False(t, params.Required)
			assert.Equal(t, "Created", params.Extensions["x-go-name"])

		case "foo_slice":
			assert.Equal(t, "a FooSlice has foos which are strings", params.Description)
			assert.Equal(t, "FooSlice", params.Extensions["x-go-name"])
			assert.Equal(t, "query", params.In)
			assert.Equal(t, "array", params.Type)
			assert.False(t, params.Required)
			assert.True(t, params.UniqueItems)
			assert.Equal(t, "pipe", params.CollectionFormat)
			assert.NotNil(t, params.Items, "foo_slice should have had an items property")
			assert.EqualValues(t, 3, *params.MinItems, "'foo_slice' should have had 3 min items")
			assert.EqualValues(t, 10, *params.MaxItems, "'foo_slice' should have had 10 max items")
			itprop := params.Items
			assert.EqualValues(t, 3, *itprop.MinLength, "'foo_slice.items.minLength' should have been 3")
			assert.EqualValues(t, 10, *itprop.MaxLength, "'foo_slice.items.maxLength' should have been 10")
			assert.EqualValues(t, "\\w+", itprop.Pattern, "'foo_slice.items.pattern' should have \\w+")

		case "items":
			assert.Equal(t, "Items", params.Extensions["x-go-name"])
			assert.Equal(t, "body", params.In)
			assert.NotNil(t, params.Schema)
			aprop := params.Schema
			assert.Equal(t, "array", aprop.Type[0])
			assert.NotNil(t, aprop.Items)
			assert.NotNil(t, aprop.Items.Schema)
			itprop := aprop.Items.Schema
			assert.Len(t, itprop.Properties, 4)
			assert.Len(t, itprop.Required, 3)
			assertProperty(t, itprop, "number", "id", "int32", "ID")
			iprop, ok := itprop.Properties["id"]
			assert.True(t, ok)
			assert.Equal(t, "ID of this no model instance.\nids in this application start at 11 and are smaller than 1000", iprop.Description)
			assert.EqualValues(t, 1000, *iprop.Maximum)
			assert.True(t, iprop.ExclusiveMaximum, "'id' should have had an exclusive maximum")
			assert.NotNil(t, iprop.Minimum)
			assert.EqualValues(t, 10, *iprop.Minimum)
			assert.True(t, iprop.ExclusiveMinimum, "'id' should have had an exclusive minimum")

			assertRef(t, itprop, "pet", "Pet", "#/definitions/Pet")
			iprop, ok = itprop.Properties["pet"]
			assert.True(t, ok)
			assert.Equal(t, "The Pet to add to this NoModel items bucket.\nPets can appear more than once in the bucket", iprop.Description)

			assertProperty(t, itprop, "number", "quantity", "int16", "Quantity")
			iprop, ok = itprop.Properties["quantity"]
			assert.True(t, ok)
			assert.Equal(t, "The amount of pets to add to this bucket.", iprop.Description)
			assert.EqualValues(t, 1, *iprop.Minimum)
			assert.EqualValues(t, 10, *iprop.Maximum)

			assertProperty(t, itprop, "string", "notes", "", "Notes")
			iprop, ok = itprop.Properties["notes"]
			assert.True(t, ok)
			assert.Equal(t, "Notes to add to this item.\nThis can be used to add special instructions.", iprop.Description)

		default:
			assert.Fail(t, "unkown property: "+params.Name)
		}
	}
}

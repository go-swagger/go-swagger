package codescan

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/go-openapi/spec"
)

const epsilon = 1e-9

func TestSchemaBuilder_Struct_Tag(t *testing.T) {
	sctx := loadPetstorePkgsCtx(t)
	var td *entityDecl
	for k := range sctx.app.Models {
		if k.Name != "Tag" {
			continue
		}
		td = sctx.app.Models[k]
		break
	}
	require.NotNil(t, td)

	prs := &schemaBuilder{
		ctx:  sctx,
		decl: td,
	}
	result := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(result))
}

func TestSchemaBuilder_Struct_Pet(t *testing.T) {
	// Debug = true
	// defer func() { Debug = false }()

	sctx := loadPetstorePkgsCtx(t)
	var td *entityDecl
	for k := range sctx.app.Models {
		if k.Name != "Pet" {
			continue
		}
		td = sctx.app.Models[k]
		break
	}
	require.NotNil(t, td)

	prs := &schemaBuilder{
		ctx:  sctx,
		decl: td,
	}
	result := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(result))
}

func TestSchemaBuilder_Struct_Order(t *testing.T) {
	// Debug = true
	// defer func() { Debug = false }()

	sctx := loadPetstorePkgsCtx(t)
	var td *entityDecl
	for k := range sctx.app.Models {
		if k.Name != "Order" {
			continue
		}
		td = sctx.app.Models[k]
		break
	}
	require.NotNil(t, td)

	prs := &schemaBuilder{
		ctx:  sctx,
		decl: td,
	}
	result := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(result))
}

func TestSchemaBuilder(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "NoModel")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))
	schema := models["NoModel"]

	assert.Equal(t, spec.StringOrArray([]string{"object"}), schema.Type)
	assert.Equal(t, "NoModel is a struct without an annotation.", schema.Title)
	assert.Equal(t, "NoModel exists in a package\nbut is not annotated with the swagger model annotations\nso it should now show up in a test.", schema.Description)
	assert.Len(t, schema.Required, 3)
	assert.Len(t, schema.Properties, 12)

	assertProperty(t, &schema, "integer", "id", "int64", "ID")
	prop, ok := schema.Properties["id"]
	assert.Equal(t, "ID of this no model instance.\nids in this application start at 11 and are smaller than 1000", prop.Description)
	assert.True(t, ok, "should have had an 'id' property")
	assert.InDelta(t, 1000.00, *prop.Maximum, epsilon)
	assert.True(t, prop.ExclusiveMaximum, "'id' should have had an exclusive maximum")
	assert.NotNil(t, prop.Minimum)
	assert.InDelta(t, 10.00, *prop.Minimum, epsilon)
	assert.True(t, prop.ExclusiveMinimum, "'id' should have had an exclusive minimum")
	assert.Equal(t, 11, prop.Default, "ID default value is incorrect")

	assertProperty(t, &schema, "string", "NoNameOmitEmpty", "", "")
	prop, ok = schema.Properties["NoNameOmitEmpty"]
	assert.Equal(t, "A field which has omitempty set but no name", prop.Description)
	assert.True(t, ok, "should have had an 'NoNameOmitEmpty' property")

	assertProperty(t, &schema, "string", "noteb64", "byte", "Note")
	prop, ok = schema.Properties["noteb64"]
	assert.True(t, ok, "should have a 'noteb64' property")
	assert.Nil(t, prop.Items)

	assertProperty(t, &schema, "integer", "score", "int32", "Score")
	prop, ok = schema.Properties["score"]
	assert.Equal(t, "The Score of this model", prop.Description)
	assert.True(t, ok, "should have had a 'score' property")
	assert.InDelta(t, 45.00, *prop.Maximum, epsilon)
	assert.False(t, prop.ExclusiveMaximum, "'score' should not have had an exclusive maximum")
	assert.NotNil(t, prop.Minimum)
	assert.InDelta(t, 3.00, *prop.Minimum, epsilon)
	assert.False(t, prop.ExclusiveMinimum, "'score' should not have had an exclusive minimum")
	assert.EqualValues(t, 27, prop.Example)

	expectedNameExtensions := spec.Extensions{
		"x-go-name": "Name",
		"x-property-array": []any{
			"value1",
			"value2",
		},
		"x-property-array-obj": []any{
			map[string]any{
				"name":  "obj",
				"value": "field",
			},
		},
		"x-property-value": "value",
	}

	assertProperty(t, &schema, "string", "name", "", "Name")
	prop, ok = schema.Properties["name"]
	assert.True(t, ok)
	assert.Equal(t, "Name of this no model instance", prop.Description)
	require.NotNil(t, prop.MinLength)
	require.NotNil(t, prop.MaxLength)
	assert.Equal(t, int64(4), *prop.MinLength)
	assert.Equal(t, int64(50), *prop.MaxLength)
	assert.Equal(t, "[A-Za-z0-9-.]*", prop.Pattern)
	assert.Equal(t, expectedNameExtensions, prop.Extensions)

	assertProperty(t, &schema, "string", "created", "date-time", "Created")
	prop, ok = schema.Properties["created"]
	assert.Equal(t, "Created holds the time when this entry was created", prop.Description)
	assert.True(t, ok, "should have a 'created' property")
	assert.True(t, prop.ReadOnly, "'created' should be read only")

	assertProperty(t, &schema, "string", "gocreated", "date-time", "GoTimeCreated")
	prop, ok = schema.Properties["gocreated"]
	assert.Equal(t, "GoTimeCreated holds the time when this entry was created in go time.Time", prop.Description)
	assert.True(t, ok, "should have a 'gocreated' property")

	assertArrayProperty(t, &schema, "string", "foo_slice", "", "FooSlice")
	prop, ok = schema.Properties["foo_slice"]
	assert.Equal(t, "a FooSlice has foos which are strings", prop.Description)
	assert.True(t, ok, "should have a 'foo_slice' property")
	require.NotNil(t, prop.Items, "foo_slice should have had an items property")
	require.NotNil(t, prop.Items.Schema, "foo_slice.items should have had a schema property")
	assert.True(t, prop.UniqueItems, "'foo_slice' should have unique items")
	assert.Equal(t, int64(3), *prop.MinItems, "'foo_slice' should have had 3 min items")
	assert.Equal(t, int64(10), *prop.MaxItems, "'foo_slice' should have had 10 max items")
	itprop := prop.Items.Schema
	assert.Equal(t, int64(3), *itprop.MinLength, "'foo_slice.items.minLength' should have been 3")
	assert.Equal(t, int64(10), *itprop.MaxLength, "'foo_slice.items.maxLength' should have been 10")
	assert.Equal(t, "\\w+", itprop.Pattern, "'foo_slice.items.pattern' should have \\w+")

	assertArrayProperty(t, &schema, "string", "time_slice", "date-time", "TimeSlice")
	prop, ok = schema.Properties["time_slice"]
	assert.Equal(t, "a TimeSlice is a slice of times", prop.Description)
	assert.True(t, ok, "should have a 'time_slice' property")
	require.NotNil(t, prop.Items, "time_slice should have had an items property")
	require.NotNil(t, prop.Items.Schema, "time_slice.items should have had a schema property")
	assert.True(t, prop.UniqueItems, "'time_slice' should have unique items")
	assert.Equal(t, int64(3), *prop.MinItems, "'time_slice' should have had 3 min items")
	assert.Equal(t, int64(10), *prop.MaxItems, "'time_slice' should have had 10 max items")

	assertArrayProperty(t, &schema, "array", "bar_slice", "", "BarSlice")
	prop, ok = schema.Properties["bar_slice"]
	assert.Equal(t, "a BarSlice has bars which are strings", prop.Description)
	assert.True(t, ok, "should have a 'bar_slice' property")
	require.NotNil(t, prop.Items, "bar_slice should have had an items property")
	require.NotNil(t, prop.Items.Schema, "bar_slice.items should have had a schema property")
	assert.True(t, prop.UniqueItems, "'bar_slice' should have unique items")
	assert.Equal(t, int64(3), *prop.MinItems, "'bar_slice' should have had 3 min items")
	assert.Equal(t, int64(10), *prop.MaxItems, "'bar_slice' should have had 10 max items")

	itprop = prop.Items.Schema
	require.NotNil(t, itprop)
	assert.Equal(t, int64(4), *itprop.MinItems, "'bar_slice.items.minItems' should have been 4")
	assert.Equal(t, int64(9), *itprop.MaxItems, "'bar_slice.items.maxItems' should have been 9")

	itprop2 := itprop.Items.Schema
	require.NotNil(t, itprop2)
	assert.Equal(t, int64(5), *itprop2.MinItems, "'bar_slice.items.items.minItems' should have been 5")
	assert.Equal(t, int64(8), *itprop2.MaxItems, "'bar_slice.items.items.maxItems' should have been 8")

	itprop3 := itprop2.Items.Schema
	require.NotNil(t, itprop3)
	assert.Equal(t, int64(3), *itprop3.MinLength, "'bar_slice.items.items.items.minLength' should have been 3")
	assert.Equal(t, int64(10), *itprop3.MaxLength, "'bar_slice.items.items.items.maxLength' should have been 10")
	assert.Equal(t, "\\w+", itprop3.Pattern, "'bar_slice.items.items.items.pattern' should have \\w+")

	assertArrayProperty(t, &schema, "array", "deep_time_slice", "", "DeepTimeSlice")
	prop, ok = schema.Properties["deep_time_slice"]
	assert.Equal(t, "a DeepSlice has bars which are time", prop.Description)
	assert.True(t, ok, "should have a 'deep_time_slice' property")
	require.NotNil(t, prop.Items, "deep_time_slice should have had an items property")
	require.NotNil(t, prop.Items.Schema, "deep_time_slice.items should have had a schema property")
	assert.True(t, prop.UniqueItems, "'deep_time_slice' should have unique items")
	assert.Equal(t, int64(3), *prop.MinItems, "'deep_time_slice' should have had 3 min items")
	assert.Equal(t, int64(10), *prop.MaxItems, "'deep_time_slice' should have had 10 max items")
	itprop = prop.Items.Schema
	require.NotNil(t, itprop)
	assert.Equal(t, int64(4), *itprop.MinItems, "'deep_time_slice.items.minItems' should have been 4")
	assert.Equal(t, int64(9), *itprop.MaxItems, "'deep_time_slice.items.maxItems' should have been 9")

	itprop2 = itprop.Items.Schema
	require.NotNil(t, itprop2)
	assert.Equal(t, int64(5), *itprop2.MinItems, "'deep_time_slice.items.items.minItems' should have been 5")
	assert.Equal(t, int64(8), *itprop2.MaxItems, "'deep_time_slice.items.items.maxItems' should have been 8")

	itprop3 = itprop2.Items.Schema
	require.NotNil(t, itprop3)

	assertArrayProperty(t, &schema, "object", "items", "", "Items")
	prop, ok = schema.Properties["items"]
	assert.True(t, ok, "should have an 'items' slice")
	assert.NotNil(t, prop.Items, "items should have had an items property")
	assert.NotNil(t, prop.Items.Schema, "items.items should have had a schema property")
	itprop = prop.Items.Schema
	assert.Len(t, itprop.Properties, 5)
	assert.Len(t, itprop.Required, 4)
	assertProperty(t, itprop, "integer", "id", "int32", "ID")
	iprop, ok := itprop.Properties["id"]
	assert.True(t, ok)
	assert.Equal(t, "ID of this no model instance.\nids in this application start at 11 and are smaller than 1000", iprop.Description)
	require.NotNil(t, iprop.Maximum)
	assert.InDelta(t, 1000.00, *iprop.Maximum, epsilon)
	assert.True(t, iprop.ExclusiveMaximum, "'id' should have had an exclusive maximum")
	require.NotNil(t, iprop.Minimum)
	assert.InDelta(t, 10.00, *iprop.Minimum, epsilon)
	assert.True(t, iprop.ExclusiveMinimum, "'id' should have had an exclusive minimum")
	assert.Equal(t, 11, iprop.Default, "ID default value is incorrect")

	assertRef(t, itprop, "pet", "Pet", "#/definitions/pet")
	iprop, ok = itprop.Properties["pet"]
	assert.True(t, ok)
	if itprop.Ref.String() != "" {
		assert.Equal(t, "The Pet to add to this NoModel items bucket.\nPets can appear more than once in the bucket", iprop.Description)
	}

	assertProperty(t, itprop, "integer", "quantity", "int16", "Quantity")
	iprop, ok = itprop.Properties["quantity"]
	assert.True(t, ok)
	assert.Equal(t, "The amount of pets to add to this bucket.", iprop.Description)
	assert.InDelta(t, 1.00, *iprop.Minimum, epsilon)
	assert.InDelta(t, 10.00, *iprop.Maximum, epsilon)

	assertProperty(t, itprop, "string", "expiration", "date-time", "Expiration")
	iprop, ok = itprop.Properties["expiration"]
	assert.True(t, ok)
	assert.Equal(t, "A dummy expiration date.", iprop.Description)

	assertProperty(t, itprop, "string", "notes", "", "Notes")
	iprop, ok = itprop.Properties["notes"]
	assert.True(t, ok)
	assert.Equal(t, "Notes to add to this item.\nThis can be used to add special instructions.", iprop.Description)

	decl2 := getClassificationModel(sctx, "StoreOrder")
	require.NotNil(t, decl2)
	require.NoError(t, (&schemaBuilder{decl: decl2, ctx: sctx}).Build(models))
	msch, ok := models["order"]
	pn := "github.com/go-swagger/go-swagger/fixtures/goparsing/classification/models"
	assert.True(t, ok)
	assert.Equal(t, pn, msch.Extensions["x-go-package"])
	assert.Equal(t, "StoreOrder", msch.Extensions["x-go-name"])
}

func TestSchemaBuilder_AddExtensions(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	models := make(map[string]spec.Schema)
	decl := getClassificationModel(sctx, "StoreOrder")
	require.NotNil(t, decl)
	require.NoError(t, (&schemaBuilder{decl: decl, ctx: sctx}).Build(models))

	msch, ok := models["order"]
	pn := "github.com/go-swagger/go-swagger/fixtures/goparsing/classification/models"
	assert.True(t, ok)
	assert.Equal(t, pn, msch.Extensions["x-go-package"])
	assert.Equal(t, "StoreOrder", msch.Extensions["x-go-name"])
	assert.Equal(t, "StoreOrder represents an order in this application.", msch.Title)
}

func TestTextMarhalCustomType(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "TextMarshalModel")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))
	schema := models["TextMarshalModel"]
	assertProperty(t, &schema, "string", "id", "uuid", "ID")
	assertArrayProperty(t, &schema, "string", "ids", "uuid", "IDs")
	assertProperty(t, &schema, "string", "struct", "", "Struct")
	assertProperty(t, &schema, "string", "map", "", "Map")
	assertMapProperty(t, &schema, "string", "mapUUID", "uuid", "MapUUID")
	assertRef(t, &schema, "url", "URL", "#/definitions/URL")
	assertProperty(t, &schema, "string", "time", "date-time", "Time")
	assertProperty(t, &schema, "string", "structStrfmt", "date-time", "StructStrfmt")
	assertProperty(t, &schema, "string", "structStrfmtPtr", "date-time", "StructStrfmtPtr")
	assertProperty(t, &schema, "string", "customUrl", "url", "CustomURL")
}

func TestEmbeddedTypes(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "ComplexerOne")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))
	schema := models["ComplexerOne"]
	assertProperty(t, &schema, "integer", "age", "int32", "Age")
	assertProperty(t, &schema, "integer", "id", "int64", "ID")
	assertProperty(t, &schema, "string", "createdAt", "date-time", "CreatedAt")
	assertProperty(t, &schema, "string", "extra", "", "Extra")
	assertProperty(t, &schema, "string", "name", "", "Name")
	assertProperty(t, &schema, "string", "notes", "", "Notes")
}

func TestParsePrimitiveSchemaProperty(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "PrimateModel")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["PrimateModel"]
	assertProperty(t, &schema, "boolean", "a", "", "A")
	assertProperty(t, &schema, "integer", "b", "int32", "B")
	assertProperty(t, &schema, "string", "c", "", "C")
	assertProperty(t, &schema, "integer", "d", "int64", "D")
	assertProperty(t, &schema, "integer", "e", "int8", "E")
	assertProperty(t, &schema, "integer", "f", "int16", "F")
	assertProperty(t, &schema, "integer", "g", "int32", "G")
	assertProperty(t, &schema, "integer", "h", "int64", "H")
	assertProperty(t, &schema, "integer", "i", "uint64", "I")
	assertProperty(t, &schema, "integer", "j", "uint8", "J")
	assertProperty(t, &schema, "integer", "k", "uint16", "K")
	assertProperty(t, &schema, "integer", "l", "uint32", "L")
	assertProperty(t, &schema, "integer", "m", "uint64", "M")
	assertProperty(t, &schema, "number", "n", "float", "N")
	assertProperty(t, &schema, "number", "o", "double", "O")
	assertProperty(t, &schema, "integer", "p", "uint8", "P")
	assertProperty(t, &schema, "integer", "q", "uint64", "Q")
}

func TestParseStringFormatSchemaProperty(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "FormattedModel")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["FormattedModel"]
	assertProperty(t, &schema, "string", "a", "byte", "A")
	assertProperty(t, &schema, "string", "b", "creditcard", "B")
	assertProperty(t, &schema, "string", "c", "date", "C")
	assertProperty(t, &schema, "string", "d", "date-time", "D")
	assertProperty(t, &schema, "string", "e", "duration", "E")
	assertProperty(t, &schema, "string", "f", "email", "F")
	assertProperty(t, &schema, "string", "g", "hexcolor", "G")
	assertProperty(t, &schema, "string", "h", "hostname", "H")
	assertProperty(t, &schema, "string", "i", "ipv4", "I")
	assertProperty(t, &schema, "string", "j", "ipv6", "J")
	assertProperty(t, &schema, "string", "k", "isbn", "K")
	assertProperty(t, &schema, "string", "l", "isbn10", "L")
	assertProperty(t, &schema, "string", "m", "isbn13", "M")
	assertProperty(t, &schema, "string", "n", "rgbcolor", "N")
	assertProperty(t, &schema, "string", "o", "ssn", "O")
	assertProperty(t, &schema, "string", "p", "uri", "P")
	assertProperty(t, &schema, "string", "q", "uuid", "Q")
	assertProperty(t, &schema, "string", "r", "uuid3", "R")
	assertProperty(t, &schema, "string", "s", "uuid4", "S")
	assertProperty(t, &schema, "string", "t", "uuid5", "T")
	assertProperty(t, &schema, "string", "u", "mac", "U")
}

func TestStringStructTag(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "JSONString")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	sch := models["jsonString"]
	assertProperty(t, &sch, "string", "someInt", "int64", "SomeInt")
	assertProperty(t, &sch, "string", "someInt8", "int8", "SomeInt8")
	assertProperty(t, &sch, "string", "someInt16", "int16", "SomeInt16")
	assertProperty(t, &sch, "string", "someInt32", "int32", "SomeInt32")
	assertProperty(t, &sch, "string", "someInt64", "int64", "SomeInt64")
	assertProperty(t, &sch, "string", "someUint", "uint64", "SomeUint")
	assertProperty(t, &sch, "string", "someUint8", "uint8", "SomeUint8")
	assertProperty(t, &sch, "string", "someUint16", "uint16", "SomeUint16")
	assertProperty(t, &sch, "string", "someUint32", "uint32", "SomeUint32")
	assertProperty(t, &sch, "string", "someUint64", "uint64", "SomeUint64")
	assertProperty(t, &sch, "string", "someFloat64", "double", "SomeFloat64")
	assertProperty(t, &sch, "string", "someString", "", "SomeString")
	assertProperty(t, &sch, "string", "someBool", "", "SomeBool")
	assertProperty(t, &sch, "string", "SomeDefaultInt", "int64", "")

	prop, ok := sch.Properties["somethingElse"]
	if assert.True(t, ok) {
		assert.NotEqual(t, "string", prop.Type)
	}
}

func TestPtrFieldStringStructTag(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "JSONPtrString")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	sch := models["jsonPtrString"]
	assertProperty(t, &sch, "string", "someInt", "int64", "SomeInt")
	assertProperty(t, &sch, "string", "someInt8", "int8", "SomeInt8")
	assertProperty(t, &sch, "string", "someInt16", "int16", "SomeInt16")
	assertProperty(t, &sch, "string", "someInt32", "int32", "SomeInt32")
	assertProperty(t, &sch, "string", "someInt64", "int64", "SomeInt64")
	assertProperty(t, &sch, "string", "someUint", "uint64", "SomeUint")
	assertProperty(t, &sch, "string", "someUint8", "uint8", "SomeUint8")
	assertProperty(t, &sch, "string", "someUint16", "uint16", "SomeUint16")
	assertProperty(t, &sch, "string", "someUint32", "uint32", "SomeUint32")
	assertProperty(t, &sch, "string", "someUint64", "uint64", "SomeUint64")
	assertProperty(t, &sch, "string", "someFloat64", "double", "SomeFloat64")
	assertProperty(t, &sch, "string", "someString", "", "SomeString")
	assertProperty(t, &sch, "string", "someBool", "", "SomeBool")

	prop, ok := sch.Properties["somethingElse"]
	if assert.True(t, ok) {
		assert.NotEqual(t, "string", prop.Type)
	}
}

func TestIgnoredStructField(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "IgnoredFields")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	sch := models["ignoredFields"]
	assertProperty(t, &sch, "string", "someIncludedField", "", "SomeIncludedField")
	assertProperty(t, &sch, "string", "someErroneouslyIncludedField", "", "SomeErroneouslyIncludedField")
	assert.Len(t, sch.Properties, 2)
}

func TestParseStructFields(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "SimpleComplexModel")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["SimpleComplexModel"]
	assertProperty(t, &schema, "object", "emb", "", "Emb")
	eSchema := schema.Properties["emb"]
	assertProperty(t, &eSchema, "integer", "cid", "int64", "CID")
	assertProperty(t, &eSchema, "string", "baz", "", "Baz")

	assertRef(t, &schema, "top", "Top", "#/definitions/Something")
	assertRef(t, &schema, "notSel", "NotSel", "#/definitions/NotSelected")
}

func TestParsePointerFields(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "Pointdexter")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["Pointdexter"]

	assertProperty(t, &schema, "integer", "id", "int64", "ID")
	assertProperty(t, &schema, "string", "name", "", "Name")
	assertProperty(t, &schema, "object", "emb", "", "Emb")
	assertProperty(t, &schema, "string", "t", "uuid5", "T")
	eSchema := schema.Properties["emb"]
	assertProperty(t, &eSchema, "integer", "cid", "int64", "CID")
	assertProperty(t, &eSchema, "string", "baz", "", "Baz")

	assertRef(t, &schema, "top", "Top", "#/definitions/Something")
	assertRef(t, &schema, "notSel", "NotSel", "#/definitions/NotSelected")
}

func TestEmbeddedStarExpr(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "EmbeddedStarExpr")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["EmbeddedStarExpr"]

	assertProperty(t, &schema, "integer", "embeddedMember", "int64", "EmbeddedMember")
	assertProperty(t, &schema, "integer", "notEmbedded", "int64", "NotEmbedded")
}

func TestArrayOfPointers(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "Cars")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["cars"]
	assertProperty(t, &schema, "array", "cars", "", "Cars")
}

func TestOverridingOneIgnore(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "OverridingOneIgnore")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["OverridingOneIgnore"]

	assertProperty(t, &schema, "integer", "id", "int64", "ID")
	assertProperty(t, &schema, "string", "name", "", "Name")
	assert.Len(t, schema.Properties, 2)
}

func TestParseSliceFields(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "SliceAndDice")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["SliceAndDice"]

	assertArrayProperty(t, &schema, "integer", "ids", "int64", "IDs")
	assertArrayProperty(t, &schema, "string", "names", "", "Names")
	assertArrayProperty(t, &schema, "string", "uuids", "uuid", "UUIDs")
	assertArrayProperty(t, &schema, "object", "embs", "", "Embs")
	eSchema := schema.Properties["embs"].Items.Schema
	assertArrayProperty(t, eSchema, "integer", "cid", "int64", "CID")
	assertArrayProperty(t, eSchema, "string", "baz", "", "Baz")

	assertArrayRef(t, &schema, "tops", "Tops", "#/definitions/Something")
	assertArrayRef(t, &schema, "notSels", "NotSels", "#/definitions/NotSelected")

	assertArrayProperty(t, &schema, "integer", "ptrIds", "int64", "PtrIDs")
	assertArrayProperty(t, &schema, "string", "ptrNames", "", "PtrNames")
	assertArrayProperty(t, &schema, "string", "ptrUuids", "uuid", "PtrUUIDs")
	assertArrayProperty(t, &schema, "object", "ptrEmbs", "", "PtrEmbs")
	eSchema = schema.Properties["ptrEmbs"].Items.Schema
	assertArrayProperty(t, eSchema, "integer", "ptrCid", "int64", "PtrCID")
	assertArrayProperty(t, eSchema, "string", "ptrBaz", "", "PtrBaz")

	assertArrayRef(t, &schema, "ptrTops", "PtrTops", "#/definitions/Something")
	assertArrayRef(t, &schema, "ptrNotSels", "PtrNotSels", "#/definitions/NotSelected")
}

func TestParseMapFields(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "MapTastic")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["MapTastic"]

	assertMapProperty(t, &schema, "integer", "ids", "int64", "IDs")
	assertMapProperty(t, &schema, "string", "names", "", "Names")
	assertMapProperty(t, &schema, "string", "uuids", "uuid", "UUIDs")
	assertMapProperty(t, &schema, "object", "embs", "", "Embs")
	eSchema := schema.Properties["embs"].AdditionalProperties.Schema
	assertMapProperty(t, eSchema, "integer", "cid", "int64", "CID")
	assertMapProperty(t, eSchema, "string", "baz", "", "Baz")

	assertMapRef(t, &schema, "tops", "Tops", "#/definitions/Something")
	assertMapRef(t, &schema, "notSels", "NotSels", "#/definitions/NotSelected")

	assertMapProperty(t, &schema, "integer", "ptrIds", "int64", "PtrIDs")
	assertMapProperty(t, &schema, "string", "ptrNames", "", "PtrNames")
	assertMapProperty(t, &schema, "string", "ptrUuids", "uuid", "PtrUUIDs")
	assertMapProperty(t, &schema, "object", "ptrEmbs", "", "PtrEmbs")
	eSchema = schema.Properties["ptrEmbs"].AdditionalProperties.Schema
	assertMapProperty(t, eSchema, "integer", "ptrCid", "int64", "PtrCID")
	assertMapProperty(t, eSchema, "string", "ptrBaz", "", "PtrBaz")

	assertMapRef(t, &schema, "ptrTops", "PtrTops", "#/definitions/Something")
	assertMapRef(t, &schema, "ptrNotSels", "PtrNotSels", "#/definitions/NotSelected")
}

func TestInterfaceField(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "Interfaced")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["Interfaced"]
	assertProperty(t, &schema, "", "custom_data", "", "CustomData")
}

func TestAliasedTypes(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "OtherTypes")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))

	schema := models["OtherTypes"]
	assertRef(t, &schema, "named", "Named", "#/definitions/SomeStringType")
	assertRef(t, &schema, "numbered", "Numbered", "#/definitions/SomeIntType")
	assertProperty(t, &schema, "string", "dated", "date-time", "Dated")
	assertRef(t, &schema, "timed", "Timed", "#/definitions/SomeTimedType")
	assertRef(t, &schema, "petted", "Petted", "#/definitions/SomePettedType")
	assertRef(t, &schema, "somethinged", "Somethinged", "#/definitions/SomethingType")
	assertRef(t, &schema, "strMap", "StrMap", "#/definitions/SomeStringMap")
	assertRef(t, &schema, "strArrMap", "StrArrMap", "#/definitions/SomeArrayStringMap")

	assertRef(t, &schema, "manyNamed", "ManyNamed", "#/definitions/SomeStringsType")
	assertRef(t, &schema, "manyNumbered", "ManyNumbered", "#/definitions/SomeIntsType")
	assertArrayProperty(t, &schema, "string", "manyDated", "date-time", "ManyDated")
	assertRef(t, &schema, "manyTimed", "ManyTimed", "#/definitions/SomeTimedsType")
	assertRef(t, &schema, "manyPetted", "ManyPetted", "#/definitions/SomePettedsType")
	assertRef(t, &schema, "manySomethinged", "ManySomethinged", "#/definitions/SomethingsType")

	assertArrayRef(t, &schema, "nameds", "Nameds", "#/definitions/SomeStringType")
	assertArrayRef(t, &schema, "numbereds", "Numbereds", "#/definitions/SomeIntType")
	assertArrayProperty(t, &schema, "string", "dateds", "date-time", "Dateds")
	assertArrayRef(t, &schema, "timeds", "Timeds", "#/definitions/SomeTimedType")
	assertArrayRef(t, &schema, "petteds", "Petteds", "#/definitions/SomePettedType")
	assertArrayRef(t, &schema, "somethingeds", "Somethingeds", "#/definitions/SomethingType")

	assertRef(t, &schema, "modsNamed", "ModsNamed", "#/definitions/modsSomeStringType")
	assertRef(t, &schema, "modsNumbered", "ModsNumbered", "#/definitions/modsSomeIntType")
	assertProperty(t, &schema, "string", "modsDated", "date-time", "ModsDated")
	assertRef(t, &schema, "modsTimed", "ModsTimed", "#/definitions/modsSomeTimedType")
	assertRef(t, &schema, "modsPetted", "ModsPetted", "#/definitions/modsSomePettedType")

	assertArrayRef(t, &schema, "modsNameds", "ModsNameds", "#/definitions/modsSomeStringType")
	assertArrayRef(t, &schema, "modsNumbereds", "ModsNumbereds", "#/definitions/modsSomeIntType")
	assertArrayProperty(t, &schema, "string", "modsDateds", "date-time", "ModsDateds")
	assertArrayRef(t, &schema, "modsTimeds", "ModsTimeds", "#/definitions/modsSomeTimedType")
	assertArrayRef(t, &schema, "modsPetteds", "ModsPetteds", "#/definitions/modsSomePettedType")

	assertRef(t, &schema, "manyModsNamed", "ManyModsNamed", "#/definitions/modsSomeStringsType")
	assertRef(t, &schema, "manyModsNumbered", "ManyModsNumbered", "#/definitions/modsSomeIntsType")
	assertArrayProperty(t, &schema, "string", "manyModsDated", "date-time", "ManyModsDated")
	assertRef(t, &schema, "manyModsTimed", "ManyModsTimed", "#/definitions/modsSomeTimedsType")
	assertRef(t, &schema, "manyModsPetted", "ManyModsPetted", "#/definitions/modsSomePettedsType")
	assertRef(t, &schema, "manyModsPettedPtr", "ManyModsPettedPtr", "#/definitions/modsSomePettedsPtrType")

	assertProperty(t, &schema, "string", "namedAlias", "", "NamedAlias")
	assertProperty(t, &schema, "integer", "numberedAlias", "int64", "NumberedAlias")
	assertArrayProperty(t, &schema, "string", "namedsAlias", "", "NamedsAlias")
	assertArrayProperty(t, &schema, "integer", "numberedsAlias", "int64", "NumberedsAlias")
}

func TestAliasedModels(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)

	names := []string{
		"SomeStringType",
		"SomeIntType",
		"SomeTimeType",
		"SomeTimedType",
		"SomePettedType",
		"SomethingType",
		"SomeStringsType",
		"SomeIntsType",
		"SomeTimesType",
		"SomeTimedsType",
		"SomePettedsType",
		"SomethingsType",
		"SomeObject",
		"SomeStringMap",
		"SomeIntMap",
		"SomeTimeMap",
		"SomeTimedMap",
		"SomePettedMap",
		"SomeSomethingMap",
	}

	defs := make(map[string]spec.Schema)
	for _, nm := range names {
		decl := getClassificationModel(sctx, nm)
		require.NotNil(t, decl)

		prs := &schemaBuilder{
			decl: decl,
			ctx:  sctx,
		}
		require.NoError(t, prs.Build(defs))
	}

	for k := range defs {
		for i, b := range names {
			if b == k {
				// remove the entry from the collection
				names = append(names[:i], names[i+1:]...)
			}
		}
	}
	if assert.Empty(t, names) {
		// single value types
		assertDefinition(t, defs, "SomeStringType", "string", "", "")
		assertDefinition(t, defs, "SomeIntType", "integer", "int64", "")
		assertDefinition(t, defs, "SomeTimeType", "string", "date-time", "")
		assertDefinition(t, defs, "SomeTimedType", "string", "date-time", "")
		assertRefDefinition(t, defs, "SomePettedType", "#/definitions/pet", "")
		assertRefDefinition(t, defs, "SomethingType", "#/definitions/Something", "")

		// slice types
		assertArrayDefinition(t, defs, "SomeStringsType", "string", "", "")
		assertArrayDefinition(t, defs, "SomeIntsType", "integer", "int64", "")
		assertArrayDefinition(t, defs, "SomeTimesType", "string", "date-time", "")
		assertArrayDefinition(t, defs, "SomeTimedsType", "string", "date-time", "")
		assertArrayWithRefDefinition(t, defs, "SomePettedsType", "#/definitions/pet", "")
		assertArrayWithRefDefinition(t, defs, "SomethingsType", "#/definitions/Something", "")

		// map types
		assertMapDefinition(t, defs, "SomeObject", "object", "", "")
		assertMapDefinition(t, defs, "SomeStringMap", "string", "", "")
		assertMapDefinition(t, defs, "SomeIntMap", "integer", "int64", "")
		assertMapDefinition(t, defs, "SomeTimeMap", "string", "date-time", "")
		assertMapDefinition(t, defs, "SomeTimedMap", "string", "date-time", "")
		assertMapWithRefDefinition(t, defs, "SomePettedMap", "#/definitions/pet", "")
		assertMapWithRefDefinition(t, defs, "SomeSomethingMap", "#/definitions/Something", "")
	}
}

func TestAliasedTopLevelModels(t *testing.T) {
	t.Run("with options: no scan models, with aliases as ref", func(t *testing.T) {
		t.Run("with goparsing/spec", func(t *testing.T) {
			sctx, err := newScanCtx(&Options{
				Packages: []string{
					"github.com/go-swagger/go-swagger/fixtures/goparsing/spec",
				},
				ScanModels: false,
				RefAliases: true,
			})
			require.NoError(t, err)

			t.Run("should find User definition in source", func(t *testing.T) {
				_, hasUser := sctx.FindDecl("github.com/go-swagger/go-swagger/fixtures/goparsing/spec", "User")
				require.True(t, hasUser)
			})

			var decl *entityDecl
			t.Run("should find Customer definition in source", func(t *testing.T) {
				var hasCustomer bool
				decl, hasCustomer = sctx.FindDecl("github.com/go-swagger/go-swagger/fixtures/goparsing/spec", "Customer")
				require.True(t, hasCustomer)
			})

			t.Run("with schema builder", func(t *testing.T) {
				require.NotNil(t, decl)
				builder := &schemaBuilder{
					ctx:  sctx,
					decl: decl,
				}

				t.Run("should build model for Customer", func(t *testing.T) {
					models := make(map[string]spec.Schema)
					require.NoError(t, builder.Build(models))

					assertRefDefinition(t, models, "Customer", "#/definitions/User", "")
				})

				t.Run("should have discovered models for User and Customer", func(t *testing.T) {
					require.Len(t, builder.postDecls, 2)
					foundUserIndex := -1
					foundCustomerIndex := -1

					for i, discoveredDecl := range builder.postDecls {
						switch discoveredDecl.Obj().Name() {
						case "User":
							foundUserIndex = i
						case "Customer":
							foundCustomerIndex = i
						}
					}
					require.GreaterOrEqual(t, foundUserIndex, 0)
					require.GreaterOrEqual(t, foundCustomerIndex, 0)

					userBuilder := &schemaBuilder{
						ctx:  sctx,
						decl: builder.postDecls[foundUserIndex],
					}

					t.Run("should build model for User", func(t *testing.T) {
						models := make(map[string]spec.Schema)
						require.NoError(t, userBuilder.Build(models))

						require.Contains(t, models, "User")

						user := models["User"]
						assert.True(t, user.Type.Contains("object"))

						userProperties := user.Properties
						require.Contains(t, userProperties, "name")
					})
				})
			})
		})
	})

	t.Run("with options: no scan models, without aliases as ref", func(t *testing.T) {
		t.Run("with goparsing/spec", func(t *testing.T) {
			sctx, err := newScanCtx(&Options{
				Packages: []string{
					"github.com/go-swagger/go-swagger/fixtures/goparsing/spec",
				},
				ScanModels: false,
				RefAliases: false,
			})
			require.NoError(t, err)

			t.Run("should find User definition in source", func(t *testing.T) {
				_, hasUser := sctx.FindDecl("github.com/go-swagger/go-swagger/fixtures/goparsing/spec", "User")
				require.True(t, hasUser)
			})

			var decl *entityDecl
			t.Run("should find Customer definition in source", func(t *testing.T) {
				var hasCustomer bool
				decl, hasCustomer = sctx.FindDecl("github.com/go-swagger/go-swagger/fixtures/goparsing/spec", "Customer")
				require.True(t, hasCustomer)
			})

			t.Run("with schema builder", func(t *testing.T) {
				require.NotNil(t, decl)
				builder := &schemaBuilder{
					ctx:  sctx,
					decl: decl,
				}

				t.Run("should build model for Customer", func(t *testing.T) {
					models := make(map[string]spec.Schema)
					require.NoError(t, builder.Build(models))

					require.Contains(t, models, "Customer")
					customer := models["Customer"]
					require.NotContains(t, models, "User")

					assert.True(t, customer.Type.Contains("object"))

					customerProperties := customer.Properties
					assert.Contains(t, customerProperties, "name")
					assert.NotEmpty(t, customer.Title)
				})

				t.Run("should have discovered only Customer", func(t *testing.T) {
					require.Len(t, builder.postDecls, 1)
					discovered := builder.postDecls[0]
					assert.Equal(t, "Customer", discovered.Obj().Name())
				})
			})
		})
	})
}

func TestAliasedSchemas(t *testing.T) {
	t.Setenv("SWAGGER_GENERATE_EXTENSION", "true")

	fixturesPath := filepath.Join("..", "fixtures", "goparsing", "go123", "aliased", "schema")
	var sp *spec.Swagger
	t.Run("end-to-end source scan should succeed", func(t *testing.T) {
		var err error
		sp, err = Run(&Options{
			WorkDir:    fixturesPath,
			BuildTags:  "testscanner", // fixture code is excluded from normal build
			ScanModels: true,
			RefAliases: true,
		})
		require.NoError(t, err)
	})

	if enableSpecOutput {
		// for debugging, output the resulting spec as YAML
		yml, err := marshalToYAMLFormat(sp)
		require.NoError(t, err)

		_, _ = os.Stdout.Write(yml)
	}

	shouldHaveExt := func(t *testing.T, sch spec.Schema, ext string) {
		t.Helper()
		pkg, hasExt := sch.Extensions.GetString(ext)
		assert.True(t, hasExt)
		assert.NotEmpty(t, pkg)
	}
	shouldHaveGoPackageExt := func(t *testing.T, sch spec.Schema) {
		t.Helper()
		shouldHaveExt(t, sch, "x-go-package")
	}
	shouldHaveTitle := func(t *testing.T, sch spec.Schema) {
		t.Helper()
		assert.NotEmpty(t, sch.Title)
	}
	shouldNotHaveTitle := func(t *testing.T, sch spec.Schema) {
		t.Helper()
		assert.Empty(t, sch.Title)
	}

	t.Run("type aliased to any should yield an empty schema", func(t *testing.T) {
		anything, ok := sp.Definitions["Anything"]
		require.True(t, ok)

		shouldHaveGoPackageExt(t, anything)
		shouldHaveTitle(t, anything)

		// after stripping extension and title, should be empty
		anything.VendorExtensible = spec.VendorExtensible{}
		anything.Title = ""
		assert.Equal(t, spec.Schema{}, anything)
	})

	t.Run("type aliased to an empty struct should yield an empty object", func(t *testing.T) {
		empty, ok := sp.Definitions["Empty"]
		require.True(t, ok)

		shouldHaveGoPackageExt(t, empty)
		shouldHaveTitle(t, empty)

		// after stripping extension and title, should be empty
		empty.VendorExtensible = spec.VendorExtensible{}
		empty.Title = ""
		emptyObject := &spec.Schema{}
		emptyObject = emptyObject.Typed("object", "").WithProperties(map[string]spec.Schema{})
		assert.Equal(t, *emptyObject, empty)
	})

	t.Run("struct fields defined as any or interface{} should yield properties with an empty schema", func(t *testing.T) {
		extended, ok := sp.Definitions["ExtendedID"]
		require.True(t, ok)

		t.Run("struct with an embedded alias should render as allOf", func(t *testing.T) {
			require.Len(t, extended.AllOf, 2)
			shouldHaveTitle(t, extended)

			foundAliased := false
			foundProps := false
			for idx, member := range extended.AllOf {
				isProps := len(member.Properties) > 0
				isAlias := member.Ref.String() != ""

				switch {
				case isProps:
					props := member
					t.Run("with property of type any", func(t *testing.T) {
						evenMore, ok := props.Properties["EvenMore"]
						require.True(t, ok)
						assert.Equal(t, spec.Schema{}, evenMore)
					})

					t.Run("with property of type interface{}", func(t *testing.T) {
						evenMore, ok := props.Properties["StillMore"]
						require.True(t, ok)
						assert.Equal(t, spec.Schema{}, evenMore)
					})

					t.Run("non-aliased properties remain unaffected", func(t *testing.T) {
						more, ok := props.Properties["more"]
						require.True(t, ok)

						shouldHaveExt(t, more, "x-go-name") // because we have a struct tag
						shouldNotHaveTitle(t, more)

						// after stripping extension and title, should be empty
						more.VendorExtensible = spec.VendorExtensible{}

						strSchema := &spec.Schema{}
						strSchema = strSchema.Typed("string", "")
						assert.Equal(t, *strSchema, more)
					})
					foundProps = true
				case isAlias:
					assertIsRef(t, &member, "#/definitions/Empty")
					foundAliased = true
				default:
					assert.Failf(t, "embedded members in struct are not as expected", "unexpected member in allOf: %d", idx)
				}
			}
			require.True(t, foundProps)
			require.True(t, foundAliased)
		})
	})

	t.Run("aliased primitive types remain unaffected", func(t *testing.T) {
		uuid, ok := sp.Definitions["UUID"]
		require.True(t, ok)

		shouldHaveGoPackageExt(t, uuid)
		shouldHaveTitle(t, uuid)

		// after strip extension, should be equal to integer with format
		uuid.VendorExtensible = spec.VendorExtensible{}
		uuid.Title = ""
		intSchema := &spec.Schema{}
		intSchema = intSchema.Typed("integer", "int64")
		assert.Equal(t, *intSchema, uuid)
	})

	t.Run("with struct having fields aliased to any or interface{}", func(t *testing.T) {
		order, ok := sp.Definitions["order"]
		require.True(t, ok)

		t.Run("field defined on an alias should produce a ref", func(t *testing.T) {
			t.Run("with alias to any", func(t *testing.T) {
				_, ok = order.Properties["DeliveryOption"]
				require.True(t, ok)
				assertRef(t, &order, "DeliveryOption", "", "#/definitions/Anything") // points to an alias to any
			})

			t.Run("with alias to primitive type", func(t *testing.T) {
				_, ok = order.Properties["id"]
				require.True(t, ok)
				assertRef(t, &order, "id", "", "#/definitions/UUID") // points to an alias to any
			})

			t.Run("with alias to struct type", func(t *testing.T) {
				_, ok = order.Properties["extended_id"]
				require.True(t, ok)
				assertRef(t, &order, "extended_id", "", "#/definitions/ExtendedID") // points to an alias to any
			})

			t.Run("inside anonymous array", func(t *testing.T) {
				items, ok := order.Properties["items"]
				require.True(t, ok)

				require.NotNil(t, items)
				require.NotNil(t, items.Items)

				assert.True(t, items.Type.Contains("array"))
				t.Run("field as any should render as empty object", func(t *testing.T) {
					require.NotNil(t, items.Items.Schema)
					itemsSchema := items.Items.Schema
					assert.True(t, itemsSchema.Type.Contains("object"))

					require.Contains(t, itemsSchema.Properties, "extra_options")
					extraOptions := itemsSchema.Properties["extra_options"]
					shouldHaveExt(t, extraOptions, "x-go-name")

					extraOptions.VendorExtensible = spec.VendorExtensible{}
					empty := spec.Schema{}
					assert.Equal(t, empty, extraOptions)
				})
			})
		})

		t.Run("struct field defined as any should produce an empty schema", func(t *testing.T) {
			extras, ok := order.Properties["Extras"]
			require.True(t, ok)
			assert.Equal(t, spec.Schema{}, extras)
		})

		t.Run("struct field defined as interface{} should produce an empty schema", func(t *testing.T) {
			extras, ok := order.Properties["MoreExtras"]
			require.True(t, ok)
			assert.Equal(t, spec.Schema{}, extras)
		})
	})

	t.Run("type redefinitions and syntactic aliases to any should render the same", func(t *testing.T) {
		whatnot, ok := sp.Definitions["whatnot"]
		require.True(t, ok)
		// after strip extension, should be empty
		whatnot.VendorExtensible = spec.VendorExtensible{}
		assert.Equal(t, spec.Schema{}, whatnot)

		whatnotAlias, ok := sp.Definitions["whatnot_alias"]
		require.True(t, ok)
		// after strip extension, should be empty
		whatnotAlias.VendorExtensible = spec.VendorExtensible{}
		assert.Equal(t, spec.Schema{}, whatnotAlias)

		whatnot2, ok := sp.Definitions["whatnot2"]
		require.True(t, ok)
		// after strip extension, should be empty
		whatnot2.VendorExtensible = spec.VendorExtensible{}
		assert.Equal(t, spec.Schema{}, whatnot2)

		whatnot2Alias, ok := sp.Definitions["whatnot2_alias"]
		require.True(t, ok)
		// after strip extension, should be empty
		whatnot2Alias.VendorExtensible = spec.VendorExtensible{}
		assert.Equal(t, spec.Schema{}, whatnot2Alias)
	})

	t.Run("alias to another alias is resolved as a ref", func(t *testing.T) {
		void, ok := sp.Definitions["void"]
		require.True(t, ok)

		assertIsRef(t, &void, "#/definitions/Empty") // points to another alias
	})

	t.Run("type redefinition to anonymous is not an alias and is resolved as an object", func(t *testing.T) {
		empty, ok := sp.Definitions["empty_redefinition"]
		require.True(t, ok)

		shouldHaveGoPackageExt(t, empty)
		shouldNotHaveTitle(t, empty)

		// after stripping extension and title, should be empty
		empty.VendorExtensible = spec.VendorExtensible{}
		emptyObject := &spec.Schema{}
		emptyObject = emptyObject.Typed("object", "").WithProperties(map[string]spec.Schema{})
		assert.Equal(t, *emptyObject, empty)
	})

	t.Run("alias to a named interface should render as a $ref", func(t *testing.T) {
		iface, ok := sp.Definitions["iface_alias"]
		require.True(t, ok)

		assertIsRef(t, &iface, "#/definitions/iface") // points to an interface
	})

	t.Run("interface redefinition is not an alias and should render as a $ref", func(t *testing.T) {
		iface, ok := sp.Definitions["iface_redefinition"]
		require.True(t, ok)

		assertIsRef(t, &iface, "#/definitions/iface") // points to an interface
	})

	t.Run("anonymous interface should render a schema", func(t *testing.T) {
		iface, ok := sp.Definitions["anonymous_iface"]
		require.True(t, ok)

		require.NotEmpty(t, iface.Properties)
		require.Contains(t, iface.Properties, "String")
	})

	t.Run("anonymous struct should render as an anonymous schema", func(t *testing.T) {
		obj, ok := sp.Definitions["anonymous_struct"]
		require.True(t, ok)

		require.NotEmpty(t, obj.Properties)
		require.Contains(t, obj.Properties, "A")

		a := obj.Properties["A"]
		assert.True(t, a.Type.Contains("object"))
		require.Contains(t, a.Properties, "B")
		b := a.Properties["B"]
		assert.True(t, b.Type.Contains("integer"))
	})

	t.Run("standalone model with a tag should be rendered", func(t *testing.T) {
		shouldSee, ok := sp.Definitions["ShouldSee"]
		require.True(t, ok)
		assert.True(t, shouldSee.Type.Contains("boolean"))
	})

	t.Run("standalone model without a tag should not be rendered", func(t *testing.T) {
		_, ok := sp.Definitions["ShouldNotSee"]
		require.False(t, ok)

		_, ok = sp.Definitions["ShouldNotSeeSlice"]
		require.False(t, ok)

		_, ok = sp.Definitions["ShouldNotSeeMap"]
		require.False(t, ok)
	})

	t.Run("with aliases in slices and arrays", func(t *testing.T) {
		t.Run("slice redefinition should render as schema", func(t *testing.T) {
			t.Run("with anonymous slice", func(t *testing.T) {
				slice, ok := sp.Definitions["slice_type"] // []any
				require.True(t, ok)
				assert.True(t, slice.Type.Contains("array"))
				require.NotNil(t, slice.Items)
				require.NotNil(t, slice.Items.Schema)

				assert.Equal(t, &spec.Schema{}, slice.Items.Schema)
			})

			t.Run("with anonymous struct", func(t *testing.T) {
				slice, ok := sp.Definitions["slice_of_structs"] // type X = []struct{}
				require.True(t, ok)
				assert.True(t, slice.Type.Contains("array"))

				require.NotNil(t, slice.Items)
				require.NotNil(t, slice.Items.Schema)

				emptyObject := &spec.Schema{}
				emptyObject = emptyObject.Typed("object", "").WithProperties(map[string]spec.Schema{})
				assert.Equal(t, emptyObject, slice.Items.Schema)
			})
		})

		t.Run("alias to anonymous slice should render as schema", func(t *testing.T) {
			t.Run("with anonymous slice", func(t *testing.T) {
				slice, ok := sp.Definitions["slice_alias"] // type X = []any
				require.True(t, ok)
				assert.True(t, slice.Type.Contains("array"))

				require.NotNil(t, slice.Items)
				require.NotNil(t, slice.Items.Schema)

				assert.Equal(t, &spec.Schema{}, slice.Items.Schema)
			})

			t.Run("with anonymous struct", func(t *testing.T) {
				slice, ok := sp.Definitions["slice_of_structs_alias"] // type X = []struct{}
				require.True(t, ok)
				assert.True(t, slice.Type.Contains("array"))
				require.NotNil(t, slice.Items)
				require.NotNil(t, slice.Items.Schema)

				emptyObject := &spec.Schema{}
				emptyObject = emptyObject.Typed("object", "").WithProperties(map[string]spec.Schema{})
				assert.Equal(t, emptyObject, slice.Items.Schema)
			})
		})

		t.Run("alias to named alias to anonymous slice should render as ref", func(t *testing.T) {
			slice, ok := sp.Definitions["slice_to_slice"] // type X = Slice
			require.True(t, ok)
			assertIsRef(t, &slice, "#/definitions/slice_type") // points to a named alias
		})
	})

	t.Run("with aliases in interfaces", func(t *testing.T) {
		t.Run("should render anonymous interface as a schema", func(t *testing.T) {
			iface, ok := sp.Definitions["anonymous_iface"] // e.g. type X interface{ String() string}
			require.True(t, ok)

			require.True(t, iface.Type.Contains("object"))
			require.Contains(t, iface.Properties, "String")
			prop := iface.Properties["String"]
			require.True(t, prop.Type.Contains("string"))
			assert.Len(t, iface.Properties, 1)
		})

		t.Run("alias to an anonymous interface should render as a $ref", func(t *testing.T) {
			iface, ok := sp.Definitions["anonymous_iface_alias"]
			require.True(t, ok)

			assertIsRef(t, &iface, "#/definitions/anonymous_iface") // points to an anonymous interface
		})

		t.Run("named interface should render as a schema", func(t *testing.T) {
			iface, ok := sp.Definitions["iface"]
			require.True(t, ok)

			require.True(t, iface.Type.Contains("object"))
			require.Contains(t, iface.Properties, "Get")
			prop := iface.Properties["Get"]
			require.True(t, prop.Type.Contains("string"))
			assert.Len(t, iface.Properties, 1)
		})

		t.Run("named interface with embedded types should render as allOf", func(t *testing.T) {
			iface, ok := sp.Definitions["iface_embedded"]
			require.True(t, ok)

			require.Len(t, iface.AllOf, 2)
			foundEmbedded := false
			foundMethod := false
			for idx, member := range iface.AllOf {
				require.True(t, member.Type.Contains("object"))
				require.NotEmpty(t, member.Properties)
				require.Len(t, member.Properties, 1)
				propGet, isEmbedded := member.Properties["Get"]
				propMethod, isMethod := member.Properties["Dump"]

				switch {
				case isEmbedded:
					assert.True(t, propGet.Type.Contains("string"))
					foundEmbedded = true
				case isMethod:
					assert.True(t, propMethod.Type.Contains("array"))
					foundMethod = true
				default:
					assert.Failf(t, "embedded members in interface are not as expected", "unexpected member in allOf: %d", idx)
				}
			}
			require.True(t, foundEmbedded)
			require.True(t, foundMethod)
		})

		t.Run("named interface with embedded anonymous interface should render as allOf", func(t *testing.T) {
			iface, ok := sp.Definitions["iface_embedded_anonymous"]
			require.True(t, ok)

			require.Len(t, iface.AllOf, 2)
			foundEmbedded := false
			foundAnonymous := false
			for idx, member := range iface.AllOf {
				require.True(t, member.Type.Contains("object"))
				require.NotEmpty(t, member.Properties)
				require.Len(t, member.Properties, 1)
				propGet, isEmbedded := member.Properties["String"]
				propAnonymous, isAnonymous := member.Properties["Error"]

				switch {
				case isEmbedded:
					assert.True(t, propGet.Type.Contains("string"))
					foundEmbedded = true
				case isAnonymous:
					assert.True(t, propAnonymous.Type.Contains("string"))
					foundAnonymous = true
				default:
					assert.Failf(t, "embedded members in interface are not as expected", "unexpected member in allOf: %d", idx)
				}
			}
			require.True(t, foundEmbedded)
			require.True(t, foundAnonymous)
		})

		t.Run("composition of empty interfaces is rendered as an empty schema", func(t *testing.T) {
			iface, ok := sp.Definitions["iface_embedded_empty"]
			require.True(t, ok)

			iface.VendorExtensible = spec.VendorExtensible{}
			assert.Equal(t, spec.Schema{}, iface)
		})

		t.Run("interface embedded with an alias should be rendered as allOf, with a ref", func(t *testing.T) {
			iface, ok := sp.Definitions["iface_embedded_with_alias"]
			require.True(t, ok)

			require.Len(t, iface.AllOf, 3)
			foundEmbedded := false
			foundEmbeddedAnon := false
			foundRef := false
			for idx, member := range iface.AllOf {
				propGet, isEmbedded := member.Properties["String"]
				propAnonymous, isAnonymous := member.Properties["Dump"]
				isRef := member.Ref.String() != ""

				switch {
				case isEmbedded:
					require.True(t, member.Type.Contains("object"))
					require.Len(t, member.Properties, 1)
					assert.True(t, propGet.Type.Contains("string"))
					foundEmbedded = true
				case isAnonymous:
					require.True(t, member.Type.Contains("object"))
					require.Len(t, member.Properties, 1)
					assert.True(t, propAnonymous.Type.Contains("array"))
					foundEmbeddedAnon = true
				case isRef:
					require.Empty(t, member.Properties)
					assertIsRef(t, &member, "#/definitions/iface_alias")
					foundRef = true
				default:
					assert.Failf(t, "embedded members in interface are not as expected", "unexpected member in allOf: %d", idx)
				}
			}
			require.True(t, foundEmbedded)
			require.True(t, foundEmbeddedAnon)
			require.True(t, foundRef)
		})
	})

	t.Run("with aliases in embedded types", func(t *testing.T) {
		t.Run("embedded alias should render as a $ref", func(t *testing.T) {
			iface, ok := sp.Definitions["embedded_with_alias"]
			require.True(t, ok)

			require.Len(t, iface.AllOf, 3)
			foundAnything := false
			foundUUID := false
			foundProps := false
			for idx, member := range iface.AllOf {
				isProps := len(member.Properties) > 0
				isRef := member.Ref.String() != ""

				switch {
				case isProps:
					require.True(t, member.Type.Contains("object"))
					require.Len(t, member.Properties, 3)
					assert.Contains(t, member.Properties, "EvenMore")
					foundProps = true
				case isRef:
					switch member.Ref.String() {
					case "#/definitions/Anything":
						foundAnything = true
					case "#/definitions/UUID":
						foundUUID = true
					default:
						assert.Failf(t,
							"embedded members in interface are not as expected", "unexpected $ref for member (%v): %d",
							member.Ref, idx,
						)
					}
				default:
					assert.Failf(t, "embedded members in interface are not as expected", "unexpected member in allOf: %d", idx)
				}
			}
			require.True(t, foundAnything)
			require.True(t, foundUUID)
			require.True(t, foundProps)
		})
	})
}

func TestSpecialSchemas(t *testing.T) {
	t.Setenv("SWAGGER_GENERATE_EXTENSION", "true")

	fixturesPath := filepath.Join("..", "fixtures", "goparsing", "go123", "special")
	var sp *spec.Swagger

	t.Run("end-to-end source scan should succeed", func(t *testing.T) {
		var err error
		sp, err = Run(&Options{
			WorkDir:    fixturesPath,
			BuildTags:  "testscanner", // fixture code is excluded from normal build
			ScanModels: true,
			RefAliases: true,
		})
		require.NoError(t, err)
	})

	if enableSpecOutput {
		// for debugging, output the resulting spec as YAML
		yml, err := marshalToYAMLFormat(sp)
		require.NoError(t, err)

		_, _ = os.Stdout.Write(yml)
	}

	t.Run("top-level primitive declaration should render just fine", func(t *testing.T) {
		primitive, ok := sp.Definitions["primitive"]
		require.True(t, ok)

		require.True(t, primitive.Type.Contains("string"))
	})

	t.Run("alias to unsafe pointer at top level should render empty", func(t *testing.T) {
		uptr, ok := sp.Definitions["unsafe_pointer_alias"]
		require.True(t, ok)
		var empty spec.Schema
		uptr.VendorExtensible = spec.VendorExtensible{}
		require.Equal(t, empty, uptr)
	})

	t.Run("alias to uintptr at top level should render as integer", func(t *testing.T) {
		uptr, ok := sp.Definitions["upointer_alias"]
		require.True(t, ok)
		require.True(t, uptr.Type.Contains("integer"))
		require.Equal(t, "uint64", uptr.Format)
	})

	t.Run("top-level map[string]... should render just fine", func(t *testing.T) {
		gomap, ok := sp.Definitions["go_map"]
		require.True(t, ok)
		require.True(t, gomap.Type.Contains("object"))
		require.NotNil(t, gomap.AdditionalProperties)

		mapSchema := gomap.AdditionalProperties.Schema
		require.NotNil(t, mapSchema)
		require.True(t, mapSchema.Type.Contains("integer"))
		require.Equal(t, "uint16", mapSchema.Format)
	})

	t.Run("untagged struct referenced by a tagged model should be discovered", func(t *testing.T) {
		gostruct, ok := sp.Definitions["GoStruct"]
		require.True(t, ok)
		require.True(t, gostruct.Type.Contains("object"))
		require.NotEmpty(t, gostruct.Properties)

		t.Run("pointer property should render just fine", func(t *testing.T) {
			a, ok := gostruct.Properties["A"]
			require.True(t, ok)
			require.True(t, a.Type.Contains("number"))
			require.Equal(t, "float", a.Format)
		})
	})

	t.Run("tagged unsupported map type should render empty", func(t *testing.T) {
		idx, ok := sp.Definitions["index_map"]
		require.True(t, ok)
		var empty spec.Schema
		idx.VendorExtensible = spec.VendorExtensible{}
		require.Equal(t, empty, idx)
	})

	t.Run("redefinition of the builtin error type should render as a string", func(t *testing.T) {
		goerror, ok := sp.Definitions["go_error"]
		require.True(t, ok)
		require.True(t, goerror.Type.Contains("string"))

		t.Run("a type based on the error builtin should be decorated with a x-go-type: error extension", func(t *testing.T) {
			val, hasExt := goerror.Extensions.GetString("x-go-type")
			assert.True(t, hasExt)
			assert.Equal(t, "error", val)
		})
	})

	t.Run("with SpecialTypes struct", func(t *testing.T) {
		t.Run("in spite of all the pitfalls, the struct should be rendered", func(t *testing.T) {
			special, ok := sp.Definitions["special_types"]
			require.True(t, ok)
			require.True(t, special.Type.Contains("object"))
			props := special.Properties
			require.NotEmpty(t, props)
			require.Empty(t, special.AllOf)

			t.Run("property pointer to struct should render as a ref", func(t *testing.T) {
				ptr, ok := props["PtrStruct"]
				require.True(t, ok)
				assertIsRef(t, &ptr, "#/definitions/GoStruct")
			})

			t.Run("property as time.Time should render as a formatted string", func(t *testing.T) {
				str, ok := props["ShouldBeStringTime"]
				require.True(t, ok)
				require.True(t, str.Type.Contains("string"))
				require.Equal(t, "date-time", str.Format)
			})

			t.Run("property as *time.Time should also render as a formatted string", func(t *testing.T) {
				str, ok := props["ShouldAlsoBeStringTime"]
				require.True(t, ok)
				require.True(t, str.Type.Contains("string"))
				require.Equal(t, "date-time", str.Format)
			})

			t.Run("property as builtin error should render as a string", func(t *testing.T) {
				goerror, ok := props["Err"]
				require.True(t, ok)
				require.True(t, goerror.Type.Contains("string"))

				t.Run("a type based on the error builtin should be decorated with a x-go-type: error extension", func(t *testing.T) {
					val, hasExt := goerror.Extensions.GetString("x-go-type")
					assert.True(t, hasExt)
					assert.Equal(t, "error", val)
				})
			})

			t.Run("type recognized as a text marshaler should render as a string", func(t *testing.T) {
				m, ok := props["Marshaler"]
				require.True(t, ok)
				require.True(t, m.Type.Contains("string"))

				t.Run("a type based on the encoding.TextMarshaler decorated with a x-go-type extension", func(t *testing.T) {
					val, hasExt := m.Extensions.GetString("x-go-type")
					assert.True(t, hasExt)
					assert.Equal(t, "github.com/go-swagger/go-swagger/fixtures/goparsing/go123/special.IsATextMarshaler", val)
				})
			})

			t.Run("a json.RawMessage should be recognized and render as an object (yes this is wrong)", func(t *testing.T) {
				m, ok := props["Message"]
				require.True(t, ok)
				require.True(t, m.Type.Contains("object"))
			})

			t.Run("type time.Duration is not recognized as a special type and should just render as a ref", func(t *testing.T) {
				d, ok := props["Duration"]
				require.True(t, ok)
				assertIsRef(t, &d, "#/definitions/Duration")

				t.Run("discovered definition should be an integer", func(t *testing.T) {
					duration, ok := sp.Definitions["Duration"]
					require.True(t, ok)
					require.True(t, duration.Type.Contains("integer"))
					require.Equal(t, "int64", duration.Format)

					t.Run("time.Duration schema should be decorated with a x-go-package: time", func(t *testing.T) {
						val, hasExt := duration.Extensions.GetString("x-go-package")
						assert.True(t, hasExt)
						assert.Equal(t, "time", val)
					})
				})
			})

			t.Run("with strfmt types", func(t *testing.T) {
				t.Run("a strfmt.Date should be recognized and render as a formatted string", func(t *testing.T) {
					d, ok := props["FormatDate"]
					require.True(t, ok)
					require.True(t, d.Type.Contains("string"))
					require.Equal(t, "date", d.Format)
				})

				t.Run("a strfmt.DateTime should be recognized and render as a formatted string", func(t *testing.T) {
					d, ok := props["FormatTime"]
					require.True(t, ok)
					require.True(t, d.Type.Contains("string"))
					require.Equal(t, "date-time", d.Format)
				})

				t.Run("a strfmt.UUID should be recognized and render as a formatted string", func(t *testing.T) {
					u, ok := props["FormatUUID"]
					require.True(t, ok)
					require.True(t, u.Type.Contains("string"))
					require.Equal(t, "uuid", u.Format)
				})

				t.Run("a pointer to strfmt.UUID should be recognized and render as a formatted string", func(t *testing.T) {
					u, ok := props["PtrFormatUUID"]
					require.True(t, ok)
					require.True(t, u.Type.Contains("string"))
					require.Equal(t, "uuid", u.Format)
				})
			})

			t.Run("a property which is a map should render just fine, with a ref", func(t *testing.T) {
				mm, ok := props["Map"]
				require.True(t, ok)
				require.True(t, mm.Type.Contains("object"))
				require.NotNil(t, mm.AdditionalProperties)
				mapSchema := mm.AdditionalProperties.Schema
				require.NotNil(t, mapSchema)
				assertIsRef(t, mapSchema, "#/definitions/GoStruct")
			})

			t.Run(`with the "WhatNot" anonymous inner struct`, func(t *testing.T) {
				t.Run("should render as an anonymous schema, in spite of all the unsupported things", func(t *testing.T) {
					wn, ok := props["WhatNot"]
					require.True(t, ok)
					require.True(t, wn.Type.Contains("object"))
					require.NotEmpty(t, wn.Properties)

					markedProps := make([]string, 0)

					for _, unsupportedProp := range []string{
						"AA", // complex128
						"A",  // complex64
						"B",  // chan int
						"C",  // func()
						"D",  // func() string
						"E",  // unsafe.Pointer
					} {
						t.Run("with property "+unsupportedProp, func(t *testing.T) {
							prop, ok := wn.Properties[unsupportedProp]
							require.True(t, ok)
							markedProps = append(markedProps, unsupportedProp)

							t.Run("unsupported type in property should render as an empty schema", func(t *testing.T) {
								var empty spec.Schema
								require.Equal(t, empty, prop)
							})
						})
					}

					for _, supportedProp := range []string{
						"F", // uintptr
						"G",
						"H",
						"I",
						"J",
						"K",
					} {
						t.Run("with property "+supportedProp, func(t *testing.T) {
							prop, ok := wn.Properties[supportedProp]
							require.True(t, ok)
							markedProps = append(markedProps, supportedProp)

							switch supportedProp {
							case "F":
								t.Run("uintptr should render as integer", func(t *testing.T) {
									require.True(t, prop.Type.Contains("integer"))
									require.Equal(t, "uint64", prop.Format)
								})
							case "G", "H":
								t.Run(
									"math/big types are not recognized as special types and as TextMarshalers they render as string",
									func(t *testing.T) {
										require.True(t, prop.Type.Contains("string"))
									})
							case "I":
								t.Run("go array should render as a json array", func(t *testing.T) {
									require.True(t, prop.Type.Contains("array"))
									require.NotNil(t, prop.Items)
									itemsSchema := prop.Items.Schema
									require.NotNil(t, itemsSchema)

									require.True(t, itemsSchema.Type.Contains("integer"))
									// [5]byte is not recognized an array of bytes, but of uint8
									// (internally this is the same for go)
									require.Equal(t, "uint8", itemsSchema.Format)
								})
							case "J", "K":
								t.Run("reflect types should render just fine", func(t *testing.T) {
									var dest string
									if supportedProp == "J" {
										dest = "Type"
									} else {
										dest = "Value"
									}
									assertIsRef(t, &prop, "#/definitions/"+dest)

									t.Run("the $ref should exist", func(t *testing.T) {
										deref, ok := sp.Definitions[dest]
										require.True(t, ok)
										val, hasExt := deref.Extensions.GetString("x-go-package")
										assert.True(t, hasExt)
										assert.Equal(t, "reflect", val)
									})
								})
							}
						})
					}

					t.Run("we should not have any property left in WhatNot", func(t *testing.T) {
						for _, key := range markedProps {
							delete(wn.Properties, key)
						}

						require.Empty(t, wn.Properties)
					})

					t.Run("surprisingly, a tagged unexported top-level definition can be rendered", func(t *testing.T) {
						unexported, ok := sp.Definitions["unexported"]
						require.True(t, ok)
						require.True(t, unexported.Type.Contains("object"))
					})

					t.Run("the IsATextMarshaler type is not identified as a discovered type and is not rendered", func(t *testing.T) {
						_, ok := sp.Definitions["IsATextMarshaler"]
						require.False(t, ok)
					})

					t.Run("a top-level go array should render just fine", func(t *testing.T) {
						// Notice that the semantics of fixed length are lost in this mapping
						goarray, ok := sp.Definitions["go_array"]
						require.True(t, ok)
						require.True(t, goarray.Type.Contains("array"))
						require.NotNil(t, goarray.Items)
						itemsSchema := goarray.Items.Schema
						require.NotNil(t, itemsSchema)
						require.True(t, itemsSchema.Type.Contains("integer"))
						require.Equal(t, "int64", itemsSchema.Format)
					})
				})
			})
		})
	})

	t.Run("with generic types", func(t *testing.T) {
		// NOTE: codescan does not really support generic types.
		// This test just makes sure generic definitions don't crash the scanner.
		//
		// The general approach of the scanner is to make an empty schema out of anything
		// it doesn't understand.

		// generic_constraint
		t.Run("generic type constraint should render like an interface", func(t *testing.T) {
			generic, ok := sp.Definitions["generic_constraint"]
			require.True(t, ok)
			require.Len(t, generic.AllOf, 1) // scanner only understood one member, and skipped the ~uint16 member is doesn't understand
			member := generic.AllOf[0]
			require.True(t, member.Type.Contains("object"))
			require.Len(t, member.Properties, 1)
			prop, ok := member.Properties["Uint"]
			require.True(t, ok)
			require.True(t, prop.Type.Contains("integer"))
			require.Equal(t, "uint16", prop.Format)
		})

		// numerical_constraint
		t.Run("generic type constraint with union type should render an empty schema", func(t *testing.T) {
			generic, ok := sp.Definitions["numerical_constraint"]
			require.True(t, ok)
			var empty spec.Schema
			generic.VendorExtensible = spec.VendorExtensible{}
			require.Equal(t, empty, generic)
		})

		// generic_map
		t.Run("generic map should render an empty schema", func(t *testing.T) {
			generic, ok := sp.Definitions["generic_map"]
			require.True(t, ok)
			var empty spec.Schema
			generic.VendorExtensible = spec.VendorExtensible{}
			require.Equal(t, empty, generic)
		})

		// generic_map_alias
		t.Run("generic map alias to an anonymous generic type should render an empty schema", func(t *testing.T) {
			generic, ok := sp.Definitions["generic_map_alias"]
			require.True(t, ok)
			var empty spec.Schema
			generic.VendorExtensible = spec.VendorExtensible{}
			require.Equal(t, empty, generic)
		})

		// generic_indirect
		t.Run("generic map alias to a named generic type should render a ref", func(t *testing.T) {
			generic, ok := sp.Definitions["generic_indirect"]
			require.True(t, ok)
			assertIsRef(t, &generic, "#/definitions/generic_map_alias")
		})

		// generic_slice
		t.Run("generic slice should render as an array of empty schemas", func(t *testing.T) {
			generic, ok := sp.Definitions["generic_slice"]
			require.True(t, ok)
			require.True(t, generic.Type.Contains("array"))
			require.NotNil(t, generic.Items)
			itemsSchema := generic.Items.Schema
			require.NotNil(t, itemsSchema)
			var empty spec.Schema
			require.Equal(t, &empty, itemsSchema)
		})

		// union_alias:
		t.Run("alias to type constraint should render a ref", func(t *testing.T) {
			generic, ok := sp.Definitions["union_alias"]
			require.True(t, ok)
			assertIsRef(t, &generic, "#/definitions/numerical_constraint")
		})
	})
}

func TestEmbeddedAllOf(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "AllOfModel")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))
	schema := models["AllOfModel"]

	require.Len(t, schema.AllOf, 3)
	asch := schema.AllOf[0]
	assertProperty(t, &asch, "integer", "age", "int32", "Age")
	assertProperty(t, &asch, "integer", "id", "int64", "ID")
	assertProperty(t, &asch, "string", "name", "", "Name")

	asch = schema.AllOf[1]
	assert.Equal(t, "#/definitions/withNotes", asch.Ref.String())

	asch = schema.AllOf[2]
	assertProperty(t, &asch, "string", "createdAt", "date-time", "CreatedAt")
	assertProperty(t, &asch, "integer", "did", "int64", "DID")
	assertProperty(t, &asch, "string", "cat", "", "Cat")
}

func TestPointersAreNullableByDefaultWhenSetXNullableForPointersIsSet(t *testing.T) {
	assertModel := func(sctx *scanCtx, packagePath, modelName string) {
		decl, _ := sctx.FindDecl(packagePath, modelName)
		require.NotNil(t, decl)
		prs := &schemaBuilder{
			ctx:  sctx,
			decl: decl,
		}
		models := make(map[string]spec.Schema)
		require.NoError(t, prs.Build(models))

		schema := models[modelName]
		require.Len(t, schema.Properties, 5)

		require.Contains(t, schema.Properties, "Value1")
		assert.Equal(t, true, schema.Properties["Value1"].Extensions["x-nullable"])
		require.Contains(t, schema.Properties, "Value2")
		assert.NotContains(t, schema.Properties["Value2"].Extensions, "x-nullable")
		require.Contains(t, schema.Properties, "Value3")
		assert.Equal(t, false, schema.Properties["Value3"].Extensions["x-nullable"])
		require.Contains(t, schema.Properties, "Value4")
		assert.NotContains(t, schema.Properties["Value4"].Extensions, "x-nullable")
		assert.Equal(t, false, schema.Properties["Value4"].Extensions["x-isnullable"])
		require.Contains(t, schema.Properties, "Value5")
		assert.NotContains(t, schema.Properties["Value5"].Extensions, "x-nullable")
	}

	packagePath := "github.com/go-swagger/go-swagger/fixtures/enhancements/pointers-nullable-by-default"
	sctx, err := newScanCtx(&Options{Packages: []string{packagePath}, SetXNullableForPointers: true})
	require.NoError(t, err)

	assertModel(sctx, packagePath, "Item")
	assertModel(sctx, packagePath, "ItemInterface")
}

func TestPointersAreNotNullableByDefaultWhenSetXNullableForPointersIsNotSet(t *testing.T) {
	assertModel := func(sctx *scanCtx, packagePath, modelName string) {
		decl, _ := sctx.FindDecl(packagePath, modelName)
		require.NotNil(t, decl)
		prs := &schemaBuilder{
			ctx:  sctx,
			decl: decl,
		}
		models := make(map[string]spec.Schema)
		require.NoError(t, prs.Build(models))

		schema := models[modelName]
		require.Len(t, schema.Properties, 5)

		require.Contains(t, schema.Properties, "Value1")
		assert.NotContains(t, schema.Properties["Value1"].Extensions, "x-nullable")
		require.Contains(t, schema.Properties, "Value2")
		assert.NotContains(t, schema.Properties["Value2"].Extensions, "x-nullable")
		require.Contains(t, schema.Properties, "Value3")
		assert.Equal(t, false, schema.Properties["Value3"].Extensions["x-nullable"])
		require.Contains(t, schema.Properties, "Value4")
		assert.NotContains(t, schema.Properties["Value4"].Extensions, "x-nullable")
		assert.Equal(t, false, schema.Properties["Value4"].Extensions["x-isnullable"])
		require.Contains(t, schema.Properties, "Value5")
		assert.NotContains(t, schema.Properties["Value5"].Extensions, "x-nullable")
	}

	packagePath := "github.com/go-swagger/go-swagger/fixtures/enhancements/pointers-nullable-by-default"
	sctx, err := newScanCtx(&Options{Packages: []string{packagePath}})
	require.NoError(t, err)

	assertModel(sctx, packagePath, "Item")
	assertModel(sctx, packagePath, "ItemInterface")
}

func TestSwaggerTypeNamed(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "NamedWithType")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))
	schema := models["namedWithType"]

	assertProperty(t, &schema, "object", "some_map", "", "SomeMap")
}

func TestSwaggerTypeNamedWithGenerics(t *testing.T) {
	tests := map[string]func(t *testing.T, models map[string]spec.Schema){
		"NamedStringResults": func(t *testing.T, models map[string]spec.Schema) {
			schema := models["namedStringResults"]
			assertArrayProperty(t, &schema, "string", "matches", "", "Matches")
		},
		"NamedStoreOrderResults": func(t *testing.T, models map[string]spec.Schema) {
			schema := models["namedStoreOrderResults"]
			assertArrayRef(t, &schema, "matches", "Matches", "#/definitions/order")
		},
		"NamedStringSlice": func(t *testing.T, models map[string]spec.Schema) {
			assertArrayDefinition(t, models, "namedStringSlice", "string", "", "NamedStringSlice")
		},
		"NamedStoreOrderSlice": func(t *testing.T, models map[string]spec.Schema) {
			assertArrayWithRefDefinition(t, models, "namedStoreOrderSlice", "#/definitions/order", "NamedStoreOrderSlice")
		},
		"NamedStringMap": func(t *testing.T, models map[string]spec.Schema) {
			assertMapDefinition(t, models, "namedStringMap", "string", "", "NamedStringMap")
		},
		"NamedStoreOrderMap": func(t *testing.T, models map[string]spec.Schema) {
			assertMapWithRefDefinition(t, models, "namedStoreOrderMap", "#/definitions/order", "NamedStoreOrderMap")
		},
		"NamedMapOfStoreOrderSlices": func(t *testing.T, models map[string]spec.Schema) {
			assertMapDefinition(t, models, "namedMapOfStoreOrderSlices", "array", "", "NamedMapOfStoreOrderSlices")
			arraySchema := models["namedMapOfStoreOrderSlices"].AdditionalProperties.Schema
			assertArrayWithRefDefinition(t, map[string]spec.Schema{
				"array": *arraySchema,
			}, "array", "#/definitions/order", "")
		},
	}

	for testName, testFunc := range tests {
		t.Run(testName, func(t *testing.T) {
			sctx := loadClassificationPkgsCtx(t)
			decl := getClassificationModel(sctx, testName)
			require.NotNil(t, decl)
			prs := &schemaBuilder{
				ctx:  sctx,
				decl: decl,
			}
			models := make(map[string]spec.Schema)
			require.NoError(t, prs.Build(models))
			testFunc(t, models)
		})
	}
}

func TestSwaggerTypeStruct(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	decl := getClassificationModel(sctx, "NullString")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))
	schema := models["NullString"]

	assert.True(t, schema.Type.Contains("string"))
}

func TestStructDiscriminators(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)

	models := make(map[string]spec.Schema)
	for _, tn := range []string{"BaseStruct", "Giraffe", "Gazelle"} {
		decl := getClassificationModel(sctx, tn)
		require.NotNil(t, decl)
		prs := &schemaBuilder{
			ctx:  sctx,
			decl: decl,
		}
		require.NoError(t, prs.Build(models))
	}

	schema := models["animal"]

	assert.Equal(t, "BaseStruct", schema.Extensions["x-go-name"])
	assert.Equal(t, "jsonClass", schema.Discriminator)

	sch := models["gazelle"]
	assert.Len(t, sch.AllOf, 2)
	cl, _ := sch.Extensions.GetString("x-class")
	assert.Equal(t, "a.b.c.d.E", cl)
	cl, _ = sch.Extensions.GetString("x-go-name")
	assert.Equal(t, "Gazelle", cl)

	sch = models["giraffe"]
	assert.Len(t, sch.AllOf, 2)
	cl, _ = sch.Extensions.GetString("x-class")
	assert.Empty(t, cl)
	cl, _ = sch.Extensions.GetString("x-go-name")
	assert.Equal(t, "Giraffe", cl)

	// sch = noModelDefs["lion"]

	// b, _ := json.MarshalIndent(sch, "", "  ")
	// fmt.Println(string(b))
}

func TestInterfaceDiscriminators(t *testing.T) {
	sctx := loadClassificationPkgsCtx(t)
	models := make(map[string]spec.Schema)
	for _, tn := range []string{"BaseStruct", "Identifiable", "WaterType", "Fish", "TeslaCar", "ModelS", "ModelX", "ModelA", "Cars"} {
		decl := getClassificationModel(sctx, tn)
		require.NotNil(t, decl)

		prs := &schemaBuilder{
			ctx:  sctx,
			decl: decl,
		}
		require.NoError(t, prs.Build(models))
	}

	schema, ok := models["fish"]

	if assert.True(t, ok) && assert.Len(t, schema.AllOf, 5) {
		sch := schema.AllOf[3]
		assert.Len(t, sch.Properties, 1)
		assertProperty(t, &sch, "string", "colorName", "", "ColorName")

		sch = schema.AllOf[2]
		assert.Equal(t, "#/definitions/extra", sch.Ref.String())

		sch = schema.AllOf[0]
		assert.Len(t, sch.Properties, 1)
		assertProperty(t, &sch, "integer", "id", "int64", "ID")

		sch = schema.AllOf[1]
		assert.Equal(t, "#/definitions/water", sch.Ref.String())

		sch = schema.AllOf[4]
		assert.Len(t, sch.Properties, 2)
		assertProperty(t, &sch, "string", "name", "", "Name")
		assertProperty(t, &sch, "string", "jsonClass", "", "StructType")
		assert.Equal(t, "jsonClass", sch.Discriminator)
	}

	schema, ok = models["modelS"]
	if assert.True(t, ok) {
		assert.Len(t, schema.AllOf, 2)
		cl, _ := schema.Extensions.GetString("x-class")
		assert.Equal(t, "com.tesla.models.ModelS", cl)
		cl, _ = schema.Extensions.GetString("x-go-name")
		assert.Equal(t, "ModelS", cl)

		sch := schema.AllOf[0]
		assert.Equal(t, "#/definitions/TeslaCar", sch.Ref.String())
		sch = schema.AllOf[1]
		assert.Len(t, sch.Properties, 1)
		assertProperty(t, &sch, "string", "edition", "", "Edition")
	}

	schema, ok = models["modelA"]
	if assert.True(t, ok) {
		cl, _ := schema.Extensions.GetString("x-go-name")
		assert.Equal(t, "ModelA", cl)

		sch, ok := schema.Properties["Tesla"]
		if assert.True(t, ok) {
			assert.Equal(t, "#/definitions/TeslaCar", sch.Ref.String())
		}

		assertProperty(t, &schema, "integer", "doors", "int64", "Doors")
	}
}

func TestAddExtension(t *testing.T) {
	ve := &spec.VendorExtensible{
		Extensions: make(spec.Extensions),
	}

	key := "x-go-name"
	value := "Name"
	addExtension(ve, key, value)
	assert.Equal(t, value, ve.Extensions[key].(string))

	key2 := "x-go-package"
	value2 := "schema"
	t.Setenv("SWAGGER_GENERATE_EXTENSION", "true")
	addExtension(ve, key2, value2)
	assert.Equal(t, value2, ve.Extensions[key2].(string))

	key3 := "x-go-class"
	value3 := "Spec"
	t.Setenv("SWAGGER_GENERATE_EXTENSION", "false")
	addExtension(ve, key3, value3)
	assert.Nil(t, ve.Extensions[key3])
}

func getClassificationModel(sctx *scanCtx, nm string) *entityDecl {
	decl, ok := sctx.FindDecl("github.com/go-swagger/go-swagger/fixtures/goparsing/classification/models", nm)
	if !ok {
		return nil
	}
	return decl
}

func assertArrayProperty(t testing.TB, schema *spec.Schema, typeName, jsonName, format, goName string) {
	prop := schema.Properties[jsonName]
	assert.NotEmpty(t, prop.Type)
	assert.True(t, prop.Type.Contains("array"))
	assert.NotNil(t, prop.Items)
	if typeName != "" {
		assert.Equal(t, typeName, prop.Items.Schema.Type[0])
	}
	assert.Equal(t, goName, prop.Extensions["x-go-name"])
	assert.Equal(t, format, prop.Items.Schema.Format)
}

func assertArrayRef(t testing.TB, schema *spec.Schema, jsonName, goName, fragment string) {
	assertArrayProperty(t, schema, "", jsonName, "", goName)
	psch := schema.Properties[jsonName].Items.Schema
	assert.Equal(t, fragment, psch.Ref.String())
}

func assertProperty(t testing.TB, schema *spec.Schema, typeName, jsonName, format, goName string) {
	if typeName == "" {
		assert.Empty(t, schema.Properties[jsonName].Type)
	} else if assert.NotEmpty(t, schema.Properties[jsonName].Type) {
		assert.Equal(t, typeName, schema.Properties[jsonName].Type[0])
	}
	if goName == "" {
		assert.Nil(t, schema.Properties[jsonName].Extensions["x-go-name"])
	} else {
		assert.Equal(t, goName, schema.Properties[jsonName].Extensions["x-go-name"])
	}
	assert.Equal(t, format, schema.Properties[jsonName].Format)
}

func assertRef(t testing.TB, schema *spec.Schema, jsonName, _, fragment string) {
	t.Helper()

	assert.Empty(t, schema.Properties[jsonName].Type)
	psch := schema.Properties[jsonName]
	assert.Equal(t, fragment, psch.Ref.String())
}

func assertIsRef(t testing.TB, schema *spec.Schema, fragment string) {
	t.Helper()

	assert.Equal(t, fragment, schema.Ref.String())
}

func assertDefinition(t testing.TB, defs map[string]spec.Schema, defName, typeName, formatName, goName string) {
	t.Helper()

	schema, ok := defs[defName]
	if assert.True(t, ok) {
		if assert.NotEmpty(t, schema.Type) {
			assert.Equal(t, typeName, schema.Type[0])
			if goName != "" {
				assert.Equal(t, goName, schema.Extensions["x-go-name"])
			} else {
				assert.Nil(t, schema.Extensions["x-go-name"])
			}
			assert.Equal(t, formatName, schema.Format)
		}
	}
}

func assertMapDefinition(t testing.TB, defs map[string]spec.Schema, defName, typeName, formatName, goName string) {
	schema, ok := defs[defName]
	if assert.True(t, ok) {
		if assert.NotEmpty(t, schema.Type) {
			assert.Equal(t, "object", schema.Type[0])
			adl := schema.AdditionalProperties
			if assert.NotNil(t, adl) && assert.NotNil(t, adl.Schema) {
				if len(adl.Schema.Type) > 0 {
					assert.Equal(t, typeName, adl.Schema.Type[0])
				}
				assert.Equal(t, formatName, adl.Schema.Format)
			}
			if goName != "" {
				assert.Equal(t, goName, schema.Extensions["x-go-name"])
			} else {
				assert.Nil(t, schema.Extensions["x-go-name"])
			}
		}
	}
}

func assertMapWithRefDefinition(t testing.TB, defs map[string]spec.Schema, defName, refURL, goName string) {
	schema, ok := defs[defName]
	if assert.True(t, ok) {
		if assert.NotEmpty(t, schema.Type) {
			assert.Equal(t, "object", schema.Type[0])
			adl := schema.AdditionalProperties
			if assert.NotNil(t, adl) && assert.NotNil(t, adl.Schema) {
				if assert.NotZero(t, adl.Schema.Ref) {
					assert.Equal(t, refURL, adl.Schema.Ref.String())
				}
			}
			if goName != "" {
				assert.Equal(t, goName, schema.Extensions["x-go-name"])
			} else {
				assert.Nil(t, schema.Extensions["x-go-name"])
			}
		}
	}
}

func assertArrayDefinition(t testing.TB, defs map[string]spec.Schema, defName, typeName, formatName, goName string) {
	schema, ok := defs[defName]
	if assert.True(t, ok) {
		if assert.NotEmpty(t, schema.Type) {
			assert.Equal(t, "array", schema.Type[0])
			adl := schema.Items
			if assert.NotNil(t, adl) && assert.NotNil(t, adl.Schema) {
				assert.Equal(t, typeName, adl.Schema.Type[0])
				assert.Equal(t, formatName, adl.Schema.Format)
			}
			if goName != "" {
				assert.Equal(t, goName, schema.Extensions["x-go-name"])
			} else {
				assert.Nil(t, schema.Extensions["x-go-name"])
			}
		}
	}
}

func assertArrayWithRefDefinition(t testing.TB, defs map[string]spec.Schema, defName, refURL, goName string) {
	schema, ok := defs[defName]
	if assert.True(t, ok) {
		if assert.NotEmpty(t, schema.Type) {
			assert.Equal(t, "array", schema.Type[0])
			adl := schema.Items
			if assert.NotNil(t, adl) && assert.NotNil(t, adl.Schema) {
				if assert.NotZero(t, adl.Schema.Ref) {
					assert.Equal(t, refURL, adl.Schema.Ref.String())
				}
			}
			if goName != "" {
				assert.Equal(t, goName, schema.Extensions["x-go-name"])
			} else {
				assert.Nil(t, schema.Extensions["x-go-name"])
			}
		}
	}
}

func assertRefDefinition(t testing.TB, defs map[string]spec.Schema, defName, refURL, goName string) {
	schema, ok := defs[defName]
	if assert.True(t, ok) {
		if assert.NotZero(t, schema.Ref) {
			url := schema.Ref.String()
			assert.Equal(t, refURL, url)
			if goName != "" {
				assert.Equal(t, goName, schema.Extensions["x-go-name"])
			} else {
				assert.Nil(t, schema.Extensions["x-go-name"])
			}
		}
	}
}

func assertMapProperty(t testing.TB, schema *spec.Schema, typeName, jsonName, format, goName string) {
	prop := schema.Properties[jsonName]
	assert.NotEmpty(t, prop.Type)
	assert.True(t, prop.Type.Contains("object"))
	assert.NotNil(t, prop.AdditionalProperties)
	if typeName != "" {
		assert.Equal(t, typeName, prop.AdditionalProperties.Schema.Type[0])
	}
	assert.Equal(t, goName, prop.Extensions["x-go-name"])
	assert.Equal(t, format, prop.AdditionalProperties.Schema.Format)
}

func assertMapRef(t testing.TB, schema *spec.Schema, jsonName, goName, fragment string) {
	assertMapProperty(t, schema, "", jsonName, "", goName)
	psch := schema.Properties[jsonName].AdditionalProperties.Schema
	assert.Equal(t, fragment, psch.Ref.String())
}

func marshalToYAMLFormat(swspec any) ([]byte, error) {
	b, err := json.Marshal(swspec)
	if err != nil {
		return nil, err
	}

	var jsonObj any
	if err := yaml.Unmarshal(b, &jsonObj); err != nil {
		return nil, err
	}

	return yaml.Marshal(jsonObj)
}

func TestEmbeddedDescriptionAndTags(t *testing.T) {
	packagePath := "github.com/go-swagger/go-swagger/fixtures/bugs/3125/minimal"
	sctx, err := newScanCtx(&Options{
		Packages:    []string{packagePath},
		DescWithRef: true,
	})
	require.NoError(t, err)
	decl, _ := sctx.FindDecl(packagePath, "Item")
	require.NotNil(t, decl)
	prs := &schemaBuilder{
		ctx:  sctx,
		decl: decl,
	}
	models := make(map[string]spec.Schema)
	require.NoError(t, prs.Build(models))
	schema := models["Item"]

	assert.Equal(t, []string{"value1", "value2"}, schema.Required)
	require.Len(t, schema.Properties, 2)

	require.Contains(t, schema.Properties, "value1")
	assert.Equal(t, "Nullable value", schema.Properties["value1"].Description)
	assert.Equal(t, true, schema.Properties["value1"].Extensions["x-nullable"])

	require.Contains(t, schema.Properties, "value2")
	assert.Equal(t, "Non-nullable value", schema.Properties["value2"].Description)
	assert.NotContains(t, schema.Properties["value2"].Extensions, "x-nullable")
	assert.Equal(t, `{"value": 42}`, schema.Properties["value2"].Example)
}

func TestIssue2540(t *testing.T) {
	t.Run("should produce example and default for top level declaration only",
		testIssue2540(false, `{
		"Book": {
      "description": "At this moment, a book is only described by its publishing date\nand author.",
      "type": "object",
      "title": "Book holds all relevant information about a book.",
			"example": "{ \"Published\": 2026, \"Author\": \"Fred\" }",
      "default": "{ \"Published\": 1900, \"Author\": \"Unknown\" }",
      "properties": {
        "Author": {
          "$ref": "#/definitions/Author"
        },
        "Published": {
          "type": "integer",
          "format": "int64",
          "minimum": 0,
          "example": 2021
        }
      }
    }
  }`),
	)
	t.Run("should produce example and default for top level declaration and embedded $ref field",
		testIssue2540(true, `{
		"Book": {
      "description": "At this moment, a book is only described by its publishing date\nand author.",
      "type": "object",
      "title": "Book holds all relevant information about a book.",
			"example": "{ \"Published\": 2026, \"Author\": \"Fred\" }",
      "default": "{ \"Published\": 1900, \"Author\": \"Unknown\" }",
      "properties": {
        "Author": {
          "$ref": "#/definitions/Author",
          "example": "{ \"Name\": \"Tolkien\" }"
        },
        "Published": {
          "type": "integer",
          "format": "int64",
          "minimum": 0,
          "example": 2021
        }
      }
    }
  }`),
	)
}

func testIssue2540(descWithRef bool, expectedJSON string) func(*testing.T) {
	return func(t *testing.T) {
		t.Setenv("SWAGGER_GENERATE_EXTENSION", "false")
		packagePath := "github.com/go-swagger/go-swagger/fixtures/bugs/2540/foo"
		sctx, err := newScanCtx(&Options{
			Packages:    []string{packagePath},
			DescWithRef: descWithRef,
		})
		require.NoError(t, err)

		decl, _ := sctx.FindDecl(packagePath, "Book")
		require.NotNil(t, decl)
		prs := &schemaBuilder{
			ctx:  sctx,
			decl: decl,
		}

		models := make(map[string]spec.Schema)
		require.NoError(t, prs.Build(models))

		b, err := json.Marshal(models)
		require.NoError(t, err)
		assert.JSONEq(t, expectedJSON, string(b))
	}
}

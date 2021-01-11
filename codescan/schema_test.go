package codescan

import (
	"os"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	assert.EqualValues(t, 1000, *prop.Maximum)
	assert.True(t, prop.ExclusiveMaximum, "'id' should have had an exclusive maximum")
	assert.NotNil(t, prop.Minimum)
	assert.EqualValues(t, 10, *prop.Minimum)
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
	assert.EqualValues(t, 45, *prop.Maximum)
	assert.False(t, prop.ExclusiveMaximum, "'score' should not have had an exclusive maximum")
	assert.NotNil(t, prop.Minimum)
	assert.EqualValues(t, 3, *prop.Minimum)
	assert.False(t, prop.ExclusiveMinimum, "'score' should not have had an exclusive minimum")
	assert.Equal(t, 27, prop.Example)

	expectedNameExtensions := spec.Extensions{
		"x-go-name": "Name",
		"x-property-array": []interface{}{
			"value1",
			"value2",
		},
		"x-property-array-obj": []interface{}{
			map[string]interface{}{
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
	assert.EqualValues(t, 4, *prop.MinLength)
	assert.EqualValues(t, 50, *prop.MaxLength)
	assert.Equal(t, "[A-Za-z0-9-.]*", prop.Pattern)
	assert.EqualValues(t, expectedNameExtensions, prop.Extensions)

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
	assert.NotNil(t, prop.Items, "foo_slice should have had an items property")
	assert.NotNil(t, prop.Items.Schema, "foo_slice.items should have had a schema property")
	assert.True(t, prop.UniqueItems, "'foo_slice' should have unique items")
	assert.EqualValues(t, 3, *prop.MinItems, "'foo_slice' should have had 3 min items")
	assert.EqualValues(t, 10, *prop.MaxItems, "'foo_slice' should have had 10 max items")
	itprop := prop.Items.Schema
	assert.EqualValues(t, 3, *itprop.MinLength, "'foo_slice.items.minLength' should have been 3")
	assert.EqualValues(t, 10, *itprop.MaxLength, "'foo_slice.items.maxLength' should have been 10")
	assert.EqualValues(t, "\\w+", itprop.Pattern, "'foo_slice.items.pattern' should have \\w+")

	assertArrayProperty(t, &schema, "string", "time_slice", "date-time", "TimeSlice")
	prop, ok = schema.Properties["time_slice"]
	assert.Equal(t, "a TimeSlice is a slice of times", prop.Description)
	assert.True(t, ok, "should have a 'time_slice' property")
	assert.NotNil(t, prop.Items, "time_slice should have had an items property")
	assert.NotNil(t, prop.Items.Schema, "time_slice.items should have had a schema property")
	assert.True(t, prop.UniqueItems, "'time_slice' should have unique items")
	assert.EqualValues(t, 3, *prop.MinItems, "'time_slice' should have had 3 min items")
	assert.EqualValues(t, 10, *prop.MaxItems, "'time_slice' should have had 10 max items")

	assertArrayProperty(t, &schema, "array", "bar_slice", "", "BarSlice")
	prop, ok = schema.Properties["bar_slice"]
	assert.Equal(t, "a BarSlice has bars which are strings", prop.Description)
	assert.True(t, ok, "should have a 'bar_slice' property")
	assert.NotNil(t, prop.Items, "bar_slice should have had an items property")
	assert.NotNil(t, prop.Items.Schema, "bar_slice.items should have had a schema property")
	assert.True(t, prop.UniqueItems, "'bar_slice' should have unique items")
	assert.EqualValues(t, 3, *prop.MinItems, "'bar_slice' should have had 3 min items")
	assert.EqualValues(t, 10, *prop.MaxItems, "'bar_slice' should have had 10 max items")
	itprop = prop.Items.Schema
	if assert.NotNil(t, itprop) {
		assert.EqualValues(t, 4, *itprop.MinItems, "'bar_slice.items.minItems' should have been 4")
		assert.EqualValues(t, 9, *itprop.MaxItems, "'bar_slice.items.maxItems' should have been 9")
		itprop2 := itprop.Items.Schema
		if assert.NotNil(t, itprop2) {
			assert.EqualValues(t, 5, *itprop2.MinItems, "'bar_slice.items.items.minItems' should have been 5")
			assert.EqualValues(t, 8, *itprop2.MaxItems, "'bar_slice.items.items.maxItems' should have been 8")
			itprop3 := itprop2.Items.Schema
			if assert.NotNil(t, itprop3) {
				assert.EqualValues(t, 3, *itprop3.MinLength, "'bar_slice.items.items.items.minLength' should have been 3")
				assert.EqualValues(t, 10, *itprop3.MaxLength, "'bar_slice.items.items.items.maxLength' should have been 10")
				assert.EqualValues(t, "\\w+", itprop3.Pattern, "'bar_slice.items.items.items.pattern' should have \\w+")
			}
		}
	}

	assertArrayProperty(t, &schema, "array", "deep_time_slice", "", "DeepTimeSlice")
	prop, ok = schema.Properties["deep_time_slice"]
	assert.Equal(t, "a DeepSlice has bars which are time", prop.Description)
	assert.True(t, ok, "should have a 'deep_time_slice' property")
	assert.NotNil(t, prop.Items, "deep_time_slice should have had an items property")
	assert.NotNil(t, prop.Items.Schema, "deep_time_slice.items should have had a schema property")
	assert.True(t, prop.UniqueItems, "'deep_time_slice' should have unique items")
	assert.EqualValues(t, 3, *prop.MinItems, "'deep_time_slice' should have had 3 min items")
	assert.EqualValues(t, 10, *prop.MaxItems, "'deep_time_slice' should have had 10 max items")
	itprop = prop.Items.Schema
	if assert.NotNil(t, itprop) {
		assert.EqualValues(t, 4, *itprop.MinItems, "'deep_time_slice.items.minItems' should have been 4")
		assert.EqualValues(t, 9, *itprop.MaxItems, "'deep_time_slice.items.maxItems' should have been 9")
		itprop2 := itprop.Items.Schema
		if assert.NotNil(t, itprop2) {
			assert.EqualValues(t, 5, *itprop2.MinItems, "'deep_time_slice.items.items.minItems' should have been 5")
			assert.EqualValues(t, 8, *itprop2.MaxItems, "'deep_time_slice.items.items.maxItems' should have been 8")
			itprop3 := itprop2.Items.Schema
			assert.NotNil(t, itprop3)
		}
	}

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
	assert.EqualValues(t, 1000, *iprop.Maximum)
	assert.True(t, iprop.ExclusiveMaximum, "'id' should have had an exclusive maximum")
	assert.NotNil(t, iprop.Minimum)
	assert.EqualValues(t, 10, *iprop.Minimum)
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
	assert.EqualValues(t, 1, *iprop.Minimum)
	assert.EqualValues(t, 10, *iprop.Maximum)

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
	assertProperty(t, &schema, "object", "custom_data", "", "CustomData")
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
	assert.Equal(t, schema.Discriminator, "jsonClass")

	sch := models["gazelle"]
	assert.Len(t, sch.AllOf, 2)
	cl, _ := sch.Extensions.GetString("x-class")
	assert.Equal(t, "a.b.c.d.E", cl)
	cl, _ = sch.Extensions.GetString("x-go-name")
	assert.Equal(t, "Gazelle", cl)

	sch = models["giraffe"]
	assert.Len(t, sch.AllOf, 2)
	cl, _ = sch.Extensions.GetString("x-class")
	assert.Equal(t, "", cl)
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
		sch := schema.AllOf[0]
		assert.Len(t, sch.Properties, 1)
		assertProperty(t, &sch, "string", "colorName", "", "ColorName")

		sch = schema.AllOf[1]
		assert.Equal(t, "#/definitions/extra", sch.Ref.String())

		sch = schema.AllOf[2]
		assert.Len(t, sch.Properties, 1)
		assertProperty(t, &sch, "integer", "id", "int64", "ID")

		sch = schema.AllOf[3]
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
	_ = os.Setenv("SWAGGER_GENERATE_EXTENSION", "true")
	addExtension(ve, key2, value2)
	assert.Equal(t, value2, ve.Extensions[key2].(string))

	key3 := "x-go-class"
	value3 := "Spec"
	_ = os.Setenv("SWAGGER_GENERATE_EXTENSION", "false")
	addExtension(ve, key3, value3)
	assert.Equal(t, nil, ve.Extensions[key3])
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
		assert.Equal(t, nil, schema.Properties[jsonName].Extensions["x-go-name"])
	} else {
		assert.Equal(t, goName, schema.Properties[jsonName].Extensions["x-go-name"])
	}
	assert.Equal(t, format, schema.Properties[jsonName].Format)
}

func assertRef(t testing.TB, schema *spec.Schema, jsonName, _, fragment string) {
	assert.Empty(t, schema.Properties[jsonName].Type)
	psch := schema.Properties[jsonName]
	assert.Equal(t, fragment, psch.Ref.String())
}

func assertDefinition(t testing.TB, defs map[string]spec.Schema, defName, typeName, formatName, goName string) {
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

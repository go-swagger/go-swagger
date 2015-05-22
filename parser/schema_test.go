package parser

import (
	goparser "go/parser"
	"testing"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func TestSchemaParser(t *testing.T) {

	prog := classifierProgram(t)
	docFile := "../fixtures/goparsing/classification/models/nomodel.go"
	fileTree, err := goparser.ParseFile(prog.Fset, docFile, nil, goparser.ParseComments)
	if err != nil {
		t.FailNow()
	}

	sp := schemaParser(prog)
	definitions := make(map[string]spec.Schema)
	sp.Parse(fileTree, definitions)
	schema := definitions["NoModel"]

	assert.Equal(t, spec.StringOrArray([]string{"object"}), schema.Type)
	assert.Equal(t, "NoModel is a struct that exists in a package\nbut is not annotated with the swagger model annotations\nso it should now show up in a test", schema.Title)
	//assert.Equal(t, "this model is not explictly mentioned in the import paths\nbut because it it transitively required by the order\nit should also be collected.", schema.Description)
}

func TestParsePrimitiveSchemaProperty(t *testing.T) {
	prog := classifierProgram(t)
	docFile := "../fixtures/goparsing/classification/models/nomodel.go"
	fileTree, err := goparser.ParseFile(prog.Fset, docFile, nil, goparser.ParseComments)
	if err != nil {
		t.FailNow()
	}
	sp := schemaParser(prog)
	definitions := make(map[string]spec.Schema)
	err = sp.Parse(fileTree, definitions)
	assert.NoError(t, err)
	schema := definitions["PrimateModel"]
	assertProperty(t, &schema, "boolean", "a", "", "A")
	assertProperty(t, &schema, "string", "b", "", "B")
	assertProperty(t, &schema, "string", "c", "", "C")
	assertProperty(t, &schema, "number", "d", "int64", "D")
	assertProperty(t, &schema, "number", "e", "int8", "E")
	assertProperty(t, &schema, "number", "f", "int16", "F")
	assertProperty(t, &schema, "number", "g", "int32", "G")
	assertProperty(t, &schema, "number", "h", "int64", "H")
	assertProperty(t, &schema, "number", "i", "uint64", "I")
	assertProperty(t, &schema, "number", "j", "uint8", "J")
	assertProperty(t, &schema, "number", "k", "uint16", "K")
	assertProperty(t, &schema, "number", "l", "uint32", "L")
	assertProperty(t, &schema, "number", "m", "uint64", "M")
	assertProperty(t, &schema, "number", "n", "float", "N")
	assertProperty(t, &schema, "number", "o", "double", "O")

	//pretty.Println(schema)
}

func TestParseStringFormatSchemaProperty(t *testing.T) {
	prog := classifierProgram(t)
	docFile := "../fixtures/goparsing/classification/models/nomodel.go"
	fileTree, err := goparser.ParseFile(prog.Fset, docFile, nil, goparser.ParseComments)
	if err != nil {
		t.FailNow()
	}
	sp := schemaParser(prog)
	definitions := make(map[string]spec.Schema)
	err = sp.Parse(fileTree, definitions)
	assert.NoError(t, err)
	schema := definitions["FormattedModel"]
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

	//pretty.Println(schema)
}

func assertProperty(t *testing.T, schema *spec.Schema, typeName, jsonName, format, goName string) {
	assert.NotEmpty(t, schema.Properties[jsonName].Type)
	assert.Equal(t, typeName, schema.Properties[jsonName].Type[0])
	assert.Equal(t, goName, schema.Properties[jsonName].Extensions["x-go-name"])
	assert.Equal(t, format, schema.Properties[jsonName].Format)
}

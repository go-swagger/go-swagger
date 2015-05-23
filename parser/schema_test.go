package parser

import (
	"fmt"
	goparser "go/parser"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"golang.org/x/tools/go/loader"

	"github.com/casualjim/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

var classificationProg *loader.Program
var noModelDefs map[string]spec.Schema

func init() {
	classificationProg = classifierProgram()
	docFile := "../fixtures/goparsing/classification/models/nomodel.go"
	fileTree, err := goparser.ParseFile(classificationProg.Fset, docFile, nil, goparser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	sp := schemaParser(classificationProg)
	noModelDefs = make(map[string]spec.Schema)
	err = sp.Parse(fileTree, noModelDefs)
	if err != nil {
		log.Fatal(err)
	}
}

func TestSchemaValueExtractors(t *testing.T) {
	strfmts := []string{
		"// +swagger:strfmt ",
		"* +swagger:strfmt ",
		"* +swagger:strfmt ",
		" +swagger:strfmt ",
		"+swagger:strfmt ",
		"// +swagger:strfmt    ",
		"* +swagger:strfmt     ",
		"* +swagger:strfmt    ",
		" +swagger:strfmt     ",
		"+swagger:strfmt      ",
	}
	models := []string{
		"// +swagger:model ",
		"* +swagger:model ",
		"* +swagger:model ",
		" +swagger:model ",
		"+swagger:model ",
		"// +swagger:model    ",
		"* +swagger:model     ",
		"* +swagger:model    ",
		" +swagger:model     ",
		"+swagger:model      ",
	}
	validParams := []string{
		"yada123",
		"date",
		"date-time",
		"long-combo-1-with-combo-2-and-a-3rd-one-too",
	}
	invalidParams := []string{
		"1-yada-3",
		"1-2-3",
		"-yada-3",
		"-2-3",
		"*blah",
		"blah*",
	}

	verifySwaggerOneArgSwaggerTag(t, rxStrFmt, strfmts, validParams, append(invalidParams, "", "  ", " "))
	verifySwaggerOneArgSwaggerTag(t, rxModelOverride, models, append(validParams, "", "  ", " "), invalidParams)

	verifyMinMax(t, rxMinimum, "min", []string{"", ">", "="})
	verifyMinMax(t, rxMaximum, "max", []string{"", "<", "="})
	verifyNumeric2Words(t, rxMultipleOf, "multiple", "of")

	verifyIntegerMinMaxManyWords(t, rxMinLength, "min", []string{"len", "length"})
	verifyIntegerMinMaxManyWords(t, rxMaxLength, "max", []string{"len", "length"})
	// pattern
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	patArgs := []string{"^\\w+$", "[A-Za-z0-9-.]*"}
	patNames := []string{"pattern", "Pattern"}
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, nm := range patNames {
				for _, es2 := range extraSpaces {
					for _, es3 := range extraSpaces {
						for _, arg := range patArgs {
							line := strings.Join([]string{pref, es1, nm, es2, ":", es3, arg}, "")
							matches := rxPattern.FindStringSubmatch(line)
							assert.Len(t, matches, 2)
							assert.Equal(t, arg, matches[1])
						}
					}
				}
			}
		}
	}

	verifyIntegerMinMaxManyWords(t, rxMinItems, "min", []string{"items"})
	verifyIntegerMinMaxManyWords(t, rxMaxItems, "max", []string{"items"})
	verifyBoolean(t, rxUnique, []string{"unique"}, nil)

	verifyBoolean(t, rxReadOnly, []string{"read"}, []string{"only"})
	verifyBoolean(t, rxRequired, []string{"required"}, nil)
}

func makeMinMax(lower string) (res []string) {
	for _, a := range []string{"", "imum"} {
		res = append(res, lower+a, strings.Title(lower)+a)
	}
	return
}

func verifyBoolean(t *testing.T, matcher *regexp.Regexp, names, names2 []string) {
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	validArgs := []string{"true", "false"}
	invalidArgs := []string{"TRUE", "FALSE", "t", "f", "1", "0", "True", "False", "true*", "false*"}
	var nms []string
	for _, nm := range names {
		nms = append(nms, nm, strings.Title(nm))
	}

	var nms2 []string
	for _, nm := range names2 {
		nms2 = append(nms2, nm, strings.Title(nm))
	}

	var rnms []string
	if len(nms2) > 0 {
		for _, nm := range nms {
			for _, es := range append(extraSpaces, "-") {
				for _, nm2 := range nms2 {
					rnms = append(rnms, strings.Join([]string{nm, es, nm2}, ""))
				}
			}
		}
	} else {
		rnms = nms
	}

	var cnt int
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, nm := range rnms {
				for _, es2 := range extraSpaces {
					for _, es3 := range extraSpaces {
						for _, vv := range validArgs {
							line := strings.Join([]string{pref, es1, nm, es2, ":", es3, vv}, "")
							matches := matcher.FindStringSubmatch(line)
							assert.Len(t, matches, 2)
							assert.Equal(t, vv, matches[1])
							cnt++
						}
						for _, iv := range invalidArgs {
							line := strings.Join([]string{pref, es1, nm, es2, ":", es3, iv}, "")
							matches := matcher.FindStringSubmatch(line)
							assert.Empty(t, matches)
							cnt++
						}
					}
				}
			}
		}
	}
	var nm2 string
	if len(names2) > 0 {
		nm2 = " " + names2[0]
	}
	fmt.Printf("tested %d %s%s combinations\n", cnt, names[0], nm2)
}

func verifyIntegerMinMaxManyWords(t *testing.T, matcher *regexp.Regexp, name1 string, words []string) {
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	validNumericArgs := []string{"0", "1234"}
	invalidNumericArgs := []string{"1A3F", "2e10", "*12", "12*", "-1235", "0.0", "1234.0394", "-2948.484"}

	var names []string
	for _, w := range words {
		names = append(names, w, strings.Title(w))
	}

	var cnt int
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, nm1 := range makeMinMax(name1) {
				for _, es2 := range append(extraSpaces, "-") {
					for _, nm2 := range names {
						for _, es3 := range extraSpaces {
							for _, es4 := range extraSpaces {
								for _, vv := range validNumericArgs {
									line := strings.Join([]string{pref, es1, nm1, es2, nm2, es3, ":", es4, vv}, "")
									matches := matcher.FindStringSubmatch(line)
									//fmt.Printf("matching %q, matches (%d): %v\n", line, len(matches), matches)
									assert.Len(t, matches, 2)
									assert.Equal(t, vv, matches[1])
									cnt++
								}
								for _, iv := range invalidNumericArgs {
									line := strings.Join([]string{pref, es1, nm1, es2, nm2, es3, ":", es4, iv}, "")
									matches := matcher.FindStringSubmatch(line)
									assert.Empty(t, matches)
									cnt++
								}
							}
						}
					}
				}
			}
		}
	}
	var nm2 string
	if len(words) > 0 {
		nm2 = " " + words[0]
	}
	fmt.Printf("tested %d %s%s combinations\n", cnt, name1, nm2)
}

func verifyNumeric2Words(t *testing.T, matcher *regexp.Regexp, name1, name2 string) {
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	validNumericArgs := []string{"0", "1234", "-1235", "0.0", "1234.0394", "-2948.484"}
	invalidNumericArgs := []string{"1A3F", "2e10", "*12", "12*"}

	var cnt int
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, es2 := range extraSpaces {
				for _, es3 := range extraSpaces {
					for _, es4 := range extraSpaces {
						for _, vv := range validNumericArgs {
							lines := []string{
								strings.Join([]string{pref, es1, name1, es2, name2, es3, ":", es4, vv}, ""),
								strings.Join([]string{pref, es1, strings.Title(name1), es2, strings.Title(name2), es3, ":", es4, vv}, ""),
								strings.Join([]string{pref, es1, strings.Title(name1), es2, name2, es3, ":", es4, vv}, ""),
								strings.Join([]string{pref, es1, name1, es2, strings.Title(name2), es3, ":", es4, vv}, ""),
							}
							for _, line := range lines {
								matches := matcher.FindStringSubmatch(line)
								//fmt.Printf("matching %q, matches (%d): %v\n", line, len(matches), matches)
								assert.Len(t, matches, 2)
								assert.Equal(t, vv, matches[1])
								cnt++
							}
						}
						for _, iv := range invalidNumericArgs {
							lines := []string{
								strings.Join([]string{pref, es1, name1, es2, name2, es3, ":", es4, iv}, ""),
								strings.Join([]string{pref, es1, strings.Title(name1), es2, strings.Title(name2), es3, ":", es4, iv}, ""),
								strings.Join([]string{pref, es1, strings.Title(name1), es2, name2, es3, ":", es4, iv}, ""),
								strings.Join([]string{pref, es1, name1, es2, strings.Title(name2), es3, ":", es4, iv}, ""),
							}
							for _, line := range lines {
								matches := matcher.FindStringSubmatch(line)
								//fmt.Printf("matching %q, matches (%d): %v\n", line, len(matches), matches)
								assert.Empty(t, matches)
								cnt++
							}
						}
					}
				}
			}
		}
	}
	fmt.Printf("tested %d %s %s combinations\n", cnt, name1, name2)
}

func verifyMinMax(t *testing.T, matcher *regexp.Regexp, name string, operators []string) {
	extraSpaces := []string{"", " ", "  ", "     "}
	prefixes := []string{"//", "*", ""}
	validNumericArgs := []string{"0", "1234", "-1235", "0.0", "1234.0394", "-2948.484"}
	invalidNumericArgs := []string{"1A3F", "2e10", "*12", "12*"}

	var cnt int
	for _, pref := range prefixes {
		for _, es1 := range extraSpaces {
			for _, wrd := range makeMinMax(name) {
				for _, es2 := range extraSpaces {
					for _, es3 := range extraSpaces {
						for _, op := range operators {
							for _, es4 := range extraSpaces {
								for _, vv := range validNumericArgs {
									line := strings.Join([]string{pref, es1, wrd, es2, ":", es3, op, es4, vv}, "")
									matches := matcher.FindStringSubmatch(line)
									//fmt.Printf("matching %q, matches (%d): %v\n", line, len(matches), matches)
									assert.Len(t, matches, 3)
									assert.Equal(t, vv, matches[2])
									cnt++
								}
								for _, iv := range invalidNumericArgs {
									line := strings.Join([]string{pref, es1, wrd, es2, ":", es3, op, es4, iv}, "")
									matches := matcher.FindStringSubmatch(line)
									assert.Empty(t, matches)
									cnt++
								}
							}
						}
					}
				}
			}
		}
	}
	fmt.Printf("tested %d %s combinations\n", cnt, name)
}

func verifySwaggerOneArgSwaggerTag(t *testing.T, matcher *regexp.Regexp, prefixes, validParams, invalidParams []string) {
	for _, pref := range prefixes {
		for _, param := range validParams {
			line := pref + param
			matches := matcher.FindStringSubmatch(line)
			assert.Len(t, matches, 2)
			assert.Equal(t, strings.TrimSpace(param), matches[1])
		}
	}

	for _, pref := range prefixes {
		for _, param := range invalidParams {
			line := pref + param
			matches := matcher.FindStringSubmatch(line)
			assert.Empty(t, matches)
		}
	}
}

func TestSchemaParser(t *testing.T) {
	schema := noModelDefs["NoModel"]

	assert.Equal(t, spec.StringOrArray([]string{"object"}), schema.Type)
	assert.Equal(t, "NoModel is a struct that exists in a package\nbut is not annotated with the swagger model annotations\nso it should now show up in a test", schema.Title)
	//assert.Equal(t, "this model is not explictly mentioned in the import paths\nbut because it it transitively required by the order\nit should also be collected.", schema.Description)
	assert.Len(t, schema.Required, 3)

	assertProperty(t, &schema, "number", "id", "int64", "ID")
	prop, ok := schema.Properties["id"]
	assert.True(t, ok, "should have had an 'id' property")
	assert.EqualValues(t, 1000, *prop.Maximum)
	assert.True(t, prop.ExclusiveMaximum, "'id' should have had an exclusive maximum")
	assert.NotNil(t, prop.Minimum)
	assert.EqualValues(t, 10, *prop.Minimum)
	assert.True(t, prop.ExclusiveMinimum, "'id' should have had an exclusive minimum")

	assertProperty(t, &schema, "number", "score", "int32", "Score")
	prop, ok = schema.Properties["score"]
	assert.True(t, ok, "should have had a 'score' property")
	assert.EqualValues(t, 45, *prop.Maximum)
	assert.False(t, prop.ExclusiveMaximum, "'score' should not have had an exclusive maximum")
	assert.NotNil(t, prop.Minimum)
	assert.EqualValues(t, 3, *prop.Minimum)
	assert.False(t, prop.ExclusiveMinimum, "'score' should not have had an exclusive minimum")

	assertProperty(t, &schema, "string", "created", "date-time", "Created")
	prop, ok = schema.Properties["created"]
	assert.True(t, ok, "should have a 'created' property")
	assert.True(t, prop.ReadOnly, "'created' should be read only")

	definitions := make(map[string]spec.Schema)
	sp := schemaParser(classificationProg)
	pn := "github.com/casualjim/go-swagger/fixtures/goparsing/classification/models"
	pnr := "../fixtures/goparsing/classification/models"
	pkg := classificationProg.Package(pnr)
	assert.NotNil(t, pkg)

	fnd := false
	for _, fil := range pkg.Files {
		nm := filepath.Base(classificationProg.Fset.File(fil.Pos()).Name())
		if nm == "order.go" {
			fnd = true
			sp.Parse(fil, definitions)
			break
		}
	}
	assert.True(t, fnd)
	msch, ok := definitions["order"]
	assert.True(t, ok)
	assert.Equal(t, pn, msch.Extensions["x-go-package"])
	assert.Equal(t, "StoreOrder", msch.Extensions["x-go-name"])
}

func TestParsePrimitiveSchemaProperty(t *testing.T) {
	schema := noModelDefs["PrimateModel"]
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
}

func TestParseStringFormatSchemaProperty(t *testing.T) {
	schema := noModelDefs["FormattedModel"]
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
}

func assertProperty(t *testing.T, schema *spec.Schema, typeName, jsonName, format, goName string) {
	if typeName == "" {
		assert.Empty(t, schema.Properties[jsonName].Type)
	} else {
		assert.NotEmpty(t, schema.Properties[jsonName].Type)
		assert.Equal(t, typeName, schema.Properties[jsonName].Type[0])
	}
	assert.Equal(t, goName, schema.Properties[jsonName].Extensions["x-go-name"])
	assert.Equal(t, format, schema.Properties[jsonName].Format)
}

func assertRef(t *testing.T, schema *spec.Schema, jsonName, goName, fragment string) {
	assertProperty(t, schema, "", jsonName, "", goName)
	psch := schema.Properties[jsonName]
	assert.Equal(t, fragment, psch.Ref.String())
}

func TestParseStructFields(t *testing.T) {
	schema := noModelDefs["SimpleComplexModel"]
	assertProperty(t, &schema, "object", "emb", "", "Emb")
	eSchema := schema.Properties["emb"]
	assertProperty(t, &eSchema, "number", "cid", "int64", "CID")
	assertProperty(t, &eSchema, "string", "baz", "", "Baz")

	assertRef(t, &schema, "top", "Top", "#/definitions/Something")
	assertRef(t, &schema, "notSel", "NotSel", "#/definitions/NotSelected")
}

func TestParsePointerFields(t *testing.T) {
	schema := noModelDefs["Pointdexter"]

	assertProperty(t, &schema, "number", "id", "int64", "ID")
	assertProperty(t, &schema, "string", "name", "", "Name")
	assertProperty(t, &schema, "object", "emb", "", "Emb")
	assertProperty(t, &schema, "string", "t", "uuid5", "T")
	eSchema := schema.Properties["emb"]
	assertProperty(t, &eSchema, "number", "cid", "int64", "CID")
	assertProperty(t, &eSchema, "string", "baz", "", "Baz")

	assertRef(t, &schema, "top", "Top", "#/definitions/Something")
	assertRef(t, &schema, "notSel", "NotSel", "#/definitions/NotSelected")
}

func assertArrayProperty(t *testing.T, schema *spec.Schema, typeName, jsonName, format, goName string) {
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

func assertArrayRef(t *testing.T, schema *spec.Schema, jsonName, goName, fragment string) {
	assertArrayProperty(t, schema, "", jsonName, "", goName)
	psch := schema.Properties[jsonName].Items.Schema
	assert.Equal(t, fragment, psch.Ref.String())
}

func TestParseSliceFields(t *testing.T) {
	schema := noModelDefs["SliceAndDice"]

	assertArrayProperty(t, &schema, "number", "ids", "int64", "IDs")
	assertArrayProperty(t, &schema, "string", "names", "", "Names")
	assertArrayProperty(t, &schema, "string", "uuids", "uuid", "UUIDs")
	assertArrayProperty(t, &schema, "object", "embs", "", "Embs")
	eSchema := schema.Properties["embs"].Items.Schema
	assertArrayProperty(t, eSchema, "number", "cid", "int64", "CID")
	assertArrayProperty(t, eSchema, "string", "baz", "", "Baz")

	assertArrayRef(t, &schema, "tops", "Tops", "#/definitions/Something")
	assertArrayRef(t, &schema, "notSels", "NotSels", "#/definitions/NotSelected")

	assertArrayProperty(t, &schema, "number", "ptrIds", "int64", "PtrIDs")
	assertArrayProperty(t, &schema, "string", "ptrNames", "", "PtrNames")
	assertArrayProperty(t, &schema, "string", "ptrUuids", "uuid", "PtrUUIDs")
	assertArrayProperty(t, &schema, "object", "ptrEmbs", "", "PtrEmbs")
	eSchema = schema.Properties["ptrEmbs"].Items.Schema
	assertArrayProperty(t, eSchema, "number", "ptrCid", "int64", "PtrCID")
	assertArrayProperty(t, eSchema, "string", "ptrBaz", "", "PtrBaz")

	assertArrayRef(t, &schema, "ptrTops", "PtrTops", "#/definitions/Something")
	assertArrayRef(t, &schema, "ptrNotSels", "PtrNotSels", "#/definitions/NotSelected")
}

func assertMapProperty(t *testing.T, schema *spec.Schema, typeName, jsonName, format, goName string) {
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

func assertMapRef(t *testing.T, schema *spec.Schema, jsonName, goName, fragment string) {
	assertMapProperty(t, schema, "", jsonName, "", goName)
	psch := schema.Properties[jsonName].AdditionalProperties.Schema
	assert.Equal(t, fragment, psch.Ref.String())
}

func TestParseMapFields(t *testing.T) {
	schema := noModelDefs["MapTastic"]

	assertMapProperty(t, &schema, "number", "ids", "int64", "IDs")
	assertMapProperty(t, &schema, "string", "names", "", "Names")
	assertMapProperty(t, &schema, "string", "uuids", "uuid", "UUIDs")
	assertMapProperty(t, &schema, "object", "embs", "", "Embs")
	eSchema := schema.Properties["embs"].AdditionalProperties.Schema
	assertMapProperty(t, eSchema, "number", "cid", "int64", "CID")
	assertMapProperty(t, eSchema, "string", "baz", "", "Baz")

	assertMapRef(t, &schema, "tops", "Tops", "#/definitions/Something")
	assertMapRef(t, &schema, "notSels", "NotSels", "#/definitions/NotSelected")

	assertMapProperty(t, &schema, "number", "ptrIds", "int64", "PtrIDs")
	assertMapProperty(t, &schema, "string", "ptrNames", "", "PtrNames")
	assertMapProperty(t, &schema, "string", "ptrUuids", "uuid", "PtrUUIDs")
	assertMapProperty(t, &schema, "object", "ptrEmbs", "", "PtrEmbs")
	eSchema = schema.Properties["ptrEmbs"].AdditionalProperties.Schema
	assertMapProperty(t, eSchema, "number", "ptrCid", "int64", "PtrCID")
	assertMapProperty(t, eSchema, "string", "ptrBaz", "", "PtrBaz")

	assertMapRef(t, &schema, "ptrTops", "PtrTops", "#/definitions/Something")
	assertMapRef(t, &schema, "ptrNotSels", "PtrNotSels", "#/definitions/NotSelected")
}

package validate

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
	. "github.com/smartystreets/goconvey/convey"
)

type schemaTestT struct {
	Description string       `json:"description"`
	Schema      *spec.Schema `json:"schema"`
	Tests       []struct {
		Description string      `json:"description"`
		Data        interface{} `json:"data"`
		Valid       bool        `json:"valid"`
	}
}

var jsonSchemaFixturesPath = filepath.Join("..", "fixtures", "jsonschema_suite")

var ints = []interface{}{
	1,
	int8(1),
	int16(1),
	int(1),
	int32(1),
	int64(1),
	uint8(1),
	uint16(1),
	uint(1),
	uint32(1),
	uint64(1),
	5.0,
	float32(5.0),
	float64(5.0),
}

var notInts = []interface{}{
	5.1,
	float32(5.1),
	float64(5.1),
}

var notNumbers = []interface{}{
	map[string]string{},
	struct{}{},
	time.Time{},
	"yada",
}

var enabled = []string{
	"minLength",
	"maxLength",
	"pattern",
	"type",
	"minimum",
	"maximum",
	"multipleOf",
	"enum",
	"default",
	"dependencies",
	"items",
	"maxItems",
	"maxProperties",
	"minItems",
	"minProperties",
	"patternProperties",
	"required",
	"additionalItems",
	"uniqueItems",
	"properties",
	"additionalProperties",
	"allOf",
	"not",
	"oneOf",
	"anyOf",

	// These still fail
	// Ref is not implemented yet, so these should not pass yet.
	// "definitions",
	// "ref",
	// "refRemote",
}

func isEnabled(nm string) bool {
	return util.ContainsStringsCI(enabled, nm)
}

func TestJSONSchemaSuite(t *testing.T) {

	Convey("The JSON Schema test suite", t, func() {
		files, _ := ioutil.ReadDir(jsonSchemaFixturesPath)

		for _, f := range files {
			if f.IsDir() {
				continue
			}
			fileName := f.Name()
			specName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			if isEnabled(specName) {

				Convey("for "+specName, func() {
					b, _ := ioutil.ReadFile(filepath.Join(jsonSchemaFixturesPath, fileName))

					var testDescriptions []schemaTestT
					json.Unmarshal(b, &testDescriptions)

					for _, testDescription := range testDescriptions {

						Convey(testDescription.Description, func() {

							validator := newSchemaValidator(testDescription.Schema, "data")

							for _, test := range testDescription.Tests {

								Convey(test.Description, func() {

									result := validator.Validate(test.Data)
									So(result, ShouldNotBeNil)

									if test.Valid {
										So(result.Errors, ShouldBeEmpty)

									} else {
										So(result.Errors, ShouldNotBeEmpty)
									}

								})
							}
						})
					}
				})
			}
		}
	})
	// Convey("go types for int", t, func() {
	// 	for _, i := range ints {
	// 		runSchemaTest("type/schema_0.json", i, true)
	// 	}
	// 	for _, i := range notInts {
	// 		runSchemaTest("type/schema_0.json", i, false)
	// 	}
	// 	for _, i := range notNumbers {
	// 		runSchemaTest("type/schema_0.json", i, false)
	// 	}
	// })

	// Convey("go types for number", t, func() {
	// 	for _, i := range ints {
	// 		runSchemaTest("type/schema_1.json", i, true)
	// 	}
	// 	for _, i := range notInts {
	// 		runSchemaTest("type/schema_1.json", i, true)
	// 	}
	// 	for _, i := range notNumbers {
	// 		runSchemaTest("type/schema_1.json", i, false)
	// 	}
	// })
}

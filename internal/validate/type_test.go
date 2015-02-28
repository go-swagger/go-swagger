package validate

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
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

var jsonSchemaFixturesPath = filepath.Join("..", "..", "fixtures", "jsonschema_suite")

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
	"ref",
	"definitions",
	"refRemote",
	"format",
}

type noopResCache struct {
}

func (n *noopResCache) Get(key string) (interface{}, bool) {
	return nil, false
}
func (n *noopResCache) Set(string, interface{}) {

}

func isEnabled(nm string) bool {
	return util.ContainsStringsCI(enabled, nm)
}

func TestJSONSchemaSuite(t *testing.T) {
	go func() {
		err := http.ListenAndServe(":1234", http.FileServer(http.Dir(jsonSchemaFixturesPath+"/remotes")))
		if err != nil {
			panic(err.Error())
		}
	}()
	Convey("The JSON Schema test suite", t, func() {
		files, err := ioutil.ReadDir(jsonSchemaFixturesPath)
		if err != nil {
			t.Fatal(err)
		}

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

							So(spec.ExpandSchema(testDescription.Schema, nil, nil /*new(noopResCache)*/), ShouldBeNil)

							validator := NewSchemaValidator(testDescription.Schema, nil, "data", strfmt.Default)

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
}

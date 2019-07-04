//+build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1571/gen-fixture-simple-tuple-minimal/models"
	"github.com/stretchr/testify/assert"
)

func Test_TupleThing(t *testing.T) {
	base := "tupleThing-data"
	cwd, _ := os.Getwd()
	cwd = filepath.Join(cwd, "json-data")
	schemaSource := filepath.Join(cwd, "tupleThing.json")
	// read schema
	jsonSchema, _ := ioutil.ReadFile(schemaSource)
	schema := new(spec.Schema)
	err := json.Unmarshal(jsonSchema, schema)
	if !assert.NoError(t, err) {
		t.FailNow()
		return
	}
	filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		fixture := info.Name()
		//t.Logf("Found: %s", fixture)
		if !info.IsDir() && strings.HasPrefix(fixture, base) {
			// read fixture
			buf, _ := ioutil.ReadFile(filepath.Join("json-data", fixture))

			t.Logf("INFO:Fixture: %s: %s", fixture, string(buf))
			input := []interface{}{}
			erm := json.Unmarshal(buf, &input)
			if !assert.NoErrorf(t, erm, "ERROR:Error unmarshaling fixture: %v", erm) {
				t.FailNow()
				return erm
			}

			// run validate.AgainstSchema
			//bb, _ := json.MarshalIndent(schema, "", " ")
			//t.Log(string(bb))
			erj := validate.AgainstSchema(schema, input, strfmt.Default)
			if erj == nil {
				t.Logf("INFO:Validation AgainstSchema for %s returned: valid", fixture)
			} else {
				t.Logf("INFO AgainstSchema for %s returned: invalid, with %v", fixture, erj)
			}
			//bb, _ = json.MarshalIndent(schema, "", " ")
			//t.Log(string(bb))
			// unmarshall into model
			var erv error
			model := models.TupleThing{}
			eru := model.UnmarshalJSON(buf)
			if assert.NoErrorf(t, eru, "ERROR:Error unmarshaling struct: %v", eru) {
				// run model validation
				erv = model.Validate(strfmt.Default)
				if erv == nil {
					t.Logf("INFO:Validation for %s returned: valid", fixture)

				} else {
					t.Logf("INFO:Validation for %s returned: invalid, with: %v", fixture, erv)
				}
			} else {
				t.FailNow()
				return eru
			}
			// marshall the model back to json
			bbb, erm := model.MarshalJSON()
			if assert.NoErrorf(t, erm, "ERROR:Error marshaling: %v", erm) {
				t.Logf("INFO:Data marshalled as: %s", string(bbb))
			}
			// compare validation methods
			if erv != nil && erj == nil || erv == nil && erj != nil {
				t.Logf("ERROR:Our validators returned different results for: %s", fixture)
				if fixture == strings.Join([]string{base, "2"}, "-")+".json" {
					t.Logf("WARNING: expected failure - see issue #1486")
				} else {
					t.Fail()
				}
			}
		}
		return nil
	})

}

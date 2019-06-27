//+build ignore

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-swagger/go-swagger/fixtures/bugs/1548/gen-fixture-1548/models"
	"github.com/stretchr/testify/assert"
)

func Test_Base64Thing(t *testing.T) {
	base := "base64Thing-data"
	cwd, _ := os.Getwd()
	// read schema
	filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		fixture := info.Name()
		if !info.IsDir() && strings.HasPrefix(fixture, base) {
			// read fixture
			buf, _ := ioutil.ReadFile(fixture)

			t.Logf("INFO:Fixture: %s: %s", fixture, string(buf))
			var input interface{}
			erm := json.Unmarshal(buf, &input)
			if !assert.NoError(t, erm, "ERROR:Error unmarshaling fixture: %v", erm) {
				t.FailNow()
				return erm
			}

			var erv error
			model := models.Base64Model{}
			eru := json.Unmarshal(buf, &model)
			if fixture != base+"-3.json" {
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
			} else {
				t.Logf("INFO: expected error= invalid base 64 string")
			}
			// marshall the model back to json
			bbb, erm := json.Marshal(model)
			if assert.NoErrorf(t, erm, "ERROR:Error marshaling: %v", erm) {
				t.Logf("INFO:Data internal representation: %s", string(model.Prop1))
				t.Logf("INFO:Data marshalled as: %s", string(bbb))
			}
		}
		return nil
	})

}

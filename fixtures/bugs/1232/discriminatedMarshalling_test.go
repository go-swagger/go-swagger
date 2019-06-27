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
	"github.com/go-swagger/go-swagger/fixtures/bugs/1232/gen-fixture-1232/models"
	"github.com/stretchr/testify/assert"
)

func Test_Pet(t *testing.T) {
	base := "pet-data"
	cwd, _ := os.Getwd()
	filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		fixture := info.Name()
		if !info.IsDir() && strings.HasPrefix(fixture, base) {
			// read fixture
			buf, _ := ioutil.ReadFile(fixture)

			t.Logf("Fixture: %s", string(buf))
			input := []interface{}{}
			json.Unmarshal(buf, input)

			// unmarshall into model
			model := models.TupleThing{}
			err = model.UnmarshalJSON(buf)
			if assert.NoError(t, err) {
				err = model.Validate(strfmt.Default)
				if err == nil {
					t.Logf("Validation for %s returned: valid", fixture)

				} else {
					t.Logf("Validation for %s returned: invalid, with: %v", fixture, err)
				}
			}
		}
		return nil
	})

}

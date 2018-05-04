// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package post

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

var defaulterFixturesPath = filepath.Join("..", "fixtures", "defaulting")

func TestDefaulter(t *testing.T) {
	schema, err := defaulterFixture()
	assert.NoError(t, err)

	validator := validate.NewSchemaValidator(schema, nil, "", strfmt.Default)
	x := defaulterFixtureInput()
	t.Logf("Before: %v", x)

	r := validator.Validate(x)
	assert.False(t, r.HasErrors(), fmt.Sprintf("unexpected validation error: %v", r.AsError()))

	ApplyDefaults(r)
	t.Logf("After: %v", x)
	var expected interface{}
	err = json.Unmarshal([]byte(`{
		"existing": 100,
		"int": 42,
		"str": "Hello",
		"obj": {"foo": "bar"},
		"nested": {"inner": 7},
		"all": {"foo": 42, "bar": 42},
		"any": {"foo": 42},
		"one": {"bar": 42}
	}`), &expected)
	assert.NoError(t, err)
	assert.Equal(t, expected, x)
}

func TestDefaulterSimple(t *testing.T) {
	schema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Properties: map[string]spec.Schema{
				"int": {
					SchemaProps: spec.SchemaProps{
						Default: float64(42),
					},
				},
				"str": {
					SchemaProps: spec.SchemaProps{
						Default: "Hello",
					},
				},
			},
		},
	}
	validator := validate.NewSchemaValidator(&schema, nil, "", strfmt.Default)
	x := map[string]interface{}{}
	t.Logf("Before: %v", x)
	r := validator.Validate(x)
	assert.False(t, r.HasErrors(), fmt.Sprintf("unexpected validation error: %v", r.AsError()))

	ApplyDefaults(r)
	t.Logf("After: %v", x)
	var expected interface{}
	err := json.Unmarshal([]byte(`{
		"int": 42,
		"str": "Hello"
	}`), &expected)
	assert.NoError(t, err)
	assert.Equal(t, expected, x)
}

func BenchmarkDefaulting(b *testing.B) {
	b.ReportAllocs()

	schema, err := defaulterFixture()
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		validator := validate.NewSchemaValidator(schema, nil, "", strfmt.Default)
		x := defaulterFixtureInput()
		r := validator.Validate(x)
		assert.False(b, r.HasErrors(), fmt.Sprintf("unexpected validation error: %v", r.AsError()))
		ApplyDefaults(r)
	}
}

func defaulterFixtureInput() map[string]interface{} {
	return map[string]interface{}{
		"existing": float64(100),
		"nested":   map[string]interface{}{},
		"all":      map[string]interface{}{},
		"any":      map[string]interface{}{},
		"one":      map[string]interface{}{},
	}
}

func defaulterFixture() (*spec.Schema, error) {
	fname := filepath.Join(defaulterFixturesPath, "schema.json")
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	var schema spec.Schema
	if err := json.Unmarshal(b, &schema); err != nil {
		return nil, err
	}

	return &schema, spec.ExpandSchema(&schema, nil, nil /*new(noopResCache)*/)
}

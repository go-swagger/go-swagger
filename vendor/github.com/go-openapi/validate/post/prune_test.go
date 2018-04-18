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
	validate "github.com/go-openapi/validate"
)

var pruneFixturesPath = filepath.Join("..", "fixtures", "pruning")

func TestPrune(t *testing.T) {
	schema, err := pruningFixture()
	assert.NoError(t, err)

	x := map[string]interface{}{
		"foo": 42,
		"bar": 42,
		"x":   42,
		"nested": map[string]interface{}{
			"x": 42,
			"inner": map[string]interface{}{
				"foo": 42,
				"bar": 42,
				"x":   42,
			},
		},
		"all": map[string]interface{}{
			"foo": 42,
			"bar": 42,
			"x":   42,
		},
		"any": map[string]interface{}{
			"foo": 42,
			"bar": 42,
			"x":   42,
		},
		"one": map[string]interface{}{
			"bar": 42,
			"x":   42,
		},
		"array": []interface{}{
			map[string]interface{}{
				"foo": 42,
				"bar": 123,
			},
			map[string]interface{}{
				"x": 42,
				"y": 123,
			},
		},
	}
	t.Logf("Before: %v", x)

	validator := validate.NewSchemaValidator(schema, nil, "", strfmt.Default)
	r := validator.Validate(x)
	assert.False(t, r.HasErrors(), fmt.Sprintf("unexpected validation error: %v", r.AsError()))

	Prune(r)
	t.Logf("After: %v", x)
	expected := map[string]interface{}{
		"foo": 42,
		"bar": 42,
		"nested": map[string]interface{}{
			"inner": map[string]interface{}{
				"foo": 42,
				"bar": 42,
			},
		},
		"all": map[string]interface{}{
			"foo": 42,
			"bar": 42,
		},
		"any": map[string]interface{}{
			// intentionally only list one: the first matching
			"foo": 42,
		},
		"one": map[string]interface{}{
			"bar": 42,
		},
		"array": []interface{}{
			map[string]interface{}{
				"foo": 42,
			},
			map[string]interface{}{},
		},
	}
	assert.Equal(t, expected, x)
}

func pruningFixture() (*spec.Schema, error) {
	fname := filepath.Join(pruneFixturesPath, "schema.json")
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

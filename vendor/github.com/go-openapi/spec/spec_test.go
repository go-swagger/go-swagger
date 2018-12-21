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

package spec_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

// mimics what the go-openapi/load does
var (
	yamlLoader = swag.YAMLDoc
	rex        = regexp.MustCompile(`"\$ref":\s*"(.+)"`)
)

func loadOrFail(t *testing.T, path string) *spec.Swagger {
	raw, erl := yamlLoader(path)
	if erl != nil {
		t.Logf("can't load fixture %s: %v", path, erl)
		t.FailNow()
		return nil
	}
	swspec := new(spec.Swagger)
	if err := json.Unmarshal(raw, swspec); err != nil {
		t.FailNow()
		return nil
	}
	return swspec
}

// Test unitary fixture for dev and bug fixing
func Test_Issue1429(t *testing.T) {
	prevPathLoader := spec.PathLoader
	defer func() {
		spec.PathLoader = prevPathLoader
	}()
	spec.PathLoader = yamlLoader
	path := filepath.Join("fixtures", "bugs", "1429", "swagger.yaml")

	// load and full expand
	sp := loadOrFail(t, path)
	err := spec.ExpandSpec(sp, &spec.ExpandOptions{RelativeBase: path, SkipSchemas: false})
	if !assert.NoError(t, err) {
		t.FailNow()
		return
	}

	// assert well expanded
	if !assert.Truef(t, (sp.Paths != nil && sp.Paths.Paths != nil), "expected paths to be available in fixture") {
		t.FailNow()
		return
	}
	for _, pi := range sp.Paths.Paths {
		for _, param := range pi.Get.Parameters {
			if assert.NotNilf(t, param.Schema, "expected param schema not to be nil") {
				// all param fixtures are body param with schema
				// all $ref expanded
				assert.Equal(t, "", param.Schema.Ref.String())
			}
		}
		for code, response := range pi.Get.Responses.StatusCodeResponses {
			// all response fixtures are with StatusCodeResponses, but 200
			if code == 200 {
				assert.Nilf(t, response.Schema, "expected response schema to be nil")
				continue
			}
			if assert.NotNilf(t, response.Schema, "expected response schema not to be nil") {
				assert.Equal(t, "", response.Schema.Ref.String())
			}
		}
	}
	for _, def := range sp.Definitions {
		assert.Equal(t, "", def.Ref.String())
	}

	// reload and SkipSchemas: true
	sp = loadOrFail(t, path)
	err = spec.ExpandSpec(sp, &spec.ExpandOptions{RelativeBase: path, SkipSchemas: true})
	if !assert.NoError(t, err) {
		t.FailNow()
		return
	}

	// assert well resolved
	if !assert.Truef(t, (sp.Paths != nil && sp.Paths.Paths != nil), "expected paths to be available in fixture") {
		t.FailNow()
		return
	}
	for _, pi := range sp.Paths.Paths {
		for _, param := range pi.Get.Parameters {
			if assert.NotNilf(t, param.Schema, "expected param schema not to be nil") {
				// all param fixtures are body param with schema
				if param.Name == "plainRequest" {
					// this one is expanded
					assert.Equal(t, "", param.Schema.Ref.String())
					continue
				}
				if param.Name == "nestedBody" {
					// this one is local
					assert.True(t, strings.HasPrefix(param.Schema.Ref.String(), "#/definitions/"))
					continue
				}
				if param.Name == "remoteRequest" {
					assert.Contains(t, param.Schema.Ref.String(), "remote/remote.yaml#/")
					continue
				}
				assert.Contains(t, param.Schema.Ref.String(), "responses.yaml#/")
			}
		}
		for code, response := range pi.Get.Responses.StatusCodeResponses {
			// all response fixtures are with StatusCodeResponses, but 200
			if code == 200 {
				assert.Nilf(t, response.Schema, "expected response schema to be nil")
				continue
			}
			if code == 204 {
				assert.Contains(t, response.Schema.Ref.String(), "remote/remote.yaml#/")
				continue
			}
			if code == 404 {
				assert.Equal(t, "", response.Schema.Ref.String())
				continue
			}
			assert.Containsf(t, response.Schema.Ref.String(), "responses.yaml#/", "expected remote ref at resp. %d", code)
		}
	}
	for _, def := range sp.Definitions {
		assert.Contains(t, def.Ref.String(), "responses.yaml#/")
	}
}

func Test_MoreLocalExpansion(t *testing.T) {
	prevPathLoader := spec.PathLoader
	defer func() {
		spec.PathLoader = prevPathLoader
	}()
	spec.PathLoader = yamlLoader
	path := filepath.Join("fixtures", "local_expansion", "spec2.yaml")

	// load and full expand
	sp := loadOrFail(t, path)
	err := spec.ExpandSpec(sp, &spec.ExpandOptions{RelativeBase: path, SkipSchemas: false})
	if !assert.NoError(t, err) {
		t.FailNow()
		return
	}
	// asserts all $ref expanded
	jazon, _ := json.MarshalIndent(sp, "", " ")
	assert.NotContains(t, jazon, `"$ref"`)
	//t.Log(string(jazon))
}

func Test_Issue69(t *testing.T) {
	// this checks expansion for the dapperbox spec (circular ref issues)

	path := filepath.Join("fixtures", "bugs", "69", "dapperbox.json")

	// expand with relative path
	// load and expand
	sp := loadOrFail(t, path)
	err := spec.ExpandSpec(sp, &spec.ExpandOptions{RelativeBase: path, SkipSchemas: false})
	if !assert.NoError(t, err) {
		t.FailNow()
		return
	}
	// asserts all $ref expanded
	jazon, _ := json.MarshalIndent(sp, "", " ")

	// assert all $ref match  "$ref": "#/definitions/something"
	m := rex.FindAllStringSubmatch(string(jazon), -1)
	if assert.NotNil(t, m) {
		for _, matched := range m {
			subMatch := matched[1]
			assert.True(t, strings.HasPrefix(subMatch, "#/definitions/"),
				"expected $ref to be inlined, got: %s", matched[0])
		}
	}
}

func Test_Issue1621(t *testing.T) {
	prevPathLoader := spec.PathLoader
	defer func() {
		spec.PathLoader = prevPathLoader
	}()
	spec.PathLoader = yamlLoader
	path := filepath.Join("fixtures", "bugs", "1621", "fixture-1621.yaml")

	// expand with relative path
	// load and expand
	sp := loadOrFail(t, path)

	err := spec.ExpandSpec(sp, &spec.ExpandOptions{RelativeBase: path, SkipSchemas: false})
	if !assert.NoError(t, err) {
		t.FailNow()
		return
	}
	// asserts all $ref expanded
	jazon, _ := json.MarshalIndent(sp, "", " ")
	m := rex.FindAllStringSubmatch(string(jazon), -1)
	assert.Nil(t, m)
}

func Test_Issue1614(t *testing.T) {

	path := filepath.Join("fixtures", "bugs", "1614", "gitea.json")

	// expand with relative path
	// load and expand
	sp := loadOrFail(t, path)
	err := spec.ExpandSpec(sp, &spec.ExpandOptions{RelativeBase: path, SkipSchemas: false})
	if !assert.NoError(t, err) {
		t.FailNow()
		return
	}
	// asserts all $ref expanded
	jazon, _ := json.MarshalIndent(sp, "", " ")

	// assert all $ref maches  "$ref": "#/definitions/something"
	m := rex.FindAllStringSubmatch(string(jazon), -1)
	if assert.NotNil(t, m) {
		for _, matched := range m {
			subMatch := matched[1]
			assert.True(t, strings.HasPrefix(subMatch, "#/definitions/"),
				"expected $ref to be inlined, got: %s", matched[0])
		}
	}

	// now with option CircularRefAbsolute
	sp = loadOrFail(t, path)
	err = spec.ExpandSpec(sp, &spec.ExpandOptions{RelativeBase: path, SkipSchemas: false, AbsoluteCircularRef: true})
	if !assert.NoError(t, err) {
		t.FailNow()
		return
	}
	// asserts all $ref expanded
	jazon, _ = json.MarshalIndent(sp, "", " ")

	// assert all $ref maches  "$ref": "{file path}#/definitions/something"
	refPath, _ := os.Getwd()
	refPath = filepath.Join(refPath, path)
	m = rex.FindAllStringSubmatch(string(jazon), -1)
	if assert.NotNil(t, m) {
		for _, matched := range m {
			subMatch := matched[1]
			assert.True(t, strings.HasPrefix(subMatch, refPath+"#/definitions/"),
				"expected $ref to be inlined, got: %s", matched[0])
		}
	}
}

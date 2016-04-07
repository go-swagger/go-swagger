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

package spec

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	testingutil "github.com/go-swagger/go-swagger/internal/testing"
	"github.com/go-swagger/go-swagger/jsonpointer"
	"github.com/go-swagger/go-swagger/swag"
	"github.com/stretchr/testify/assert"
)

func TestSpecExpansion(t *testing.T) {
	spec := new(Swagger)
	// resolver, err := defaultSchemaLoader(spec, nil, nil)
	// assert.NoError(t, err)

	err := expandSpec(spec)
	assert.NoError(t, err)

	specDoc, err := swag.JSONDoc("../fixtures/expansion/all-the-things.json")
	assert.NoError(t, err)

	spec = new(Swagger)
	err = json.Unmarshal(specDoc, spec)
	assert.NoError(t, err)

	pet := spec.Definitions["pet"]
	errorModel := spec.Definitions["errorModel"]
	petResponse := spec.Responses["petResponse"]
	petResponse.Schema = &pet
	stringResponse := spec.Responses["stringResponse"]
	tagParam := spec.Parameters["tag"]
	idParam := spec.Parameters["idParam"]

	err = expandSpec(spec)
	assert.NoError(t, err)

	assert.Equal(t, tagParam, spec.Parameters["query"])
	assert.Equal(t, petResponse, spec.Responses["petResponse"])
	assert.Equal(t, petResponse, spec.Responses["anotherPet"])
	assert.Equal(t, pet, *spec.Responses["petResponse"].Schema)
	assert.Equal(t, stringResponse, *spec.Paths.Paths["/"].Get.Responses.Default)
	assert.Equal(t, petResponse, spec.Paths.Paths["/"].Get.Responses.StatusCodeResponses[200])
	assert.Equal(t, pet, *spec.Paths.Paths["/pets"].Get.Responses.StatusCodeResponses[200].Schema.Items.Schema)
	assert.Equal(t, errorModel, *spec.Paths.Paths["/pets"].Get.Responses.Default.Schema)
	assert.Equal(t, pet, spec.Definitions["petInput"].AllOf[0])
	assert.Equal(t, spec.Definitions["petInput"], *spec.Paths.Paths["/pets"].Post.Parameters[0].Schema)
	assert.Equal(t, petResponse, spec.Paths.Paths["/pets"].Post.Responses.StatusCodeResponses[200])
	assert.Equal(t, errorModel, *spec.Paths.Paths["/pets"].Post.Responses.Default.Schema)
	pi := spec.Paths.Paths["/pets/{id}"]
	assert.Equal(t, idParam, pi.Get.Parameters[0])
	assert.Equal(t, petResponse, pi.Get.Responses.StatusCodeResponses[200])
	assert.Equal(t, errorModel, *pi.Get.Responses.Default.Schema)
	assert.Equal(t, idParam, pi.Delete.Parameters[0])
	assert.Equal(t, errorModel, *pi.Delete.Responses.Default.Schema)
}

func TestResponseExpansion(t *testing.T) {
	specDoc, err := swag.JSONDoc("../fixtures/expansion/all-the-things.json")
	assert.NoError(t, err)

	spec := new(Swagger)
	err = json.Unmarshal(specDoc, spec)
	assert.NoError(t, err)

	resolver, err := defaultSchemaLoader(spec, nil, nil)
	assert.NoError(t, err)

	resp := spec.Responses["anotherPet"]
	expected := spec.Responses["petResponse"]

	err = expandResponse(&resp, resolver)
	assert.NoError(t, err)
	assert.Equal(t, expected, resp)

	resp2 := spec.Paths.Paths["/"].Get.Responses.Default
	expected = spec.Responses["stringResponse"]

	err = expandResponse(resp2, resolver)
	assert.NoError(t, err)
	assert.Equal(t, expected, *resp2)

	resp = spec.Paths.Paths["/"].Get.Responses.StatusCodeResponses[200]
	expected = spec.Responses["petResponse"]

	err = expandResponse(&resp, resolver)
	assert.NoError(t, err)
	// assert.Equal(t, expected, resp)
}

func TestParameterExpansion(t *testing.T) {
	paramDoc, err := swag.JSONDoc("../fixtures/expansion/params.json")
	assert.NoError(t, err)

	spec := new(Swagger)
	err = json.Unmarshal(paramDoc, spec)
	assert.NoError(t, err)

	resolver, err := defaultSchemaLoader(spec, nil, nil)
	assert.NoError(t, err)

	param := spec.Parameters["query"]
	expected := spec.Parameters["tag"]

	err = expandParameter(&param, resolver)
	assert.NoError(t, err)
	assert.Equal(t, expected, param)

	param = spec.Paths.Paths["/cars/{id}"].Parameters[0]
	expected = spec.Parameters["id"]

	err = expandParameter(&param, resolver)
	assert.NoError(t, err)
	assert.Equal(t, expected, param)
}

func TestCircularRefsExpansion(t *testing.T) {
	carsDoc, err := swag.JSONDoc("../fixtures/expansion/circularRefs.json")
	assert.NoError(t, err)

	spec := new(Swagger)
	err = json.Unmarshal(carsDoc, spec)
	assert.NoError(t, err)

	resolver, err := defaultSchemaLoader(spec, nil, nil)
	assert.NoError(t, err)
	schema := spec.Definitions["car"]

	assert.NotPanics(t, func() {
		err = expandSchema(&schema, "#/definitions/car", resolver)
		assert.NoError(t, err)
	}, "Calling expand schema with circular refs, should not panic!")
}

func TestIssue415(t *testing.T) {
	doc, err := swag.YAMLDoc("../fixtures/expansion/clickmeter.yaml")
	assert.NoError(t, err)

	spec := new(Swagger)
	err = json.Unmarshal(doc, spec)
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		err = expandSpec(spec)
		assert.NoError(t, err)
	}, "Calling expand spec with response schemas that have circular refs, should not panic!")
}

func TestCircularSpecExpansion(t *testing.T) {
	doc, err := swag.YAMLDoc("../fixtures/expansion/circularSpec.yaml")
	assert.NoError(t, err)

	spec := new(Swagger)
	err = json.Unmarshal(doc, spec)
	assert.NoError(t, err)

	assert.NotPanics(t, func() {
		err = expandSpec(spec)
		assert.NoError(t, err)
	}, "Calling expand spec with circular refs, should not panic!")
}

func TestItemsExpansion(t *testing.T) {
	carsDoc, err := swag.JSONDoc("../fixtures/expansion/schemas2.json")
	assert.NoError(t, err)

	spec := new(Swagger)
	err = json.Unmarshal(carsDoc, spec)
	assert.NoError(t, err)

	resolver, err := defaultSchemaLoader(spec, nil, nil)
	assert.NoError(t, err)

	schema := spec.Definitions["car"]
	oldBrand := schema.Properties["brand"]
	assert.NotEmpty(t, oldBrand.Items.Schema.Ref.String())
	assert.NotEqual(t, spec.Definitions["brand"], oldBrand)

	err = expandSchema(&schema, "#/definitions/car", resolver)
	assert.NoError(t, err)

	newBrand := schema.Properties["brand"]
	assert.Empty(t, newBrand.Items.Schema.Ref.String())
	assert.Equal(t, spec.Definitions["brand"], *newBrand.Items.Schema)

	schema = spec.Definitions["truck"]
	assert.NotEmpty(t, schema.Items.Schema.Ref.String())

	err = expandSchema(&schema, "#/definitions/truck", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.Items.Schema.Ref.String())
	assert.Equal(t, spec.Definitions["car"], *schema.Items.Schema)

	sch := new(Schema)
	err = expandSchema(sch, "", resolver)
	assert.NoError(t, err)

	schema = spec.Definitions["batch"]
	err = expandSchema(&schema, "#/definitions/batch", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.Items.Schema.Items.Schema.Ref.String())
	assert.Equal(t, *schema.Items.Schema.Items.Schema, spec.Definitions["brand"])

	schema = spec.Definitions["batch2"]
	err = expandSchema(&schema, "#/definitions/batch2", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.Items.Schemas[0].Items.Schema.Ref.String())
	assert.Empty(t, schema.Items.Schemas[1].Items.Schema.Ref.String())
	assert.Equal(t, *schema.Items.Schemas[0].Items.Schema, spec.Definitions["brand"])
	assert.Equal(t, *schema.Items.Schemas[1].Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["allofBoth"]
	err = expandSchema(&schema, "#/definitions/allofBoth", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.AllOf[0].Items.Schema.Ref.String())
	assert.Empty(t, schema.AllOf[1].Items.Schema.Ref.String())
	assert.Equal(t, *schema.AllOf[0].Items.Schema, spec.Definitions["brand"])
	assert.Equal(t, *schema.AllOf[1].Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["anyofBoth"]
	err = expandSchema(&schema, "#/definitions/anyofBoth", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.AnyOf[0].Items.Schema.Ref.String())
	assert.Empty(t, schema.AnyOf[1].Items.Schema.Ref.String())
	assert.Equal(t, *schema.AnyOf[0].Items.Schema, spec.Definitions["brand"])
	assert.Equal(t, *schema.AnyOf[1].Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["oneofBoth"]
	err = expandSchema(&schema, "#/definitions/oneofBoth", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.OneOf[0].Items.Schema.Ref.String())
	assert.Empty(t, schema.OneOf[1].Items.Schema.Ref.String())
	assert.Equal(t, *schema.OneOf[0].Items.Schema, spec.Definitions["brand"])
	assert.Equal(t, *schema.OneOf[1].Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["notSomething"]
	err = expandSchema(&schema, "#/definitions/notSomething", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.Not.Items.Schema.Ref.String())
	assert.Equal(t, *schema.Not.Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["withAdditional"]
	err = expandSchema(&schema, "#/definitions/withAdditional", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.AdditionalProperties.Schema.Items.Schema.Ref.String())
	assert.Equal(t, *schema.AdditionalProperties.Schema.Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["withAdditionalItems"]
	err = expandSchema(&schema, "#/definitions/withAdditionalItems", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.AdditionalItems.Schema.Items.Schema.Ref.String())
	assert.Equal(t, *schema.AdditionalItems.Schema.Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["withPattern"]
	err = expandSchema(&schema, "#/definitions/withPattern", resolver)
	assert.NoError(t, err)
	prop := schema.PatternProperties["^x-ab"]
	assert.Empty(t, prop.Items.Schema.Ref.String())
	assert.Equal(t, *prop.Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["deps"]
	err = expandSchema(&schema, "#/definitions/deps", resolver)
	assert.NoError(t, err)
	prop2 := schema.Dependencies["something"]
	assert.Empty(t, prop2.Schema.Items.Schema.Ref.String())
	assert.Equal(t, *prop2.Schema.Items.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["defined"]
	err = expandSchema(&schema, "#/definitions/defined", resolver)
	assert.NoError(t, err)
	prop = schema.Definitions["something"]
	assert.Empty(t, prop.Items.Schema.Ref.String())
	assert.Equal(t, *prop.Items.Schema, spec.Definitions["tag"])
}

func TestSchemaExpansion(t *testing.T) {
	carsDoc, err := swag.JSONDoc("../fixtures/expansion/schemas1.json")
	assert.NoError(t, err)

	spec := new(Swagger)
	err = json.Unmarshal(carsDoc, spec)
	assert.NoError(t, err)

	resolver, err := defaultSchemaLoader(spec, nil, nil)
	assert.NoError(t, err)

	schema := spec.Definitions["car"]
	oldBrand := schema.Properties["brand"]
	assert.NotEmpty(t, oldBrand.Ref.String())
	assert.NotEqual(t, spec.Definitions["brand"], oldBrand)

	err = expandSchema(&schema, "#/definitions/car", resolver)
	assert.NoError(t, err)

	newBrand := schema.Properties["brand"]
	assert.Empty(t, newBrand.Ref.String())
	assert.Equal(t, spec.Definitions["brand"], newBrand)

	schema = spec.Definitions["truck"]
	assert.NotEmpty(t, schema.Ref.String())

	err = expandSchema(&schema, "#/definitions/truck", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.Ref.String())
	assert.Equal(t, spec.Definitions["car"], schema)

	sch := new(Schema)
	err = expandSchema(sch, "", resolver)
	assert.NoError(t, err)

	schema = spec.Definitions["batch"]
	err = expandSchema(&schema, "#/definitions/batch", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.Items.Schema.Ref.String())
	assert.Equal(t, *schema.Items.Schema, spec.Definitions["brand"])

	schema = spec.Definitions["batch2"]
	err = expandSchema(&schema, "#/definitions/batch2", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.Items.Schemas[0].Ref.String())
	assert.Empty(t, schema.Items.Schemas[1].Ref.String())
	assert.Equal(t, schema.Items.Schemas[0], spec.Definitions["brand"])
	assert.Equal(t, schema.Items.Schemas[1], spec.Definitions["tag"])

	schema = spec.Definitions["allofBoth"]
	err = expandSchema(&schema, "#/definitions/allofBoth", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.AllOf[0].Ref.String())
	assert.Empty(t, schema.AllOf[1].Ref.String())
	assert.Equal(t, schema.AllOf[0], spec.Definitions["brand"])
	assert.Equal(t, schema.AllOf[1], spec.Definitions["tag"])

	schema = spec.Definitions["anyofBoth"]
	err = expandSchema(&schema, "#/definitions/anyofBoth", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.AnyOf[0].Ref.String())
	assert.Empty(t, schema.AnyOf[1].Ref.String())
	assert.Equal(t, schema.AnyOf[0], spec.Definitions["brand"])
	assert.Equal(t, schema.AnyOf[1], spec.Definitions["tag"])

	schema = spec.Definitions["oneofBoth"]
	err = expandSchema(&schema, "#/definitions/oneofBoth", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.OneOf[0].Ref.String())
	assert.Empty(t, schema.OneOf[1].Ref.String())
	assert.Equal(t, schema.OneOf[0], spec.Definitions["brand"])
	assert.Equal(t, schema.OneOf[1], spec.Definitions["tag"])

	schema = spec.Definitions["notSomething"]
	err = expandSchema(&schema, "#/definitions/notSomething", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.Not.Ref.String())
	assert.Equal(t, *schema.Not, spec.Definitions["tag"])

	schema = spec.Definitions["withAdditional"]
	err = expandSchema(&schema, "#/definitions/withAdditional", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.AdditionalProperties.Schema.Ref.String())
	assert.Equal(t, *schema.AdditionalProperties.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["withAdditionalItems"]
	err = expandSchema(&schema, "#/definitions/withAdditionalItems", resolver)
	assert.NoError(t, err)
	assert.Empty(t, schema.AdditionalItems.Schema.Ref.String())
	assert.Equal(t, *schema.AdditionalItems.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["withPattern"]
	err = expandSchema(&schema, "#/definitions/withPattern", resolver)
	assert.NoError(t, err)
	prop := schema.PatternProperties["^x-ab"]
	assert.Empty(t, prop.Ref.String())
	assert.Equal(t, prop, spec.Definitions["tag"])

	schema = spec.Definitions["deps"]
	err = expandSchema(&schema, "#/definitions/deps", resolver)
	assert.NoError(t, err)
	prop2 := schema.Dependencies["something"]
	assert.Empty(t, prop2.Schema.Ref.String())
	assert.Equal(t, *prop2.Schema, spec.Definitions["tag"])

	schema = spec.Definitions["defined"]
	err = expandSchema(&schema, "#/definitions/defined", resolver)
	assert.NoError(t, err)
	prop = schema.Definitions["something"]
	assert.Empty(t, prop.Ref.String())
	assert.Equal(t, prop, spec.Definitions["tag"])

}

func TestDefaultResolutionCache(t *testing.T) {

	cache := defaultResolutionCache()

	sch, ok := cache.Get("not there")
	assert.False(t, ok)
	assert.Nil(t, sch)

	sch, ok = cache.Get("http://swagger.io/v2/schema.json")
	assert.True(t, ok)
	assert.Equal(t, swaggerSchema, sch)

	sch, ok = cache.Get("http://json-schema.org/draft-04/schema")
	assert.True(t, ok)
	assert.Equal(t, jsonSchema, sch)

	cache.Set("something", "here")
	sch, ok = cache.Get("something")
	assert.True(t, ok)
	assert.Equal(t, "here", sch)
}

func resolutionContextServer() *httptest.Server {
	var servedAt string
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// fmt.Println("got a request for", req.URL.String())
		if req.URL.Path == "/resolution.json" {

			b, _ := ioutil.ReadFile("../fixtures/specs/resolution.json")
			var ctnt map[string]interface{}
			json.Unmarshal(b, &ctnt)
			ctnt["id"] = servedAt

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(200)
			bb, _ := json.Marshal(ctnt)
			rw.Write(bb)
			return
		}
		if req.URL.Path == "/resolution2.json" {
			b, _ := ioutil.ReadFile("../fixtures/specs/resolution2.json")
			var ctnt map[string]interface{}
			json.Unmarshal(b, &ctnt)
			ctnt["id"] = servedAt

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(200)
			bb, _ := json.Marshal(ctnt)
			rw.Write(bb)
			return
		}

		if req.URL.Path == "/boolProp.json" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(200)
			b, _ := json.Marshal(map[string]interface{}{
				"type": "boolean",
			})
			rw.Write(b)
			return
		}

		if req.URL.Path == "/deeper/stringProp.json" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(200)
			b, _ := json.Marshal(map[string]interface{}{
				"type": "string",
			})
			rw.Write(b)
			return
		}

		if req.URL.Path == "/deeper/arrayProp.json" {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(200)
			b, _ := json.Marshal(map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "file",
				},
			})
			rw.Write(b)
			return
		}

		rw.WriteHeader(http.StatusNotFound)
	}))
	servedAt = server.URL
	return server
}

func TestResolveRemoteRef_RootSame(t *testing.T) {
	specs := "../fixtures/specs"
	fileserver := http.FileServer(http.Dir(specs))
	server := httptest.NewServer(fileserver)
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var result Swagger
		ref, _ := NewRef(server.URL + "/refed.json#")
		resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
		if assert.NoError(t, resolver.Resolve(&ref, &result)) {
			assertSpecs(t, result, *rootDoc)
		}
	}
}

func TestResolveRemoteRef_FromFragment(t *testing.T) {
	specs := "../fixtures/specs"
	fileserver := http.FileServer(http.Dir(specs))
	server := httptest.NewServer(fileserver)
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Schema
		ref, err := NewRef(server.URL + "/refed.json#/definitions/pet")
		if assert.NoError(t, err) {
			resolver := &schemaLoader{root: rootDoc, cache: defaultResolutionCache(), loadDoc: swag.JSONDoc}
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, []string{"id", "name"}, tgt.Required)
			}
		}
	}
}

func TestResolveRemoteRef_FromInvalidFragment(t *testing.T) {
	specs := "../fixtures/specs"
	fileserver := http.FileServer(http.Dir(specs))
	server := httptest.NewServer(fileserver)
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Schema
		ref, err := NewRef(server.URL + "/refed.json#/definitions/NotThere")
		if assert.NoError(t, err) {
			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			assert.Error(t, resolver.Resolve(&ref, &tgt))
		}
	}
}

func TestResolveRemoteRef_WithResolutionContext(t *testing.T) {
	server := resolutionContextServer()
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Schema
		ref, err := NewRef(server.URL + "/resolution.json#/definitions/bool")
		if assert.NoError(t, err) {
			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, StringOrArray([]string{"boolean"}), tgt.Type)
			}
		}
	}
}

func TestResolveRemoteRef_WithNestedResolutionContext(t *testing.T) {
	server := resolutionContextServer()
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Schema
		ref, err := NewRef(server.URL + "/resolution.json#/items/items")
		if assert.NoError(t, err) {
			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, StringOrArray([]string{"string"}), tgt.Type)
			}
		}
	}
}

func TestResolveRemoteRef_WithNestedResolutionContextWithFragment(t *testing.T) {
	server := resolutionContextServer()
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Schema
		ref, err := NewRef(server.URL + "/resolution2.json#/items/items")
		if assert.NoError(t, err) {
			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, StringOrArray([]string{"file"}), tgt.Type)
			}
		}
	}
}

func TestResolveRemoteRef_ToParameter(t *testing.T) {
	specs := "../fixtures/specs"
	fileserver := http.FileServer(http.Dir(specs))
	server := httptest.NewServer(fileserver)
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Parameter
		ref, err := NewRef(server.URL + "/refed.json#/parameters/idParam")
		if assert.NoError(t, err) {

			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, "id", tgt.Name)
				assert.Equal(t, "path", tgt.In)
				assert.Equal(t, "ID of pet to fetch", tgt.Description)
				assert.True(t, tgt.Required)
				assert.Equal(t, "integer", tgt.Type)
				assert.Equal(t, "int64", tgt.Format)
			}
		}
	}
}

func TestResolveRemoteRef_ToPathItem(t *testing.T) {
	specs := "../fixtures/specs"
	fileserver := http.FileServer(http.Dir(specs))
	server := httptest.NewServer(fileserver)
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt PathItem
		ref, err := NewRef(server.URL + "/refed.json#/paths/" + jsonpointer.Escape("/pets/{id}"))
		if assert.NoError(t, err) {

			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, rootDoc.Paths.Paths["/pets/{id}"].Get, tgt.Get)
			}
		}
	}
}

func TestResolveRemoteRef_ToResponse(t *testing.T) {
	specs := "../fixtures/specs"
	fileserver := http.FileServer(http.Dir(specs))
	server := httptest.NewServer(fileserver)
	defer server.Close()

	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Response
		ref, err := NewRef(server.URL + "/refed.json#/responses/petResponse")
		if assert.NoError(t, err) {

			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, rootDoc.Responses["petResponse"], tgt)
			}
		}
	}
}

func TestResolveLocalRef_SameRoot(t *testing.T) {
	rootDoc := new(Swagger)
	json.Unmarshal(testingutil.PetStoreJSONMessage, rootDoc)

	result := new(Swagger)
	ref, _ := NewRef("#")
	resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
	err := resolver.Resolve(&ref, result)
	if assert.NoError(t, err) {
		assert.Equal(t, rootDoc, result)
	}
}

func TestResolveLocalRef_FromFragment(t *testing.T) {
	rootDoc := new(Swagger)
	json.Unmarshal(testingutil.PetStoreJSONMessage, rootDoc)

	var tgt Schema
	ref, err := NewRef("#/definitions/Category")
	if assert.NoError(t, err) {
		resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
		err := resolver.Resolve(&ref, &tgt)
		if assert.NoError(t, err) {
			assert.Equal(t, "Category", tgt.ID)
		}
	}
}

func TestResolveLocalRef_FromInvalidFragment(t *testing.T) {
	rootDoc := new(Swagger)
	json.Unmarshal(testingutil.PetStoreJSONMessage, rootDoc)

	var tgt Schema
	ref, err := NewRef("#/definitions/NotThere")
	if assert.NoError(t, err) {
		resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
		err := resolver.Resolve(&ref, &tgt)
		assert.Error(t, err)
	}
}

func TestResolveLocalRef_Parameter(t *testing.T) {
	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Parameter
		ref, err := NewRef("#/parameters/idParam")
		if assert.NoError(t, err) {
			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, "id", tgt.Name)
				assert.Equal(t, "path", tgt.In)
				assert.Equal(t, "ID of pet to fetch", tgt.Description)
				assert.True(t, tgt.Required)
				assert.Equal(t, "integer", tgt.Type)
				assert.Equal(t, "int64", tgt.Format)
			}
		}
	}
}

func TestResolveLocalRef_PathItem(t *testing.T) {
	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt PathItem
		ref, err := NewRef("#/paths/" + jsonpointer.Escape("/pets/{id}"))
		if assert.NoError(t, err) {
			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, rootDoc.Paths.Paths["/pets/{id}"].Get, tgt.Get)
			}
		}
	}
}

func TestResolveLocalRef_Response(t *testing.T) {
	rootDoc := new(Swagger)
	b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
	if assert.NoError(t, err) && assert.NoError(t, json.Unmarshal(b, rootDoc)) {
		var tgt Response
		ref, err := NewRef("#/responses/petResponse")
		if assert.NoError(t, err) {
			resolver, _ := defaultSchemaLoader(rootDoc, nil, nil)
			if assert.NoError(t, resolver.Resolve(&ref, &tgt)) {
				assert.Equal(t, rootDoc.Responses["petResponse"], tgt)
			}
		}
	}
}

package spec

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casualjim/go-swagger/jsonpointer"
	testingutil "github.com/casualjim/go-swagger/testing"
	"github.com/casualjim/go-swagger/util"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestSchemaExpansion(t *testing.T) {
	// carsDoc, err := util.JSONDoc("../fixtures/expansion/schemas1.json")
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

// func compareSpecs(actual, spec)

func TestResolveRemoteRef(t *testing.T) {
	specs := "../fixtures/specs"
	fileserver := http.FileServer(http.Dir(specs))

	Convey("resolving a remote ref", t, func() {
		server := httptest.NewServer(fileserver)
		Reset(func() {
			server.Close()
		})

		Convey("in a swagger spec", func() {
			rootDoc := new(Swagger)
			b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
			So(err, ShouldBeNil)
			json.Unmarshal(b, rootDoc)

			Convey("resolves root to same schema", func() {
				var result Swagger
				ref, _ := NewRef(server.URL + "/refed.json#")
				resolver, _ := defaultSchemaLoader(rootDoc, nil)
				err = resolver.Resolve(&ref, &result)
				So(err, ShouldBeNil)
				compareSpecs(result, *rootDoc)
			})

			Convey("to a schema", func() {

				Convey("from a fragment", func() {
					var tgt Schema
					ref, err := NewRef(server.URL + "/refed.json#/definitions/pet")
					So(err, ShouldBeNil)
					resolver := &schemaLoader{root: rootDoc, cache: defaultResolutionCache(), loadDoc: util.JSONDoc}
					err = resolver.Resolve(&ref, &tgt)
					So(err, ShouldBeNil)
					So(tgt.Required, ShouldResemble, []string{"id", "name"})
				})

				Convey("from an invalid fragment", func() {
					var tgt Schema
					ref, err := NewRef(server.URL + "/refed.json#/definitions/NotThere")
					So(err, ShouldBeNil)

					resolver, _ := defaultSchemaLoader(rootDoc, nil)
					err = resolver.Resolve(&ref, &tgt)
					So(err, ShouldNotBeNil)
				})

				Convey("with a resolution context", func() {
					server.Close()
					server = resolutionContextServer()
					var tgt Schema
					ref, err := NewRef(server.URL + "/resolution.json#/definitions/bool")
					So(err, ShouldBeNil)

					resolver, _ := defaultSchemaLoader(rootDoc, nil)
					err = resolver.Resolve(&ref, &tgt)
					So(err, ShouldBeNil)
					So(tgt.Type, ShouldResemble, StringOrArray([]string{"boolean"}))
				})

				Convey("with a nested resolution context", func() {
					server.Close()
					server = resolutionContextServer()
					var tgt Schema
					ref, err := NewRef(server.URL + "/resolution.json#/items/items")
					So(err, ShouldBeNil)

					resolver, _ := defaultSchemaLoader(rootDoc, nil)
					// pretty.Println(resolver)
					err = resolver.Resolve(&ref, &tgt)
					So(err, ShouldBeNil)
					So(tgt.Type, ShouldResemble, StringOrArray([]string{"string"}))
				})

				Convey("with a nested resolution context with a fragment", func() {
					server.Close()
					server = resolutionContextServer()
					var tgt Schema
					ref, err := NewRef(server.URL + "/resolution2.json#/items/items")
					So(err, ShouldBeNil)

					resolver, _ := defaultSchemaLoader(rootDoc, nil)
					// pretty.Println(resolver)
					err = resolver.Resolve(&ref, &tgt)
					So(err, ShouldBeNil)
					So(tgt.Type, ShouldResemble, StringOrArray([]string{"file"}))
				})
			})

			Convey("to a parameter", func() {
				var tgt Parameter
				ref, err := NewRef(server.URL + "/refed.json#/parameters/idParam")
				So(err, ShouldBeNil)

				resolver, _ := defaultSchemaLoader(rootDoc, nil)
				err = resolver.Resolve(&ref, &tgt)
				So(err, ShouldBeNil)
				So(tgt.Name, ShouldEqual, "id")
				So(tgt.In, ShouldEqual, "path")
				So(tgt.Description, ShouldEqual, "ID of pet to fetch")
				So(tgt.Required, ShouldBeTrue)
				So(tgt.Type, ShouldEqual, "integer")
				So(tgt.Format, ShouldEqual, "int64")
			})

			Convey("to a path item object", func() {
				var tgt PathItem
				ref, err := NewRef(server.URL + "/refed.json#/paths/" + jsonpointer.Escape("/pets/{id}"))
				So(err, ShouldBeNil)

				resolver, _ := defaultSchemaLoader(rootDoc, nil)
				err = resolver.Resolve(&ref, &tgt)
				So(err, ShouldBeNil)
				So(tgt.Get, ShouldResemble, rootDoc.Paths.Paths["/pets/{id}"].Get)
			})

			Convey("to a response object", func() {
				var tgt Response
				ref, err := NewRef(server.URL + "/refed.json#/responses/petResponse")
				So(err, ShouldBeNil)

				resolver, _ := defaultSchemaLoader(rootDoc, nil)
				err = resolver.Resolve(&ref, &tgt)
				So(err, ShouldBeNil)
				So(tgt, ShouldResemble, rootDoc.Responses["petResponse"])
			})
		})
	})

}

func TestResolveLocalRef(t *testing.T) {
	rootDoc := new(Swagger)
	json.Unmarshal(testingutil.PetStoreJSONMessage, rootDoc)

	Convey("resolving local a ref", t, func() {

		Convey("in a swagger spec", func() {

			Convey("to a schema", func() {

				Convey("resolves root to same ptr instance", func() {
					var result interface{}
					ref, _ := NewRef("#")
					resolver, _ := defaultSchemaLoader(rootDoc, nil)
					err := resolver.Resolve(&ref, &result)
					So(err, ShouldBeNil)
					So(result, ShouldEqual, rootDoc)
				})

				Convey("from a fragment", func() {
					var tgt Schema
					ref, err := NewRef("#/definitions/Category")
					So(err, ShouldBeNil)

					resolver, _ := defaultSchemaLoader(rootDoc, nil)
					err = resolver.Resolve(&ref, &tgt)
					So(err, ShouldBeNil)
					So(tgt.ID, ShouldEqual, "Category")
				})

				Convey("from an invalid fragment", func() {
					var tgt Schema
					ref, err := NewRef("#/definitions/NotThere")
					So(err, ShouldBeNil)

					resolver, _ := defaultSchemaLoader(rootDoc, nil)
					err = resolver.Resolve(&ref, &tgt)
					So(err, ShouldNotBeNil)
				})

			})

			Convey("to a parameter", func() {
				rootDoc = new(Swagger)
				b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
				So(err, ShouldBeNil)
				json.Unmarshal(b, rootDoc)

				var tgt Parameter
				ref, err := NewRef("#/parameters/idParam")
				So(err, ShouldBeNil)

				resolver, _ := defaultSchemaLoader(rootDoc, nil)
				err = resolver.Resolve(&ref, &tgt)
				So(err, ShouldBeNil)
				So(tgt.Name, ShouldEqual, "id")
				So(tgt.In, ShouldEqual, "path")
				So(tgt.Description, ShouldEqual, "ID of pet to fetch")
				So(tgt.Required, ShouldBeTrue)
				So(tgt.Type, ShouldEqual, "integer")
				So(tgt.Format, ShouldEqual, "int64")
			})

			Convey("to a path item object", func() {
				rootDoc = new(Swagger)
				b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
				So(err, ShouldBeNil)
				json.Unmarshal(b, rootDoc)

				var tgt PathItem
				ref, err := NewRef("#/paths/" + jsonpointer.Escape("/pets/{id}"))
				So(err, ShouldBeNil)

				resolver, _ := defaultSchemaLoader(rootDoc, nil)
				err = resolver.Resolve(&ref, &tgt)
				So(err, ShouldBeNil)
				So(tgt.Get, ShouldEqual, rootDoc.Paths.Paths["/pets/{id}"].Get)
			})

			Convey("to a response object", func() {
				rootDoc = new(Swagger)
				b, err := ioutil.ReadFile("../fixtures/specs/refed.json")
				So(err, ShouldBeNil)
				json.Unmarshal(b, rootDoc)

				var tgt Response
				ref, err := NewRef("#/responses/petResponse")
				So(err, ShouldBeNil)

				resolver, _ := defaultSchemaLoader(rootDoc, nil)
				err = resolver.Resolve(&ref, &tgt)
				So(err, ShouldBeNil)
				So(tgt, ShouldResemble, rootDoc.Responses["petResponse"])
			})

		})
	})

}

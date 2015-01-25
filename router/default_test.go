package router

import (
	"testing"

	swagger_api "github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/testing/petstore"
	"github.com/stretchr/testify/assert"
)

func TestRouterBuilder(t *testing.T) {
	spec, api := petstore.NewAPI(t)

	assert.Len(t, spec.RequiredConsumes(), 3)
	assert.Len(t, spec.RequiredProduces(), 5)
	assert.Len(t, spec.OperationIDs(), 4)

	// context := NewContext(spec, api)
	builder := petAPIRouterBuilder(spec, api)
	getRecords := builder.records["GET"]
	postRecords := builder.records["POST"]
	deleteRecords := builder.records["DELETE"]

	assert.Len(t, getRecords, 2)
	assert.Len(t, postRecords, 1)
	assert.Len(t, deleteRecords, 1)

	assert.Empty(t, builder.records["PATCH"])
	assert.Empty(t, builder.records["OPTIONS"])
	assert.Empty(t, builder.records["HEAD"])
	assert.Empty(t, builder.records["PUT"])

	rec := postRecords[0]
	assert.Equal(t, rec.Key, "/pets")
	val := rec.Value.(*routeEntry)
	assert.Len(t, val.Consumers, 3)
	assert.Len(t, val.Producers, 5)
	assert.Len(t, val.Consumes, 3)
	assert.Len(t, val.Produces, 5)

	assert.Len(t, val.Parameters, 1)
}

func TestRouterStruct(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	router := Default(spec, api)

	methods := router.OtherMethods("post", "/pets/{id}")
	assert.Len(t, methods, 2)

	entry, ok := router.Lookup("delete", "/pets/{id}")
	assert.True(t, ok)
	assert.NotNil(t, entry)
	assert.Len(t, entry.Params, 1)
	assert.Equal(t, "id", entry.Params[0].Name)

	_, ok = router.Lookup("delete", "/pets")
	assert.False(t, ok)

	_, ok = router.Lookup("post", "/no-pets")
	assert.False(t, ok)
}

func petAPIRouterBuilder(spec *spec.Document, api *swagger_api.API) *defaultRouteBuilder {
	builder := newDefaultRouteBuilder(spec, api)
	builder.AddRoute("GET", "/pets", spec.AllPaths()["/pets"].Get)
	builder.AddRoute("POST", "/pets", spec.AllPaths()["/pets"].Post)
	builder.AddRoute("DELETE", "/pets/{id}", spec.AllPaths()["/pets/{id}"].Delete)
	builder.AddRoute("GET", "/pets/{id}", spec.AllPaths()["/pets/{id}"].Get)

	return builder
}

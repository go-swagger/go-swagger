package middleware

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/go-swagger/go-swagger/httpkit/middleware/untyped"
	"github.com/go-swagger/go-swagger/internal/testing/petstore"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/stretchr/testify/assert"
)

func terminator(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func TestRouterMiddleware(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	context := NewContext(spec, api, nil)
	mw := newRouter(context, http.HandlerFunc(terminator))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/api/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/api/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)

	methods := strings.Split(recorder.Header().Get("Allow"), ",")
	sort.Sort(sort.StringSlice(methods))
	assert.Equal(t, "GET,POST", strings.Join(methods, ","))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/nopets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code)

	spec, api = petstore.NewRootAPI(t)
	context = NewContext(spec, api, nil)
	mw = newRouter(context, http.HandlerFunc(terminator))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)

	methods = strings.Split(recorder.Header().Get("Allow"), ",")
	sort.Sort(sort.StringSlice(methods))
	assert.Equal(t, "GET,POST", strings.Join(methods, ","))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/nopets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code)

}

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
	router := DefaultRouter(spec, newRoutableUntypedAPI(spec, api, new(Context)))

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

func petAPIRouterBuilder(spec *spec.Document, api *untyped.API) *defaultRouteBuilder {
	builder := newDefaultRouteBuilder(spec, newRoutableUntypedAPI(spec, api, new(Context)))
	builder.AddRoute("GET", "/pets", spec.AllPaths()["/pets"].Get)
	builder.AddRoute("POST", "/pets", spec.AllPaths()["/pets"].Post)
	builder.AddRoute("DELETE", "/pets/{id}", spec.AllPaths()["/pets/{id}"].Delete)
	builder.AddRoute("GET", "/pets/{id}", spec.AllPaths()["/pets/{id}"].Get)

	return builder
}

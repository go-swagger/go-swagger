package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/testing/petstore"
	"github.com/stretchr/testify/assert"
)

func TestContentTypeValidation(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	context := NewContext(spec, api, nil)
	mw := context.ValidationMiddleware(http.HandlerFunc(terminator))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://localhost:8080/api/pets", nil)
	request.Header.Add("Accept", "*/*")
	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "http://localhost:8080/api/pets", nil)
	request.Header.Add("content-type", "application(")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 400, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("content-type"))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "http://localhost:8080/api/pets", nil)
	request.Header.Add("Accept", "*/*")
	request.Header.Add("content-type", "text/html")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 415, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
}

func TestResponseFormatValidation(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	context := NewContext(spec, api, nil)
	mw := context.ValidationMiddleware(http.HandlerFunc(terminator))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "http://localhost:8080/api/pets", nil)
	request.Header.Set(httputils.HeaderContentType, "application/json")
	request.Header.Set(httputils.HeaderAccept, "application/json")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "http://localhost:8080/api/pets", nil)
	request.Header.Set(httputils.HeaderContentType, "application/json")
	request.Header.Set(httputils.HeaderAccept, "application/sml")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotAcceptable, recorder.Code)
}

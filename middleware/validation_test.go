package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/httputils"
	"github.com/casualjim/go-swagger/testing/petstore"
	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	result := new(result)
	assert.True(t, result.IsValid())
	assert.False(t, result.HasErrors())
	result.AddErrors(errors.New(400, "yada"))
	assert.Len(t, result.Errors, 1)
	assert.True(t, result.HasErrors())
	assert.False(t, result.IsValid())
}

func TestContentTypeValidation(t *testing.T) {
	context := Serve(petstore.NewAPI(t))
	mw := newValidation(context, http.HandlerFunc(terminator))

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
	context := Serve(petstore.NewAPI(t))
	mw := newValidation(context, http.HandlerFunc(terminator))

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

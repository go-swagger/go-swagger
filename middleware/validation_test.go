package middleware

import (
	"bytes"
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
	mw := newValidation(context, http.HandlerFunc(terminator))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/pets", nil)
	request.Header.Add("Accept", "*/*")
	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/pets", nil)
	request.Header.Add("content-type", "application(")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 400, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("content-type"))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/pets", nil)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("content-type", "text/html")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 415, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("content-type"))
}

func TestResponseFormatValidation(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	context := NewContext(spec, api, nil)
	mw := newValidation(context, http.HandlerFunc(terminator))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/pets", bytes.NewBuffer([]byte(`{"name":"Dog"}`)))
	request.Header.Set(httputils.HeaderContentType, "application/json")
	request.Header.Set(httputils.HeaderAccept, "application/json")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code, recorder.Body.String())

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/pets", bytes.NewBuffer([]byte(`{"name":"Dog"}`)))
	request.Header.Set(httputils.HeaderContentType, "application/json")
	request.Header.Set(httputils.HeaderAccept, "application/sml")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotAcceptable, recorder.Code)
}

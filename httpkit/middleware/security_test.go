package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-swagger/go-swagger/internal/testing/petstore"
	"github.com/stretchr/testify/assert"
)

func TestSecurityMiddleware(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	context := NewContext(spec, api, nil)
	context.router = DefaultRouter(spec, context.api)
	mw := newSecureAPI(context, http.HandlerFunc(terminator))

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/pets", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 401, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets", nil)
	request.SetBasicAuth("admin", "wrong")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 401, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets", nil)
	request.SetBasicAuth("admin", "admin")

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/pets/1", nil)

	mw.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

}

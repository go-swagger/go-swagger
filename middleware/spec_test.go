package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casualjim/go-swagger/httpkit"
	"github.com/casualjim/go-swagger/internal/testing/petstore"
	"github.com/stretchr/testify/assert"
)

func TestServeSpecMiddleware(t *testing.T) {
	spec, api := petstore.NewAPI(t)
	ctx := NewContext(spec, api, nil)

	handler := specMiddleware(ctx, nil)
	// serves spec
	request, _ := http.NewRequest("GET", "/swagger.json", nil)
	request.Header.Add(httpkit.HeaderContentType, httpkit.JSONMime)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

	// returns 404 when no next handler
	request, _ = http.NewRequest("GET", "/api/pets", nil)
	request.Header.Add(httpkit.HeaderContentType, httpkit.JSONMime)
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, 404, recorder.Code)

	// forwards to next handler for other url
	handler = specMiddleware(ctx, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	request, _ = http.NewRequest("GET", "/api/pets", nil)
	request.Header.Add(httpkit.HeaderContentType, httpkit.JSONMime)
	recorder = httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	assert.Equal(t, 200, recorder.Code)

}

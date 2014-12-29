package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casualjim/go-swagger/swagger/testing/petstore"
	"github.com/stretchr/testify/assert"
)

func terminator(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func TestRouterMiddleware(t *testing.T) {
	context := NewContext(petstore.NewAPI(t))
	mw := NewRouter(context)

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://localhost:8080/api/pets", nil)

	mw(recorder, request, terminator)
	assert.Equal(t, 200, recorder.Code)

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "http://localhost:8080/api/pets", nil)

	mw(recorder, request, terminator)
	assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	assert.Equal(t, "GET,POST", recorder.Header().Get("Allow"))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "http://localhost:8080/pets", nil)

	mw(recorder, request, terminator)
	assert.Equal(t, http.StatusNotFound, recorder.Code)

}

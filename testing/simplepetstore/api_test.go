package simplepetstore

import (
	"net/http/httptest"
	"testing"

	"github.com/casualjim/go-swagger/httputils"
	"github.com/stretchr/testify/assert"
)

func TestSimplePetstoreSpec(t *testing.T) {
	handler, _ := NewPetstore()
	// Serves swagger spec document
	r, _ := httputils.JSONRequest("GET", "/swagger.json", nil)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, r)
	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, swaggerJSON, rw.Body.String())
}

func TestSimplePetstoreAllPets(t *testing.T) {
	handler, _ := NewPetstore()
	// Serves swagger spec document
	r, _ := httputils.JSONRequest("GET", "/api/pets", nil)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, r)
	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, "[{\"id\":1,\"name\":\"Dog\",\"status\":\"available\"},{\"id\":2,\"name\":\"Cat\",\"status\":\"pending\"}]\n", rw.Body.String())
}

func TestSimplePetstorePetByID(t *testing.T) {
	handler, _ := NewPetstore()

	// Serves swagger spec document
	r, _ := httputils.JSONRequest("GET", "/api/pets/1", nil)
	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, r)
	assert.Equal(t, 200, rw.Code)
	assert.Equal(t, "{\"id\":1,\"name\":\"Dog\",\"status\":\"available\"}\n", rw.Body.String())
}

// func TestSimplePetstoreAddPet(t *testing.T) {
// 	handler, _ := NewPetstore()
// 	// Serves swagger spec document
// 	r, _ := httputils.JSONRequest("POST", "/api/pets", nil)
// 	rw := httptest.NewRecorder()
// 	handler.ServeHTTP(rw, r)
// 	assert.Equal(t, 200, rw.Code)
// 	assert.Equal(t, swaggerJSON, rw.Body.String())
// }

// func TestSimplePetstoreDeletePet(t *testing.T) {
// 	handler, _ := NewPetstore()
// 	// Serves swagger spec document
// 	r, _ := httputils.JSONRequest("DELETE", "/api/pets/1", nil)
// 	rw := httptest.NewRecorder()
// 	handler.ServeHTTP(rw, r)
// 	assert.Equal(t, 204, rw.Code)
// 	assert.Equal(t, swaggerJSON, rw.Body.String())
// }

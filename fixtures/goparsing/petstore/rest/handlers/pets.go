package handlers

import (
	"net/http"

	"github.com/naoina/denco"
)

// GetPets +swagger:route GET /pets pets listPets
//
// Lists the pets known to the store.
//
// Responses:
// default: genericError
// 200: []pet
func GetPets(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// GetPetByID http handler.
//
// +swagger:route GET /pets/{id} pets getPetById
//
// Gets the details for a pet
func GetPetByID(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// CreatePet http handler.
//
// +swagger:route POST /pets pets createPet
//
// Creates a new pet in the store.
func CreatePet(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// UpdatePet http handler.
//
// +swagger:route GET /pets/{id} pets updatePet
//
// Updates the details for a pet.
func UpdatePet(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// DeletePet http handler.
//
// +swagger:route DELETE /pets/{id} pets deletePet
//
// Deletes a pet from the store.
func DeletePet(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

package handlers

import (
	"net/http"

	"github.com/naoina/denco"
)

// GetPets handles the list pets operation
func GetPets(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// GetPetByID gets a pet by id
func GetPetByID(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// CreatePet creates a pet
func CreatePet(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// UpdatePet updates a pet
func UpdatePet(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

// DeletePet deletes a pet from the store
func DeletePet(w http.ResponseWriter, r *http.Request, params denco.Params) {
	// some actual stuff should happen in here
}

package api

import (
	"net/http"
	"sync"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/middleware"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/testing"
)

// NewPetstore creates a new petstore api handler
func NewPetstore() (http.Handler, error) {
	spec, err := spec.New(testing.PetStoreJSONMessage, "")
	if err != nil {
		return nil, err
	}
	api := swagger.NewAPI(spec)

	api.RegisterOperation("getAllPets", getAllPets)
	api.RegisterOperation("createPet", createPet)
	api.RegisterOperation("deletePet", deletePet)
	api.RegisterOperation("getPetById", getPetByID)

	return middleware.Serve(spec, api), nil
}

var getAllPets = swagger.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	return pets, nil
})
var createPet = swagger.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	body := data.(map[string]interface{})["pet"]
	// var pet Pet
	// reflection.UnmarshalMap(body.(map[string]interface{}), &pet)
	// addPet(pet)
	return body, nil
})

var deletePet = swagger.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	id := data.(map[string]interface{})["id"].(int64)
	removePet(id)
	return nil, nil
})

var getPetByID = swagger.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	id := data.(map[string]interface{})["id"].(int64)
	return petByID(id)
})

// Tag the tag model
type Tag struct {
	ID   int64
	Name string
}

// Pet the pet model
type Pet struct {
	ID        int64
	Name      string
	PhotoURLs []string
	Status    string
	Tags      []Tag
}

var pets = []Pet{
	{1, "Dog", []string{}, "available", nil},
	{2, "Cat", []string{}, "pending", nil},
}

var petsLock = &sync.Mutex{}

func addPet(pet Pet) {
	petsLock.Lock()
	defer petsLock.Unlock()
	pets = append(pets, pet)
}

func removePet(id int64) {
	petsLock.Lock()
	defer petsLock.Unlock()

}

func petByID(id int64) (*Pet, error) {
	return nil, nil
}

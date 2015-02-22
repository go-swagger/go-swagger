package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/middleware"
	"github.com/casualjim/go-swagger/middleware/untyped"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
	"github.com/kr/pretty"
)

// NewPetstore creates a new petstore api handler
func NewPetstore() (http.Handler, error) {
	spec, err := spec.New(json.RawMessage([]byte(swaggerJSON)), "")
	if err != nil {
		return nil, err
	}
	api := untyped.NewAPI(spec)

	api.RegisterOperation("getAllPets", getAllPets)
	api.RegisterOperation("createPet", createPet)
	api.RegisterOperation("deletePet", deletePet)
	api.RegisterOperation("getPetById", getPetByID)

	return middleware.ServeWithUI(spec, api), nil
}

var getAllPets = swagger.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	fmt.Println("getAllPets")
	pretty.Println(data)
	return pets, nil
})
var createPet = swagger.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	fmt.Println("createPet")
	pretty.Println(data)
	body := data.(map[string]interface{})["pet"]
	var pet Pet
	if err := util.FromDynamicJSON(body, &pet); err != nil {
		return nil, err
	}
	addPet(pet)
	return body, nil
})

var deletePet = swagger.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	fmt.Println("deletePet")
	pretty.Println(data)
	id := data.(map[string]interface{})["id"].(int64)
	removePet(id)
	return nil, nil
})

var getPetByID = swagger.OperationHandlerFunc(func(data interface{}) (interface{}, error) {
	fmt.Println("getPetByID")
	pretty.Println(data)
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
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	PhotoURLs []string `json:"photoUrls,omitempty"`
	Status    string   `json:"status,omitempty"`
	Tags      []Tag    `json:"tags,omitempty"`
}

var pets = []Pet{
	{1, "Dog", []string{}, "available", nil},
	{2, "Cat", []string{}, "pending", nil},
}

var petsLock = &sync.Mutex{}
var lastPetID int64 = 2

func newPetID() int64 {
	return atomic.AddInt64(&lastPetID, 1)
}

func addPet(pet Pet) {
	petsLock.Lock()
	defer petsLock.Unlock()
	pet.ID = newPetID()
	pets = append(pets, pet)
}

func removePet(id int64) {
	petsLock.Lock()
	defer petsLock.Unlock()
	var newPets []Pet
	for _, pet := range pets {
		if pet.ID != id {
			newPets = append(newPets, pet)
		}
	}
	pets = newPets
}

func petByID(id int64) (*Pet, error) {
	for _, pet := range pets {
		if pet.ID == id {
			return &pet, nil
		}
	}
	return nil, errors.NotFound("not found: pet %d", id)
}

var swaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "version": "1.0.0",
    "title": "Swagger Petstore",
    "description": "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification",
    "termsOfService": "http://helloreverb.com/terms/",
    "contact": {
      "name": "Wordnik API Team"
    },
    "license": {
      "name": "MIT"
    }
  },
  "host": "localhost:8344",
  "basePath": "/api",
  "schemes": [
    "http"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/pets": {
      "get": {
        "description": "Returns all pets from the system that the user has access to",
        "operationId": "findPets",
        "produces": [
          "application/json",
          "application/xml",
          "text/xml",
          "text/html"
        ],
        "parameters": [
          {
            "name": "tags",
            "in": "query",
            "description": "tags to filter by",
            "required": false,
            "type": "array",
            "items": {
              "type": "string"
            },
            "collectionFormat": "csv"
          },
          {
            "name": "limit",
            "in": "query",
            "description": "maximum number of results to return",
            "required": false,
            "type": "integer",
            "format": "int32"
          }
        ],
        "responses": {
          "200": {
            "description": "pet response",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/pet"
              }
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/errorModel"
            }
          }
        }
      },
      "post": {
        "description": "Creates a new pet in the store.  Duplicates are allowed",
        "operationId": "addPet",
        "produces": [
          "application/json"
        ],
        "parameters": [
          {
            "name": "pet",
            "in": "body",
            "description": "Pet to add to the store",
            "required": true,
            "schema": {
              "$ref": "#/definitions/petInput"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "pet response",
            "schema": {
              "$ref": "#/definitions/pet"
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/errorModel"
            }
          }
        }
      }
    },
    "/pets/{id}": {
      "get": {
        "description": "Returns a user based on a single ID, if the user does not have access to the pet",
        "operationId": "findPetById",
        "produces": [
          "application/json",
          "application/xml",
          "text/xml",
          "text/html"
        ],
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "description": "ID of pet to fetch",
            "required": true,
            "type": "integer",
            "format": "int64"
          }
        ],
        "responses": {
          "200": {
            "description": "pet response",
            "schema": {
              "$ref": "#/definitions/pet"
            }
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/errorModel"
            }
          }
        }
      },
      "delete": {
        "description": "deletes a single pet based on the ID supplied",
        "operationId": "deletePet",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "description": "ID of pet to delete",
            "required": true,
            "type": "integer",
            "format": "int64"
          }
        ],
        "responses": {
          "204": {
            "description": "pet deleted"
          },
          "default": {
            "description": "unexpected error",
            "schema": {
              "$ref": "#/definitions/errorModel"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "pet": {
      "required": [
        "id",
        "name"
      ],
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        }
      }
    },
    "petInput": {
      "allOf": [
        {
          "$ref": "pet"
        },
        {
          "required": [
            "name"
          ],
          "properties": {
            "id": {
              "type": "integer",
              "format": "int64"
            }
          }
        }
      ]
    },
    "errorModel": {
      "required": [
        "code",
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`

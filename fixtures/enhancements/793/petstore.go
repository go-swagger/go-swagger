// Package petstore API
//
// The purpose of this application is to provie an application
// that is using plain go to define an API.
//
// This should demonstrate many comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http
//     Host: petstore.swagger.wordnik.com
//     BasePath: /api
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: John Doe<john.doe@example.com> http://john.doe.com
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//
// swagger:meta
package petstore

// NewPet represents a new pet within this application
//
// A new pet is preferable a puppy.
//
// swagger:model newPet
type NewPet struct {
	// the name for this pet
	// required: true
	Name string `json:"name"`
}

// Pet represents a pet within this application
//
// A pet is preferable cute and fluffy.
//
// Examples of typical pets are dogs and cats.
//
// swagger:model pet
type Pet struct {
	// swager:allOf
	NewPet
	// the id for this pet
	// required: true
	Identifier int64 `json:"id"`
}

// ErrorModel represents a generic error in this application
//
// An error will typically return this model.
//
// The code and message of the error,
// both get returned.
//
// swagger:model errorModel
type ErrorModel struct {
	// the code for this error
	// required: true
	Code int32 `json:"code"`
	// the message for this error
	// required: true
	Message string `json:"message"`
}

// ServeAPI serves the API for this record store
func ServeAPI(host, basePath string, schemes []string) (err error) {
	// swagger:operation GET /pets getPet
	//
	// Returns all pets from the system that the user has access to
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	// - text/xml
	// - text/html
	// parameters:
	// - name: tags
	//   in: query
	//   description: tags to filter by
	//   required: false
	//   type: array
	//   items:
	//     type: string
	//   collectionFormat: csv
	// - name: limit
	//   in: query
	//   description: maximum number of results to return
	//   required: false
	//   type: integer
	//   format: int32
	// responses:
	//   '200':
	//     description: pet response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/pet"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorModel"
	mountItem("GET", basePath+"/pets", nil)

	// swagger:operation POST /pets addPet
	//
	// Creates a new pet in the store.
	// Duplicates are allowed
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: pet
	//   in: body
	//   description: Pet to add to the store
	//   required: true
	//   schema:
	//       "$ref": "#/definitions/newPet"
	//
	// responses:
	//   '200':
	//     description: pet response
	//     schema:
	//       "$ref": "#/definitions/pet"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorModel"
	mountItem("POST", basePath+"/pets", nil)

	// swagger:operation GET /pets/{id} findPetById
	//
	// Returns a user based on a single ID,
	// if the user does not have access to the pet
	//
	// ---
	// produces:
	// - application/json
	// - application/xml
	// - text/xml
	// - text/html
	// parameters:
	// - name: id
	//   in: path
	//   description: ID of pet to fetch
	//   required: true
	//   type: integer
	//   format: int64
	// responses:
	//   '200':
	//     description: pet response
	//     schema:
	//       "$ref": "#/definitions/pet"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorModel"
	mountItem("GET", basePath+"/pets/{id}", nil)

	// swagger:operation DELETE /pets/{id} deletePet
	//
	// deletes a single pet based on the ID supplied
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   description: ID of pet to delete
	//   required: true
	//   type: integer
	//   format: int64
	// responses:
	//   '204':
	//     description: pet deleted
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorModel"
	mountItem("DELETE", basePath+"/pets/{id}", nil)

	return
}

// not really used but I need a method to decorate the calls to
func mountItem(method, path string, handler interface{}) {}

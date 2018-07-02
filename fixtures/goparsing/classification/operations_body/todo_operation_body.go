// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package operations

// ListPetParams the params for the list pets query
type ListPetParams struct {
	// OutOfStock when set to true only the pets that are out of stock will be returned
	OutOfStock bool
}

// ServeAPI serves the API for this record store
func ServeAPI(host, basePath string, schemes []string) error {

	// swagger:route GET /pets pets users listPets
	//
	// Lists pets filtered by some parameters.
	//
	// This will show all available pets by default.
	// You can get the pets that are out of stock
	//
	// Consumes:
	// application/json
	// application/x-protobuf
	//
	// Produces:
	// application/json
	// application/x-protobuf
	//
	// Schemes: http, https, ws, wss
	//
	// Security:
	// api_key:
	// oauth: read, write
	//
	// Responses:
	// default: body:genericError
	// 200: body:someResponse
	// 422: body:validationError
	mountItem("GET", basePath+"/pets", nil)

	/* swagger:route POST /pets pets users createPet

	Create a pet based on the parameters.

	Consumes:
	- application/json
	- application/x-protobuf

	Produces:
	- application/json
	- application/x-protobuf

	Schemes: http, https, ws, wss

	Parameters:
	+ name:        request
	  description: The request model.
	  in:          body
	  type:        petModel
	  unknown:     invalid key that will not get parsed. Added to increase coverage.
	+ name:        id
	  description: The pet id
	  in:          path
	  required:    true
	  allowEmpty:  false

	Responses:
	default: body:genericError
	200: body:someResponse
	422: body:validationError

	Security:
	api_key:
	oauth: read, write */
	mountItem("POST", basePath+"/pets", nil)

	// swagger:route GET /orders orders listOrders
	//
	// lists orders filtered by some parameters.
	//
	// Consumes:
	// application/json
	// application/x-protobuf
	//
	// Produces:
	// application/json
	// application/x-protobuf
	//
	// Schemes: http, https, ws, wss
	//
	// Security:
	// api_key:
	// oauth: orders:read, https://www.googleapis.com/auth/userinfo.email
	//
	// Parameters:
	//
	// Responses:
	// default: body:genericError
	// 200: body:someResponse
	// 422: body:validationError
	mountItem("GET", basePath+"/orders", nil)

	// swagger:route POST /orders orders createOrder
	//
	// create an order based on the parameters.
	//
	// Consumes:
	// application/json
	// application/x-protobuf
	//
	// Produces:
	// application/json
	// application/x-protobuf
	//
	// Schemes: http, https, ws, wss
	//
	// Security:
	// api_key:
	// oauth: read, write
	//
	// Parameters:
	// + name:        id
	//   description: The order id
	//   in:          invalidIn
	//   required:    false
	//   allowEmpty:  true
	//   noValue  (to increase coverage, line without colon, split result will be 1)
	// + name:        request
	//   description: The request model.
	//   in:          body
	//   type:        orderModel
	//
	// Responses:
	// default: body:genericError
	// 200: body:someResponse
	// 422: body:validationError
	mountItem("POST", basePath+"/orders", nil)

	// swagger:route GET /orders/{id} orders orderDetails
	//
	// gets the details for an order.
	//
	// Consumes:
	// application/json
	// application/x-protobuf
	//
	// Produces:
	// application/json
	// application/x-protobuf
	//
	// Schemes: http, https, ws, wss
	//
	// Security:
	// api_key:
	// oauth: read, write
	//
	// Responses:
	// default: body:genericError
	// 200: body:someResponse
	// 422: body:validationError
	mountItem("GET", basePath+"/orders/:id", nil)

	// swagger:route PUT /orders/{id} orders updateOrder
	//
	// Update the details for an order.
	//
	// When the order doesn't exist this will return an error.
	//
	// Consumes:
	// application/json
	// application/x-protobuf
	//
	// Produces:
	// application/json
	// application/x-protobuf
	//
	// Schemes: http, https, ws, wss
	//
	// Security:
	// api_key:
	// oauth: read, write
	//
	// Responses:
	// default: body:genericError
	// 200: body:someResponse
	// 422: body:validationError
	mountItem("PUT", basePath+"/orders/:id", nil)

	// swagger:route DELETE /orders/{id} deleteOrder
	//
	// delete a particular order.
	//
	// Consumes:
	// application/json
	// application/x-protobuf
	//
	// Produces:
	// application/json
	// application/x-protobuf
	//
	// Schemes: http, https, ws, wss
	//
	// Security:
	// api_key:
	// oauth: read, write
	//
	// Responses:
	// default: body:genericError
	// 200: body:someResponse
	// 422: body:validationError
	mountItem("DELETE", basePath+"/orders/:id", nil)

	// swagger:route POST /param-test params testParams
	//
	// Allow some params with constraints.
	//
	// Consumes:
	// application/json
	// application/x-protobuf
	//
	// Produces:
	// application/json
	// application/x-protobuf
	//
	// Schemes: http, https, ws, wss
	//
	// Security:
	// api_key:
	// oauth: read, write
	//
	// Parameters:
	// + name:        someNumber
	//   description: some number
	//   in:          path
	//   required:    true
	//   allowEmpty:  true
	//   type:        number
	//   max:         20
	//   min:         10
	//   default:     15
	// + name:        someQuery
	//   description: some query values
	//   in:          query
	//   type:        array
	//   minLength:   5
	//   maxLength:   20
	// + name:        someBoolean
	//   in:          path
	//   description: some boolean
	//   type:        boolean
	//   default:     true
	// + name:        constraintsOnInvalidType
	//   description: test constraints on invalid types
	//   in:          query
	//   type:        bool
	//   min:         1
	//   max:         10
	//   minLength:   1
	//   maxLength:   10
	//   format:      abcde
	//   default:     false
	// + name:        noType
	//   description: test no type
	//   min:         1
	//   max:         10
	//   minLength:   1
	//   maxLength:   10
	//   default:     something
	// + name:        request
	//   description: The request model.
	//   in:          body
	//   type:        string
	//   enum:        apple, orange, pineapple, peach, plum
	//   default:     orange
	//
	// Responses:
	// default: body:genericError
	// 200: body:someResponse
	// 422: body:validationError
	mountItem("POST", basePath+"/param-test", nil)

	return nil
}

// not really used but I need a method to decorate the calls to
func mountItem(method, path string, handler interface{}) {}

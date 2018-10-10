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

// ServeAPI serves the API for this record store
func ServeAPI(host, basePath string, schemes []string) (err error) {
	// swagger:operation GET /pets pets getPet
	//
	// List all pets
	//
	// ---
	// parameters:
	//   - name: limit
	//     in: query
	//     description: How many items to return at one time (max 100)
	//     required: false
	//     type: integer
	//     format: int32
	// consumes:
	//   - "application/json"
	//   - "application/xml"
	// produces:
	//   - "application/xml"
	//   - "application/json"
	// responses:
	//   "200":
	//     description: An paged array of pets
	//     headers:
	//       x-next:
	//         type: string
	//         description: A link to the next page of responses
	//     schema:
	//       type: array
	//       items:
	//         schema:
	//           type: object
	//           required:
	//             - id
	//             - name
	//           properties:
	//             id:
	//               type: integer
	//               format: int64
	//             name:
	//               type: string
	//   default:
	//     description: unexpected error
	//     schema:
	//       type: object
	//       required:
	//         - code
	//         - message
	//       properties:
	//         code:
	//           type: integer
	//           format: int32
	//         message:
	//           type: string
	// security:
	//   -
	//     petstore_auth:
	//       - "write:pets"
	//       - "read:pets"
	mountItem("GET", basePath+"/pets", nil)

	// swagger:operation PUT /pets/{id} pets updatePet
	//
	// Updates the details for a pet.
	//
	// Some long explanation,
	// spanning over multipele lines,
	// AKA the description.
	//
	// ---
	// consumes:
	//   - "application/json"
	//   - "application/xml"
	// produces:
	//   - "application/xml"
	//   - "application/json"
	// parameters:
	//   -
	//     in: "body"
	//     name: "body"
	//     description: "Pet object that needs to be added to the store"
	//     required: true
	//     schema:
	//       type: object
	//       required:
	//       - name
	//       properties:
	//         name:
	//           type: string
	//         age:
	//           type: integer
	//           format: int32
	//           minimum: 0
	//   -
	//     in: "path"
	//     name: "id"
	//     description: "Pet object that needs to be added to the store"
	//     required: true
	//     schema:
	//       type: string
	//       pattern: "[A-Z]{3}-[0-9]{3}"
	// responses:
	//   400:
	//     description: "Invalid ID supplied"
	//   404:
	//     description: "Pet not found"
	//   405:
	//     description: "Validation exception"
	// security:
	//   -
	//     petstore_auth:
	//       - "write:pets"
	//       - "read:pets"
	mountItem("PUT", basePath+"/pets/{id}", nil)

	// swagger:operation GET /v1/events Events getEvents
	//
	// Events
	//
	// Mitigation Events
	//
	// ---
	// consumes:
	//   - "application/json"
	//   - "application/xml"
	// produces:
	//   - "application/xml"
	//   - "application/json"
	// parameters:
	// - name: running
	//   in: query
	//   description: (boolean) Filters
	//   required: false
	//   type: boolean
	//
	// responses:
	//  '200':
	//    description: '200'
	//    schema:
	//	    "$ref": "#/definitions/ListResponse"
	//  '400':
	//    description: '400'
	//    schema:
	//      "$ref": "#/definitions/ErrorResponse"
	// security:
	//   -
	//     petstore_auth:
	//       - "write:pets"
	//       - "read:pets"
	mountItem("GET", basePath+"/events", nil)

	// no errors to return, all good
	return
}

// not really used but I need a method to decorate the calls to
func mountItem(method, path string, handler interface{}) {}

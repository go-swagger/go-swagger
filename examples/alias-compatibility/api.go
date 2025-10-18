// Package demo API
//
// Demonstrates type alias handling with the --transparent-aliases flag.
//
// This example shows how type aliases (e.g., type UserID = Identifier)
// are handled differently:
//
// Without --transparent-aliases (default post-#3227):
//   - UserID appears as a definition in swagger.json
//   - User.id references #/definitions/UserID
//
// With --transparent-aliases:
//   - UserID does NOT appear in definitions
//   - User.id references #/definitions/Identifier directly
//   - Matches pre-#3227 behavior
//
// To see the difference:
//   swagger generate spec -m -o without-flag.json
//   swagger generate spec -m --transparent-aliases -o with-flag.json
//   diff <(jq . without-flag.json) <(jq . with-flag.json)
//
//	Schemes: https
//	Host: localhost
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package demo

// Identifier represents a unique identifier
type Identifier string

// UserID is an alias to Identifier for user-specific IDs
type UserID = Identifier

// User represents a user in the system
type User struct {
	ID   UserID `json:"id"`
	Name string `json:"name"`
}

// UserResponse represents a user response
//
// swagger:response UserResponse
type UserResponse struct {
	// in: body
	Body User
}

// swagger:route GET /users/{id} getUser
//
// Get a user by ID
//
// Responses:
//   200: UserResponse
func getUser() {}

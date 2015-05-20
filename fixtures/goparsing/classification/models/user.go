package models

import "github.com/casualjim/go-swagger/strfmt"

// User represents the user for this application
//
// A user is the security principal for this aplication.
// It's also used as one of main axis for reporting.
//
// A user can have friends with whom they can share what they like.
//
// +swagger:model
type User struct {
	// the id for this user
	//
	// required: true
	// min: 1
	ID int64 `json:"id"`

	// the name for this user
	// required: true
	// min length: 3
	Name string `json:"name"`

	// the email address for this user
	//
	// required: true
	// unique: true
	Email strfmt.Email `json:"login"`

	// the friends for this user
	Friends []User `json:"friends"`
}

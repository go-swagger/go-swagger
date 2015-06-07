package models

// A User can purchase pets
// swagger:model user
type User struct {
	// The id of the user.
	//
	// required: true
	ID int64 `json:"id"`

	// The name of the user.
	//
	// required: true
	Name string `json:"name"`
}

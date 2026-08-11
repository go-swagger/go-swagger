package bbb

// User - user model
// swagger:model
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	// User type
	// min: 1
	// max: 5
	Type int `json:"user_type"`
}

package mods

// Category represents a category for a pet
// categories are things like: cat, dog, fish
//
// Even though this model is not annotated with anything,
// and it's not included in the initial imports
// it should still register because it's a required file for the pet model
type Category struct {
	// ID the id of the category
	//
	// required: true
	// min: 1
	ID int64 `json:"id"`

	// Name the name of the category
	//
	// required: true
	Name string `json:"name"`
}

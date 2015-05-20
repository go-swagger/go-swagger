package models

// NoModel is a struct that exists in a package
// but is not annotated with the swagger model annotations
// so it should now show up in a test
//
type NoModel struct {
	// ID of this no model instance
	ID int64 `json:"id"`
	// Name of this no model instance
	Name string `json:"name"`
}

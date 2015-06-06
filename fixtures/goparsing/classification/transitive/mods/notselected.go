package mods

// NotSelected is a model that is in a transitive package
//
// This model is not annotated and should not be detected for parsing.
type NotSelected struct {
	// ID the id of this not selected model
	ID int64 `json:"id"`
	// Name the name of this not selected model
	Name string `json:"name"`
}

// Notable is a model in a transitive package.
// it's used for embedding in another model
//
// +swagger:model withNotes
type Notable struct {
	Notes string `json:"notes"`

	Extra string `json:"extra"`
}

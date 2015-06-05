package models

// A Tag is an extra piece of data to provide more information about a pet.
// It is used to describe the animals available in the store.
// +swagger:model tag
type Tag struct {
	// The id of the tag.
	//
	// required: true
	ID int64 `json:"id"`

	// The value of the tag.
	//
	// required: true
	Value string `json:"value"`
}

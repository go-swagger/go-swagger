package models

// A Pet is the main product in the store.
// It is used to describe the animals available in the store.
//
// swagger:model pet
type Pet struct {
	// The id of the pet.
	//
	// required: true
	ID int64 `json:"id"`

	// The name of the pet.
	//
	// required: true
	// pattern: \w[\w-]+
	// minimum length: 3
	// maximum length: 50
	Name string `json:"name"`

	// The photo urls for the pet.
	// This only accepts jpeg or png images.
	//
	// items pattern: \.(jpe?g|png)$
	PhotoURLs []string `json:"photoUrls,omitempty"`

	// The current status of the pet in the store.
	Status string `json:"status,omitempty"`

	// Extra bits of information attached to this pet.
	//
	Tags []Tag `json:"tags,omitempty"`
}

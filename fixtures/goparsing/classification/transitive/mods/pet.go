package mods

// Pet represents a pet in our store
//
// this model is not explictly mentioned in the import paths
// but because it it transitively required by the order
// it should also be collected.
//
// +swagger:model pet
type Pet struct {
	// ID the id of this pet
	//
	// required: true
	ID int64 `json:"id"`

	// Name the name of the pet
	// this is more like the breed or race of the pet
	//
	// required: true
	// min length: 3
	Name string `json:"name"`

	// Category the category this pet belongs to.
	//
	// required: true
	Category *Category `json:"category"`
}

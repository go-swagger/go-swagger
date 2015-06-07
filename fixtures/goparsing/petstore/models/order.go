package models

import "github.com/casualjim/go-swagger/strfmt"

// An Order for one or more pets by a user.
// swagger:model order
type Order struct {
	// the ID of the order
	//
	// required: true
	ID int64 `json:"id"`

	// the id of the user who placed the order.
	//
	// required: true
	UserID int64 `json:"userId"`

	// the time at which this order was made.
	//
	// required: true
	OrderedAt strfmt.DateTime `json:"orderedAt"`

	// the items for this order
	// mininum items: 1
	Items []struct {

		// the id of the pet to order
		//
		// required: true
		PetID int64 `json:"petId"`

		// the quantity of this pet to order
		//
		// required: true
		// minimum: 1
		Quantity int32 `json:"qty"`
	} `json:"items"`
}

package petstoreapp

import "github.com/casualjim/go-swagger/strfmt"

// Order an order for one or more pets
// an order can be for one type of pet at a time.
//
// +swagger:model
type Order struct {

	// Status Order Status
	Status string `json:"status" xml:"status"`

	// Complete
	Complete bool `json:"complete" xml:"complete"`

	// ID
	ID int64 `json:"id" xml:"id"`

	// PetID
	PetID int64 `json:"petId" xml:"petId"`

	// Quantity
	Quantity int32 `json:"quantity" xml:"quantity"`

	// ShipDate
	ShipDate strfmt.DateTime `json:"shipDate" xml:"shipDate"`
}

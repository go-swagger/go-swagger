package models

import "github.com/casualjim/go-swagger/fixtures/goparsing/classification/transitive/mods"

// StoreOrder represents an order in this application.
//
// An order can either be created, processed or completed.
//
// swagger:model order
type StoreOrder struct {
	// the id for this order
	//
	// required: true
	// min: 1
	ID int64 `json:"id"`

	// the name for this user
	//
	// required: true
	// min length: 3
	UserID int64 `json:"userId"`

	// the items for this order
	Items []struct {
		ID       int32    `json:"id"`
		Pet      mods.Pet `json:"pet"`
		Quantity int16    `json:"quantity"`
	} `json:"items"`
}

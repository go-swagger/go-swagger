// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build testscanner

package schema

// UUID should be discovered through dependency analysis.
type UUID = int64

// Anything should be discovered through dependency analysis.
type Anything = any

// Empty should be discovered through dependency analysis.
type Empty = struct{}

// # StoreOrder represents an order in this application.
//
// An order can either be created, processed or completed.
//
// swagger:model order
type StoreOrder struct {
	// the id for this order
	//
	// required: true
	// min: 1
	ID UUID `json:"id"`

	EID ExtendedID `json:"extended_id"`

	// the name for this user
	//
	// required: true
	// min length: 3
	UserID int64 `json:"userId"`

	// the category of this user
	//
	// required: true
	// default: bar
	// enum: foo,bar,none
	Category string `json:"category"`

	// the items for this order
	Items []struct {
		ID           int32 `json:"id"`
		Quantity     int16 `json:"quantity"`
		ExtraOptions any   `json:"extra_options"`
	} `json:"items"`

	Extras any

	MoreExtras     interface{}
	DeliveryOption Anything
}

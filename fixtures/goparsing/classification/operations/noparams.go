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

package operations

import (
	"bytes"

	"github.com/go-openapi/strfmt"
	"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/models"
	"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/transitive/mods"
)

// MyFileParams contains the uploaded file data
// swagger:parameters myOperation
type MyFileParams struct {
	// MyFormFile desc.
	//
	// in: formData
	//
	// swagger:file
	MyFormFile *bytes.Buffer `json:"myFormFile"`
}

// MyFunc contains a struct with parameters.
func MyFunc() {
	// MyFuncFileParams contains the uploaded file data in a function.
	// swagger:parameters myFuncOperation
	type MyFuncFileParams struct {
		// MyFormFile desc.
		//
		// in: formData
		//
		// swagger:file
		MyFormFile *bytes.Buffer `json:"myFormFile"`
	}
}

// EmbeddedFileParams embeds a *MyFileParams
// swagger:parameters myOtherOperation
type EmbeddedFileParams struct {
	*MyFileParams

	// ExtraParam desc.
	// in: formData
	// required: true
	ExtraParam int `json:"extraParam"`
}

// An OrderBodyParams model.
//
// This is used for operations that want an Order as body of the request
// swagger:parameters updateOrder
// swagger:parameters createOrder
type OrderBodyParams struct {
	// The order to submit.
	//
	// in: body
	// required: true
	Order *models.StoreOrder `json:"order"`
}

// An MultipleOrderParams model.
//
// This is used for operations that want multiple orders as the body
// swagger:parameters getOrders
type MultipleOrderParams struct {
	// The orders
	// required: true
	Orders []*OrderBodyParams `json:"orders"`

	// And another thing
	// in: body
	Another []struct {
		That string `json:"that"`
	} `json:"another"`
}

// Success
// swagger:parameters getConfiguration
type SetConfiguration struct {
	// in:body
	Value map[string]string
}

// swagger:parameters putNumPlate
type NumPlates struct {
	// in: body
	NumPlates interface{} `json:"num_plates"`
}

// A ComplexerOneParams is composed of a SimpleOne and some extra fields
// swagger:parameters yetAnotherOperation
type ComplexerOneParams struct {
	SimpleOne
	mods.NotSelected
	mods.Notable
	CreatedAt strfmt.DateTime `json:"createdAt"`
	Secret    string          `json:"-"`

	// in: formData
	Informity string `json:"informity"`

	NoTagName string `json:",omitempty"`
}

// NoParams is a struct that exists in a package
// but is not annotated with the swagger params annotations
// so it should now show up in a test
//
// swagger:parameters someOperation
// swagger:parameters anotherOperation
type NoParams struct {
	// ID of this no model instance.
	// ids in this application start at 11 and are smaller than 1000
	//
	// required: true
	// minimum: > 10
	// maximum: < 1000
	// in: path
	// default: 1
	ID int64 `json:"id"`

	// The Score of this model
	//
	// required: true
	// minimum: 3
	// maximum: 45
	// multiple of: 3
	// in: query
	// default: 2
	// example: 27
	Score int32 `json:"score"`

	// Name of this no model instance
	//
	// min length: 4
	// max length: 50
	// pattern: [A-Za-z0-9-.]*
	// required: true
	// in: header
	Name string `json:"x-hdr-name"`

	// Created holds the time when this entry was created
	//
	// required: false
	// in: query
	Created strfmt.DateTime `json:"created"`

	// The Category of this model
	//
	// required: true
	// enum: foo,bar,none
	// default: bar
	// in: query
	Category string `json:"category"`

	// Type of this model
	//
	// enum: 1,3,5
	// default: 1
	// in: query
	Type int `json:"type"`

	// This is mix in enum. And actually on output should be valid form where int will be int and
	// string will also be presented.
	//
	// enum: 1,rsq,qaz
	// in: query
	BadEnum int `json:"bad_enum"`

	// a FooSlice has foos which are strings
	//
	// min items: 3
	// max items: 10
	// unique: true
	// items.minLength: 3
	// items.maxLength: 10
	// items.pattern: \w+
	// collection format: pipe
	// items.default: bar
	// in: query
	FooSlice []string `json:"foo_slice"`

	// a BarSlice has bars which are strings
	//
	// min items: 3
	// max items: 10
	// unique: true
	// items.minItems: 4
	// items.maxItems: 9
	// items.enum: bar1,bar2,bar3
	// items.default: bar2
	// items.items.minItems: 5
	// items.items.maxItems: 8
	// items.items.items.minLength: 3
	// items.items.items.maxLength: 10
	// items.items.items.pattern: \w+
	// collection format: pipe
	// in: query
	BarSlice [][][]string `json:"bar_slice"`

	// the items for this order
	//
	// in: body
	Items []struct {
		// ID of this no model instance.
		// ids in this application start at 11 and are smaller than 1000
		//
		// required: true
		// minimum: > 10
		// maximum: < 1000
		// default: 3
		ID int32 `json:"id"`

		// The Pet to add to this NoModel items bucket.
		// Pets can appear more than once in the bucket
		//
		// required: true
		Pet *mods.Pet `json:"pet"`

		// The amount of pets to add to this bucket.
		//
		// required: true
		// minimum: 1
		// maximum: 10
		Quantity int16 `json:"quantity"`

		// Notes to add to this item.
		// This can be used to add special instructions.
		//
		//
		// required: false
		Notes string `json:"notes"`
	} `json:"items"`
}

// NoParamsAlias is a struct that exists in a package
// but is not annotated with the swagger params annotations
// so it should now show up in a test
//
// swagger:parameters someAliasOperation
type NoParamsAlias struct {
	// default "in" is "query" => this params should be aliased
	// required: true
	// minimum: 1
	// maximum: 10
	IntAlias    SomeIntType    `json:"intAlias"`
	StringAlias SomeStringType `json:"stringAlias"`
	// in: path
	IntAliasPath SomeIntType `json:"intAliasPath"`
	// in: formData
	IntAliasForm SomeIntType `json:"intAliasForm"`
}

// SomeStringType is a type that refines string
type SomeStringType string

// SomeIntType is a type that refines int64
type SomeIntType int64

package models

import (
	"time"

	"github.com/casualjim/go-swagger/fixtures/goparsing/classification/transitive/mods"
	"github.com/casualjim/go-swagger/strfmt"
)

// NoModel is a struct that exists in a package
// but is not annotated with the swagger model annotations
// so it should now show up in a test
//
type NoModel struct {
	// ID of this no model instance.
	// ids in this application start at 11 and are smaller than 1000
	//
	//
	// required: true
	// minimum: > 10
	// maximum: < 1000
	ID int64 `json:"id"`

	// The Score of this model
	//
	//
	// required: true
	// minimum: 3
	// maximum: 45
	// multiple of: 3
	Score int32 `json:"score"`

	// Name of this no model instance
	//
	//
	// min length: 4
	// max length: 50
	// pattern: [A-Za-z0-9-.]*
	// required: true
	Name string `json:"name"`

	// Created holds the time when this entry was created
	//
	//
	// required: false
	// read only: true
	Created strfmt.DateTime `json:"created"`

	// a FooSlice has foos which are strings
	//
	//
	// min items: 3
	// max items: 10
	// unique: true
	// items.minLength: 3
	// items.maxLength: 10
	// items.pattern: \w+
	FooSlice []string `json:"foo_slice"`

	// the items for this order
	Items []struct {
		// ID of this no model instance.
		// ids in this application start at 11 and are smaller than 1000
		//
		//
		// required: true
		// minimum: > 10
		// maximum: < 1000
		ID int32 `json:"id"`

		// The Pet to add to this NoModel items bucket.
		// Pets can appear more than once in the bucket
		//
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

// SomeStringType is a type that refines string
type SomeStringType string

// SomeIntType is a type that refines int64
type SomeIntType int64

// SomeTimeType is a type that refines time.Time
type SomeTimeType time.Time

// A PrimateModel is a struct with nothing but primitives.
//
// It only has values 1 level deep and each of those is of a very simple
// builtin type.
type PrimateModel struct {
	A bool `json:"a"`

	B rune   `json:"b"`
	C string `json:"c"`

	D int   `json:"d"`
	E int8  `json:"e"`
	F int16 `json:"f"`
	G int32 `json:"g"`
	H int64 `json:"h"`

	I uint   `json:"i"`
	J uint8  `json:"j"`
	K uint16 `json:"k"`
	L uint32 `json:"l"`
	M uint64 `json:"m"`

	N float32 `json:"n"`
	O float64 `json:"o"`
}

// A FormattedModel is a struct with only strfmt types
//
// It only has values 1 level deep and is used for testing the conversion
type FormattedModel struct {
	A strfmt.Base64     `json:"a"`
	B strfmt.CreditCard `json:"b"`
	C strfmt.Date       `json:"c"`
	D strfmt.DateTime   `json:"d"`
	E strfmt.Duration   `json:"e"`
	F strfmt.Email      `json:"f"`
	G strfmt.HexColor   `json:"g"`
	H strfmt.Hostname   `json:"h"`
	I strfmt.IPv4       `json:"i"`
	J strfmt.IPv6       `json:"j"`
	K strfmt.ISBN       `json:"k"`
	L strfmt.ISBN10     `json:"l"`
	M strfmt.ISBN13     `json:"m"`
	N strfmt.RGBColor   `json:"n"`
	O strfmt.SSN        `json:"o"`
	P strfmt.URI        `json:"p"`
	Q strfmt.UUID       `json:"q"`
	R strfmt.UUID3      `json:"r"`
	S strfmt.UUID4      `json:"s"`
	T strfmt.UUID5      `json:"t"`
}

// A SimpleComplexModel is a struct with only other struct types
//
// It doesn't have slices or arrays etc but only complex types
// so also no primitives or string formatters
type SimpleComplexModel struct {
	Top Something `json:"top"`

	NotSel mods.NotSelected `json:"notSel"`

	Emb struct {
		CID int64  `json:"cid"`
		Baz string `json:"baz"`
	} `json:"emb"`
}

// A Something struct is used by other structs
type Something struct {
	DID int64  `json:"did"`
	Cat string `json:"cat"`
}

// Pointdexter is a struct with only pointers
type Pointdexter struct {
	ID   *int64        `json:"id"`
	Name *string       `json:"name"`
	T    *strfmt.UUID5 `json:"t"`
	Top  *Something    `json:"top"`

	NotSel *mods.NotSelected `json:"notSel"`

	Emb *struct {
		CID *int64  `json:"cid"`
		Baz *string `json:"baz"`
	} `json:"emb"`
}

// A SliceAndDice struct contains only slices
//
// the elements of the slices are structs, primitives or string formats
// there is also a pointer version of each property
type SliceAndDice struct {
	IDs     []int64            `json:"ids"`
	Names   []string           `json:"names"`
	UUIDs   []strfmt.UUID      `json:"uuids"`
	Tops    []Something        `json:"tops"`
	NotSels []mods.NotSelected `json:"notSels"`
	Embs    []struct {
		CID []int64  `json:"cid"`
		Baz []string `json:"baz"`
	} `json:"embs"`

	PtrIDs     []*int64            `json:"ptrIds"`
	PtrNames   []*string           `json:"ptrNames"`
	PtrUUIDs   []*strfmt.UUID      `json:"ptrUuids"`
	PtrTops    []*Something        `json:"ptrTops"`
	PtrNotSels []*mods.NotSelected `json:"ptrNotSels"`
	PtrEmbs    []*struct {
		PtrCID []*int64  `json:"ptrCid"`
		PtrBaz []*string `json:"ptrBaz"`
	} `json:"ptrEmbs"`
}

// A MapTastic struct contains only maps
//
// the values of the maps are structs, primitives or string formats
// there is also a pointer version of each property
type MapTastic struct {
	IDs     map[string]int64            `json:"ids"`
	Names   map[string]string           `json:"names"`
	UUIDs   map[string]strfmt.UUID      `json:"uuids"`
	Tops    map[string]Something        `json:"tops"`
	NotSels map[string]mods.NotSelected `json:"notSels"`
	Embs    map[string]struct {
		CID map[string]int64  `json:"cid"`
		Baz map[string]string `json:"baz"`
	} `json:"embs"`

	PtrIDs     map[string]*int64            `json:"ptrIds"`
	PtrNames   map[string]*string           `json:"ptrNames"`
	PtrUUIDs   map[string]*strfmt.UUID      `json:"ptrUuids"`
	PtrTops    map[string]*Something        `json:"ptrTops"`
	PtrNotSels map[string]*mods.NotSelected `json:"ptrNotSels"`
	PtrEmbs    map[string]*struct {
		PtrCID map[string]*int64  `json:"ptrCid"`
		PtrBaz map[string]*string `json:"ptrBaz"`
	} `json:"ptrEmbs"`
}

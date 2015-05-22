package models

import (
	"time"

	mods "github.com/casualjim/go-swagger/fixtures/goparsing/classification/transitive/mods"
	"github.com/casualjim/go-swagger/strfmt"
)

// NoModel is a struct that exists in a package
// but is not annotated with the swagger model annotations
// so it should now show up in a test
//
type NoModel struct {
	// ID of this no model instance
	ID int64 `json:"id"`
	// Name of this no model instance
	Name string `json:"name"`

	// the time when this entry was created
	Created strfmt.DateTime `json:"created"`

	FooArr   [3]string
	FooSlice []string

	NestedFoo          [][]string
	DeeplyNestedFooBar [][]map[string][]map[string]string

	BarIFace interface{}

	// the items for this order
	Items []struct {
		ID       int32    `json:"id"`
		Pet      mods.Pet `json:"pet"`
		Quantity int16    `json:"quantity"`
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

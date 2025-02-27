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

package models

import (
	"net/url"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/transitive/mods"
)

// NoModel is a struct without an annotation.
// NoModel exists in a package
// but is not annotated with the swagger model annotations
// so it should now show up in a test.
type NoModel struct {
	// ID of this no model instance.
	// ids in this application start at 11 and are smaller than 1000
	//
	// required: true
	// minimum: > 10
	// maximum: < 1000
	// default: 11
	ID int64 `json:"id"`

	Ignored      string `json:"-"`
	IgnoredOther string `json:"-,omitempty"`

	// A field which has omitempty set but no name
	NoNameOmitEmpty string `json:",omitempty"`

	// Note is a free form data in base64
	//
	// swagger:strfmt byte
	Note []byte `json:"noteb64,omitempty"`

	// The Score of this model
	//
	// required: true
	// minimum: 3
	// maximum: 45
	// multiple of: 3
	// example: 27
	Score int32 `json:"score"`

	// Name of this no model instance
	//
	// min length: 4
	// max length: 50
	// pattern: [A-Za-z0-9-.]*
	// required: true
	//
	// Extensions:
	// ---
	// x-property-value: value
	// x-property-array:
	//   - value1
	//   - value2
	// x-property-array-obj:
	//   - name: obj
	//     value: field
	// ---
	//
	Name string `json:"name"`

	// Created holds the time when this entry was created
	//
	// required: false
	// read only: true
	Created strfmt.DateTime `json:"created"`

	// GoTimeCreated holds the time when this entry was created in go time.Time
	//
	// required: false
	GoTimeCreated time.Time `json:"gocreated"`

	// a FooSlice has foos which are strings
	//
	// min items: 3
	// max items: 10
	// unique: true
	// items.minLength: 3
	// items.maxLength: 10
	// items.pattern: \w+
	FooSlice []string `json:"foo_slice"`

	// a TimeSlice is a slice of times
	//
	// min items: 3
	// max items: 10
	// unique: true
	TimeSlice []time.Time `json:"time_slice"`

	// a BarSlice has bars which are strings
	//
	// min items: 3
	// max items: 10
	// unique: true
	// items.minItems: 4
	// items.maxItems: 9
	// items.items.minItems: 5
	// items.items.maxItems: 8
	// items.items.items.minLength: 3
	// items.items.items.maxLength: 10
	// items.items.items.pattern: \w+
	BarSlice [][][]string `json:"bar_slice"`

	// a DeepSlice has bars which are time
	//
	// min items: 3
	// max items: 10
	// unique: true
	// items.minItems: 4
	// items.maxItems: 9
	// items.items.minItems: 5
	// items.items.maxItems: 8
	DeepTimeSlice [][][]time.Time `json:"deep_time_slice"`

	// the items for this order
	Items []struct {
		// ID of this no model instance.
		// ids in this application start at 11 and are smaller than 1000
		//
		// required: true
		// minimum: > 10
		// maximum: < 1000
		// default: 11
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

		// A dummy expiration date.
		//
		// required: true
		Expiration time.Time `json:"expiration"`

		// Notes to add to this item.
		// This can be used to add special instructions.
		//
		// required: false
		Notes string `json:"notes"`

		AlsoIgnored string `json:"-"`
	} `json:"items"`
}

// A OtherTypes struct contains type aliases
type OtherTypes struct {
	Named       SomeStringType     `json:"named"`
	Numbered    SomeIntType        `json:"numbered"`
	Dated       SomeTimeType       `json:"dated"`
	Timed       SomeTimedType      `json:"timed"`
	Petted      SomePettedType     `json:"petted"`
	Somethinged SomethingType      `json:"somethinged"`
	StrMap      SomeStringMap      `json:"strMap"`
	StrArrMap   SomeArrayStringMap `json:"strArrMap"`

	ManyNamed       SomeStringsType `json:"manyNamed"`
	ManyNumbered    SomeIntsType    `json:"manyNumbered"`
	ManyDated       SomeTimesType   `json:"manyDated"`
	ManyTimed       SomeTimedsType  `json:"manyTimed"`
	ManyPetted      SomePettedsType `json:"manyPetted"`
	ManySomethinged SomethingsType  `json:"manySomethinged"`

	Nameds       []SomeStringType `json:"nameds"`
	Numbereds    []SomeIntType    `json:"numbereds"`
	Dateds       []SomeTimeType   `json:"dateds"`
	Timeds       []SomeTimedType  `json:"timeds"`
	Petteds      []SomePettedType `json:"petteds"`
	Somethingeds []SomethingType  `json:"somethingeds"`

	ModsNamed    mods.SomeStringType `json:"modsNamed"`
	ModsNumbered mods.SomeIntType    `json:"modsNumbered"`
	ModsDated    mods.SomeTimeType   `json:"modsDated"`
	ModsTimed    mods.SomeTimedType  `json:"modsTimed"`
	ModsPetted   mods.SomePettedType `json:"modsPetted"`

	ModsNameds    []mods.SomeStringType `json:"modsNameds"`
	ModsNumbereds []mods.SomeIntType    `json:"modsNumbereds"`
	ModsDateds    []mods.SomeTimeType   `json:"modsDateds"`
	ModsTimeds    []mods.SomeTimedType  `json:"modsTimeds"`
	ModsPetteds   []mods.SomePettedType `json:"modsPetteds"`

	ManyModsNamed     mods.SomeStringsType    `json:"manyModsNamed"`
	ManyModsNumbered  mods.SomeIntsType       `json:"manyModsNumbered"`
	ManyModsDated     mods.SomeTimesType      `json:"manyModsDated"`
	ManyModsTimed     mods.SomeTimedsType     `json:"manyModsTimed"`
	ManyModsPetted    mods.SomePettedsType    `json:"manyModsPetted"`
	ManyModsPettedPtr mods.SomePettedsPtrType `json:"manyModsPettedPtr"`

	NamedAlias     SomeStringTypeAlias   `json:"namedAlias"`
	NumberedAlias  SomeIntTypeAlias      `json:"numberedAlias"`
	NamedsAlias    []SomeStringTypeAlias `json:"namedsAlias"`
	NumberedsAlias []SomeIntTypeAlias    `json:"numberedsAlias"`
}

// A SimpleOne is a model with a few simple fields
type SimpleOne struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}

// A ComplexerOne is composed of a SimpleOne and some extra fields
type ComplexerOne struct {
	SimpleOne
	mods.NotSelected
	mods.Notable
	CreatedAt strfmt.DateTime `json:"createdAt"`
}

// An OverridingOne is composed of a SimpleOne and overrides a field
type OverridingOne struct {
	SimpleOne
	Age int64
}

// An OverridingOneIgnore is composed of a SimpleOne and overrides a field to ignore it
type OverridingOneIgnore struct {
	SimpleOne
	Age int32 `json:"-"`
}

// An AllOfModel is composed out of embedded structs but it should build
// an allOf property
type AllOfModel struct {
	// swagger:allOf
	SimpleOne
	// swagger:allOf
	mods.Notable

	Something // not annotated with anything, so should be included

	CreatedAt strfmt.DateTime `json:"createdAt"`
}

// An Embedded is to be embedded in EmbeddedStarExpr
type Embedded struct {
	EmbeddedMember int64 `json:"embeddedMember"`
}

// An EmbeddedStarExpr for testing the embedded StarExpr
type EmbeddedStarExpr struct {
	*Embedded
	NotEmbedded int64 `json:"notEmbedded"`
}

// A PrimateModel is a struct with nothing but builtins.
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

	P byte `json:"p"`

	Q uintptr `json:"q"`
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
	U strfmt.MAC        `json:"u"`
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

// An Interfaced struct contains objects with interface definitions
type Interfaced struct {
	CustomData interface{} `json:"custom_data"`
}

// A BaseStruct is a struct that has subtypes.
//
// it should deserialize into one of the struct types that
// enlist for being implementations of this struct
//
// swagger:model animal
type BaseStruct struct {
	// ID of this no model instance.
	// ids in this application start at 11 and are smaller than 1000
	//
	// required: true
	// minimum: > 10
	// maximum: < 1000
	ID int64 `json:"id"`

	// Name of this no model instance
	//
	// min length: 4
	// max length: 50
	// pattern: [A-Za-z0-9-.]*
	// required: true
	Name string `json:"name"`

	// StructType the type of this polymorphic model
	//
	// discriminator: true
	StructType string `json:"jsonClass"`
}

/* TODO: implement this in the scanner

// A Lion is a struct that "subtypes" the BaseStruct
//
// it does so because it included the fields in the struct body
// The scanner assumes it will follow the rules and describes this type
// as discriminated in the swagger spec based on the discriminator value
// annotation.
//
// swagger:model lion
// swagger:discriminatorValue animal org.horrible.java.fqpn.TheLionDataObjectFactoryInstanceServiceImpl
type Lion struct {
	// ID of this no model instance.
	// ids in this application start at 11 and are smaller than 1000
	//
	// required: true
	// minimum: > 10
	// maximum: < 1000
	ID int64 `json:"id"`

	// Name of this no model instance
	//
	// min length: 4
	// max length: 50
	// pattern: [A-Za-z0-9-.]*
	// required: true
	Name string `json:"name"`

	// StructType the type of this polymorphic model
	StructType string `json:"jsonClass"`

	// Leader is true when this is the leader of its group
	//
	// default value: false
	Leader bool `json:"leader"`
}
*/

// A Giraffe is a struct that embeds BaseStruct
//
// the annotation is not necessary here because of inclusion
// of a discriminated type
// it infers the name of the x-class value from its context
//
// swagger:model giraffe
type Giraffe struct {
	// swagger:allOf
	BaseStruct

	// NeckSize the size of the neck of this giraffe
	NeckSize int64 `json:"neckSize"`
}

// A Gazelle is a struct is discriminated for BaseStruct.
//
// The struct includes the BaseStruct and that embedded value
// is annotated with the discriminator value annotation so it
// where it only requires 1 argument because it knows which
// discriminator type this belongs to
//
// swagger:model gazelle
type Gazelle struct {
	// swagger:allOf a.b.c.d.E
	BaseStruct

	// The size of the horns
	HornSize float32 `json:"hornSize"`
}

// Identifiable is an interface for things that have an ID
type Identifiable interface {
	// ID of this no model instance.
	// ids in this application start at 11 and are smaller than 1000
	//
	// required: true
	// minimum: > 10
	// maximum: < 1000
	// swagger:name id
	ID() int64
}

// WaterType is an interface describing a water type
//
// swagger:model water
type WaterType interface {
	// swagger:name sweetWater
	SweetWater() bool
	// swagger:name saltWater
	SaltWater() bool
}

// Fish represents a base type implemented as interface
// the nullary methods of this interface will be included as
//
// swagger:model fish
type Fish interface {
	Identifiable // interfaces like this are included as if they were defined directly on this type

	// embeds decorated with allOf are included as refs

	// swagger:allOf
	WaterType

	// swagger:allOf
	mods.ExtraInfo

	mods.EmbeddedColor

	Items(id, size int64) []string

	// Name of this no model instance
	//
	// min length: 4
	// max length: 50
	// pattern: [A-Za-z0-9-.]*
	// required: true
	// swagger:name name
	Name() string

	// StructType the type of this polymorphic model
	// Discriminator: true
	// swagger:name jsonClass
	StructType() string
}

// TeslaCar is a tesla car
//
// swagger:model
type TeslaCar interface {
	// The model of tesla car
	//
	// discriminated: true
	// swagger:name model
	Model() string

	// AutoPilot returns true when it supports autopilot
	// swagger:name autoPilot
	AutoPilot() bool
}

// The ModelS version of the tesla car
//
// swagger:model modelS
type ModelS struct {
	// swagger:allOf com.tesla.models.ModelS
	TeslaCar
	// The edition of this Model S
	Edition string `json:"edition"`
}

// Test proper parsing of type declaration _blocks_:
type (
	// The ModelX version of the tesla car
	//
	// swagger:model modelX
	ModelX struct {
		// swagger:allOf com.tesla.models.ModelX
		TeslaCar
		// The number of doors on this Model X
		Doors int `json:"doors"`
	}

	// The ModelA version of the tesla car
	//
	// swagger:model modelA
	ModelA struct {
		Tesla TeslaCar
		// The number of doors on this Model A
		Doors int `json:"doors"`
	}
)

// Cars is a collection of cars
//
// swagger:model cars
type Cars struct {
	Cars []*TeslaCar `json:"cars"`
}

// JSONString has fields with ",string" JSON directives.
//
// swagger:model jsonString
type JSONString struct {
	// Should be encoded as a string with string format "integer"
	SomeInt    int    `json:"someInt,string"`
	SomeInt8   int8   `json:"someInt8,string"`
	SomeInt16  int16  `json:"someInt16,string"`
	SomeInt32  int32  `json:"someInt32,string"`
	SomeInt64  int64  `json:"someInt64,string"`
	SomeUint   uint   `json:"someUint,string"`
	SomeUint8  uint8  `json:"someUint8,string"`
	SomeUint16 uint16 `json:"someUint16,string"`
	SomeUint32 uint32 `json:"someUint32,string"`
	SomeUint64 uint64 `json:"someUint64,string"`

	// Should be encoded as a string with string format "double"
	SomeFloat64 float64 `json:"someFloat64,string"`

	// Should be encoded as a string with no format
	SomeString string `json:"someString,string"`

	// Should be encoded as a string with no format
	SomeBool bool `json:"someBool,string"`

	// The ",string" directive should be ignore before the type isn't scalar
	SomethingElse Cars `json:"somethingElse,string"`

	// The ",omitempty,string" directive should be valid
	SomeDefaultInt int `json:",omitempty,string"`
}

// JSONPtrString has fields with ",string" JSON directives.
//
// swagger:model jsonPtrString
type JSONPtrString struct {
	// Should be encoded as a string with string format "integer"
	SomeInt    *int    `json:"someInt,string"`
	SomeInt8   *int8   `json:"someInt8,string"`
	SomeInt16  *int16  `json:"someInt16,string"`
	SomeInt32  *int32  `json:"someInt32,string"`
	SomeInt64  *int64  `json:"someInt64,string"`
	SomeUint   *uint   `json:"someUint,string"`
	SomeUint8  *uint8  `json:"someUint8,string"`
	SomeUint16 *uint16 `json:"someUint16,string"`
	SomeUint32 *uint32 `json:"someUint32,string"`
	SomeUint64 *uint64 `json:"someUint64,string"`

	// Should be encoded as a string with string format "double"
	SomeFloat64 *float64 `json:"someFloat64,string"`

	// Should be encoded as a string with no format
	SomeString *string `json:"someString,string"`

	// Should be encoded as a string with no format
	SomeBool *bool `json:"someBool,string"`

	// The ",string" directive should be ignore before the type isn't scalar
	SomethingElse *Cars `json:"somethingElse,string"`
}

// IgnoredFields demostrates the use of swagger:ignore on struct fields.
//
// swagger:model ignoredFields
type IgnoredFields struct {
	SomeIncludedField string `json:"someIncludedField"`

	// swagger:ignore
	SomeIgnoredField string `json:"someIgnoredField"`

	// This swagger:ignore tag won't work - it needs to be in the field's doc
	// block
	SomeErroneouslyIncludedField string `json:"someErroneouslyIncludedField"` // swagger:ignore
}

// UUID is a type that represents a UUID as a string
type UUID [16]byte

func (uuid UUID) MarshalText() ([]byte, error) {
	return []byte("hola desde UUID"), nil
}

type MarshalTextStruct struct {
	Hola string
}

func (cm MarshalTextStruct) MarshalText() ([]byte, error) {
	return []byte("hi from CustomStruct"), nil
}

type MarshalTextMap map[string]interface{}

func (cm MarshalTextMap) MarshalText() ([]byte, error) {
	return []byte("hola desde CustomMap"), nil
}

// swagger:strfmt date-time
type MarshalTextStructStrfmt struct {
	Foo string `json:"foo"`
}

func (cm MarshalTextStructStrfmt) MarshalText() ([]byte, error) {
	return []byte("hi from CustomStructStrfmt"), nil
}

// swagger:strfmt date-time
type MarshalTextStructStrfmtPtr struct {
	Foo string `json:"foo"`
}

func (cm MarshalTextStructStrfmtPtr) MarshalText() ([]byte, error) {
	return []byte("hi frome CustomStructStrfmtPtr"), nil
}

// swagger:strfmt url
type URL url.URL

// TextMarshalModel demostrates the use of MarshalText from different fields
//
// swagger:model TextMarshalModel
type TextMarshalModel struct {
	ID              UUID                        `json:"id"`
	IDs             []UUID                      `json:"ids"`
	Struct          MarshalTextStruct           `json:"struct"`
	Map             MarshalTextMap              `json:"map"`
	MapUUID         map[string]UUID             `json:"mapUUID"`
	URL             url.URL                     `json:"url"` // url.URL not has TextMarshal!
	Time            time.Time                   `json:"time"`
	StructStrfmt    MarshalTextStructStrfmt     `json:"structStrfmt"`
	StructStrfmtPtr *MarshalTextStructStrfmtPtr `json:"structStrfmtPtr"`
	CustomURL       URL                         `json:"customUrl"`
}

// swagger:type object
type SomeObjectMap interface{}

// swagger:model namedWithType
type NamedWithType struct {
	SomeMap SomeObjectMap `json:"some_map"`
}

//
// Next models are related to named types with type arguments
//

type GenericResults[T any] struct {
	RecordsMatched uint32 `json:"records_matched"`
	Matches        []T    `json:"matches"`
}

// swagger:model namedStringResults
type NamedStringResults GenericResults[string]

// swagger:model namedStoreOrderResults
type NamedStoreOrderResults GenericResults[StoreOrder]

type GenericSlice[T any] []T

// swagger:model namedStringSlice
type NamedStringSlice GenericSlice[string]

// swagger:model namedStoreOrderSlice
type NamedStoreOrderSlice GenericSlice[StoreOrder]

type GenericMap[K comparable, V any] map[K]V

// swagger:model namedStringMap
type NamedStringMap GenericMap[string, string]

// swagger:model namedStoreOrderMap
type NamedStoreOrderMap GenericMap[string, StoreOrder]

// swagger:model namedMapOfStoreOrderSlices
type NamedMapOfStoreOrderSlices GenericMap[string, GenericSlice[StoreOrder]]

//
// End of models related to named types with type arguments
//

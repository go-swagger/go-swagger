package mods

import (
	"time"

	"github.com/casualjim/go-swagger/strfmt"
)

// SomeStringsType is a type that refines []string
type SomeStringsType []string

// SomeIntsType is a type that refines []int64
type SomeIntsType []int64

// SomeTimesType is a type that refines time.Time
// +swagger:strfmt date-time
type SomeTimesType []time.Time

// SomeTimedsType is a type that refines strfmt.DateTime
type SomeTimedsType []strfmt.DateTime

// SomePettedsType is a type that refines mods.Pet
type SomePettedsType []Pet

// SomeStringType is a type that refines string
type SomeStringType string

// SomeIntType is a type that refines int64
type SomeIntType int64

// SomeTimeType is a type that refines time.Time
// +swagger:strfmt date-time
type SomeTimeType time.Time

// SomeTimedType is a type that refines strfmt.DateTime
type SomeTimedType strfmt.DateTime

// SomePettedType is a type that refines Pet
type SomePettedType Pet

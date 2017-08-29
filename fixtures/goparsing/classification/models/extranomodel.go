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
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-swagger/go-swagger/fixtures/goparsing/classification/transitive/mods"
)

// A Something struct is used by other structs
type Something struct {
	DID int64  `json:"did"`
	Cat string `json:"cat"`
}

// SomeStringType is a type that refines string
type SomeStringType string

// SomeIntType is a type that refines int64
type SomeIntType int64

// SomeTimeType is a type that refines time.Time
// swagger:strfmt date-time
type SomeTimeType time.Time

// SomeTimedType is a type that refines strfmt.DateTime
type SomeTimedType strfmt.DateTime

// SomePettedType is a type that refines mods.Pet
type SomePettedType mods.Pet

// SomethingType is a type that refines a type contained in the same package
type SomethingType Something

// SomeStringsType is a type that refines []string
type SomeStringsType []string

// SomeIntsType is a type that refines []int64
type SomeIntsType []int64

// SomeTimesType is a type that refines time.Time
// swagger:strfmt date-time
type SomeTimesType []time.Time

// SomeTimedsType is a type that refines strfmt.DateTime
type SomeTimedsType []strfmt.DateTime

// SomePettedsType is a type that refines mods.Pet
type SomePettedsType []mods.Pet

// SomethingsType is a type that refines a type contained in the same package
type SomethingsType []Something

// SomeObject is a type that refines an untyped map
type SomeObject map[string]interface{}

// SomeStringMap is a type that refines a string value map
type SomeStringMap map[string]string

// SomeArrayStringMap is a type that refines a array of strings value map
type SomeArrayStringMap map[string][]string

// SomeIntMap is a type that refines an int value map
type SomeIntMap map[string]int64

// SomeTimeMap is a type that refines a time.Time value map
// swagger:strfmt date-time
type SomeTimeMap map[string]time.Time

// SomeTimedMap is a type that refines an strfmt.DateTime value map
type SomeTimedMap map[string]strfmt.DateTime

// SomePettedMap is a type that refines a pet value map
type SomePettedMap map[string]mods.Pet

// SomeSomethingMap is a type that refines a Something value map
type SomeSomethingMap map[string]Something

// SomeStringTypeAlias is a type that refines string
// swagger:alias
type SomeStringTypeAlias string

// SomeIntTypeAlias is a type that refines int64
// swagger:alias
type SomeIntTypeAlias int64

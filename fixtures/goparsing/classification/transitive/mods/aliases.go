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

package mods

import (
	"time"

	"github.com/go-openapi/strfmt"
)

// SomeStringsType is a type that refines []string
// swagger:model modsSomeStringsType
type SomeStringsType []string

// SomeIntsType is a type that refines []int64
// swagger:model modsSomeIntsType
type SomeIntsType []int64

// SomeTimesType is a type that refines time.Time
// swagger:strfmt date-time
// swagger:model modsSomeTimesType
type SomeTimesType []time.Time

// SomeTimedsType is a type that refines strfmt.DateTime
// swagger:model modsSomeTimedsType
type SomeTimedsType []strfmt.DateTime

// SomePettedsType is a type that refines mods.Pet
// swagger:model modsSomePettedsType
type SomePettedsType []Pet

// SomePettedsPtrType is a type that refines array of mods.Pet pointers
// swagger:model modsSomePettedsPtrType
type SomePettedsPtrType []*Pet

// SomeStringType is a type that refines string
// swagger:model modsSomeStringType
type SomeStringType string

// SomeIntType is a type that refines int64
// swagger:model modsSomeIntType
type SomeIntType int64

// SomeTimeType is a type that refines time.Time
// swagger:strfmt date-time
// swagger:model modsSomeTimeType
type SomeTimeType time.Time

// SomeTimedType is a type that refines strfmt.DateTime
// swagger:model modsSomeTimedType
type SomeTimedType strfmt.DateTime

// SomePettedType is a type that refines Pet
// swagger:model modsSomePettedType
type SomePettedType Pet

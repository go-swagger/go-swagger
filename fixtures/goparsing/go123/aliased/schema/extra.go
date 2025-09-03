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

// swagger:model void
type Void = Empty

// swagger:model empty_redefinition
type EmptyRedefinition struct{}

// swagger:model anonymous_struct
type AnonymousStruct struct {
	A struct {
		B int
	}
}

// swagger:model whatnot
type WhatNot any

// swagger:model aliased_id
type AliasedID = ExtendedID

// swagger:model whatnot_alias
type WhatNotAlias = any

// swagger:model iface
type Iface interface {
	Get() string
	Set(string)
}

// swagger:model iface_alias
type IfaceAlias = Iface

// swagger:model iface_redefinition
type IfaceRedefinition Iface

// swagger:model anonymous_iface
type AnonymousIface interface{ String() string }

// swagger:model anonymous_iface_alias
type AnonymousIfaceAlias = AnonymousIface

// swagger:model whatnot2
//
// This is a type redefinition.
type WhatNot2 interface{}

// swagger:model whatnot2_alias
//
// This is a syntactic alias.
type WhatNot2Alias = interface{}

// swagger:model slice_type
type Slice []any

// swagger:model slice_alias
type SliceAlias = []any

// swagger:model slice_to_slice
type SliceToSlice = Slice

// swagger:model slice_of_structs
//
// SliceOfStructs is an slice of anonymous structs.
type SliceOfStructs []struct{}

// swagger:model slice_of_structs_alias
type SliceOfStructsAlias = []struct{}

// swagger:model
type ShouldSee bool

type ShouldNotSee bool
type ShouldNotSeeSlice []int
type ShouldNotSeeMap map[string]int

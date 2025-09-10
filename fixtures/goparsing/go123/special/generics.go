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

/*//go:build testscanner*/

package special

// swagger:model generic_constraint
type Constraint interface {
	~uint16
	interface{ Uint() uint16 }
}

// swagger:model numerical_constraint
type Numerical interface {
	~uint16 | ~int16
}

// swagger:model union_alias
type Union = Numerical

// swagger:model generic_map
type GenericMap[K comparable, V any] map[K]V

// swagger:model generic_map_alias
type GenericMapAlias[K comparable, V any] = map[K]V

// swagger:model generic_indirect
type GenericIndirect[K comparable, V any] = GenericMapAlias[K, V]

// swagger:model generic_slice
type GenericSlice[T any] []T

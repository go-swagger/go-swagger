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

// swagger:model iface_embedded
type IfaceEmbedded interface {
	Iface

	Dump() []byte
	error
}

// swagger:model iface_embedded_with_alias
type IfaceEmbeddedWithAlias interface {
	IfaceAlias
	AnonymousIface

	Dump() []byte
}

// swagger:model iface_embedded_empty
type IfaceEmpty interface {
	interface{}
	any
}

// swagger:model iface_embedded_anonymous
type IfaceEmbeddedAnonymous interface {
	interface{ String() string }
	interface{ Error() string }
	interface{}
	error
}

// ExtendedID should be discovered through dependency analysis.
type ExtendedID struct {
	Empty

	More      string `json:"more"`
	EvenMore  any
	StillMore interface{}
	_         struct{}
}

// swagger:model embedded_with_alias
type EmbeddedWithAlias struct {
	Anything
	UUID

	More      string `json:"more"`
	EvenMore  any
	StillMore interface{}
}

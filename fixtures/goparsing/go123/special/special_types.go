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

import (
	"encoding/json"
	"math/big"
	"reflect"
	"time"
	"unsafe"

	"github.com/go-openapi/strfmt"
)

// swagger:model primitive
type Primitive string

// swagger:model unsafe_pointer_alias
type Unsafe = unsafe.Pointer

// swagger:model upointer_alias
type UIntPtr = uintptr

// swagger:model go_map
type GoMap map[string]uint16

type GoStruct struct {
	A *float32
}

// swagger:model index_map
type UnsupportedMap map[int]struct{}

// swagger:model go_error
type Error error

// swagger:model special_types
type SpecialTypes struct {
	PtrStruct              *GoStruct
	ShouldBeStringTime     time.Time
	ShouldAlsoBeStringTime *time.Time
	Error                  error
	Marshaler              IsATextMarshaler
	Message                json.RawMessage
	Duration               time.Duration
	FormatDate             strfmt.Date
	FormatTime             strfmt.DateTime
	FormatUUID             strfmt.UUID
	PtrFormatUUID          *strfmt.UUID
	Err                    error
	Map                    map[string]*GoStruct

	// and what not
	WhatNot struct {
		unexported int
		AA         complex128
		A          complex64
		B          chan int
		C          func()
		D          func() string
		E          unsafe.Pointer
		F          uintptr
		G          *big.Float
		H          *big.Int
		I          [5]byte
		J          reflect.Type
		K          reflect.Value
	}
}

// swagger:model unexported
type unexported struct{}

type IsATextMarshaler struct {
	unexported string
}

func (m IsATextMarshaler) MarshalText() ([]byte, error) {
	return []byte(m.unexported), nil
}

// swagger:model go_array
type GoArray [3]int64

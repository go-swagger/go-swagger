
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

package flags

import (
	"reflect"
)

// Arg represents a positional argument on the command line.
type Arg struct {
	// The name of the positional argument (used in the help)
	Name string

	// A description of the positional argument (used in the help)
	Description string

	// Whether a positional argument is required
	Required int

	value reflect.Value
	tag   multiTag
}

func (a *Arg) isRemaining() bool {
	return a.value.Type().Kind() == reflect.Slice
}

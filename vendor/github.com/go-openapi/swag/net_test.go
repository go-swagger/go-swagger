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

package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitHostPort(t *testing.T) {
	data := []struct {
		Input string
		Host  string
		Port  int
		Err   bool
	}{
		{"localhost:3933", "localhost", 3933, false},
		{"localhost:yellow", "", -1, true},
		{"localhost", "", -1, true},
		{"localhost:", "", -1, true},
		{"localhost:3933", "localhost", 3933, false},
	}

	for _, e := range data {
		h, p, err := SplitHostPort(e.Input)
		if (!e.Err && assert.NoError(t, err)) || (e.Err && assert.Error(t, err)) {
			assert.Equal(t, e.Host, h)
			assert.Equal(t, e.Port, p)
		}
	}
}

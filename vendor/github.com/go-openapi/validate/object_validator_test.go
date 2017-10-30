// Copyright 2017 go-swagger maintainers
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

package validate

import "testing"
import "github.com/stretchr/testify/assert"

func TestItemsMustBeTypeArray(t *testing.T) {
	ov := new(objectValidator)
	dataValid := map[string]interface{}{
		"type":  "array",
		"items": "dummy",
	}
	dataInvalid := map[string]interface{}{
		"type":  "object",
		"items": "dummy",
	}
	res := ov.Validate(dataValid)
	assert.Equal(t, 0, len(res.Errors))
	res = ov.Validate(dataInvalid)
	assert.NotEqual(t, 0, len(res.Errors))
}

// +build !go1.11

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

package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseResponseInfo(t *testing.T) {
	//Parse just untagged ref
	ref, arr, isDef, desc, err := parseTags("myModel")
	assert.NoError(t, err)
	assert.Equal(t, ref, "myModel")
	assert.Equal(t, arr, 0)
	assert.False(t, isDef)
	assert.Equal(t, desc, "")
	//Parse untagged ref & desc
	ref, arr, isDef, desc, err = parseTags("myModel11 and my desc")
	assert.NoError(t, err)
	assert.Equal(t, ref, "myModel11")
	assert.Equal(t, arr, 0)
	assert.False(t, isDef)
	assert.Equal(t, desc, "and my desc")
	//Parse tagged body & desc
	ref, arr, isDef, desc, err = parseTags("body:myModel55 description:and my desc")
	assert.NoError(t, err)
	assert.Equal(t, ref, "myModel55")
	assert.Equal(t, arr, 0)
	assert.True(t, isDef)
	assert.Equal(t, desc, "and my desc")
	//Parse tagged response & desc
	ref, arr, isDef, desc, err = parseTags("response:myModel42 and my desc")
	assert.NoError(t, err)
	assert.Equal(t, ref, "myModel42")
	assert.Equal(t, arr, 0)
	assert.False(t, isDef)
	assert.Equal(t, desc, "and my desc")
	//Parse tagged array body & desc
	ref, arr, isDef, desc, err = parseTags("body:[]myModel55 description:and my desc")
	assert.NoError(t, err)
	assert.Equal(t, ref, "myModel55")
	assert.Equal(t, arr, 1)
	assert.True(t, isDef)
	assert.Equal(t, desc, "and my desc")
	//Parse description with ignored tags after (which also test *just* description)
	ref, arr, isDef, desc, err = parseTags("description:desc response:myModel28")
	assert.NoError(t, err)
	assert.Equal(t, ref, "")
	assert.Equal(t, arr, 0)
	assert.False(t, isDef)
	assert.Equal(t, desc, "desc response:myModel28")
	//Parse just description
	ref, arr, isDef, desc, err = parseTags("description:desc")
	assert.NoError(t, err)
	assert.Equal(t, ref, "")
	assert.Equal(t, arr, 0)
	assert.False(t, isDef)
	assert.Equal(t, desc, "desc")
	//Parse unrecognized tags
	_, _, _, _, err = parseTags("body:good woah:desc response:myModel28")
	assert.NotNil(t, err)
	//Parse repeated ref tag
	_, _, _, _, err = parseTags("body:m1 response:m2")
	assert.NotNil(t, err)
}

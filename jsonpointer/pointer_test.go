// Copyright 2013 sigu-399 ( https://github.com/sigu-399 )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author       sigu-399
// author-github  https://github.com/sigu-399
// author-mail    sigu.399@gmail.com
//
// repository-name  jsonpointer
// repository-desc  An implementation of JSON Pointer - Go language
//
// description    Automated tests on package.
//
// created        03-03-2013

package jsonpointer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestDocumentNBItems = 11
	TestNodeObjNBItems  = 4
	TestDocumentString  = `{
"foo": ["bar", "baz"],
"obj": { "a":1, "b":2, "c":[3,4], "d":[ {"e":9}, {"f":[50,51]} ] },
"": 0,
"a/b": 1,
"c%d": 2,
"e^f": 3,
"g|h": 4,
"i\\j": 5,
"k\"l": 6,
" ": 7,
"m~n": 8
}`
)

var testDocumentJSON interface{}

func init() {
	json.Unmarshal([]byte(TestDocumentString), &testDocumentJSON)
}

func TestEscaping(t *testing.T) {

	ins := []string{`/`, `/`, `/a~1b`, `/a~1b`, `/c%d`, `/e^f`, `/g|h`, `/i\j`, `/k"l`, `/ `, `/m~0n`}
	outs := []float64{0, 0, 1, 1, 2, 3, 4, 5, 6, 7, 8}

	for i := range ins {

		p, err := New(ins[i])
		if err != nil {
			t.Errorf("New(%v) error %v", ins[i], err.Error())
		}

		result, _, err := p.Get(testDocumentJSON)
		if err != nil {
			t.Errorf("Get(%v) error %v", ins[i], err.Error())
		}

		if result != outs[i] {
			t.Errorf("Get(%v) = %v, expect %v", ins[i], result, outs[i])
		}
	}

}

func TestFullDocument(t *testing.T) {

	in := ``

	p, err := New(in)
	if err != nil {
		t.Errorf("New(%v) error %v", in, err.Error())
	}

	result, _, err := p.Get(testDocumentJSON)
	if err != nil {
		t.Errorf("Get(%v) error %v", in, err.Error())
	}

	if len(result.(map[string]interface{})) != TestDocumentNBItems {
		t.Errorf("Get(%v) = %v, expect full document", in, result)
	}
}

func TestGetNode(t *testing.T) {

	in := `/obj`

	p, err := New(in)
	if err != nil {
		t.Errorf("New(%v) error %v", in, err.Error())
	}

	result, _, err := p.Get(testDocumentJSON)
	if err != nil {
		t.Errorf("Get(%v) error %v", in, err.Error())
	}

	if len(result.(map[string]interface{})) != TestNodeObjNBItems {
		t.Errorf("Get(%v) = %v, expect full document", in, result)
	}
}

func TestArray(t *testing.T) {

	ins := []string{`/foo/0`, `/foo/0`, `/foo/1`}
	outs := []string{"bar", "bar", "baz"}

	for i := range ins {

		p, err := New(ins[i])
		if err != nil {
			t.Errorf("New(%v) error %v", ins[i], err.Error())
		}

		result, _, err := p.Get(testDocumentJSON)
		if err != nil {
			t.Errorf("Get(%v) error %v", ins[i], err.Error())
		}

		if result != outs[i] {
			t.Errorf("Get(%v) = %v, expect %v", ins[i], result, outs[i])
		}
	}

}

func TestOtherThings(t *testing.T) {
	_, err := New("abc")
	assert.Error(t, err)

	p, err := New("")
	assert.NoError(t, err)
	assert.Equal(t, "", p.String())

	p, err = New("/obj/a")
	assert.Equal(t, "/obj/a", p.String())

	s := Escape("m~n")
	assert.Equal(t, "m~0n", s)
	s = Escape("m/n")
	assert.Equal(t, "m~1n", s)

	p, err = New("/foo/3")
	assert.NoError(t, err)
	_, _, err = p.Get(testDocumentJSON)
	assert.Error(t, err)

	p, err = New("/foo/a")
	assert.NoError(t, err)
	_, _, err = p.Get(testDocumentJSON)
	assert.Error(t, err)

	p, err = New("/notthere")
	assert.NoError(t, err)
	_, _, err = p.Get(testDocumentJSON)
	assert.Error(t, err)

	p, err = New("/invalid")
	assert.NoError(t, err)
	_, _, err = p.Get(1234)
	assert.Error(t, err)

	p, err = New("/foo/1")
	assert.NoError(t, err)
	expected := "hello"
	_, err = p.Set(testDocumentJSON, expected)
	assert.NoError(t, err)
	v, _, err := p.Get(testDocumentJSON)
	assert.NoError(t, err)
	assert.Equal(t, expected, v)

	esc := Escape("a/")
	assert.Equal(t, "a~1", esc)
	unesc := Unescape(esc)
	assert.Equal(t, "a/", unesc)

	unesc = Unescape("~01")
	assert.Equal(t, "~1", unesc)
	assert.Equal(t, "~0~1", Escape("~/"))
	assert.Equal(t, "~/", Unescape("~0~1"))
}

func TestObject(t *testing.T) {

	ins := []string{`/obj/a`, `/obj/b`, `/obj/c/0`, `/obj/c/1`, `/obj/c/1`, `/obj/d/1/f/0`}
	outs := []float64{1, 2, 3, 4, 4, 50}

	for i := range ins {

		p, err := New(ins[i])
		if err != nil {
			t.Errorf("New(%v) error %v", ins[i], err.Error())
		}

		result, _, err := p.Get(testDocumentJSON)
		if err != nil {
			t.Errorf("Get(%v) error %v", ins[i], err.Error())
		}

		if result != outs[i] {
			t.Errorf("Get(%v) = %v, expect %v", ins[i], result, outs[i])
		}
	}

}

func TestSetNode(t *testing.T) {

	jsonText := `{"a":[{"b": 1, "c": 2}], "d": 3}`

	var jsonDocument interface{}
	json.Unmarshal([]byte(jsonText), &jsonDocument)

	in := "/a/0/c"

	p, err := New(in)
	if err != nil {
		t.Errorf("New(%v) error %v", in, err.Error())
	}

	_, err = p.Set(jsonDocument, 999)
	if err != nil {
		t.Errorf("Set(%v) error %v", in, err.Error())
	}

	firstNode := jsonDocument.(map[string]interface{})
	if len(firstNode) != 2 {
		t.Errorf("Set(%s) failed", in)
	}

	sliceNode := firstNode["a"].([]interface{})
	if len(sliceNode) != 1 {
		t.Errorf("Set(%s) failed", in)
	}

	changedNode := sliceNode[0].(map[string]interface{})
	changedNodeValue := changedNode["c"].(int)

	if changedNodeValue != 999 {
		if len(sliceNode) != 1 {
			t.Errorf("Set(%s) failed", in)
		}
	}

}

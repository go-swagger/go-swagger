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
// description    Main and unique file.
//
// created        25-02-2013

package jsonpointer

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	emptyPointer     = ``
	pointerSeparator = `/`

	invalidStart = `JSON pointer must be empty or start with a "` + pointerSeparator
)

type implStruct struct {
	mode string // "SET" or "GET"

	inDocument interface{}

	setInValue interface{}

	getOutNode interface{}
	getOutKind reflect.Kind
	outError   error
}

// New creates a new json pointer for the given string
func New(jsonPointerString string) (Pointer, error) {

	var p Pointer
	err := p.parse(jsonPointerString)
	return p, err

}

// Pointer the json pointer reprsentation
type Pointer struct {
	referenceTokens []string
}

// "Constructor", parses the given string JSON pointer
func (p *Pointer) parse(jsonPointerString string) error {

	var err error

	if jsonPointerString != emptyPointer {
		if !strings.HasPrefix(jsonPointerString, pointerSeparator) {
			err = errors.New(invalidStart)
		} else {
			referenceTokens := strings.Split(jsonPointerString, pointerSeparator)
			for _, referenceToken := range referenceTokens[1:] {
				p.referenceTokens = append(p.referenceTokens, referenceToken)
			}
		}
	}

	return err
}

// Get uses the pointer to retrieve a value from a JSON document
func (p *Pointer) Get(document interface{}) (interface{}, reflect.Kind, error) {

	is := &implStruct{mode: "GET", inDocument: document}
	p.implementation(is)
	return is.getOutNode, is.getOutKind, is.outError

}

// Set uses the pointer to update a value from a JSON document
func (p *Pointer) Set(document interface{}, value interface{}) (interface{}, error) {

	is := &implStruct{mode: "SET", inDocument: document, setInValue: value}
	p.implementation(is)
	return document, is.outError

}

// Both Get and Set functions use the same implementation to avoid code duplication
func (p *Pointer) implementation(i *implStruct) {

	kind := reflect.Invalid

	// Full document when empty
	if len(p.referenceTokens) == 0 {
		i.getOutNode = i.inDocument
		i.outError = nil
		i.getOutKind = kind
		i.outError = nil
		return
	}

	node := i.inDocument

	for ti, token := range p.referenceTokens {

		decodedToken := Unescape(token)
		isLastToken := ti == len(p.referenceTokens)-1

		rValue := reflect.ValueOf(node)
		kind = rValue.Kind()

		switch kind {

		case reflect.Map:
			m := node.(map[string]interface{})
			if _, ok := m[decodedToken]; ok {
				node = m[decodedToken]
				if isLastToken && i.mode == "SET" {
					m[decodedToken] = i.setInValue
				}
			} else {
				i.outError = fmt.Errorf("object has no key '%s'", token)
				i.getOutKind = kind
				i.getOutNode = nil
				return
			}

		case reflect.Slice:
			s := node.([]interface{})
			tokenIndex, err := strconv.Atoi(token)
			if err != nil {
				i.outError = fmt.Errorf("invalid array index '%s'", token)
				i.getOutKind = kind
				i.getOutNode = nil
				return
			}
			sLength := len(s)
			if tokenIndex < 0 || tokenIndex >= sLength {
				i.outError = fmt.Errorf("index out of bounds array[0,%d] index '%d'", sLength, tokenIndex)
				i.getOutKind = kind
				i.getOutNode = nil
				return
			}

			node = s[tokenIndex]
			if isLastToken && i.mode == "SET" {
				s[tokenIndex] = i.setInValue
			}

		default:
			i.outError = fmt.Errorf("invalid token reference '%s'", token)
			i.getOutKind = kind
			i.getOutNode = nil
			return
		}

	}

	rValue := reflect.ValueOf(node)
	kind = rValue.Kind()

	i.getOutNode = node
	i.getOutKind = kind
	i.outError = nil
}

// Pointer to string representation function
func (p *Pointer) String() string {

	if len(p.referenceTokens) == 0 {
		return emptyPointer
	}

	pointerString := pointerSeparator + strings.Join(p.referenceTokens, pointerSeparator)

	return pointerString
}

// Specific JSON pointer encoding here
// ~0 => ~
// ~1 => /
// ... and vice versa

const (
	encRefTok0 = `~0`
	encRefTok1 = `~1`
	decRefTok0 = `~`
	decRefTok1 = `/`
)

// Unescape unescapes a json pointer reference token string to the original representation
func Unescape(token string) string {
	step1 := strings.Replace(token, encRefTok1, decRefTok1, -1)
	step2 := strings.Replace(step1, encRefTok0, decRefTok0, -1)
	return step2
}

// Escape escapes a pointer reference token string
func Escape(token string) string {
	step1 := strings.Replace(token, decRefTok0, encRefTok0, -1)
	step2 := strings.Replace(step1, decRefTok1, encRefTok1, -1)
	return step2
}

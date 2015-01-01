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
// repository-name  jsonreference
// repository-desc  An implementation of JSON Reference - Go language
//
// description    Main and unique file.
//
// created        26-02-2013

package jsonreference

import (
	"errors"
	"net/url"
	"strings"

	"github.com/casualjim/go-swagger/swagger/jsonpointer"
)

const (
	const_fragment_char = `#`
)

func NewJsonReference(jsonReferenceString string) (JsonReference, error) {

	var r JsonReference
	err := r.parse(jsonReferenceString)
	return r, err

}

type JsonReference struct {
	referenceUrl     *url.URL
	referencePointer jsonpointer.JsonPointer

	HasFullUrl      bool
	HasUrlPathOnly  bool
	HasFragmentOnly bool
	HasFileScheme   bool
	HasFullFilePath bool
}

func (r *JsonReference) GetUrl() *url.URL {
	return r.referenceUrl
}

func (r *JsonReference) GetPointer() *jsonpointer.JsonPointer {
	return &r.referencePointer
}

func (r *JsonReference) String() string {

	if r.referenceUrl != nil {
		return r.referenceUrl.String()
	}

	if r.HasFragmentOnly {
		return const_fragment_char + r.referencePointer.String()
	}

	return r.referencePointer.String()
}

func (r *JsonReference) IsCanonical() bool {
	return (r.HasFileScheme && r.HasFullFilePath) || (!r.HasFileScheme && r.HasFullUrl)
}

// "Constructor", parses the given string JSON reference
func (r *JsonReference) parse(jsonReferenceString string) (err error) {

	r.referenceUrl, err = url.Parse(jsonReferenceString)
	if err != nil {
		return
	}
	refUrl := r.referenceUrl

	if refUrl.Scheme != "" && refUrl.Host != "" {
		r.HasFullUrl = true
	} else {
		if refUrl.Path != "" {
			r.HasUrlPathOnly = true
		} else if refUrl.RawQuery == "" && refUrl.Fragment != "" {
			r.HasFragmentOnly = true
		}
	}

	r.HasFileScheme = refUrl.Scheme == "file"
	r.HasFullFilePath = strings.HasPrefix(refUrl.Path, "/")

	// invalid json-pointer error means url has no json-pointer fragment. simply ignore error
	r.referencePointer, _ = jsonpointer.NewJsonPointer(refUrl.Fragment)

	return
}

// Creates a new reference from a parent and a child
// If the child cannot inherit from the parent, an error is returned
func (r *JsonReference) Inherits(child JsonReference) (*JsonReference, error) {
	childUrl := child.GetUrl()
	parentUrl := r.GetUrl()
	if childUrl == nil {
		return nil, errors.New("childUrl is nil!")
	}
	if parentUrl == nil {
		return nil, errors.New("parentUrl is nil!")
	}

	ref, err := NewJsonReference(parentUrl.ResolveReference(childUrl).String())
	if err != nil {
		return nil, err
	}
	return &ref, err
}

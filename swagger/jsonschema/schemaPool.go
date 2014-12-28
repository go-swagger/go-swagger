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

// author           sigu-399
// author-github    https://github.com/sigu-399
// author-mail      sigu.399@gmail.com
//
// repository-name  jsonschema
// repository-desc  An implementation of JSON Schema, based on IETF's draft v4 - Go language.
//
// description		Defines resources pooling.
//                  Eases referencing and avoids downloading the same resource twice.
//
// created          26-02-2013

package jsonschema

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/casualjim/go-swagger/swagger/assets"
	"github.com/xeipuuv/gojsonreference"
)

type schemaPoolDocument struct {
	Document interface{}
}

type schemaPool struct {
	schemaPoolDocuments map[string]*schemaPoolDocument
	standaloneDocument  interface{}
}

func newSchemaPool() *schemaPool {

	p := &schemaPool{}
	p.schemaPoolDocuments = make(map[string]*schemaPoolDocument)
	p.standaloneDocument = nil

	return p
}

func (p *schemaPool) SetStandaloneDocument(document interface{}) {
	p.standaloneDocument = document
}

func (p *schemaPool) GetStandaloneDocument() (document interface{}) {
	return p.standaloneDocument
}

func (p *schemaPool) GetAssetDocument(reference gojsonreference.JsonReference, path string) (*schemaPoolDocument, error) {
	internalLog(fmt.Sprintf("Get document from pool (%s) :", reference.String()))

	var err error

	// It is not possible to load anything that is not canonical...
	if !reference.IsCanonical() {
		return nil, fmt.Errorf("reference must be canonical %s", reference.String())
	}

	refToURL := reference
	refToURL.GetUrl().Fragment = ""

	var spd *schemaPoolDocument

	if d, ok := p.schemaPoolDocuments[refToURL.String()]; ok {
		return d, nil
	}

	// Load the document
	b, err := assets.Asset(path)
	if err != nil {
		return nil, err
	}
	var document interface{}
	if err := json.Unmarshal(b, &document); err != nil {
		return nil, err
	}

	spd = &schemaPoolDocument{Document: document}
	// add the document to the pool for potential later use
	p.schemaPoolDocuments[refToURL.String()] = spd

	return spd, nil
}

func (p *schemaPool) GetDocument(reference gojsonreference.JsonReference) (*schemaPoolDocument, error) {

	internalLog(fmt.Sprintf("Get document from pool (%s) :", reference.String()))

	var err error

	// It is not possible to load anything that is not canonical...
	if !reference.IsCanonical() {
		return nil, fmt.Errorf("reference must be canonical %s", reference.String())
	}

	refToURL := reference
	refToURL.GetUrl().Fragment = ""

	var spd *schemaPoolDocument

	if d, ok := p.schemaPoolDocuments[refToURL.String()]; ok {
		return d, nil
	}

	// Load the document

	var document interface{}

	if reference.HasFileScheme {

		internalLog(" Loading new document from file")

		// Load from file
		filename := strings.Replace(refToURL.String(), "file://", "", -1)
		document, err = GetFileJson(filename)
		if err != nil {
			return nil, err
		}

	} else {

		internalLog(" Loading new document from http")

		// Load from HTTP
		document, err = GetHttpJson(refToURL.String())
		if err != nil {
			return nil, err
		}

	}

	spd = &schemaPoolDocument{Document: document}
	// add the document to the pool for potential later use
	p.schemaPoolDocuments[refToURL.String()] = spd

	return spd, nil
}

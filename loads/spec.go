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

package loads

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/go-swagger/go-swagger/analysis"
	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
)

// DocLoader represents a doc loader type
type DocLoader func(string) (json.RawMessage, error)

// JSONSpec loads a spec from a json document
func JSONSpec(path string) (*Document, error) {
	data, err := swag.JSONDoc(path)
	if err != nil {
		return nil, err
	}
	// convert to json
	return Analyzed(json.RawMessage(data), "")
}

// YAMLSpec loads a swagger spec document
func YAMLSpec(path string) (*Document, error) {
	data, err := swag.YAMLDoc(path)
	if err != nil {
		return nil, err
	}

	return Analyzed(data, "")
}

// Document represents a swagger spec document
type Document struct {
	// specAnalyzer
	Analyzer *analysis.Spec
	spec     *spec.Swagger
	schema   *spec.Schema
	raw      json.RawMessage
	orig     *Document
}

// Spec loads a new spec document
func Spec(path string) (*Document, error) {
	specURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(specURL.Path)
	if ext == ".yaml" || ext == ".yml" {
		return YAMLSpec(path)
	}

	return JSONSpec(path)
}

// Analyzed creates a new analyzed spec document
func Analyzed(data json.RawMessage, version string) (*Document, error) {
	if version == "" {
		version = "2.0"
	}
	if version != "2.0" {
		return nil, fmt.Errorf("spec version %q is not supported", version)
	}

	swspec := new(spec.Swagger)
	if err := json.Unmarshal(data, swspec); err != nil {
		return nil, err
	}

	d := &Document{
		Analyzer: analysis.New(swspec),
		schema:   spec.MustLoadSwagger20Schema(),
		spec:     swspec,
		raw:      data,
	}
	// d.initialize()
	d.orig = &(*d)
	d.orig.spec = &(*swspec)
	return d, nil
}

// Expanded expands the ref fields in the spec document and returns a new spec document
func (d *Document) Expanded() (*Document, error) {
	swspec := new(spec.Swagger)
	if err := json.Unmarshal(d.raw, swspec); err != nil {
		return nil, err
	}
	if err := spec.ExpandSpec(swspec); err != nil {
		return nil, err
	}

	dd := &Document{
		Analyzer: analysis.New(swspec),
		spec:     swspec,
		schema:   spec.MustLoadSwagger20Schema(),
		raw:      d.raw,
	}
	// dd.initialize()
	dd.orig = d.orig
	dd.orig.spec = &(*d.orig.spec)

	return dd, nil
}

// BasePath the base path for this spec
func (d *Document) BasePath() string {
	return d.spec.BasePath
}

// Version returns the version of this spec
func (d *Document) Version() string {
	return d.spec.Swagger
}

// Schema returns the swagger 2.0 schema
func (d *Document) Schema() *spec.Schema {
	return d.schema
}

// Spec returns the swagger spec object model
func (d *Document) Spec() *spec.Swagger {
	return d.spec
}

// Host returns the host for the API
func (d *Document) Host() string {
	return d.spec.Host
}

// Raw returns the raw swagger spec as json bytes
func (d *Document) Raw() json.RawMessage {
	return d.raw
}

// Reload reanalyzes the spec
func (d *Document) Reload() *Document {
	orig := d.orig
	sp := *d.orig.spec
	d.Analyzer = analysis.New(&sp)
	d.orig = orig
	return d
}

// ResetDefinitions gives a shallow copy with the models reset
func (d *Document) ResetDefinitions() *Document {
	defs := make(map[string]spec.Schema, len(d.orig.spec.Definitions))
	for k, v := range d.orig.spec.Definitions {
		defs[k] = v
	}

	dd := *d
	cp := *d.orig.spec
	dd.spec = &cp
	dd.schema = spec.MustLoadSwagger20Schema()
	dd.spec.Definitions = defs
	// dd.initialize()
	dd.orig = d.orig
	return dd.Reload()
}

// Pristine creates a new pristine document instance based on the input data
func (d *Document) Pristine() *Document {
	dd, _ := Analyzed(d.Raw(), d.Version())
	return dd
}

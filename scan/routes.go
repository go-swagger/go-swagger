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
	"fmt"
	"go/ast"

	"github.com/go-openapi/spec"

	"golang.org/x/tools/go/loader"
)

func opConsumesSetter(op *spec.Operation) func([]string) {
	return func(consumes []string) { op.Consumes = consumes }
}

func opProducesSetter(op *spec.Operation) func([]string) {
	return func(produces []string) { op.Produces = produces }
}

func opSchemeSetter(op *spec.Operation) func([]string) {
	return func(schemes []string) { op.Schemes = schemes }
}

func opSecurityDefsSetter(op *spec.Operation) func([]map[string][]string) {
	return func(securityDefs []map[string][]string) { op.Security = securityDefs }
}

func opResponsesSetter(op *spec.Operation) func(*spec.Response, map[int]spec.Response) {
	return func(def *spec.Response, scr map[int]spec.Response) {
		if op.Responses == nil {
			op.Responses = new(spec.Responses)
		}
		op.Responses.Default = def
		op.Responses.StatusCodeResponses = scr
	}
}

func newRoutesParser(prog *loader.Program) *routesParser {
	return &routesParser{
		program: prog,
	}
}

type routesParser struct {
	program     *loader.Program
	definitions map[string]spec.Schema
	operations  map[string]*spec.Operation
	responses   map[string]spec.Response
}

func (rp *routesParser) Parse(gofile *ast.File, target interface{}) error {
	tgt := target.(*spec.Paths)
	for _, comsec := range gofile.Comments {
		content := parsePathAnnotation(rxRoute, comsec.List)

		if content.Method == "" {
			continue // it's not, next!
		}

		pthObj := tgt.Paths[content.Path]
		op := setPathOperation(
			content.Method, content.ID,
			&pthObj, rp.operations[content.ID])

		op.Tags = content.Tags

		sp := new(sectionedParser)
		sp.setTitle = func(lines []string) { op.Summary = joinDropLast(lines) }
		sp.setDescription = func(lines []string) { op.Description = joinDropLast(lines) }
		sr := newSetResponses(rp.definitions, rp.responses, opResponsesSetter(op))
		sp.taggers = []tagParser{
			newMultiLineTagParser("Consumes", newMultilineDropEmptyParser(rxConsumes, opConsumesSetter(op)), false),
			newMultiLineTagParser("Produces", newMultilineDropEmptyParser(rxProduces, opProducesSetter(op)), false),
			newSingleLineTagParser("Schemes", newSetSchemes(opSchemeSetter(op))),
			newMultiLineTagParser("Security", newSetSecurity(rxSecuritySchemes, opSecurityDefsSetter(op)), false),
			newMultiLineTagParser("Responses", sr, false),
		}
		if err := sp.Parse(content.Remaining); err != nil {
			return fmt.Errorf("operation (%s): %v", op.ID, err)
		}

		if tgt.Paths == nil {
			tgt.Paths = make(map[string]spec.PathItem)
		}
		tgt.Paths[content.Path] = pthObj
	}

	return nil
}

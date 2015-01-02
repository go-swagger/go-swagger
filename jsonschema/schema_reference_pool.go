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
// description      Pool of referenced schemas.
//
// created          25-06-2013

package jsonschema

import (
	"fmt"
)

type schemaReferencePool struct {
	schemaPoolDocuments map[string]*jsonSchema
}

func newSchemaReferencePool() *schemaReferencePool {

	p := &schemaReferencePool{}
	p.schemaPoolDocuments = make(map[string]*jsonSchema)

	return p
}

func (p *schemaReferencePool) GetSchema(ref string) (r *jsonSchema, o bool) {

	internalLog(fmt.Sprintf("Get schema from reference pool (%s) :", ref))

	if sch, ok := p.schemaPoolDocuments[ref]; ok {
		internalLog(fmt.Sprintf(" Found"))
		return sch, true
	}

	internalLog(fmt.Sprintf(" Not found"))
	return nil, false
}

func (p *schemaReferencePool) AddSchema(ref string, sch *jsonSchema) {

	internalLog(fmt.Sprintf("Adding schema to reference pool (%s) :", ref))
	p.schemaPoolDocuments[ref] = sch
}

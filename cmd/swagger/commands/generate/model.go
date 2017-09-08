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

package generate

import (
	"errors"
	"log"
)

// Model the generate model file command
type Model struct {
	shared
	Name        []string `long:"name" short:"n" description:"the model to generate"`
	NoValidator bool     `long:"skip-validator" description:"when present will not generate a model validator"`
	NoStruct    bool     `long:"skip-struct" description:"when present will not generate the model struct"`
	DumpData    bool     `long:"dump-data" description:"when present dumps the json for the template generator instead of generating files"`
}

// Execute generates a model file
func (m *Model) Execute(args []string) error {

	if m.DumpData && len(m.Name) > 1 {
		return errors.New("only 1 model at a time is supported for dumping data")
	}

	if m.ExistingModels != "" {
		log.Println("Warning: Ignoring existing-models flag when generating models.")
	}
	s := &Server{
		shared:         m.shared,
		Models:         m.Name,
		DumpData:       m.DumpData,
		ExcludeMain:    true,
		ExcludeSpec:    true,
		SkipSupport:    true,
		SkipOperations: true,
		SkipValidation: m.NoValidator,
		SkipModels:     m.NoStruct,
	}
	return s.Execute(args)
}

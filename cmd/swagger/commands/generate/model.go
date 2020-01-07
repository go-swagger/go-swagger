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

	"github.com/go-swagger/go-swagger/generator"
)

type modelOptions struct {
	ModelPackage               string   `long:"model-package" short:"m" description:"the package to save the models" default:"models"`
	Models                     []string `long:"model" short:"M" description:"specify a model to include in generation, repeat for multiple (defaults to all)"`
	ExistingModels             string   `long:"existing-models" description:"use pre-generated models e.g. github.com/foobar/model"`
	StrictAdditionalProperties bool     `long:"strict-additional-properties" description:"disallow extra properties when additionalProperties is set to false"`
	KeepSpecOrder              bool     `long:"keep-spec-order" description:"keep schema properties order identical to spec file"`
}

func (mo modelOptions) apply(opts *generator.GenOpts) {
	opts.ModelPackage = mo.ModelPackage
	opts.Models = mo.Models
	opts.ExistingModels = mo.ExistingModels
	opts.StrictAdditionalProperties = mo.StrictAdditionalProperties
	opts.PropertiesSpecOrder = mo.KeepSpecOrder
}

// WithModels adds the model options group
type WithModels struct {
	Models modelOptions `group:"Options for model generation"`
}

// Model the generate model file command
type Model struct {
	WithShared
	WithModels

	NoStruct bool `long:"skip-struct" description:"when present will not generate the model struct"`

	Name []string `long:"name" short:"n" description:"the model to generate, repeat for multiple (defaults to all). Same as --models"`
}

func (m Model) apply(opts *generator.GenOpts) {
	m.Shared.apply(opts)
	m.Models.apply(opts)

	opts.IncludeModel = !m.NoStruct
	opts.IncludeValidator = !m.NoStruct
}

func (m Model) log(rp string) {
	log.Printf(`Generation completed!

For this generation to compile you need to have some packages in your GOPATH:

	* github.com/go-openapi/validate
	* github.com/go-openapi/strfmt

You can get these now with: go get -u -f %s/...
`, rp)
}

func (m *Model) generate(opts *generator.GenOpts) error {
	// NOTE: at the moment, the model generator (generator.GenerateDefinition)
	// is not standalone: use server generator as a proxy
	opts.IncludeSupport = false
	opts.IncludeMain = false
	return generator.GenerateServer("", append(m.Name, m.Models.Models...), nil, opts)
}

// Execute generates a model file
func (m *Model) Execute(args []string) error {

	if m.Shared.DumpData && (len(m.Name) > 1 || len(m.Models.Models) > 1) {
		return errors.New("only 1 model at a time is supported for dumping data")
	}

	if m.Models.ExistingModels != "" {
		log.Println("warning: Ignoring existing-models flag when generating models.")
	}
	return createSwagger(m)
}

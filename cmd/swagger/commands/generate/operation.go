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

type operationOptions struct {
	Operations []string `description:"specify an operation to include, repeat for multiple (defaults to all)" long:"operation"                                 short:"O"`
	Tags       []string `description:"the tags to include, if not specified defaults to all"                  group:"operations"                               long:"tags"`
	APIPackage string   `default:"operations"                                                                 description:"the package to save the operations" long:"api-package" short:"a"`
	WithEnumCI bool     `description:"allow case-insensitive enumerations"                                    long:"with-enum-ci"`

	// tags handling
	SkipTagPackages bool `description:"skips the generation of tag-based operation packages, resulting in a flat generation" long:"skip-tag-packages"`
}

func (oo operationOptions) apply(opts *generator.GenOpts) {
	opts.Operations = oo.Operations
	opts.Tags = oo.Tags
	opts.APIPackage = oo.APIPackage
	opts.AllowEnumCI = oo.WithEnumCI
	opts.SkipTagPackages = oo.SkipTagPackages
}

// WithOperations adds the operations options group.
type WithOperations struct {
	Operations operationOptions `group:"Options for operation generation"`
}

// Operation the generate operation files command.
type Operation struct {
	WithShared
	WithOperations

	clientOptions
	serverOptions
	schemeOptions
	mediaOptions

	ModelPackage string `default:"models" description:"the package to save the models" long:"model-package" short:"m"`

	NoHandler    bool `description:"when present will not generate an operation handler"       long:"skip-handler"`
	NoStruct     bool `description:"when present will not generate the parameter model struct" long:"skip-parameters"`
	NoResponses  bool `description:"when present will not generate the response model struct"  long:"skip-responses"`
	NoURLBuilder bool `description:"when present will not generate a URL builder"              long:"skip-url-builder"`

	Name []string `description:"the operations to generate, repeat for multiple (defaults to all). Same as --operations" long:"name" short:"n"`
}

func (o Operation) apply(opts *generator.GenOpts) {
	o.Shared.apply(opts)
	o.Operations.apply(opts)
	o.clientOptions.apply(opts)
	o.serverOptions.apply(opts)
	o.schemeOptions.apply(opts)
	o.mediaOptions.apply(opts)

	opts.ModelPackage = o.ModelPackage
	opts.IncludeHandler = !o.NoHandler
	opts.IncludeResponses = !o.NoResponses
	opts.IncludeParameters = !o.NoStruct
	opts.IncludeURLBuilder = !o.NoURLBuilder
}

func (o *Operation) generate(opts *generator.GenOpts) error {
	return generator.GenerateServerOperation(append(o.Name, o.Operations.Operations...), opts)
}

func (o Operation) log(_ string) {
	log.Println(`Generation completed!

For this generation to compile you need to have some packages in your go.mod:

	* github.com/go-openapi/runtime

You can get these now with: go mod tidy`)
}

// Execute generates a model file.
func (o *Operation) Execute(_ []string) error {
	if o.Shared.DumpData && len(append(o.Name, o.Operations.Operations...)) > 1 {
		return errors.New("only 1 operation at a time is supported for dumping data")
	}

	return createSwagger(o)
}

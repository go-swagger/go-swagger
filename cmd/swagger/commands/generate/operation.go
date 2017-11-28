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
	"path/filepath"

	"github.com/go-swagger/go-swagger/generator"
)

// Operation the generate operation files command
type Operation struct {
	shared
	Name           []string `long:"name" short:"n" required:"true" description:"the operations to generate, repeat for multiple"`
	Tags           []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal      string   `short:"P" long:"principal" description:"the model to use for the security principal"`
	DefaultScheme  string   `long:"default-scheme" description:"the default scheme for this API" default:"http"`
	NoHandler      bool     `long:"skip-handler" description:"when present will not generate an operation handler"`
	NoStruct       bool     `long:"skip-parameters" description:"when present will not generate the parameter model struct"`
	NoResponses    bool     `long:"skip-responses" description:"when present will not generate the response model struct"`
	NoValidator    bool     `long:"skip-validator" description:"when present will not generate a model validator"`
	NoURLBuilder   bool     `long:"skip-url-builder" description:"when present will not generate a URL builder"`
	DumpData       bool     `long:"dump-data" description:"when present dumps the json for the template generator instead of generating files"`
	SkipFlattening bool     `long:"skip-flatten" description:"skips flattening of spec prior to generation"`
	SkipValidation bool     `long:"skip-validation" description:"skips validation of spec prior to generation"`
}

// Execute generates a model file
func (o *Operation) Execute(args []string) error {
	if o.DumpData && len(o.Name) > 1 {
		return errors.New("only 1 operation at a time is supported for dumping data")
	}

	cfg, err := readConfig(string(o.ConfigFile))
	if err != nil {
		return err
	}
	setDebug(cfg)

	opts := &generator.GenOpts{
		Spec:              string(o.Spec),
		Target:            string(o.Target),
		APIPackage:        o.APIPackage,
		ModelPackage:      o.ModelPackage,
		ServerPackage:     o.ServerPackage,
		ClientPackage:     o.ClientPackage,
		Principal:         o.Principal,
		DumpData:          o.DumpData,
		DefaultScheme:     o.DefaultScheme,
		TemplateDir:       string(o.TemplateDir),
		IncludeHandler:    !o.NoHandler,
		IncludeResponses:  !o.NoResponses,
		IncludeParameters: !o.NoStruct,
		IncludeValidator:  !o.NoValidator,
		IncludeURLBuilder: !o.NoURLBuilder,
		Tags:              o.Tags,
		FlattenSpec:       !o.SkipFlattening,
		ValidateSpec:      !o.SkipValidation,
	}

	if err = opts.EnsureDefaults(false); err != nil {
		return err
	}

	if err = configureOptsFromConfig(cfg, opts); err != nil {
		return err
	}

	if err = generator.GenerateServerOperation(o.Name, opts); err != nil {
		return err
	}

	var basepath, rp, targetAbs string

	basepath, err = filepath.Abs(".")
	if err != nil {
		return err
	}
	targetAbs, err = filepath.Abs(opts.Target)
	if err != nil {
		return err
	}
	rp, err = filepath.Rel(basepath, targetAbs)
	if err != nil {
		return err
	}

	log.Printf(`Generation completed!

For this generation to compile you need to have some packages in your GOPATH:

  * github.com/go-openapi/runtime

You can get these now with: go get -u -f %s/...
`, rp)

	return nil
}

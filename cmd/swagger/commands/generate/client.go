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
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/go-swagger/go-swagger/generator"
)

// Client the command to generate a swagger client
type Client struct {
	shared
	Name            string   `long:"name" short:"A" description:"the name of the application, defaults to a mangled value of info.title"`
	Operations      []string `long:"operation" short:"O" description:"specify an operation to include, repeat for multiple"`
	Tags            []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal       string   `long:"principal" short:"P" description:"the model to use for the security principal"`
	Models          []string `long:"model" short:"M" description:"specify a model to include, repeat for multiple"`
	DefaultScheme   string   `long:"default-scheme" description:"the default scheme for this client" default:"http"`
	DefaultProduces string   `long:"default-produces" description:"the default mime type that API operations produce" default:"application/json"`
	SkipModels      bool     `long:"skip-models" description:"no models will be generated when this flag is specified"`
	SkipOperations  bool     `long:"skip-operations" description:"no operations will be generated when this flag is specified"`
	DumpData        bool     `long:"dump-data" description:"when present dumps the json for the template generator instead of generating files"`
	SkipValidation  bool     `long:"skip-validation" description:"skips validation of spec prior to generation"`
	SkipFlattening  bool     `long:"skip-flatten" description:"skips flattening of spec prior to generation"`
}

// Execute runs this command
func (c *Client) Execute(args []string) error {
	cfg, err := readConfig(string(c.ConfigFile))
	if err != nil {
		return err
	}
	setDebug(cfg)

	if c.ExistingModels != "" {
		c.SkipModels = true
	}

	var bytebuffer []byte
	var copyrightstr string
	copyrightfile := string(c.CopyrightFile)
	if copyrightfile != "" {
		//Read the Copyright from file path in opts
		bytebuffer, err = ioutil.ReadFile(copyrightfile)
		if err != nil {
			return err
		}
		copyrightstr = string(bytebuffer)
	} else {
		copyrightstr = ""
	}

	opts := &generator.GenOpts{
		Spec: string(c.Spec),

		Target:            string(c.Target),
		APIPackage:        c.APIPackage,
		ModelPackage:      c.ModelPackage,
		ServerPackage:     c.ServerPackage,
		ClientPackage:     c.ClientPackage,
		Principal:         c.Principal,
		DefaultScheme:     c.DefaultScheme,
		DefaultProduces:   c.DefaultProduces,
		IncludeModel:      !c.SkipModels,
		IncludeValidator:  !c.SkipModels,
		IncludeHandler:    !c.SkipOperations,
		IncludeParameters: !c.SkipOperations,
		IncludeResponses:  !c.SkipOperations,
		ValidateSpec:      !c.SkipValidation,
		FlattenSpec:       !c.SkipFlattening,
		Tags:              c.Tags,
		IncludeSupport:    true,
		TemplateDir:       string(c.TemplateDir),
		DumpData:          c.DumpData,
		ExistingModels:    c.ExistingModels,
		Copyright:         copyrightstr,
	}

	if err = opts.EnsureDefaults(true); err != nil {
		return err
	}

	if err = configureOptsFromConfig(cfg, opts); err != nil {
		return err
	}

	if err = generator.GenerateClient(c.Name, c.Models, c.Operations, opts); err != nil {
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
  * golang.org/x/net/context
  * golang.org/x/net/context/ctxhttp

You can get these now with: go get -u -f %s/...
`, rp)

	return nil
}

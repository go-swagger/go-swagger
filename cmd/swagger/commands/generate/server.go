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
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-swagger/go-swagger/generator"
)

// Server the command to generate an entire server application
type Server struct {
	shared
	Name           string   `long:"name" short:"A" description:"the name of the application, defaults to a mangled value of info.title"`
	Operations     []string `long:"operation" short:"O" description:"specify an operation to include, repeat for multiple"`
	Tags           []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal      string   `long:"principal" short:"P" description:"the model to use for the security principal"`
	DefaultScheme  string   `long:"default-scheme" description:"the default scheme for this API" default:"http"`
	Models         []string `long:"model" short:"M" description:"specify a model to include, repeat for multiple"`
	SkipModels     bool     `long:"skip-models" description:"no models will be generated when this flag is specified"`
	SkipOperations bool     `long:"skip-operations" description:"no operations will be generated when this flag is specified"`
	SkipSupport    bool     `long:"skip-support" description:"no supporting files will be generated when this flag is specified"`
	ExcludeMain    bool     `long:"exclude-main" description:"exclude main function, so just generate the library"`
	ExcludeSpec    bool     `long:"exclude-spec" description:"don't embed the swagger specification"`
	WithContext    bool     `long:"with-context" description:"handlers get a context as first arg"`
	DumpData       bool     `long:"dump-data" description:"when present dumps the json for the template generator instead of generating files"`
}

// Execute runs this command
func (s *Server) Execute(args []string) error {
	cfg, err := readConfig(string(s.ConfigFile))
	if err != nil {
		return err
	}
	setDebug(cfg)

	opts := &generator.GenOpts{
		Spec:              string(s.Spec),
		Target:            string(s.Target),
		APIPackage:        s.APIPackage,
		ModelPackage:      s.ModelPackage,
		ServerPackage:     s.ServerPackage,
		ClientPackage:     s.ClientPackage,
		Principal:         s.Principal,
		DefaultScheme:     s.DefaultScheme,
		IncludeModel:      !s.SkipModels,
		IncludeValidator:  !s.SkipModels,
		IncludeHandler:    !s.SkipOperations,
		IncludeParameters: !s.SkipOperations,
		IncludeResponses:  !s.SkipOperations,
		IncludeURLBuilder: !s.SkipOperations,
		IncludeMain:       !s.ExcludeMain,
		IncludeSupport:    !s.SkipSupport,
		ExcludeSpec:       s.ExcludeSpec,
		TemplateDir:       string(s.TemplateDir),
		WithContext:       s.WithContext,
		DumpData:          s.DumpData,
		Models:            s.Models,
		Operations:        s.Operations,
		Tags:              s.Tags,
		Name:              s.Name,
	}

	if e := opts.EnsureDefaults(false); e != nil {
		return e
	}

	if e := configureOptsFromConfig(cfg, opts); e != nil {
		return e
	}

	if e := generator.GenerateServer(s.Name, s.Models, s.Operations, opts); e != nil {
		return e
	}

	rp, err := filepath.Rel(".", opts.Target)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, `Generation completed!

For this generation to compile you need to have some packages in your GOPATH:

  * github.com/go-openapi/runtime
  * github.com/tylerb/graceful
  * github.com/jessevdk/go-flags
  * golang.org/x/net/context

You can get these now with: go get -u -f %s/...
`, rp)

	return nil
}

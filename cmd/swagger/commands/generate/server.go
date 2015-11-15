
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
	"github.com/go-swagger/go-swagger/generator"
	"github.com/jessevdk/go-flags"
)

type shared struct {
	Spec          flags.Filename `long:"spec" short:"f" description:"the spec file to use" default:"./swagger.json"`
	APIPackage    string         `long:"api-package" short:"a" description:"the package to save the operations" default:"operations"`
	ModelPackage  string         `long:"model-package" short:"m" description:"the package to save the models" default:"models"`
	ServerPackage string         `long:"server-package" short:"s" description:"the package to save the server specific code" default:"restapi"`
	ClientPackage string         `long:"client-package" short:"c" description:"the package to save the client specific code" default:"client"`
	Target        flags.Filename `long:"target" short:"t" default:"./" description:"the base directory for generating the files"`
	// TemplateDir  flags.Filename `long:"template-dir"`

}

// Server the command to generate an entire server application
type Server struct {
	shared
	Name           string   `long:"name" short:"A" description:"the name of the application, defaults to a mangled value of info.title"`
	Operations     []string `long:"operation" short:"O" description:"specify an operation to include, repeat for multiple"`
	Tags           []string `long:"tags" description:"the tags to include, if not specified defaults to all"`
	Principal      string   `long:"principal" short:"P" description:"the model to use for the security principal"`
	Models         []string `long:"model" short:"M" description:"specify a model to include, repeat for multiple"`
	SkipModels     bool     `long:"skip-models" description:"no models will be generated when this flag is specified"`
	SkipOperations bool     `long:"skip-operations" description:"no operations will be generated when this flag is specified"`
	SkipSupport    bool     `long:"skip-support" description:"no supporting files will be generated when this flag is specified"`
}

// Execute runs this command
func (s *Server) Execute(args []string) error {
	opts := generator.GenOpts{
		Spec:          string(s.Spec),
		Target:        string(s.Target),
		APIPackage:    s.APIPackage,
		ModelPackage:  s.ModelPackage,
		ServerPackage: s.ServerPackage,
		ClientPackage: s.ClientPackage,
		Principal:     s.Principal,
	}

	if !s.SkipModels && (len(s.Models) > 0 || len(s.Operations) == 0) {
		if err := generator.GenerateDefinition(s.Models, true, true, opts); err != nil {
			return err
		}
	}

	if !s.SkipOperations && (len(s.Operations) > 0 || len(s.Models) == 0) {
		if err := generator.GenerateServerOperation(s.Operations, s.Tags, true, true, true, opts); err != nil {
			return err
		}
	}

	if !s.SkipSupport {
		if err := generator.GenerateSupport(s.Name, s.Models, s.Operations, opts); err != nil {
			return err
		}
	}

	return nil
}

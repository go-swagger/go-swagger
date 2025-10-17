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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-swagger/go-swagger/codescan"

	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
)

// SpecFile command to generate a swagger spec from a go application.
type SpecFile struct {
	WorkDir                 string         `default:"."                                                                                                  description:"the base path to use" long:"work-dir" short:"w"`
	BuildTags               string         `default:""                                                                                                   description:"build tags"           long:"tags"     short:"t"`
	ScanModels              bool           `description:"includes models that were annotated with 'swagger:model'"                                       long:"scan-models"                 short:"m"`
	Compact                 bool           `description:"when present, doesn't prettify the json"                                                        long:"compact"`
	Output                  flags.Filename `description:"the file to write to"                                                                           long:"output"                      short:"o"`
	Input                   flags.Filename `description:"an input swagger file with which to merge"                                                      long:"input"                       short:"i"`
	Include                 []string       `description:"include packages matching pattern"                                                              long:"include"                     short:"c"`
	Exclude                 []string       `description:"exclude packages matching pattern"                                                              long:"exclude"                     short:"x"`
	IncludeTags             []string       `description:"include routes having specified tags (can be specified many times)"                             long:"include-tag"                 short:""`
	ExcludeTags             []string       `description:"exclude routes having specified tags (can be specified many times)"                             long:"exclude-tag"                 short:""`
	ExcludeDeps             bool           `description:"exclude all dependencies of project"                                                            long:"exclude-deps"                short:""`
	SetXNullableForPointers bool           `description:"set x-nullable extension to true automatically for fields of pointer types without 'omitempty'" long:"nullable-pointers"           short:"n"`
	RefAliases              bool           `description:"transform aliased types into $ref rather than expanding their definition"                       long:"ref-aliases"                 short:"r"`
	DescWithRef             bool           `description:"allow descriptions to flow alongside $ref"                                                      long:"allow-desc-with-ref"         short:""`
	Format                  string         `choice:"yaml"                                                                                                choice:"json"                      default:"json"  description:"the format for the spec document" long:"format"`
}

// Execute runs this command.
func (s *SpecFile) Execute(args []string) error {
	if len(args) == 0 { // by default consider all the paths under the working directory
		args = []string{"./..."}
	}

	input, err := loadSpec(string(s.Input))
	if err != nil {
		return err
	}

	var opts codescan.Options
	opts.Packages = args
	opts.WorkDir = s.WorkDir
	opts.InputSpec = input
	opts.ScanModels = s.ScanModels
	opts.BuildTags = s.BuildTags
	opts.Include = s.Include
	opts.Exclude = s.Exclude
	opts.IncludeTags = s.IncludeTags
	opts.ExcludeTags = s.ExcludeTags
	opts.ExcludeDeps = s.ExcludeDeps
	opts.SetXNullableForPointers = s.SetXNullableForPointers
	opts.RefAliases = s.RefAliases
	opts.DescWithRef = s.DescWithRef
	swspec, err := codescan.Run(&opts)
	if err != nil {
		return err
	}

	return writeToFile(swspec, !s.Compact, s.Format, string(s.Output))
}

func loadSpec(input string) (*spec.Swagger, error) {
	if fi, err := os.Stat(input); err == nil {
		if fi.IsDir() {
			return nil, fmt.Errorf("expected %q to be a file not a directory", input)
		}
		sp, err := loads.Spec(input)
		if err != nil {
			return nil, err
		}
		return sp.Spec(), nil
	}
	return nil, nil
}

var defaultWriter io.Writer = os.Stdout

func writeToFile(swspec *spec.Swagger, pretty bool, format string, output string) error {
	var b []byte
	var err error

	if strings.HasSuffix(output, "yml") || strings.HasSuffix(output, "yaml") || format == "yaml" {
		b, err = marshalToYAMLFormat(swspec)
	} else {
		b, err = marshalToJSONFormat(swspec, pretty)
	}

	if err != nil {
		return err
	}

	switch output {
	case "", "-":
		_, e := fmt.Fprintf(defaultWriter, "%s\n", b)
		return e
	default:
		return os.WriteFile(output, b, 0o644) //#nosec
	}

	// #nosec
}

func marshalToJSONFormat(swspec *spec.Swagger, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(swspec, "", "  ")
	}
	return json.Marshal(swspec)
}

func marshalToYAMLFormat(swspec *spec.Swagger) ([]byte, error) {
	b, err := json.Marshal(swspec)
	if err != nil {
		return nil, err
	}

	var jsonObj any
	if err := yaml.Unmarshal(b, &jsonObj); err != nil {
		return nil, err
	}

	return yaml.Marshal(jsonObj)
}

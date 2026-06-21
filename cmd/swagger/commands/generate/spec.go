// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/codescan"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/SladkyCitron/slogcolor"
	"github.com/jessevdk/go-flags"
	"go.yaml.in/yaml/v3"
)

const (
	allFromCurrent = "./..."
)

// SpecFile command to generate a swagger spec from a go application.
type SpecFile struct {
	WorkDir                        string         `default:"."                                                                                                                                                      description:"the base path to use" long:"work-dir" short:"w"`
	BuildTags                      string         `default:""                                                                                                                                                       description:"build tags"           long:"tags"     short:"t"`
	ScanModels                     bool           `description:"includes models that were annotated with 'swagger:model'"                                                                                           long:"scan-models"                 short:"m"`
	Compact                        bool           `description:"when present, doesn't prettify the json"                                                                                                            long:"compact"`
	Output                         flags.Filename `description:"the file to write to"                                                                                                                               long:"output"                      short:"o"`
	Input                          flags.Filename `description:"an input swagger file with which to merge"                                                                                                          long:"input"                       short:"i"`
	Include                        []string       `description:"include packages matching pattern"                                                                                                                  long:"include"                     short:"c"`
	Exclude                        []string       `description:"exclude packages matching pattern"                                                                                                                  long:"exclude"                     short:"x"`
	IncludeTags                    []string       `description:"include routes having specified tags (can be specified many times)"                                                                                 long:"include-tag"                 short:""`
	ExcludeTags                    []string       `description:"exclude routes having specified tags (can be specified many times)"                                                                                 long:"exclude-tag"                 short:""`
	ExcludeDeps                    bool           `description:"exclude all dependencies of project"                                                                                                                long:"exclude-deps"                short:""`
	SetXNullableForPointers        bool           `description:"set x-nullable extension to true automatically for fields of pointer types without 'omitempty'"                                                     long:"nullable-pointers"           short:"n"`
	RefAliases                     bool           `description:"transform aliased types into $ref rather than expanding their definition"                                                                           long:"ref-aliases"                 short:"r"`
	TransparentAliases             bool           `description:"treat type aliases as completely transparent, never creating definitions for them"                                                                  long:"transparent-aliases"         short:""`
	SkipExtensions                 bool           `description:"skip generation of x-go-* go-swagger extensions"                                                                                                    long:"skip-extensions"             short:""`
	SkipEnumDescriptions           bool           `description:"controls whether descriptions of enum values in field are preserved in the main description"                                                        long:"skip-enum-desc"              short:""`
	DescWithRef                    bool           `description:"allow descriptions to flow alongside $ref"                                                                                                          long:"allow-desc-with-ref"         short:""`
	Format                         string         `choice:"yaml"                                                                                                                                                    choice:"json"                      default:"json"  description:"the format for the spec document" long:"format"`
	EmitXGoType                    bool           `description:"controls whether special extension x-go-type is emitted"                                                                                            long:"emit-x-go-type"              short:""`
	EmitHierarchicalNames          bool           `description:"controls how name conflicts are handled - this enables the last resort, failsafe method using nested definitions"                                   long:"emit-hierarchical-defs"      short:""`
	SingleLineCommentAsDescription bool           `description:"controls how single line comments are handled. Default (false): as title. When true, title is skipped and only description is hydrated"             long:"single-line-comment-desc"    short:""`
	EnableAllOfCompounding         bool           `description:"controls compounded validations & descriptions with $ref. Default is to drop. When enabled, construct a allOf compound that preserves all siblings" long:"enable-allof-compounding"    short:""`
	Colorized                      bool           `description:"enable colorized diagnostics on stderr"                                                                                                             long:"colorized"                   short:""`
	Quiet                          bool           `description:"mute diagnostics on stderr"                                                                                                                         long:"quiet"                       short:"q"`
}

// Execute runs this command.
func (s *SpecFile) Execute(args []string) error {
	if len(args) == 0 { // by default consider all the paths under the working directory
		args = []string{allFromCurrent}
	}

	var input *spec.Swagger
	if len(s.Input) > 0 {
		// load an external spec to merge into
		swspec, err := loadSpec(string(s.Input))
		if err != nil {
			return err
		}
		input = swspec
	}

	skipExt := s.SkipExtensions || os.Getenv("SWAGGER_GENERATE_EXTENSION") == "false"
	debug := os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != ""
	var logger *slog.Logger
	switch {
	case s.Quiet:
		logger = noopLogger()
	case s.Colorized:
		logger = colorizedLogger()
	default:
		logger = slog.Default()
	}

	opts := codescan.Options{
		Packages:                       args,
		WorkDir:                        s.WorkDir,
		InputSpec:                      input,
		ScanModels:                     s.ScanModels,
		BuildTags:                      s.BuildTags,
		Include:                        s.Include,
		Exclude:                        s.Exclude,
		IncludeTags:                    s.IncludeTags,
		ExcludeTags:                    s.ExcludeTags,
		ExcludeDeps:                    s.ExcludeDeps,
		SetXNullableForPointers:        s.SetXNullableForPointers,
		RefAliases:                     s.RefAliases,
		TransparentAliases:             s.TransparentAliases,
		EmitRefSiblings:                s.DescWithRef,
		SkipExtensions:                 skipExt,
		SkipEnumDescriptions:           s.SkipEnumDescriptions,
		EmitXGoType:                    s.EmitXGoType,
		SingleLineCommentAsDescription: s.SingleLineCommentAsDescription,
		SkipAllOfCompounding:           !s.EnableAllOfCompounding,
		EmitHierarchicalNames:          s.EmitHierarchicalNames,
		OnDiagnostic: func(diag codescan.Diagnostic) {
			if diag.Severity == codescan.SeverityHint && !debug {
				return
			}

			var l func(string, ...any)
			switch diag.Severity {
			case codescan.SeverityError:
				l = logger.Error
			case codescan.SeverityWarning:
				l = logger.Warn
			case codescan.SeverityHint:
				fallthrough
			default:
				l = logger.Info
			}

			pth := diag.Pos.Filename
			wkdir, err := filepath.Abs(s.WorkDir)
			if err == nil {
				pth, _ = filepath.Rel(wkdir, diag.Pos.Filename)
			}

			l(diag.Message,
				slog.String("severity", diag.Severity.String()),
				slog.String("diagnostic", string(diag.Code)),
				slog.String("file", pth),
				slog.Int("line", diag.Pos.Line),
				slog.Int("column", diag.Pos.Column),
			)
		},
	}

	swspec, err := codescan.Run(&opts)
	if err != nil {
		return err
	}

	return writeToFile(swspec, !s.Compact, s.Format, string(s.Output))
}

func loadSpec(input string) (*spec.Swagger, error) {
	fi, err := os.Stat(input)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, fmt.Errorf("expected %q to be a file not a directory", input)
	}

	sp, err := loads.Spec(input)
	if err != nil {
		return nil, err
	}

	return sp.Spec(), nil
}

var defaultWriter io.Writer = os.Stdout

const generatedFileMode os.FileMode = 0o644

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
		return os.WriteFile(output, b, generatedFileMode) //#nosec
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

func noopLogger() *slog.Logger {
	return slog.New(&noopHandler{})
}

type noopHandler struct{}

func (h noopHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (h noopHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h noopHandler) WithAttrs([]slog.Attr) slog.Handler {
	return h
}

func (h noopHandler) WithGroup(string) slog.Handler {
	return h
}

func colorizedLogger() *slog.Logger {
	return slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions))
}

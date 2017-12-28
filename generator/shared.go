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

package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"

	swaggererrors "github.com/go-openapi/errors"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"golang.org/x/tools/imports"
)

//go:generate go-bindata -mode 420 -modtime 1482416923 -pkg=generator -ignore=.*\.sw? ./templates/...

// LanguageOpts to describe a language to the code generator
type LanguageOpts struct {
	ReservedWords    []string
	BaseImportFunc   func(string) string
	reservedWordsSet map[string]struct{}
	initialized      bool
	formatFunc       func(string, []byte) ([]byte, error)
}

// Init the language option
func (l *LanguageOpts) Init() {
	if !l.initialized {
		l.initialized = true
		l.reservedWordsSet = make(map[string]struct{})
		for _, rw := range l.ReservedWords {
			l.reservedWordsSet[rw] = struct{}{}
		}
	}
}

// MangleName makes sure a reserved word gets a safe name
func (l *LanguageOpts) MangleName(name, suffix string) string {
	if _, ok := l.reservedWordsSet[swag.ToFileName(name)]; !ok {
		return name
	}
	return strings.Join([]string{name, suffix}, "_")
}

// MangleVarName makes sure a reserved word gets a safe name
func (l *LanguageOpts) MangleVarName(name string) string {
	nm := swag.ToVarName(name)
	if _, ok := l.reservedWordsSet[nm]; !ok {
		return nm
	}
	return nm + "Var"
}

// FormatContent formats a file with a language specific formatter
func (l *LanguageOpts) FormatContent(name string, content []byte) ([]byte, error) {
	if l.formatFunc != nil {
		return l.formatFunc(name, content)
	}
	return content, nil
}

func (l *LanguageOpts) baseImport(tgt string) string {
	if l.BaseImportFunc != nil {
		return l.BaseImportFunc(tgt)
	}
	return ""
}

var golang = GoLangOpts()

// GoLangOpts for rendering items as golang code
func GoLangOpts() *LanguageOpts {
	opts := new(LanguageOpts)
	opts.ReservedWords = []string{
		"break", "default", "func", "interface", "select",
		"case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch",
		"const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var",
	}
	opts.formatFunc = func(ffn string, content []byte) ([]byte, error) {
		opts := new(imports.Options)
		opts.TabIndent = true
		opts.TabWidth = 2
		opts.Fragment = true
		opts.Comments = true
		return imports.Process(ffn, content, opts)
	}
	opts.BaseImportFunc = func(tgt string) string {
		// On Windows, filepath.Abs("") behaves differently than on Unix.
		// Windows: yields an error, since Abs() does not know the volume.
		// UNIX: returns current working directory
		if tgt == "" {
			tgt = "."
		}
		tgtAbsPath, err := filepath.Abs(tgt)
		if err != nil {
			log.Fatalf("could not evaluate base import path with target \"%s\": %v", tgt, err)
		}
		var tgtAbsPathExtended string
		tgtAbsPathExtended, err = filepath.EvalSymlinks(tgtAbsPath)
		if err != nil {
			log.Fatalf("could not evaluate base import path with target \"%s\" (with symlink resolution): %v", tgtAbsPath, err)
		}

		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = filepath.Join(os.Getenv("HOME"), "go")
		}

		var pth string
		for _, gp := range filepath.SplitList(gopath) {
			// EvalSymLinks also calls the Clean
			gopathExtended, err := filepath.EvalSymlinks(gp)
			if err != nil {
				log.Fatalln(err)
			}
			gopathExtended = filepath.Join(gopathExtended, "src")
			gp = filepath.Join(gp, "src")

			// At this stage we have expanded and unexpanded target path. GOPATH is fully expanded.
			// Expanded means symlink free.
			// We compare both types of targetpath<s> with gopath.
			// If any one of them coincides with gopath , it is imperative that
			// target path lies inside gopath. How?
			// 		- Case 1: Irrespective of symlinks paths coincide. Both non-expanded paths.
			// 		- Case 2: Symlink in target path points to location inside GOPATH. (Expanded Target Path)
			//    - Case 3: Symlink in target path points to directory outside GOPATH (Unexpanded target path)

			// Case 1: - Do nothing case. If non-expanded paths match just genrate base import path as if
			//				   there are no symlinks.

			// Case 2: - Symlink in target path points to location inside GOPATH. (Expanded Target Path)
			//					 First if will fail. Second if will succeed.

			// Case 3: - Symlink in target path points to directory outside GOPATH (Unexpanded target path)
			// 					 First if will succeed and break.

			//compares non expanded path for both
			if ok, relativepath := checkPrefixAndFetchRelativePath(tgtAbsPath, gp); ok {
				pth = relativepath
				break
			}

			// Compares non-expanded target path
			if ok, relativepath := checkPrefixAndFetchRelativePath(tgtAbsPath, gopathExtended); ok {
				pth = relativepath
				break
			}

			// Compares expanded target path.
			if ok, relativepath := checkPrefixAndFetchRelativePath(tgtAbsPathExtended, gopathExtended); ok {
				pth = relativepath
				break
			}

		}

		if pth == "" {
			log.Fatalln("target must reside inside a location in the $GOPATH/src")
		}
		return pth
	}
	opts.Init()
	return opts
}

// Debug when the env var DEBUG or SWAGGER_DEBUG is not empty
// the generators will be very noisy about what they are doing
var Debug = os.Getenv("DEBUG") != "" || os.Getenv("SWAGGER_DEBUG") != ""

func findSwaggerSpec(nm string) (string, error) {
	specs := []string{"swagger.json", "swagger.yml", "swagger.yaml"}
	if nm != "" {
		specs = []string{nm}
	}
	var name string
	for _, nn := range specs {
		f, err := os.Stat(nn)
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		if err != nil && os.IsNotExist(err) {
			continue
		}
		if f.IsDir() {
			return "", fmt.Errorf("%s is a directory", nn)
		}
		name = nn
		break
	}
	if name == "" {
		return "", errors.New("couldn't find a swagger spec")
	}
	return name, nil
}

// DefaultSectionOpts for a given opts, this is used when no config file is passed
// and uses the embedded templates when no local override can be found
func DefaultSectionOpts(gen *GenOpts, client bool) {
	sec := gen.Sections
	if len(sec.Models) == 0 {
		sec.Models = []TemplateOpts{
			{
				Name:     "definition",
				Source:   "asset:model",
				Target:   "{{ joinFilePath .Target .ModelPackage }}",
				FileName: "{{ (snakize (pascalize .Name)) }}.go",
			},
		}
	}

	if len(sec.Operations) == 0 {
		if client {
			sec.Operations = []TemplateOpts{
				{
					Name:     "parameters",
					Source:   "asset:clientParameter",
					Target:   "{{ joinFilePath .Target .ClientPackage .Package }}",
					FileName: "{{ (snakize (pascalize .Name)) }}_parameters.go",
				},
				{
					Name:     "responses",
					Source:   "asset:clientResponse",
					Target:   "{{ joinFilePath .Target .ClientPackage .Package }}",
					FileName: "{{ (snakize (pascalize .Name)) }}_responses.go",
				},
			}

		} else {
			ops := []TemplateOpts{}
			if gen.IncludeParameters {
				ops = append(ops, TemplateOpts{
					Name:     "parameters",
					Source:   "asset:serverParameter",
					Target:   "{{ if eq (len .Tags) 1 }}{{ joinFilePath .Target .ServerPackage .APIPackage .Package  }}{{ else }}{{ joinFilePath .Target .ServerPackage .Package  }}{{ end }}",
					FileName: "{{ (snakize (pascalize .Name)) }}_parameters.go",
				})
			}
			if gen.IncludeURLBuilder {
				ops = append(ops, TemplateOpts{
					Name:     "urlbuilder",
					Source:   "asset:serverUrlbuilder",
					Target:   "{{ if eq (len .Tags) 1 }}{{ joinFilePath .Target .ServerPackage .APIPackage .Package  }}{{ else }}{{ joinFilePath .Target .ServerPackage .Package  }}{{ end }}",
					FileName: "{{ (snakize (pascalize .Name)) }}_urlbuilder.go",
				})
			}
			if gen.IncludeResponses {
				ops = append(ops, TemplateOpts{
					Name:     "responses",
					Source:   "asset:serverResponses",
					Target:   "{{ if eq (len .Tags) 1 }}{{ joinFilePath .Target .ServerPackage .APIPackage .Package  }}{{ else }}{{ joinFilePath .Target .ServerPackage .Package  }}{{ end }}",
					FileName: "{{ (snakize (pascalize .Name)) }}_responses.go",
				})
			}
			if gen.IncludeHandler {
				ops = append(ops, TemplateOpts{
					Name:     "handler",
					Source:   "asset:serverOperation",
					Target:   "{{ if eq (len .Tags) 1 }}{{ joinFilePath .Target .ServerPackage .APIPackage .Package  }}{{ else }}{{ joinFilePath .Target .ServerPackage .Package  }}{{ end }}",
					FileName: "{{ (snakize (pascalize .Name)) }}.go",
				})
			}
			sec.Operations = ops
		}
	}

	if len(sec.OperationGroups) == 0 {
		if client {
			sec.OperationGroups = []TemplateOpts{
				{
					Name:     "client",
					Source:   "asset:clientClient",
					Target:   "{{ joinFilePath .Target .ClientPackage .Name }}",
					FileName: "{{ (snakize (pascalize .Name)) }}_client.go",
				},
			}
		} else {
			sec.OperationGroups = []TemplateOpts{}
		}
	}

	if len(sec.Application) == 0 {
		if client {
			sec.Application = []TemplateOpts{
				{
					Name:     "facade",
					Source:   "asset:clientFacade",
					Target:   "{{ joinFilePath .Target .ClientPackage }}",
					FileName: "{{ .Name }}Client.go",
				},
			}
		} else {
			sec.Application = []TemplateOpts{
				{
					Name:       "configure",
					Source:     "asset:serverConfigureapi",
					Target:     "{{ joinFilePath .Target .ServerPackage }}",
					FileName:   "configure_{{ (snakize (pascalize .Name)) }}.go",
					SkipExists: true,
				},
				{
					Name:     "main",
					Source:   "asset:serverMain",
					Target:   "{{ joinFilePath .Target \"cmd\" (dasherize (pascalize .Name)) }}-server",
					FileName: "main.go",
				},
				{
					Name:     "embedded_spec",
					Source:   "asset:swaggerJsonEmbed",
					Target:   "{{ joinFilePath .Target .ServerPackage }}",
					FileName: "embedded_spec.go",
				},
				{
					Name:     "server",
					Source:   "asset:serverServer",
					Target:   "{{ joinFilePath .Target .ServerPackage }}",
					FileName: "server.go",
				},
				{
					Name:     "builder",
					Source:   "asset:serverBuilder",
					Target:   "{{ joinFilePath .Target .ServerPackage .Package }}",
					FileName: "{{ snakize (pascalize .Name) }}_api.go",
				},
				{
					Name:     "doc",
					Source:   "asset:serverDoc",
					Target:   "{{ joinFilePath .Target .ServerPackage }}",
					FileName: "doc.go",
				},
			}
		}
	}
	gen.Sections = sec

}

// TemplateOpts allows
type TemplateOpts struct {
	Name       string `mapstructure:"name"`
	Source     string `mapstructure:"source"`
	Target     string `mapstructure:"target"`
	FileName   string `mapstructure:"file_name"`
	SkipExists bool   `mapstructure:"skip_exists"`
	SkipFormat bool   `mapstructure:"skip_format"`
}

// SectionOpts allows for specifying options to customize the templates used for generation
type SectionOpts struct {
	Application     []TemplateOpts `mapstructure:"application"`
	Operations      []TemplateOpts `mapstructure:"operations"`
	OperationGroups []TemplateOpts `mapstructure:"operation_groups"`
	Models          []TemplateOpts `mapstructure:"models"`
}

// GenOpts the options for the generator
type GenOpts struct {
	IncludeModel       bool
	IncludeValidator   bool
	IncludeHandler     bool
	IncludeParameters  bool
	IncludeResponses   bool
	IncludeURLBuilder  bool
	IncludeMain        bool
	IncludeSupport     bool
	ExcludeSpec        bool
	DumpData           bool
	WithContext        bool
	ValidateSpec       bool
	FlattenSpec        bool
	FlattenDefinitions bool
	defaultsEnsured    bool

	Spec              string
	APIPackage        string
	ModelPackage      string
	ServerPackage     string
	ClientPackage     string
	Principal         string
	Target            string
	Sections          SectionOpts
	LanguageOpts      *LanguageOpts
	TypeMapping       map[string]string
	Imports           map[string]string
	DefaultScheme     string
	DefaultProduces   string
	DefaultConsumes   string
	TemplateDir       string
	Operations        []string
	Models            []string
	Tags              []string
	Name              string
	FlagStrategy      string
	CompatibilityMode string
	ExistingModels    string
	Copyright         string
}

// TargetPath returns the target generation path relative to the server package
// Method used by templates, e.g. with {{ .TargetPath }}
// Errors are not fatal: an empty string is returned instead
func (g *GenOpts) TargetPath() string {
	tgtAbs, err := filepath.Abs(g.Target)
	if err != nil {
		log.Printf("could not evaluate target generation path \"%s\": you must create the target directory beforehand: %v", g.Target, err)
		return ""
	}
	tgtAbs = filepath.ToSlash(tgtAbs)
	srvrAbs, err := filepath.Abs(g.ServerPackage)
	if err != nil {
		log.Printf("could not evaluate target server path \"%s\": %v", g.ServerPackage, err)
		return ""
	}
	srvrAbs = filepath.ToSlash(srvrAbs)
	tgtRel, err := filepath.Rel(srvrAbs, tgtAbs)
	if err != nil {
		log.Printf("Target path \"%s\" and server path \"%s\" are not related. You shouldn't specify an absolute path in --server-package: %v", g.Target, g.ServerPackage, err)
		return ""
	}
	tgtRel = filepath.ToSlash(tgtRel)
	return tgtRel
}

// SpecPath returns the path to the spec relative to the server package
// Method used by templates, e.g. with {{ .SpecPath }}
// Errors are not fatal: an empty string is returned instead
func (g *GenOpts) SpecPath() string {
	if strings.HasPrefix(g.Spec, "http://") || strings.HasPrefix(g.Spec, "https://") {
		return g.Spec
	}
	// Local specifications
	specAbs, err := filepath.Abs(g.Spec)
	if err != nil {
		log.Printf("could not evaluate target generation path \"%s\": you must create the target directory beforehand: %v", g.Spec, err)
		return ""
	}
	specAbs = filepath.ToSlash(specAbs)
	srvrAbs, err := filepath.Abs(g.ServerPackage)
	if err != nil {
		log.Printf("could not evaluate target server path \"%s\": %v", g.ServerPackage, err)
		return ""
	}
	srvrAbs = filepath.ToSlash(srvrAbs)
	specRel, err := filepath.Rel(srvrAbs, specAbs)
	if err != nil {
		log.Printf("Specification path \"%s\" and server path \"%s\" are not related. You shouldn't specify an absolute path in --server-package: %v", g.Spec, g.ServerPackage, err)
		return ""
	}
	specRel = filepath.ToSlash(specRel)
	return specRel
}

// EnsureDefaults for these gen opts
func (g *GenOpts) EnsureDefaults(client bool) error {
	if g.defaultsEnsured {
		return nil
	}
	DefaultSectionOpts(g, client)
	if g.LanguageOpts == nil {
		g.LanguageOpts = GoLangOpts()
	}
	g.defaultsEnsured = true
	return nil
}

func (g *GenOpts) location(t *TemplateOpts, data interface{}) (string, string, error) {
	v := reflect.Indirect(reflect.ValueOf(data))
	fld := v.FieldByName("Name")
	var name string
	if fld.IsValid() {
		log.Println("name field", fld.String())
		name = fld.String()
	}

	fldpack := v.FieldByName("Package")
	pkg := g.APIPackage
	if fldpack.IsValid() {
		log.Println("package field", fldpack.String())
		pkg = fldpack.String()
	}

	var tags []string
	tagsF := v.FieldByName("Tags")
	if tagsF.IsValid() {
		tags = tagsF.Interface().([]string)
	}

	pthTpl, err := template.New(t.Name + "-target").Funcs(FuncMap).Parse(t.Target)
	if err != nil {
		return "", "", err
	}

	fNameTpl, err := template.New(t.Name + "-filename").Funcs(FuncMap).Parse(t.FileName)
	if err != nil {
		return "", "", err
	}

	d := struct {
		Name, Package, APIPackage, ServerPackage, ClientPackage, ModelPackage, Target string
		Tags                                                                          []string
	}{
		Name:          name,
		Package:       pkg,
		APIPackage:    g.APIPackage,
		ServerPackage: g.ServerPackage,
		ClientPackage: g.ClientPackage,
		ModelPackage:  g.ModelPackage,
		Target:        g.Target,
		Tags:          tags,
	}

	// pretty.Println(data)
	var pthBuf bytes.Buffer
	if e := pthTpl.Execute(&pthBuf, d); e != nil {
		return "", "", e
	}

	var fNameBuf bytes.Buffer
	if e := fNameTpl.Execute(&fNameBuf, d); e != nil {
		return "", "", e
	}
	return pthBuf.String(), fileName(fNameBuf.String()), nil
}

func (g *GenOpts) render(t *TemplateOpts, data interface{}) ([]byte, error) {
	var templ *template.Template
	if strings.HasPrefix(strings.ToLower(t.Source), "asset:") {
		tt, err := templates.Get(strings.TrimPrefix(t.Source, "asset:"))
		if err != nil {
			return nil, err
		}
		templ = tt
	}

	if templ == nil {
		// try to load template from disk, in TemplateDir if specified
		var templateFile string
		if g.TemplateDir != "" {
			templateFile = filepath.Join(g.TemplateDir, t.Source)
		} else {
			templateFile = t.Source
		}
		content, err := ioutil.ReadFile(templateFile)
		if err != nil {
			return nil, fmt.Errorf("error while opening %s template file: %v", templateFile, err)
		}
		tt, err := template.New(t.Source).Funcs(FuncMap).Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("template parsing failed on template %s: %v", t.Name, err)
			//return nil, err
		}
		templ = tt
	}
	if templ == nil {
		return nil, fmt.Errorf("template %q not found", t.Source)
	}

	var tBuf bytes.Buffer
	if err := templ.Execute(&tBuf, data); err != nil {
		return nil, fmt.Errorf("template execution failed for template %s: %v", t.Name, err)
	}
	//if Debug {
	log.Printf("executed template %s", t.Source)
	//}

	return tBuf.Bytes(), nil
}

// Render template and write generated source code
// generated code is reformatted ("linted"), which gives an
// additional level of checking. If this step fails, the generated
// is still dumped, for template debugging purposes.
func (g *GenOpts) write(t *TemplateOpts, data interface{}) error {
	dir, fname, err := g.location(t, data)
	if err != nil {
		return fmt.Errorf("failed to resolve template location for template %s: %v", t.Name, err)
	}

	if t.SkipExists && fileExists(dir, fname) {
		if Debug {
			log.Printf("skipping generation of %s because it already exists and skip_exist directive is set for %s", filepath.Join(dir, fname), t.Name)
		}
		return nil
	}

	log.Printf("creating generated file %q in %q as %s", fname, dir, t.Name)
	content, err := g.render(t, data)
	if err != nil {
		return fmt.Errorf("failed rendering template data for %s: %v", t.Name, err)
	}

	if dir != "" {
		_, exists := os.Stat(dir)
		if os.IsNotExist(exists) {
			if Debug {
				log.Printf("creating directory %q for \"%s\"", dir, t.Name)
			}
			// Directory settings consistent with file privileges.
			// Environment's umask may alter this setup
			if e := os.MkdirAll(dir, 0755); e != nil {
				return e
			}
		}
	}

	// Conditionally format the code, unless the user wants to skip
	formatted := content
	var writeerr error

	if !t.SkipFormat {
		formatted, err = g.LanguageOpts.FormatContent(fname, content)
		if err != nil {
			log.Printf("source formatting failed on template-generated source (%q for %s). Check that your template produces valid code", filepath.Join(dir, fname), t.Name)
			writeerr = ioutil.WriteFile(filepath.Join(dir, fname), content, 0644)
			if writeerr != nil {
				return fmt.Errorf("failed to write (unformatted) file %q in %q: %v", fname, dir, writeerr)
			}
			log.Printf("unformatted generated source %q has been dumped for template debugging purposes. DO NOT build on this source!", fname)
			return fmt.Errorf("source formatting on generated source %q failed: %v", t.Name, err)
		}
	}

	writeerr = ioutil.WriteFile(filepath.Join(dir, fname), formatted, 0644)
	if writeerr != nil {
		return fmt.Errorf("failed to write file %q in %q: %v", fname, dir, writeerr)
	}
	return err
}

func fileName(in string) string {
	ext := filepath.Ext(in)
	return swag.ToFileName(strings.TrimSuffix(in, ext)) + ext
}

func (g *GenOpts) shouldRenderApp(t *TemplateOpts, app *GenApp) bool {
	switch swag.ToFileName(swag.ToGoName(t.Name)) {
	case "main":
		return g.IncludeMain
	case "embedded_spec":
		return !g.ExcludeSpec
	default:
		return true
	}
}

func (g *GenOpts) shouldRenderOperations() bool {
	return g.IncludeHandler || g.IncludeParameters || g.IncludeResponses
}

func (g *GenOpts) renderApplication(app *GenApp) error {
	log.Printf("rendering %d templates for application %s", len(g.Sections.Application), app.Name)
	for _, templ := range g.Sections.Application {
		if !g.shouldRenderApp(&templ, app) {
			continue
		}
		if err := g.write(&templ, app); err != nil {
			return err
		}
	}
	return nil
}

func (g *GenOpts) renderOperationGroup(gg *GenOperationGroup) error {
	log.Printf("rendering %d templates for operation group %s", len(g.Sections.OperationGroups), g.Name)
	for _, templ := range g.Sections.OperationGroups {
		if !g.shouldRenderOperations() {
			continue
		}

		if err := g.write(&templ, gg); err != nil {
			return err
		}
	}
	return nil
}

func (g *GenOpts) renderOperation(gg *GenOperation) error {
	log.Printf("rendering %d templates for operation %s", len(g.Sections.Operations), g.Name)
	for _, templ := range g.Sections.Operations {
		if !g.shouldRenderOperations() {
			continue
		}

		if err := g.write(&templ, gg); err != nil {
			return err
		}
	}
	return nil
}

func (g *GenOpts) renderDefinition(gg *GenDefinition) error {
	log.Printf("rendering %d templates for model %s", len(g.Sections.Models), gg.Name)
	for _, templ := range g.Sections.Models {
		if !g.IncludeModel {
			continue
		}

		if err := g.write(&templ, gg); err != nil {
			return err
		}
	}
	return nil
}

func validateSpec(path string, doc *loads.Document) (err error) {
	if doc == nil {
		if path, doc, err = loadSpec(path); err != nil {
			return err
		}
	}

	result := validate.Spec(doc, strfmt.Default)
	if result == nil {
		return nil
	}

	str := fmt.Sprintf("The swagger spec at %q is invalid against swagger specification %s. see errors :\n", path, doc.Version())
	for _, desc := range result.(*swaggererrors.CompositeError).Errors {
		str += fmt.Sprintf("- %s\n", desc)
	}
	return errors.New(str)
}

func loadSpec(specFile string) (string, *loads.Document, error) {
	// find swagger spec document, verify it exists
	specPath := specFile
	var err error
	if !strings.HasPrefix(specPath, "http") {
		specPath, err = findSwaggerSpec(specFile)
		if err != nil {
			return "", nil, err
		}
	}

	// load swagger spec
	specDoc, err := loads.Spec(specPath)
	if err != nil {
		return "", nil, err
	}
	return specPath, specDoc, nil
}

func fileExists(target, name string) bool {
	_, err := os.Stat(filepath.Join(target, name))
	return !os.IsNotExist(err)
}

func gatherModels(specDoc *loads.Document, modelNames []string) (map[string]spec.Schema, error) {
	models, mnc := make(map[string]spec.Schema), len(modelNames)
	defs := specDoc.Spec().Definitions

	if mnc > 0 {
		var unknownModels []string
		for _, m := range modelNames {
			_, ok := defs[m]
			if !ok {
				unknownModels = append(unknownModels, m)
			}
		}
		if len(unknownModels) != 0 {
			return nil, fmt.Errorf("unknown models: %s", strings.Join(unknownModels, ", "))
		}
	}
	for k, v := range defs {
		if mnc == 0 {
			models[k] = v
		}
		for _, nm := range modelNames {
			if k == nm {
				models[k] = v
			}
		}
	}
	return models, nil
}

func appNameOrDefault(specDoc *loads.Document, name, defaultName string) string {
	if strings.TrimSpace(name) == "" {
		if specDoc.Spec().Info != nil && strings.TrimSpace(specDoc.Spec().Info.Title) != "" {
			name = specDoc.Spec().Info.Title
		} else {
			name = defaultName
		}
	}
	return strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(swag.ToGoName(name), "Test"), "API"), "Test")
}

func containsString(names []string, name string) bool {
	for _, nm := range names {
		if nm == name {
			return true
		}
	}
	return false
}

type opRef struct {
	Method string
	Path   string
	Key    string
	ID     string
	Op     *spec.Operation
}

type opRefs []opRef

func (o opRefs) Len() int           { return len(o) }
func (o opRefs) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o opRefs) Less(i, j int) bool { return o[i].Key < o[j].Key }

func gatherOperations(specDoc *analysis.Spec, operationIDs []string) map[string]opRef {
	var oprefs opRefs

	for method, pathItem := range specDoc.Operations() {
		for path, operation := range pathItem {
			// nm := ensureUniqueName(operation.ID, method, path, operations)
			vv := *operation
			oprefs = append(oprefs, opRef{
				Key:    swag.ToGoName(strings.ToLower(method) + " " + path),
				Method: method,
				Path:   path,
				ID:     vv.ID,
				Op:     &vv,
			})
		}
	}

	sort.Sort(oprefs)

	operations := make(map[string]opRef)
	for _, opr := range oprefs {
		nm := opr.ID
		if nm == "" {
			nm = opr.Key
		}

		oo, found := operations[nm]
		if found && oo.Method != opr.Method && oo.Path != opr.Path {
			nm = opr.Key
		}
		if len(operationIDs) == 0 || containsString(operationIDs, opr.ID) || containsString(operationIDs, nm) {
			opr.ID = nm
			opr.Op.ID = nm
			operations[nm] = opr
		}
	}

	return operations
}

func pascalize(arg string) string {
	if len(arg) == 0 || arg[0] > '9' {
		return swag.ToGoName(arg)
	}
	if arg[0] == '+' {
		return swag.ToGoName("Plus " + arg[1:])
	}
	if arg[0] == '-' {
		return swag.ToGoName("Minus " + arg[1:])
	}

	return swag.ToGoName("Nr " + arg)
}

func pruneEmpty(in []string) (out []string) {
	for _, v := range in {
		if v != "" {
			out = append(out, v)
		}
	}
	return
}

func trimBOM(in string) string {
	return strings.Trim(in, "\xef\xbb\xbf")
}

func validateAndFlattenSpec(opts *GenOpts, specDoc *loads.Document) (*loads.Document, error) {

	var err error

	// Validate if needed
	if opts.ValidateSpec {
		if err := validateSpec(opts.Spec, specDoc); err != nil {
			return specDoc, err
		}
	}

	// Restore spec to original
	opts.Spec, specDoc, err = loadSpec(opts.Spec)
	if err != nil {
		return nil, err
	}

	absBasePath := specDoc.SpecFilePath()
	if !filepath.IsAbs(absBasePath) {
		cwd, _ := os.Getwd()
		absBasePath = filepath.Join(cwd, absBasePath)
	}

	/********************************************************************************************/
	/* Either flatten or expand should be called here before moving on the code generation part */
	/********************************************************************************************/
	if opts.FlattenSpec {
		flattenOpts := analysis.FlattenOpts{
			Expand: false,
			// BasePath must be absolute. This is guaranteed because opts.Spec is absolute
			BasePath: absBasePath,
			Spec:     analysis.New(specDoc.Spec()),
		}
		err = analysis.Flatten(flattenOpts)
	} else {
		err = spec.ExpandSpec(specDoc.Spec(), &spec.ExpandOptions{
			RelativeBase: absBasePath,
			SkipSchemas:  false,
		})
	}
	if err != nil {
		return nil, err
	}

	return specDoc, nil
}

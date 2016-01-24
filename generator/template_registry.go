package generator

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"bitbucket.org/pkg/inflect"
	"github.com/go-swagger/go-swagger/swag"
)

// FuncMap is a map with default functions for use n the templates.
// These are available in every template
var FuncMap template.FuncMap = map[string]interface{}{
	"pascalize": func(arg string) string {
		if len(arg) == 0 || arg[0] > '9' {
			return swag.ToGoName(arg)
		}

		return swag.ToGoName("Nr " + arg)
	},
	"camelize":  swag.ToJSONName,
	"humanize":  swag.ToHumanNameLower,
	"snakize":   swag.ToFileName,
	"dasherize": swag.ToCommandName,
	"pluralizeFirstWord": func(arg string) string {
		sentence := strings.Split(arg, " ")
		if len(sentence) == 1 {
			return inflect.Pluralize(arg)
		}

		return inflect.Pluralize(sentence[0]) + " " + strings.Join(sentence[1:], " ")
	},
	"json": asJSON,
	"hasInsecure": func(arg []string) bool {
		return swag.ContainsStringsCI(arg, "http") || swag.ContainsStringsCI(arg, "ws")
	},
	"hasSecure": func(arg []string) bool {
		return swag.ContainsStringsCI(arg, "https") || swag.ContainsStringsCI(arg, "wss")
	},
	"stripPackage": func(str, pkg string) string {
		parts := strings.Split(str, ".")
		strlen := len(parts)
		if strlen > 0 {
			return parts[strlen-1]
		}
		return str
	},
	"dropPackage": func(str string) string {
		parts := strings.Split(str, ".")
		strlen := len(parts)
		if strlen > 0 {
			return parts[strlen-1]
		}
		return str
	},
	"upper": func(str string) string {
		return strings.ToUpper(str)
	},
}

func NewTemplateRegistry() *TemplateRegistry {
	return &TemplateRegistry{
		funcs:            FuncMap,
		files:            make(map[string][]byte),
		fileDependencies: make(map[string][]string),
		templates:        make(map[string]TemplateDefinition),
		compiled:         make(map[string]*template.Template),
	}
}

type TemplateDefinition struct {
	Dependencies []string
	Files        []string
}

type TemplateRegistry struct {
	funcs            template.FuncMap
	files            map[string][]byte
	fileDependencies map[string][]string
	templates        map[string]TemplateDefinition
	compiled         map[string]*template.Template
}

func (t *TemplateRegistry) LoadDefaults() {
	for name, asset := range assets {
		t.AddFile(name, asset)
	}

	for name, template := range builtinTemplates {
		t.AddTemplate(name, template)
	}
}

func (t *TemplateRegistry) LoadDir(templatePath string) error {

	return filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".gotmpl") {
			assetName := strings.TrimPrefix(path, templatePath)
			if data, err := ioutil.ReadFile(path); err == nil {
				t.AddFile(assetName, data)
			}
		}
		return err
	})

}

func (t *TemplateRegistry) AddFile(name string, data []byte) {
	if t.files == nil {
		t.files = make(map[string][]byte)
	}

	t.files[name] = data

	assetDeps, found := t.fileDependencies[name]

	if !found {
		return
	}

	for _, v := range assetDeps {
		delete(t.compiled, v)
	}
}

func (t *TemplateRegistry) addAssetDependency(templateName string, definition TemplateDefinition) {

	if len(definition.Files) > 0 {
		for _, file := range definition.Files {
			t.fileDependencies[file] = append(t.fileDependencies[file], templateName)
		}

	}

	if len(definition.Dependencies) > 0 {
		for _, dep := range definition.Dependencies {
			t.addAssetDependency(dep, t.templates[dep])
		}
	}

}

func (t *TemplateRegistry) AddTemplate(name string, definition TemplateDefinition) {

	if t.templates == nil {
		t.templates = make(map[string]TemplateDefinition)
	}

	t.templates[name] = definition

	t.addAssetDependency(name, definition)

	if t.compiled == nil {
		t.compiled = make(map[string]*template.Template)
	}

	delete(t.compiled, name)
}

func (t *TemplateRegistry) parseDep(name string, templ *template.Template) (*template.Template, error) {

	def, found := t.templates[name]

	if !found {
		return templ, errors.New("Not found, " + name)
	}

	log.Println("Creating template ", name)
	templ = templ.New(name)
	if len(def.Files) > 0 {
		for _, file := range def.Files {
			if _, found := t.files[file]; !found {
				panic("Asset not loaded " + file)
			}
			templ = template.Must(templ.Parse(string(t.files[file])))
		}

	}
	if len(def.Dependencies) > 0 {

		for _, dep := range def.Dependencies {

			var err error
			templ, err = t.parseDep(dep, templ)

			if err != nil {
				return templ, err
			}
		}
	}

	return templ, nil
}

func (t *TemplateRegistry) MustGet(name string) *template.Template {

	if template, found := t.compiled[name]; found && template != nil {
		return template
	}

	definition, found := t.templates[name]

	if !found {
		panic("tried to load template " + name)
	}

	templ := template.New(name).Funcs(t.funcs)
	for _, dep := range definition.Dependencies {
		templ = template.Must(t.parseDep(dep, templ))
	}

	if len(definition.Files) > 0 {
		for _, file := range definition.Files {
			if _, found := t.files[file]; !found {
				panic("Asset not loaded " + file)
			}
			templ = template.Must(templ.Parse(string(t.files[file])))
		}

	}

	t.compiled[name] = templ

	log.Println(name, templ.DefinedTemplates())

	return templ

}

func (t *TemplateRegistry) AddFunction(name string, f interface{}) {
	if t.funcs == nil {
		t.funcs = make(map[string]interface{})
	}

	t.funcs[name] = f
}

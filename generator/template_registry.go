package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"text/template/parse"

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

func NewRepository(funcs template.FuncMap) *Repository {
	repo := Repository{
		resolved:  make(map[string]bool),
		templates: make(map[string]*template.Template),
		funcs:     funcs,
	}

	if repo.funcs == nil {

		repo.funcs = make(template.FuncMap)
	}

	return &repo
}

type Repository struct {
	templates map[string]*template.Template
	resolved  map[string]bool
	funcs     template.FuncMap
}

func (t *Repository) LoadDir(templatePath string) error {

	return filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {

		if strings.HasSuffix(path, ".gotmpl") {
			assetName := strings.TrimPrefix(path, templatePath)
			if data, err := ioutil.ReadFile(path); err == nil {
				t.AddFile(assetName, string(data))
			}
		}
		if err != nil {
			return err
		}

		return nil
	})

}

func (t *Repository) AddFile(name, data string) error {

	name = strings.TrimSuffix(filepath.Base(name), ".gotmpl")

	templ, err := template.New(name).Funcs(t.funcs).Parse(data)

	if err != nil {
		return err
	}

	// Add each defined tempalte into the cache
	for _, template := range templ.Templates() {

		t.templates[template.Name()] = template.Lookup(template.Name())
	}

	log.Println(name, templ.DefinedTemplates())
	return nil
}

func findDependencies(n parse.Node) []string {

	var deps []string

	if n == nil {
		return deps
	}
	switch node := n.(type) {
	case *parse.ListNode:
		if node != nil && node.Nodes != nil {
			for _, nn := range node.Nodes {
				deps = append(deps, findDependencies(nn)...)
			}
		}
	case *parse.IfNode:
		deps = append(deps, findDependencies(node.BranchNode.List)...)
		deps = append(deps, findDependencies(node.BranchNode.ElseList)...)
	case *parse.RangeNode:
		deps = append(deps, findDependencies(node.BranchNode.List)...)
		deps = append(deps, findDependencies(node.BranchNode.ElseList)...)
	case *parse.WithNode:
		deps = append(deps, findDependencies(node.BranchNode.List)...)
		deps = append(deps, findDependencies(node.BranchNode.ElseList)...)
	case *parse.TemplateNode:
		deps = append(deps, node.Name)
	}

	return deps

}

func (t *Repository) addDependencies(templ *template.Template) (*template.Template, error) {

	deps := findDependencies(templ.Tree.Root)

	t.resolved[templ.Name()] = true
	log.Println("Adding dependencies", templ.Name(), deps)

	for _, dep := range deps {

		log.Printf("Checking %s from %s", dep, templ.Name())

		if dep == "" {
			continue
		}
		tt := templ.Lookup(dep)

		// Check if we have it
		if tt == nil {
			tt = t.templates[dep]

			// Still dont have it return an error
			if tt == nil {
				return templ, fmt.Errorf("Could not find template %s", dep)
			}

			if !t.resolved[tt.Name()] {
				var err error
				tt, err = t.addDependencies(tt)

				if err != nil {
					return templ, err
				}
			}

			log.Println(dep, templ.DefinedTemplates(), tt.DefinedTemplates())
			// Did it get resolved when loading deps?
			if loaded := templ.Lookup(dep); loaded != nil {
				continue
			}

			var err error
			templ, err = tt.AddParseTree(templ.Name(), templ.Tree)

			if err != nil {
				return templ, fmt.Errorf("Dependency Error: %v", err)
			}

		}

		if !t.resolved[tt.Name()] {
			var err error
			tt, err = t.addDependencies(tt)

			if err != nil {
				return templ, err
			}
		}
	}

	log.Println("Loaded deps:", templ.Name(), templ.DefinedTemplates())
	return templ, nil
}

func (t *Repository) Get(name string) (*template.Template, error) {

	log.Println("Getting", name)
	templ, found := t.templates[name]

	if !found {
		return templ, fmt.Errorf("Template doesn't exist", name)
	}

	return t.addDependencies(templ)
}

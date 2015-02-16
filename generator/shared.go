package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
	"golang.org/x/tools/imports"
)

//go:generate go-bindata -pkg=generator ./templates/...

var reservedGoWords = []string{
	"break", "default", "func", "interface", "select",
	"case", "defer", "go", "map", "struct",
	"chan", "else", "goto", "package", "switch",
	"const", "fallthrough", "if", "range", "type",
	"continue", "for", "import", "return", "var",
}

var defaultGoImports = []string{
	"bool", "int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"float32", "float64", "interface{}", "string",
	"byte", "rune",
}

func findSwaggerSpec(name string) (string, error) {
	f, err := os.Stat(name)
	if err != nil {
		return "", err
	}
	if f.IsDir() {
		return "", fmt.Errorf("%s is a directory", name)
	}
	return name, nil
}

// GenOpts the options for the generator
type GenOpts struct {
	Spec         string
	APIPackage   string
	ModelPackage string
	Target       string
	TypeMapping  map[string]string
	Imports      map[string]string
}

type generatorOptions struct {
	ModelPackage    string
	TargetDirectory string
}

type propertyDescriptor struct {
	PropertyName      string `json:"propertyName,omitempty"`
	ParamName         string `json:"paramName,omitempty"`
	Path              string `json:"path,omitempty"` // language escaped string or expression
	ValueExpression   string `json:"valueExpression,omitempty"`
	IndexVar          string `json:"indexVar,omitempty"`
	IsPrimitive       bool   `json:"isPrimitive,omitempty"`       // plain old primitive type
	IsCustomFormatter bool   `json:"isCustomFormatter,omitempty"` // custom format or default format
	IsContainer       bool   `json:"isContainer,omitempty"`       // slice

}

type commonValidations struct {
	propertyDescriptor
	Required         bool          `json:"required,omitempty"`
	Type             string        `json:"type,omitempty"`
	Format           string        `json:"format,omitempty"`
	Items            *spec.Items   `json:"items,omitempty"`
	Default          interface{}   `json:"default,omitempty"`
	Maximum          *float64      `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64      `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int64        `json:"maxLength,omitempty"`
	MinLength        *int64        `json:"minLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         *int64        `json:"maxItems,omitempty"`
	MinItems         *int64        `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	MultipleOf       *float64      `json:"multipleOf,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
}

type genValidations struct {
	Type                string  `json:"type,omitempty"`
	Required            bool    `json:"required,omitempty"`
	DefaultValue        string  `json:"defaultValue,omitempty"`
	MaxLength           int64   `json:"maxLength,omitempty"`
	MinLength           int64   `json:"minLength,omitempty"`
	Pattern             string  `json:"pattern,omitempty"`
	MultipleOf          float64 `json:"multipleOf,omitempty"`
	Minimum             float64 `json:"minimum,omitempty"`
	Maximum             float64 `json:"maximum,omitempty"`
	ExclusiveMinimum    bool    `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum    bool    `json:"exclusiveMaximum,omitempty"`
	Enum                string  `json:"enum,omitempty"`
	HasValidations      bool    `json:"hasValidations,omitempty"`
	Format              string  `json:"format,empty"`
	MinItems            int64   `json:"minItems,omitempty"`
	MaxItems            int64   `json:"maxItems,omitempty"`
	UniqueItems         bool    `json:"uniqueItems,omitempty"`
	HasSliceValidations bool    `json:"hasSliceValidations,omitempty"`
	NeedsSize           bool    `json:"needsSize,omitempty"`
}

func loadSpec(specFile string) (string, *spec.Document, error) {
	// find swagger spec document, verify it exists
	specPath, err := findSwaggerSpec(specFile)
	if err != nil {
		return "", nil, err
	}

	// load swagger spec
	specDoc, err := spec.Load(specPath)
	if err != nil {
		return "", nil, err
	}
	return specPath, specDoc, nil
}

func writeToFile(target, name string, content []byte) error {
	opts := new(imports.Options)
	opts.TabIndent = true
	opts.TabWidth = 2
	opts.Fragment = true
	opts.Comments = true

	ffn := util.ToFileName(name) + ".go"
	res, err := imports.Process(ffn, content, opts)

	if err != nil {
		log.Println(err)
		return writeFile(target, ffn, content)
	}

	return writeFile(target, ffn, res)
}

func writeFile(target, ffn string, content []byte) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(target, ffn), content, 0644)
}

func commentedLines(str string) string {
	lines := strings.Split(str, "\n")
	var commented []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			if !strings.HasPrefix(strings.TrimSpace(line), "//") {
				commented = append(commented, "// "+line)
			} else {
				commented = append(commented, line)
			}
		}
	}
	return strings.Join(commented, "\n")
}

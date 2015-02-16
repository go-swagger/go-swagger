package generator

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

// Processor creates codegen models for the various stages of generation
// this allows users to plug in naming conventions specific to their language
// the default of this implementation is to follow java naming conventions
// but for the go generator it has a custom params factory
type Processor interface {
	EscapeKeyword(string) string
	GetPropertyTypeDeclaration(*Property) string
	GetStringTypeDeclaration(string) string
	APIFilename(string) string
	ModelFileName(string) string
	ModelImport(string) string
	APIImport(string) string
}

// Converter converts a swagger spec document into codegen operations and models
type Converter interface {
	FromOperation(string, string, *spec.Operation) Operation
	FromModel(string, *spec.Schema) Model
	FromParameter(*spec.Parameter) Parameter
	FromHeader(*spec.Header) Property
	FromResponse(int, *spec.Schema) Response
}

type goProcessor struct {
	opts *Options
}

func (g *goProcessor) EscapeKeyword(keyword string) string {
	return "_" + keyword
}
func (g *goProcessor) GetPropertyTypeDeclaration(property *Property) string {
	return property.Name
}
func (g *goProcessor) GetStringTypeDeclaration(name string) string {
	return name
}
func (g *goProcessor) APIFilename(name string) string {
	return util.ToFileName(name)
}
func (g *goProcessor) ModelFileName(name string) string {
	return util.ToFileName(name)
}
func (g *goProcessor) ModelImport(name string) string {
	return filepath.Join(g.opts.BaseImport, g.opts.ModelPackage)
}
func (g *goProcessor) APIImport(name string) string {
	return filepath.Join(g.opts.BaseImport, g.opts.APIPackage)
}

type goConverter struct {
	opts      *Options
	processor Processor
}

func (g *goConverter) FromOperation(method, path string, operation *spec.Operation) Operation {
	var cg Operation

	return cg
}

func (g *goConverter) FromModel(name string, schema *spec.Schema) Model {
	if util.ContainsStringsCI(g.opts.ReservedWords, name) {
		name = g.processor.EscapeKeyword(name)
	}
	var model Model
	model.ClassName = util.ToGoName(name)
	model.ClassVarName = util.ToGoName(name)
	model.Name = name
	model.Description = schema.Description
	b, _ := json.MarshalIndent(schema, "", "  ")
	model.ModelJSON = string(b)
	model.ExternalDocs = schema.ExternalDocs
	model.DefaultValue = fmt.Sprintf("%v", schema.Default)

	return model
}
func (g *goConverter) FromParameter(parameter *spec.Parameter) Parameter {
	return Parameter{}
}
func (g *goConverter) FromHeader(header *spec.Header) Property {
	return Property{}
}
func (g *goConverter) FromResponse(code string, schema *spec.Schema) Response {
	var r Response
	r.Code = code
	if code == "default" {
		r.Code = "0"
	}
	r.Message = schema.Description
	r.Schema = schema

	return r
}

// GoProcessor creates a go processor for the specified options
func GoProcessor(opts *Options) Processor {
	return &goProcessor{opts}
}

// gen generates the code based on the provided options
type gen struct {
	Spec    *spec.Document
	Options *Options
	// Processor
	Processor Processor
	Converter Converter
}

// Generate generates code based on the options
func Generate(doc *spec.Document, opts *Options, processor Processor, converter Converter) error {
	// generate models in the models package at the specified directory, this includes validators
	// generate parameters in the api/parameters package at the specified directory, this includes validators and for a server
	// that also means a request binder
	// generate operations in the api package at the specified directory
	// generate the supporting files in the api package at the specified directory
	return nil
}

// Options language specific options for a generator
type Options struct {
	// BaseImport the import clause to prefix other import paths
	BaseImport string
	// ReservedWords for a language, these will be escaped by the processor
	ReservedWords []string
	// OutputDir the destination folder
	OutputDir string
	// TemplateDir the template directory
	TemplateDir string
	// FileSuffix the file suffix to use, defaults to swagger
	FileSuffix string
	// HelpString the additional help to show for this generator
	HelpString string
	// Name the name of the generator
	Name string
	// APIPackage the package for the apis
	APIPackage string
	// ModelPackage the package for the models
	ModelPackage string
	// APIDir the directory for the api handlers
	APIDir string
	// ModelDir the directory for the model files
	ModelDir string
	// DefaultIncludes the default imports
	DefaultIncludes []string
	// TypeMapping maps a swagger type name to a language specific type name
	// only required when the default algorithm isn't sufficient
	TypeMapping map[string]string
	// InstantiationTypes the type to use when creating a new instance of the
	// specified type, this allows for languages like java to use an interface
	// and a concrete type.
	InstantiationTypes map[string]string
	// ImportMapping maps a language specific type name to an import clause
	ImportMapping map[string]string
	// APITemplateFiles the files to use when generting an api handler and their extension/suffix
	APITemplateFiles map[string]string
	// ModelTemplateFiles the files to use when generating models and their extenstion/suffix
	ModelTemplateFiles map[string]string
	// AdditionalProperties additional items to have present in the template
	AdditionalProperties map[string]interface{}
}

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

var goTypeMapping = map[string]string{
	"array":    "[]",
	"map":      "map",
	"List":     "[]",
	"boolean":  "bool",
	"int":      "int32",
	"float":    "float32",
	"number":   "inf.Dec",
	"DateTime": "swagger.DateTime",
	"long":     "int64",
	"short":    "int16",
	"char":     "rune",
	"double":   "float64",
	"object":   "interface{}",
	"integer":  "int32",
}

var goImports = map[string]string{
	"inf.Dec":   "speter.net/go/exp/math/dec/inf",
	"big.Int":   "math/big",
	"swagger.*": "github.com/casualjim/go-swagger",
}

// NewGoOptions creates the options for generating go code
func NewGoOptions() *Options {
	opts := Options{
		ReservedWords:        reservedGoWords,
		FileSuffix:           "swagger",
		DefaultIncludes:      defaultGoImports,
		TypeMapping:          goTypeMapping,
		ImportMapping:        goImports,
		APIPackage:           "operations",
		APIDir:               "operations",
		ModelPackage:         "models",
		ModelDir:             "models",
		APITemplateFiles:     make(map[string]string),
		ModelTemplateFiles:   make(map[string]string),
		AdditionalProperties: make(map[string]interface{}),
		InstantiationTypes:   make(map[string]string),
	}

	return &opts
}

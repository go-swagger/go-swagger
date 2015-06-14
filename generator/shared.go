package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
	"golang.org/x/tools/imports"
)

//go:generate go-bindata -pkg=generator -ignore=.*\.sw? ./templates/...

var (
	modelTemplate          *template.Template
	modelValidatorTemplate *template.Template
	operationTemplate      *template.Template
	parameterTemplate      *template.Template
	builderTemplate        *template.Template
	mainTemplate           *template.Template
	configureAPITemplate   *template.Template
	clientTemplate         *template.Template
	clientParamTemplate    *template.Template
	clientResponseTemplate *template.Template
	clientFacadeTemplate   *template.Template
)

func init() {

	funcs := make(template.FuncMap)
	funcs["pascalize"] = swag.ToGoName
	funcs["camelize"] = swag.ToJSONName
	funcs["humanize"] = swag.ToHumanNameLower
	funcs["dasherize"] = swag.ToCommandName

	// partial templates
	pv, _ := Asset("templates/validation/primitive.gotmpl")
	cv, _ := Asset("templates/validation/customformat.gotmpl")
	validatorTempl := template.Must(template.New("primitivevalidator").Funcs(funcs).Parse(string(pv)))
	validatorTempl = template.Must(validatorTempl.New("customformatvalidator").Parse(string(cv)))

	// common templates
	bv, _ := Asset("templates/modelvalidator.gotmpl")
	modelValidatorTemplate = template.Must(validatorTempl.Clone())
	modelValidatorTemplate = template.Must(modelValidatorTemplate.New("modelvalidator").Parse(string(bv)))

	sf, _ := Asset("templates/structfield.gotmpl")
	bm, _ := Asset("templates/model.gotmpl")
	modelTemplate = template.Must(template.New("structfield").Parse(string(sf)))
	modelTemplate = template.Must(modelTemplate.New("model").Funcs(funcs).Parse(string(bm)))

	// server templates
	bp, _ := Asset("templates/server/parameter.gotmpl")
	parameterTemplate = template.Must(validatorTempl.Clone())
	parameterTemplate = template.Must(parameterTemplate.New("parameter").Parse(string(bp)))

	bo, _ := Asset("templates/server/operation.gotmpl")
	operationTemplate = template.Must(template.New("operation").Funcs(funcs).Parse(string(bo)))

	bu, _ := Asset("templates/server/builder.gotmpl")
	builderTemplate = template.Must(template.New("builder").Funcs(funcs).Parse(string(bu)))

	bc, _ := Asset("templates/server/configureapi.gotmpl")
	configureAPITemplate = template.Must(template.New("configureapi").Funcs(funcs).Parse(string(bc)))

	mn, _ := Asset("templates/server/main.gotmpl")
	mainTemplate = template.Must(template.New("main").Funcs(funcs).Parse(string(mn)))

	// Client templates
	cp, _ := Asset("templates/client/parameter.gotmpl")
	clientParamTemplate = template.Must(validatorTempl.Clone())
	clientParamTemplate = template.Must(clientParamTemplate.New("parameter").Parse(string(cp)))

	ch, _ := Asset("templates/client/response.gotmpl")
	clientResponseTemplate = template.Must(validatorTempl.Clone())
	clientResponseTemplate = template.Must(clientResponseTemplate.New("response").Parse(string(ch)))

	cr, _ := Asset("templates/client/client.gotmpl")
	clientTemplate = template.Must(template.New("client").Funcs(funcs).Parse(string(cr)))

	cf, _ := Asset("templates/client/facade.gotmpl")
	clientFacadeTemplate = template.Must(template.New("facade").Funcs(funcs).Parse(string(cf)))
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
	Spec          string
	APIPackage    string
	ModelPackage  string
	ServerPackage string
	ClientPackage string
	Principal     string
	Target        string
	TypeMapping   map[string]string
	Imports       map[string]string
	DumpData      bool
}

type generatorOptions struct {
	ModelPackage    string
	TargetDirectory string
}

type propertyDescriptor struct {
	PropertyName      string //`json:"propertyName,omitempty"`
	ParamName         string //`json:"paramName,omitempty"`
	Path              string //`json:"path,omitempty"` // language escaped string or expression
	ValueExpression   string //`json:"valueExpression,omitempty"`
	IndexVar          string //`json:"indexVar,omitempty"`
	IsPrimitive       bool   //`json:"isPrimitive,omitempty"`       // plain old primitive type
	IsCustomFormatter bool   //`json:"isCustomFormatter,omitempty"` // custom format or default format
	IsContainer       bool   //`json:"isContainer,omitempty"`       // slice
	IsMap             bool   // json:"isMap,omitempty"
}

type commonValidations struct {
	propertyDescriptor
	Required         bool          //`json:"required,omitempty"`
	Type             string        //`json:"type,omitempty"`
	Format           string        //`json:"format,omitempty"`
	Items            *spec.Items   //`json:"items,omitempty"`
	Default          interface{}   //`json:"default,omitempty"`
	Maximum          *float64      //`json:"maximum,omitempty"`
	ExclusiveMaximum bool          //`json:"exclusiveMaximum,omitempty"`
	Minimum          *float64      //`json:"minimum,omitempty"`
	ExclusiveMinimum bool          //`json:"exclusiveMinimum,omitempty"`
	MaxLength        *int64        //`json:"maxLength,omitempty"`
	MinLength        *int64        //`json:"minLength,omitempty"`
	Pattern          string        //`json:"pattern,omitempty"`
	MaxItems         *int64        //`json:"maxItems,omitempty"`
	MinItems         *int64        //`json:"minItems,omitempty"`
	UniqueItems      bool          //`json:"uniqueItems,omitempty"`
	MultipleOf       *float64      //`json:"multipleOf,omitempty"`
	Enum             []interface{} //`json:"enum,omitempty"`
}

type genValidations struct {
	Type                string  //`json:"type,omitempty"`
	Required            bool    //`json:"required,omitempty"`
	DefaultValue        string  //`json:"defaultValue,omitempty"`
	MaxLength           int64   //`json:"maxLength,omitempty"`
	MinLength           int64   //`json:"minLength,omitempty"`
	Pattern             string  //`json:"pattern,omitempty"`
	MultipleOf          float64 //`json:"multipleOf,omitempty"`
	Minimum             float64 //`json:"minimum,omitempty"`
	Maximum             float64 //`json:"maximum,omitempty"`
	ExclusiveMinimum    bool    //`json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum    bool    //`json:"exclusiveMaximum,omitempty"`
	Enum                string  //`json:"enum,omitempty"`
	HasValidations      bool    //`json:"hasValidations,omitempty"`
	Format              string  //`json:"format,empty"`
	MinItems            int64   //`json:"minItems,omitempty"`
	MaxItems            int64   //`json:"maxItems,omitempty"`
	UniqueItems         bool    //`json:"uniqueItems,omitempty"`
	HasSliceValidations bool    //`json:"hasSliceValidations,omitempty"`
	NeedsSize           bool    //`json:"needsSize,omitempty"`
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

func fileExists(target, name string) bool {
	ffn := swag.ToFileName(name) + ".go"
	_, err := os.Stat(filepath.Join(target, ffn))
	return !os.IsNotExist(err)
}

func writeToFileIfNotExist(target, name string, content []byte) error {
	if fileExists(target, name) {
		return nil
	}
	return writeToFile(target, name, content)
}

func formatGoFile(ffn string, content []byte) ([]byte, error) {
	opts := new(imports.Options)
	opts.TabIndent = true
	opts.TabWidth = 2
	opts.Fragment = true
	opts.Comments = true

	return imports.Process(ffn, content, opts)
}

func writeToFile(target, name string, content []byte) error {
	ffn := swag.ToFileName(name) + ".go"
	res, err := formatGoFile(ffn, content)
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

func gatherModels(specDoc *spec.Document, modelNames []string) map[string]spec.Schema {
	models, mnc := make(map[string]spec.Schema), len(modelNames)
	for k, v := range specDoc.Spec().Definitions {
		for _, nm := range modelNames {
			if mnc == 0 || k == nm {
				models[k] = v
			}
		}
	}
	return models
}

func appNameOrDefault(specDoc *spec.Document, name, defaultName string) string {
	if name == "" {
		if specDoc.Spec().Info != nil && specDoc.Spec().Info.Title != "" {
			name = specDoc.Spec().Info.Title
		} else {
			name = defaultName
		}
	}
	return swag.ToGoName(name)
}

func gatherOperations(specDoc *spec.Document, operationIDs []string) map[string]spec.Operation {
	operations := make(map[string]spec.Operation)
	if len(operationIDs) == 0 {
		for _, k := range specDoc.OperationIDs() {
			if op, ok := specDoc.OperationForName(k); ok {
				operations[k] = *op
			}
		}
	} else {
		for _, k := range specDoc.OperationIDs() {
			for _, nm := range operationIDs {
				if k == nm {
					if op, ok := specDoc.OperationForName(k); ok {
						operations[k] = *op
					}
				}
			}
		}
	}
	return operations
}

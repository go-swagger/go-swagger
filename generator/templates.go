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
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"bitbucket.org/pkg/inflect"
	"github.com/go-swagger/go-swagger/swag"
)

//go:generate go-bindata -pkg=generator -ignore=.*\.sw? ./templates/...

// fwiw, don't get attached to this, still requires a better abstraction

var (
	modelTemplate          *template.Template
	modelValidatorTemplate *template.Template
	operationTemplate      *template.Template
	parameterTemplate      *template.Template
	responsesTemplate      *template.Template
	builderTemplate        *template.Template
	mainTemplate           *template.Template
	mainDocTemplate        *template.Template
	embeddedSpecTemplate   *template.Template
	configureAPITemplate   *template.Template
	clientTemplate         *template.Template
	clientParamTemplate    *template.Template
	clientResponseTemplate *template.Template
	clientFacadeTemplate   *template.Template
)

var assets = map[string][]byte{

	"validation/primitive.gotmpl":           MustAsset("templates/validation/primitive.gotmpl"),
	"validation/customformat.gotmpl":        MustAsset("templates/validation/customformat.gotmpl"),
	"docstring.gotmpl":                      MustAsset("templates/docstring.gotmpl"),
	"validation/structfield.gotmpl":         MustAsset("templates/validation/structfield.gotmpl"),
	"modelvalidator.gotmpl":                 MustAsset("templates/modelvalidator.gotmpl"),
	"structfield.gotmpl":                    MustAsset("templates/structfield.gotmpl"),
	"tupleserializer.gotmpl":                MustAsset("templates/tupleserializer.gotmpl"),
	"additionalpropertiesserializer.gotmpl": MustAsset("templates/additionalpropertiesserializer.gotmpl"),
	"schematype.gotmpl":                     MustAsset("templates/schematype.gotmpl"),
	"schemabody.gotmpl":                     MustAsset("templates/schemabody.gotmpl"),
	"schema.gotmpl":                         MustAsset("templates/schema.gotmpl"),
	"schemavalidator.gotmpl":                MustAsset("templates/schemavalidator.gotmpl"),
	"model.gotmpl":                          MustAsset("templates/model.gotmpl"),
	"header.gotmpl":                         MustAsset("templates/header.gotmpl"),
	"swagger_json_embed.gotmpl":             MustAsset("templates/swagger_json_embed.gotmpl"),

	"server/parameter.gotmpl":    MustAsset("templates/server/parameter.gotmpl"),
	"server/responses.gotmpl":    MustAsset("templates/server/responses.gotmpl"),
	"server/operation.gotmpl":    MustAsset("templates/server/operation.gotmpl"),
	"server/builder.gotmpl":      MustAsset("templates/server/builder.gotmpl"),
	"server/configureapi.gotmpl": MustAsset("templates/server/configureapi.gotmpl"),
	"server/main.gotmpl":         MustAsset("templates/server/main.gotmpl"),
	"server/doc.gotmpl":          MustAsset("templates/server/doc.gotmpl"),

	"client/parameter.gotmpl": MustAsset("templates/client/parameter.gotmpl"),
	"client/response.gotmpl":  MustAsset("templates/client/response.gotmpl"),
	"client/client.gotmpl":    MustAsset("templates/client/client.gotmpl"),
	"client/facade.gotmpl":    MustAsset("templates/client/facade.gotmpl"),
}

var (
	notNumberExp = regexp.MustCompile("[^0-9]")
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

func loadCustomTemplates(templatePath, prefix string) (bool, error) {

	recompile := false

	files, err := ioutil.ReadDir(templatePath)

	if err != nil {
		return false, err
	}

	for _, file := range files {
		templateName := file.Name()

		if prefix != "" {
			templateName = prefix + "/" + file.Name()
		}

		if !file.IsDir() {

			if _, exists := assets[templateName]; exists {
				if data, err := ioutil.ReadFile(filepath.Join(templatePath, file.Name())); err == nil {
					log.Printf("Using custom template for %s\n", templateName)
					assets[templateName] = data
					recompile = true
				}

			}
		} else {
			recompile, _ = loadCustomTemplates(filepath.Join(templatePath, file.Name()), templateName)
		}

	}

	return recompile, nil
}

func compileTemplates() {
	log.Println("compiling templates")
	// partial templates
	validatorTempl := template.Must(template.New("primitivevalidator").Funcs(FuncMap).Parse(string(assets["validation/primitive.gotmpl"])))
	validatorTempl = template.Must(validatorTempl.New("customformatvalidator").Parse(string(assets["validation/customformat.gotmpl"])))

	modelTemplate = makeModelTemplate()
	// common templates
	bv, _ := Asset("templates/modelvalidator.gotmpl") // about to be gobbled up by the model template
	modelValidatorTemplate = template.Must(validatorTempl.Clone())
	modelValidatorTemplate = template.Must(modelValidatorTemplate.New("modelvalidator.gotmpl").Parse(string(bv)))

	// server templates
	parameterTemplate = makeModelTemplate()
	//parameterTemplate = template.Must(parameterTemplate.New("docstring.gotmpl").Parse(string(assets["docstring.gotmpl"])))
	//parameterTemplate = template.Must(parameterTemplate.New("validationDocString").Parse(string(assets["validation/structfield.gotmpl"])))
	//parameterTemplate = template.Must(parameterTemplate.New("schType").Parse(string(assets["schematype.gotmpl"])))
	//parameterTemplate = template.Must(parameterTemplate.New("body").Parse(string(assets["schemabody.gotmpl"])))
	parameterTemplate = template.Must(parameterTemplate.New("parameter").Parse(string(assets["server/parameter.gotmpl"])))

	responsesTemplate = makeModelTemplate()
	responsesTemplate = template.Must(responsesTemplate.New("responses").Parse(string(assets["server/responses.gotmpl"])))

	operationTemplate = makeModelTemplate()
	operationTemplate = template.Must(operationTemplate.New("operation").Parse(string(assets["server/operation.gotmpl"])))
	builderTemplate = template.Must(template.New("builder").Funcs(FuncMap).Parse(string(assets["server/builder.gotmpl"])))
	configureAPITemplate = template.Must(template.New("configureapi").Funcs(FuncMap).Parse(string(assets["server/configureapi.gotmpl"])))
	mainTemplate = template.Must(template.New("main").Funcs(FuncMap).Parse(string(assets["server/main.gotmpl"])))
	mainDocTemplate = template.Must(template.New("meta").Funcs(FuncMap).Parse(string(assets["server/doc.gotmpl"])))

	embeddedSpecTemplate = template.Must(template.New("embedded_spec").Funcs(FuncMap).Parse(string(assets["swagger_json_embed.gotmpl"])))

	// Client templates
	clientParamTemplate = makeModelTemplate()
	//clientParamTemplate = template.Must(clientParamTemplate.New("docstring.gotmpl").Parse(string(assets["docstring.gotmpl"])))
	//clientParamTemplate = template.Must(clientParamTemplate.New("validationDocString").Parse(string(assets["validation/structfield.gotmpl"])))
	//clientParamTemplate = template.Must(clientParamTemplate.New("schType").Parse(string(assets["schematype.gotmpl"])))
	//clientParamTemplate = template.Must(clientParamTemplate.New("body").Parse(string(assets["schemabody.gotmpl"])))
	clientParamTemplate = template.Must(clientParamTemplate.New("parameter").Parse(string(assets["client/parameter.gotmpl"])))

	clientResponseTemplate = makeModelTemplate() // template.Must(validatorTempl.Clone())
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("docstring.gotmpl").Parse(string(assets["docstring.gotmpl"])))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("validationDocString").Parse(string(assets["validation/structfield.gotmpl"])))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("schType").Parse(string(assets["schematype.gotmpl"])))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("body").Parse(string(assets["schemabody.gotmpl"])))
	clientResponseTemplate = template.Must(clientResponseTemplate.New("response").Parse(string(assets["client/response.gotmpl"])))

	clientTemplate = template.Must(template.New("docstring.gotmpl").Funcs(FuncMap).Parse(string(assets["docstring.gotmpl"])))
	clientTemplate = template.Must(clientTemplate.New("validationDocString").Parse(string(assets["validation/structfield.gotmpl"])))
	clientTemplate = template.Must(clientTemplate.New("schType").Parse(string(assets["schematype.gotmpl"])))
	clientTemplate = template.Must(clientTemplate.New("body").Parse(string(assets["schemabody.gotmpl"])))
	clientTemplate = template.Must(clientTemplate.New("client").Parse(string(assets["client/client.gotmpl"])))

	clientFacadeTemplate = template.Must(template.New("docstring.gotmpl").Funcs(FuncMap).Parse(string(assets["docstring.gotmpl"])))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("validationDocString").Parse(string(assets["validation/structfield.gotmpl"])))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("schType").Parse(string(assets["schematype.gotmpl"])))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("body").Parse(string(assets["schemabody.gotmpl"])))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("facade").Parse(string(assets["client/facade.gotmpl"])))
}

func init() {

	compileTemplates()

}

func makeModelTemplate() *template.Template {
	templ := template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assets["docstring.gotmpl"])))
	templ = template.Must(templ.New("primitivevalidator").Parse(string(assets["validation/primitive.gotmpl"])))
	templ = template.Must(templ.New("customformatvalidator").Parse(string(assets["validation/customformat.gotmpl"])))
	templ = template.Must(templ.New("validationDocString").Parse(string(assets["validation/structfield.gotmpl"])))
	templ = template.Must(templ.New("schematype").Parse(string(assets["schematype.gotmpl"])))
	templ = template.Must(templ.New("body").Parse(string(assets["schemabody.gotmpl"])))
	templ = template.Must(templ.New("schema").Parse(string(assets["schema.gotmpl"])))
	templ = template.Must(templ.New("schemavalidations").Parse(string(assets["schemavalidator.gotmpl"])))
	templ = template.Must(templ.New("header").Parse(string(assets["header.gotmpl"])))
	templ = template.Must(templ.New("fields").Parse(string(assets["structfield.gotmpl"])))
	templ = template.Must(templ.New("tupleSerializer").Parse(string(assets["tupleserializer.gotmpl"])))
	templ = template.Must(templ.New("additionalpropertiesserializer.gotmpl").Parse(string(assets["additionalpropertiesserializer.gotmpl"])))
	templ = template.Must(templ.New("model").Parse(string(assets["model.gotmpl"])))
	return templ
}

func asJSON(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

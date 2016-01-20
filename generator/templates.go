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

	"PrimitiveValidation":       MustAsset("templates/validation/primitive.gotmpl"),
	"CustomFormatValidation":    MustAsset("templates/validation/customformat.gotmpl"),
	"DocString":                 MustAsset("templates/docstring.gotmpl"),
	"StuctFieldValidation":      MustAsset("templates/validation/structfield.gotmpl"),
	"ModelValidator":            MustAsset("templates/modelvalidator.gotmpl"),
	"SchemaStructField":         MustAsset("templates/structfield.gotmpl"),
	"SchemaTupleSerializer":     MustAsset("templates/tupleserializer.gotmpl"),
	"AdditionalPropsSerializer": MustAsset("templates/additionalpropertiesserializer.gotmpl"),
	"SchemaType":                MustAsset("templates/schematype.gotmpl"),
	"SchemaBody":                MustAsset("templates/schemabody.gotmpl"),
	"Schema":                    MustAsset("templates/schema.gotmpl"),
	"SchemaValidator":           MustAsset("templates/schemavalidator.gotmpl"),
	"SchemaStruct":              MustAsset("templates/model.gotmpl"),
	"Header":                    MustAsset("templates/header.gotmpl"),
	"EmbeddedSpec":              MustAsset("templates/swagger_json_embed.gotmpl"),

	"ServerParameter":    MustAsset("templates/server/parameter.gotmpl"),
	"ServerResponses":    MustAsset("templates/server/responses.gotmpl"),
	"ServerOperation":    MustAsset("templates/server/operation.gotmpl"),
	"ServerBuilder":      MustAsset("templates/server/builder.gotmpl"),
	"ServerConfigureAPI": MustAsset("templates/server/configureapi.gotmpl"),
	"ServerMain":         MustAsset("templates/server/main.gotmpl"),
	"ServerMainDoc":      MustAsset("templates/server/doc.gotmpl"),

	"ClientParameter": MustAsset("templates/client/parameter.gotmpl"),
	"ClientResponse":  MustAsset("templates/client/response.gotmpl"),
	"ClientClient":    MustAsset("templates/client/client.gotmpl"),
	"ClientFacade":    MustAsset("templates/client/facade.gotmpl"),
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

func loadCustomTemplates(templatePath string) {

}

func init() {

	// partial templates
	validatorTempl := template.Must(template.New("primitivevalidator").Funcs(FuncMap).Parse(string(assets["PrimitiveValidation"])))
	validatorTempl = template.Must(validatorTempl.New("customformatvalidator").Parse(string(assets["CustomFormatValidation"])))

	modelTemplate = makeModelTemplate()
	// common templates
	bv, _ := Asset("templates/modelvalidator.gotmpl") // about to be gobbled up by the model template
	modelValidatorTemplate = template.Must(validatorTempl.Clone())
	modelValidatorTemplate = template.Must(modelValidatorTemplate.New("modelvalidator").Parse(string(bv)))

	// server templates
	parameterTemplate = makeModelTemplate()
	//parameterTemplate = template.Must(parameterTemplate.New("docstring").Parse(string(assets["DocString"])))
	//parameterTemplate = template.Must(parameterTemplate.New("validationDocString").Parse(string(assets["StuctFieldValidation"])))
	//parameterTemplate = template.Must(parameterTemplate.New("schType").Parse(string(assets["SchemaType"])))
	//parameterTemplate = template.Must(parameterTemplate.New("body").Parse(string(assets["SchemaBody"])))
	parameterTemplate = template.Must(parameterTemplate.New("parameter").Parse(string(assets["ServerParameter"])))

	responsesTemplate = makeModelTemplate()
	responsesTemplate = template.Must(responsesTemplate.New("responses").Parse(string(assets["ServerResponses"])))

	operationTemplate = makeModelTemplate()
	operationTemplate = template.Must(operationTemplate.New("operation").Parse(string(assets["ServerOperation"])))
	builderTemplate = template.Must(template.New("builder").Funcs(FuncMap).Parse(string(assets["ServerBuilder"])))
	configureAPITemplate = template.Must(template.New("configureapi").Funcs(FuncMap).Parse(string(assets["ServerConfigureAPI"])))
	mainTemplate = template.Must(template.New("main").Funcs(FuncMap).Parse(string(assets["ServerMain"])))
	mainDocTemplate = template.Must(template.New("meta").Funcs(FuncMap).Parse(string(assets["ServerMainDoc"])))

	embeddedSpecTemplate = template.Must(template.New("embedded_spec").Funcs(FuncMap).Parse(string(assets["EmbeddedSpec"])))

	// Client templates
	clientParamTemplate = makeModelTemplate()
	//clientParamTemplate = template.Must(clientParamTemplate.New("docstring").Parse(string(assets["DocString"])))
	//clientParamTemplate = template.Must(clientParamTemplate.New("validationDocString").Parse(string(assets["StuctFieldValidation"])))
	//clientParamTemplate = template.Must(clientParamTemplate.New("schType").Parse(string(assets["SchemaType"])))
	//clientParamTemplate = template.Must(clientParamTemplate.New("body").Parse(string(assets["SchemaBody"])))
	clientParamTemplate = template.Must(clientParamTemplate.New("parameter").Parse(string(assets["ClientParameter"])))

	clientResponseTemplate = makeModelTemplate() // template.Must(validatorTempl.Clone())
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("docstring").Parse(string(assets["DocString"])))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("validationDocString").Parse(string(assets["StuctFieldValidation"])))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("schType").Parse(string(assets["SchemaType"])))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("body").Parse(string(assets["SchemaBody"])))
	clientResponseTemplate = template.Must(clientResponseTemplate.New("response").Parse(string(assets["ClientResponse"])))

	clientTemplate = template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assets["DocString"])))
	clientTemplate = template.Must(clientTemplate.New("validationDocString").Parse(string(assets["StuctFieldValidation"])))
	clientTemplate = template.Must(clientTemplate.New("schType").Parse(string(assets["SchemaType"])))
	clientTemplate = template.Must(clientTemplate.New("body").Parse(string(assets["SchemaBody"])))
	clientTemplate = template.Must(clientTemplate.New("client").Parse(string(assets["ClientClient"])))

	clientFacadeTemplate = template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assets["DocString"])))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("validationDocString").Parse(string(assets["StuctFieldValidation"])))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("schType").Parse(string(assets["SchemaType"])))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("body").Parse(string(assets["SchemaBody"])))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("facade").Parse(string(assets["ClientFacade"])))

}

func makeModelTemplate() *template.Template {
	templ := template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assets["DocString"])))
	templ = template.Must(templ.New("primitivevalidator").Parse(string(assets["PrimitiveValidation"])))
	templ = template.Must(templ.New("customformatvalidator").Parse(string(assets["CustomFormatValidation"])))
	templ = template.Must(templ.New("validationDocString").Parse(string(assets["StuctFieldValidation"])))
	templ = template.Must(templ.New("schemaType").Parse(string(assets["SchemaType"])))
	templ = template.Must(templ.New("body").Parse(string(assets["SchemaBody"])))
	templ = template.Must(templ.New("schema").Parse(string(assets["Schema"])))
	templ = template.Must(templ.New("schemavalidations").Parse(string(assets["SchemaValidator"])))
	templ = template.Must(templ.New("header").Parse(string(assets["Header"])))
	templ = template.Must(templ.New("fields").Parse(string(assets["SchemaStructField"])))
	templ = template.Must(templ.New("tupleSerializer").Parse(string(assets["SchemaTupleSerializer"])))
	templ = template.Must(templ.New("additionalPropsSerializer").Parse(string(assets["AdditionalPropsSerializer"])))
	templ = template.Must(templ.New("model").Parse(string(assets["SchemaStruct"])))
	return templ
}

func asJSON(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

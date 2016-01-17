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

var (
	assetPrimitiveValidation       = MustAsset("templates/validation/primitive.gotmpl")
	assetCustomFormatValidation    = MustAsset("templates/validation/customformat.gotmpl")
	assetDocString                 = MustAsset("templates/docstring.gotmpl")
	assetStuctFieldValidation      = MustAsset("templates/validation/structfield.gotmpl")
	assetModelValidator            = MustAsset("templates/modelvalidator.gotmpl")
	assetSchemaStructField         = MustAsset("templates/structfield.gotmpl")
	assetSchemaTupleSerializer     = MustAsset("templates/tupleserializer.gotmpl")
	assetAdditionalPropsSerializer = MustAsset("templates/additionalpropertiesserializer.gotmpl")
	assetSchemaType                = MustAsset("templates/schematype.gotmpl")
	assetSchemaBody                = MustAsset("templates/schemabody.gotmpl")
	assetSchema                    = MustAsset("templates/schema.gotmpl")
	assetSchemaValidator           = MustAsset("templates/schemavalidator.gotmpl")
	assetSchemaStruct              = MustAsset("templates/model.gotmpl")
	assetHeader                    = MustAsset("templates/header.gotmpl")
	assetEmbeddedSpec              = MustAsset("templates/swagger_json_embed.gotmpl")

	assetServerParameter    = MustAsset("templates/server/parameter.gotmpl")
	assetServerResponses    = MustAsset("templates/server/responses.gotmpl")
	assetServerOperation    = MustAsset("templates/server/operation.gotmpl")
	assetServerBuilder      = MustAsset("templates/server/builder.gotmpl")
	assetServerConfigureAPI = MustAsset("templates/server/configureapi.gotmpl")
	assetServerMain         = MustAsset("templates/server/main.gotmpl")
	assetServerMainDoc      = MustAsset("templates/server/doc.gotmpl")

	assetClientParameter = MustAsset("templates/client/parameter.gotmpl")
	assetClientResponse  = MustAsset("templates/client/response.gotmpl")
	assetClientClient    = MustAsset("templates/client/client.gotmpl")
	assetClientFacade    = MustAsset("templates/client/facade.gotmpl")
)

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

func init() {

	// partial templates
	validatorTempl := template.Must(template.New("primitivevalidator").Funcs(FuncMap).Parse(string(assetPrimitiveValidation)))
	validatorTempl = template.Must(validatorTempl.New("customformatvalidator").Parse(string(assetCustomFormatValidation)))

	modelTemplate = makeModelTemplate()
	// common templates
	bv, _ := Asset("templates/modelvalidator.gotmpl") // about to be gobbled up by the model template
	modelValidatorTemplate = template.Must(validatorTempl.Clone())
	modelValidatorTemplate = template.Must(modelValidatorTemplate.New("modelvalidator").Parse(string(bv)))

	// server templates
	parameterTemplate = makeModelTemplate()
	//parameterTemplate = template.Must(parameterTemplate.New("docstring").Parse(string(assetDocString)))
	//parameterTemplate = template.Must(parameterTemplate.New("validationDocString").Parse(string(assetStuctFieldValidation)))
	//parameterTemplate = template.Must(parameterTemplate.New("schType").Parse(string(assetSchemaType)))
	//parameterTemplate = template.Must(parameterTemplate.New("body").Parse(string(assetSchemaBody)))
	parameterTemplate = template.Must(parameterTemplate.New("parameter").Parse(string(assetServerParameter)))

	responsesTemplate = makeModelTemplate()
	responsesTemplate = template.Must(responsesTemplate.New("responses").Parse(string(assetServerResponses)))

	operationTemplate = makeModelTemplate()
	operationTemplate = template.Must(operationTemplate.New("operation").Parse(string(assetServerOperation)))
	builderTemplate = template.Must(template.New("builder").Funcs(FuncMap).Parse(string(assetServerBuilder)))
	configureAPITemplate = template.Must(template.New("configureapi").Funcs(FuncMap).Parse(string(assetServerConfigureAPI)))
	mainTemplate = template.Must(template.New("main").Funcs(FuncMap).Parse(string(assetServerMain)))
	mainDocTemplate = template.Must(template.New("meta").Funcs(FuncMap).Parse(string(assetServerMainDoc)))

	embeddedSpecTemplate = template.Must(template.New("embedded_spec").Funcs(FuncMap).Parse(string(assetEmbeddedSpec)))

	// Client templates
	clientParamTemplate = makeModelTemplate()
	//clientParamTemplate = template.Must(clientParamTemplate.New("docstring").Parse(string(assetDocString)))
	//clientParamTemplate = template.Must(clientParamTemplate.New("validationDocString").Parse(string(assetStuctFieldValidation)))
	//clientParamTemplate = template.Must(clientParamTemplate.New("schType").Parse(string(assetSchemaType)))
	//clientParamTemplate = template.Must(clientParamTemplate.New("body").Parse(string(assetSchemaBody)))
	clientParamTemplate = template.Must(clientParamTemplate.New("parameter").Parse(string(assetClientParameter)))

	clientResponseTemplate = makeModelTemplate() // template.Must(validatorTempl.Clone())
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("docstring").Parse(string(assetDocString)))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("validationDocString").Parse(string(assetStuctFieldValidation)))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("schType").Parse(string(assetSchemaType)))
	// clientResponseTemplate = template.Must(clientResponseTemplate.New("body").Parse(string(assetSchemaBody)))
	clientResponseTemplate = template.Must(clientResponseTemplate.New("response").Parse(string(assetClientResponse)))

	clientTemplate = template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assetDocString)))
	clientTemplate = template.Must(clientTemplate.New("validationDocString").Parse(string(assetStuctFieldValidation)))
	clientTemplate = template.Must(clientTemplate.New("schType").Parse(string(assetSchemaType)))
	clientTemplate = template.Must(clientTemplate.New("body").Parse(string(assetSchemaBody)))
	clientTemplate = template.Must(clientTemplate.New("client").Parse(string(assetClientClient)))

	clientFacadeTemplate = template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assetDocString)))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("validationDocString").Parse(string(assetStuctFieldValidation)))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("schType").Parse(string(assetSchemaType)))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("body").Parse(string(assetSchemaBody)))
	clientFacadeTemplate = template.Must(clientFacadeTemplate.New("facade").Parse(string(assetClientFacade)))

}

func makeModelTemplate() *template.Template {
	templ := template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assetDocString)))
	templ = template.Must(templ.New("primitivevalidator").Parse(string(assetPrimitiveValidation)))
	templ = template.Must(templ.New("customformatvalidator").Parse(string(assetCustomFormatValidation)))
	templ = template.Must(templ.New("validationDocString").Parse(string(assetStuctFieldValidation)))
	templ = template.Must(templ.New("schemaType").Parse(string(assetSchemaType)))
	templ = template.Must(templ.New("body").Parse(string(assetSchemaBody)))
	templ = template.Must(templ.New("schema").Parse(string(assetSchema)))
	templ = template.Must(templ.New("schemavalidations").Parse(string(assetSchemaValidator)))
	templ = template.Must(templ.New("header").Parse(string(assetHeader)))
	templ = template.Must(templ.New("fields").Parse(string(assetSchemaStructField)))
	templ = template.Must(templ.New("tupleSerializer").Parse(string(assetSchemaTupleSerializer)))
	templ = template.Must(templ.New("additionalPropsSerializer").Parse(string(assetAdditionalPropsSerializer)))
	templ = template.Must(templ.New("model").Parse(string(assetSchemaStruct)))
	return templ
}

func asJSON(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

package generator

import (
	"encoding/json"
	"text/template"

	"github.com/go-swagger/go-swagger/swag"
)

//go:generate go-bindata -pkg=generator -ignore=.*\.sw? ./templates/...

// fwiw, don't get attached to this, still requires a better abstraction

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

	assetServerParameter    = MustAsset("templates/server/parameter.gotmpl")
	assetServerOperation    = MustAsset("templates/server/operation.gotmpl")
	assetServerBuilder      = MustAsset("templates/server/builder.gotmpl")
	assetServerConfigureAPI = MustAsset("templates/server/configureapi.gotmpl")
	assetServerMain         = MustAsset("templates/server/main.gotmpl")

	assetClientParameter = MustAsset("templates/client/parameter.gotmpl")
	assetClientResponse  = MustAsset("templates/client/response.gotmpl")
	assetClientClient    = MustAsset("templates/client/client.gotmpl")
	assetClientFacade    = MustAsset("templates/client/facade.gotmpl")
)

// FuncMap is a map with default functions for use n the templates.
// These are available in every template
var FuncMap template.FuncMap = map[string]interface{}{
	"pascalize": swag.ToGoName,
	"camelize":  swag.ToJSONName,
	"humanize":  swag.ToHumanNameLower,
	"snakize":   swag.ToFileName,
	"dasherize": swag.ToCommandName,
	"json":      asJSON,
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
	parameterTemplate = template.Must(validatorTempl.Clone())
	parameterTemplate = template.Must(parameterTemplate.New("docstring").Parse(string(assetDocString)))
	parameterTemplate = template.Must(parameterTemplate.New("validationDocString").Parse(string(assetStuctFieldValidation)))
	parameterTemplate = template.Must(parameterTemplate.New("schType").Parse(string(assetSchemaType)))
	parameterTemplate = template.Must(parameterTemplate.New("body").Parse(string(assetSchemaBody)))
	parameterTemplate = template.Must(parameterTemplate.New("parameter").Parse(string(assetServerParameter)))

	operationTemplate = template.Must(template.New("operation").Funcs(FuncMap).Parse(string(assetServerOperation)))
	builderTemplate = template.Must(template.New("builder").Funcs(FuncMap).Parse(string(assetServerBuilder)))
	configureAPITemplate = template.Must(template.New("configureapi").Funcs(FuncMap).Parse(string(assetServerConfigureAPI)))
	mainTemplate = template.Must(template.New("main").Funcs(FuncMap).Parse(string(assetServerMain)))

	// Client templates
	clientParamTemplate = template.Must(validatorTempl.Clone())
	clientParamTemplate = template.Must(clientParamTemplate.New("docstring").Parse(string(assetDocString)))
	clientParamTemplate = template.Must(clientParamTemplate.New("validationDocString").Parse(string(assetStuctFieldValidation)))
	clientParamTemplate = template.Must(clientParamTemplate.New("schType").Parse(string(assetSchemaType)))
	clientParamTemplate = template.Must(clientParamTemplate.New("body").Parse(string(assetSchemaBody)))
	clientParamTemplate = template.Must(clientParamTemplate.New("parameter").Parse(string(assetClientParameter)))
	clientResponseTemplate = template.Must(validatorTempl.Clone())
	clientResponseTemplate = template.Must(clientResponseTemplate.New("response").Parse(string(assetClientResponse)))
	clientTemplate = template.Must(template.New("client").Funcs(FuncMap).Parse(string(assetClientClient)))
	clientFacadeTemplate = template.Must(template.New("facade").Funcs(FuncMap).Parse(string(assetClientFacade)))
}

func makeModelTemplate() *template.Template {
	templ := template.Must(template.New("docstring").Funcs(FuncMap).Parse(string(assetDocString)))
	templ = template.Must(templ.New("primitivevalidator").Parse(string(assetPrimitiveValidation)))
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

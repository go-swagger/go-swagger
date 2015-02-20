package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/util"
)

var (
	builderTemplate *template.Template
)

func init() {
	bv, _ := Asset("templates/server/builder.gotmpl")
	builderTemplate = template.Must(template.New("builder").Parse(string(bv)))
}

// GenerateSupport generates the supporting files for an API
func GenerateSupport(name string, modelNames, operationIDs []string, opts GenOpts) error {
	// Load the spec
	_, specDoc, err := loadSpec(opts.Spec)
	if err != nil {
		return err
	}

	models := make(map[string]spec.Schema)
	if len(modelNames) == 0 {
		for k, v := range specDoc.Spec().Definitions {
			models[k] = v
		}
	} else {
		for k, v := range specDoc.Spec().Definitions {
			for _, nm := range modelNames {
				if k == nm {
					models[k] = v
				}
			}
		}
	}

	operations := make(map[string]spec.Operation)
	if len(modelNames) == 0 {
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

	generator := appGenerator{
		Name:       name,
		SpecDoc:    specDoc,
		Models:     models,
		Operations: operations,
		Target:     opts.Target,
		// Package:       filepath.Base(opts.Target),
		DumpData:      opts.DumpData,
		Package:       opts.APIPackage,
		APIPackage:    opts.APIPackage,
		ModelsPackage: opts.ModelPackage,
		Principal:     opts.Principal,
	}

	return generator.Generate()
}

type appGenerator struct {
	Name          string
	SpecDoc       *spec.Document
	Package       string
	APIPackage    string
	ModelsPackage string
	Principal     string
	Models        map[string]spec.Schema
	Operations    map[string]spec.Operation
	Target        string
	DumpData      bool
}

func (a *appGenerator) Generate() error {
	app := makeCodegenApp(a.Name, a.Package, a.Target, a.ModelsPackage, a.APIPackage, a.Principal, a.SpecDoc, a.Models, a.Operations)

	if a.DumpData {
		bb, _ := json.MarshalIndent(util.ToDynamicJSON(app), "", " ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}
	return a.generateAPIBuilder(&app)
}

func (a *appGenerator) generateAPIBuilder(app *genApp) error {
	buf := bytes.NewBuffer(nil)
	if err := builderTemplate.Execute(buf, app); err != nil {
		return err
	}
	log.Println("rendered builder template:", app.AppName)
	return writeToFile(filepath.Join(a.Target, app.Package), app.AppName+"Builder", buf.Bytes())
}

var mediaTypeNames = map[string]string{
	"application/json":        "json",
	"application/x-yaml":      "yaml",
	"application/x-protobuf":  "protobuf",
	"application/x-capnproto": "capnproto",
	"application/x-thrift":    "thrift",
	"application/xml":         "xml",
	"text/xml":                "xml",
	"text/x-markdown":         "markdown",
	"text/html":               "html",
	"text/csv":                "csv",
	"text/tsv":                "tsv",
	"text/javascript":         "js",
	"text/css":                "css",
}

func makeCodegenApp(name, pkg, target, modelPackage, apiPackage, principal string, specDoc *spec.Document, models map[string]spec.Schema, operations map[string]spec.Operation) genApp {
	sw := specDoc.Spec()
	receiver := strings.ToLower(name[:1])
	appName := util.ToGoName(name)

	var consumes []genSerializer
	for _, cons := range specDoc.RequiredConsumes() {
		cn := mediaTypeNames[cons]
		consumes = append(consumes, genSerializer{
			AppName:        appName,
			ReceiverName:   receiver,
			ClassName:      util.ToGoName(cn),
			HumanClassName: util.ToHumanNameLower(cn),
			Name:           util.ToJSONName(cn),
			MediaType:      cons,
		})
	}

	var produces []genSerializer
	for _, prod := range specDoc.RequiredProduces() {
		pn := mediaTypeNames[prod]
		produces = append(produces, genSerializer{
			AppName:        appName,
			ReceiverName:   receiver,
			ClassName:      util.ToGoName(pn),
			HumanClassName: util.ToHumanNameLower(pn),
			Name:           util.ToJSONName(pn),
			MediaType:      prod,
		})
	}

	var security []genSecurityScheme
	for _, scheme := range specDoc.RequiredSchemes() {
		if req, ok := specDoc.Spec().SecurityDefinitions[scheme]; ok {
			security = append(security, genSecurityScheme{
				AppName:        appName,
				ReceiverName:   receiver,
				ClassName:      util.ToGoName(req.Name),
				HumanClassName: util.ToHumanNameLower(req.Name),
				Name:           util.ToJSONName(req.Name),
				IsBasicAuth:    strings.ToLower(req.Type) == "basic",
				IsAPIKeyAuth:   strings.ToLower(req.Type) == "apiKey",
			})
		}
	}

	var genMods []genModel
	for mn, m := range models {
		mod := *makeCodegenModel(
			mn,
			modelPackage,
			m,
			specDoc,
		)
		mod.ReceiverName = receiver
		genMods = append(genMods, mod)
	}

	var genOps []genOperation
	for on, o := range operations {
		authed := len(specDoc.SecurityRequirementsFor(&o)) > 0
		ap := apiPackage
		if apiPackage == pkg {
			ap = ""
		}
		op := makeCodegenOperation(
			on,
			ap,
			modelPackage,
			principal,
			o,
			authed,
		)
		op.ReceiverName = receiver
		genOps = append(genOps, op)
	}

	return genApp{
		Package:             pkg,
		ReceiverName:        receiver,
		AppName:             util.ToGoName(name),
		HumanAppName:        util.ToHumanNameLower(name),
		Name:                util.ToJSONName(name),
		ExternalDocs:        sw.ExternalDocs,
		Info:                sw.Info,
		Consumes:            consumes,
		Produces:            produces,
		SecurityDefinitions: security,
		Models:              genMods,
		Operations:          genOps,
	}
}

type genApp struct {
	Package             string
	ReceiverName        string
	AppName             string
	HumanAppName        string
	Name                string
	Info                *spec.Info
	ExternalDocs        *spec.ExternalDocumentation
	Imports             []string
	Consumes            []genSerializer
	Produces            []genSerializer
	SecurityDefinitions []genSecurityScheme
	Models              []genModel
	Operations          []genOperation
}

type genSerializer struct {
	ReceiverName   string
	AppName        string
	ClassName      string
	HumanClassName string
	Name           string
	MediaType      string
}

type genSecurityScheme struct {
	AppName        string
	ClassName      string
	HumanClassName string
	Name           string
	ReceiverName   string
	IsBasicAuth    bool
	IsAPIKeyAuth   bool
}

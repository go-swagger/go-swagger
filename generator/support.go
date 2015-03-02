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
	builderTemplate      *template.Template
	mainTemplate         *template.Template
	configureAPITemplate *template.Template
)

func init() {
	bv, _ := Asset("templates/server/builder.gotmpl")
	builderTemplate = template.Must(template.New("builder").Parse(string(bv)))

	bm, _ := Asset("templates/server/main.gotmpl")
	mainTemplate = template.Must(template.New("main").Parse(string(bm)))

	bc, _ := Asset("templates/server/configureapi.gotmpl")
	configureAPITemplate = template.Must(template.New("configureapi").Parse(string(bc)))
}

// GenerateSupport generates the supporting files for an API
func GenerateSupport(name string, modelNames, operationIDs []string, includeUI bool, opts GenOpts) error {
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

	if name == "" {
		if specDoc.Spec().Info != nil && specDoc.Spec().Info.Title != "" {
			name = util.ToGoName(specDoc.Spec().Info.Title)
		} else {
			name = "swagger"
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
		ServerPackage: opts.ServerPackage,
		ClientPackage: opts.ClientPackage,
		Principal:     opts.Principal,
		IncludeUI:     includeUI,
	}

	return generator.Generate()
}

type appGenerator struct {
	Name          string
	SpecDoc       *spec.Document
	Package       string
	APIPackage    string
	ModelsPackage string
	ServerPackage string
	ClientPackage string
	Principal     string
	Models        map[string]spec.Schema
	Operations    map[string]spec.Operation
	Target        string
	DumpData      bool
	IncludeUI     bool
}

func baseImport(tgt string) string {
	p, err := filepath.Abs(tgt)
	if err != nil {
		log.Fatalln(err)
	}

	var pth string
	for _, gp := range filepath.SplitList(os.Getenv("GOPATH")) {
		pp := filepath.Join(gp, "src")
		if strings.HasPrefix(p, pp) {
			pth = strings.TrimPrefix(p, pp+"/")
			break
		}
	}

	if pth == "" {
		log.Fatalln("target must reside inside a location in the gopath")
	}
	return pth
}

func (a *appGenerator) Generate() error {
	app := a.makeCodegenApp()

	if a.DumpData {
		bb, _ := json.MarshalIndent(util.ToDynamicJSON(app), "", "  ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}

	if err := a.generateAPIBuilder(&app); err != nil {
		return err
	}
	app.DefaultImports = append(app.DefaultImports, filepath.Join(baseImport(a.Target), a.ServerPackage, a.APIPackage))

	if err := a.generateConfigureAPI(&app); err != nil {
		return err
	}

	if err := a.generateMain(&app); err != nil {
		return err
	}

	return nil
}

func (a *appGenerator) generateConfigureAPI(app *genApp) error {
	pth := filepath.Join(a.Target, "cmd", util.ToCommandName(app.AppName+"Server"))
	nm := "Configure" + app.AppName
	if fileExists(pth, nm) {
		log.Println("skipped (already exists) configure api template:", app.Package+".Configure"+app.AppName)
		return nil
	}

	buf := bytes.NewBuffer(nil)
	if err := configureAPITemplate.Execute(buf, app); err != nil {
		return err
	}
	log.Println("rendered configure api template:", app.Package+".Configure"+app.AppName)
	return writeToFileIfNotExist(pth, nm, buf.Bytes())
}

func (a *appGenerator) generateMain(app *genApp) error {
	buf := bytes.NewBuffer(nil)
	if err := mainTemplate.Execute(buf, app); err != nil {
		return err
	}
	log.Println("rendered main template:", "server."+app.AppName)
	return writeToFile(filepath.Join(a.Target, "cmd", util.ToCommandName(app.AppName+"Server")), "main", buf.Bytes())
}

func (a *appGenerator) generateAPIBuilder(app *genApp) error {
	buf := bytes.NewBuffer(nil)
	if err := builderTemplate.Execute(buf, app); err != nil {
		return err
	}
	log.Println("rendered builder template:", app.Package+"."+app.AppName)
	return writeToFile(filepath.Join(a.Target, a.ServerPackage, app.Package), app.AppName+"Api", buf.Bytes())
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

var knownProducers = map[string]string{
	"json": "swagger.JSONProducer",
	"yaml": "swagger.YAMLProducer",
}

var knownConsumers = map[string]string{
	"json": "swagger.JSONConsumer",
	"yaml": "swagger.YAMLConsumer",
}

func getSerializer(sers []genSerGroup, ext string) (*genSerGroup, bool) {
	for i := range sers {
		s := &sers[i]
		if s.Name == ext {
			return s, true
		}
	}
	return nil, false
}

// func makeCodegenApp(operations map[string]spec.Operation, includeUI bool) genApp {
func (a *appGenerator) makeCodegenApp() genApp {
	sw := a.SpecDoc.Spec()
	// app := makeCodegenApp(a.Operations, a.IncludeUI)
	receiver := strings.ToLower(a.Name[:1])
	appName := util.ToGoName(a.Name)
	var defaultImports []string

	jsonb, _ := json.MarshalIndent(a.SpecDoc.Spec(), "", "  ")

	consumesJSON := false
	var consumes []genSerGroup
	for _, cons := range a.SpecDoc.RequiredConsumes() {
		cn, ok := mediaTypeNames[cons]
		if !ok {
			continue
		}
		nm := util.ToJSONName(cn)
		if nm == "json" {
			consumesJSON = true
		}

		if ser, ok := getSerializer(consumes, cn); ok {
			ser.AllSerializers = append(ser.AllSerializers, genSerializer{
				AppName:        ser.AppName,
				ReceiverName:   ser.ReceiverName,
				ClassName:      ser.ClassName,
				HumanClassName: ser.HumanClassName,
				Name:           ser.Name,
				MediaType:      cons,
				Implementation: knownConsumers[nm],
			})
			continue
		}

		ser := genSerializer{
			AppName:        appName,
			ReceiverName:   receiver,
			ClassName:      util.ToGoName(cn),
			HumanClassName: util.ToHumanNameLower(cn),
			Name:           nm,
			MediaType:      cons,
			Implementation: knownConsumers[nm],
		}

		consumes = append(consumes, genSerGroup{
			AppName:        ser.AppName,
			ReceiverName:   ser.ReceiverName,
			ClassName:      ser.ClassName,
			HumanClassName: ser.HumanClassName,
			Name:           ser.Name,
			MediaType:      cons,
			AllSerializers: []genSerializer{ser},
			Implementation: ser.Implementation,
		})
	}

	producesJSON := false
	var produces []genSerGroup
	for _, prod := range a.SpecDoc.RequiredProduces() {
		pn, ok := mediaTypeNames[prod]
		if !ok {
			continue
		}
		nm := util.ToJSONName(pn)
		if nm == "json" {
			producesJSON = true
		}

		if ser, ok := getSerializer(produces, pn); ok {
			ser.AllSerializers = append(ser.AllSerializers, genSerializer{
				AppName:        ser.AppName,
				ReceiverName:   ser.ReceiverName,
				ClassName:      ser.ClassName,
				HumanClassName: ser.HumanClassName,
				Name:           ser.Name,
				MediaType:      prod,
				Implementation: knownProducers[nm],
			})
			continue
		}
		ser := genSerializer{
			AppName:        appName,
			ReceiverName:   receiver,
			ClassName:      util.ToGoName(pn),
			HumanClassName: util.ToHumanNameLower(pn),
			Name:           nm,
			MediaType:      prod,
			Implementation: knownProducers[nm],
		}
		produces = append(produces, genSerGroup{
			AppName:        ser.AppName,
			ReceiverName:   ser.ReceiverName,
			ClassName:      ser.ClassName,
			HumanClassName: ser.HumanClassName,
			Name:           ser.Name,
			MediaType:      prod,
			Implementation: ser.Implementation,
			AllSerializers: []genSerializer{ser},
		})
	}

	var security []genSecurityScheme
	for _, scheme := range a.SpecDoc.RequiredSchemes() {
		if req, ok := a.SpecDoc.Spec().SecurityDefinitions[scheme]; ok {
			if req.Type == "basic" || req.Type == "apiKey" {
				security = append(security, genSecurityScheme{
					AppName:        appName,
					ReceiverName:   receiver,
					ClassName:      util.ToGoName(req.Name),
					HumanClassName: util.ToHumanNameLower(req.Name),
					Name:           util.ToJSONName(req.Name),
					IsBasicAuth:    strings.ToLower(req.Type) == "basic",
					IsAPIKeyAuth:   strings.ToLower(req.Type) == "apikey",
					Principal:      a.Principal,
					Source:         req.In,
				})
			}
		}
	}

	var genMods []genModel
	defaultImports = append(defaultImports, filepath.Join(baseImport(a.Target), a.ModelsPackage))
	for mn, m := range a.Models {
		mod := *makeCodegenModel(
			mn,
			a.ModelsPackage,
			m,
			a.SpecDoc,
		)
		mod.ReceiverName = receiver
		genMods = append(genMods, mod)
	}

	var genOps []genOperation
	tns := make(map[string]struct{})
	for on, o := range a.Operations {
		authed := len(a.SpecDoc.SecurityRequirementsFor(&o)) > 0
		ap := a.APIPackage
		if a.APIPackage == a.Package {
			ap = ""
		}
		if len(o.Tags) > 0 {
			for _, tag := range o.Tags {
				tns[tag] = struct{}{}
				op := makeCodegenOperation(on, tag, a.ModelsPackage, a.Principal, a.Target, o, authed)
				op.ReceiverName = receiver
				genOps = append(genOps, op)
			}
		} else {
			op := makeCodegenOperation(on, ap, a.ModelsPackage, a.Principal, a.Target, o, authed)
			op.ReceiverName = receiver
			genOps = append(genOps, op)
		}
	}
	for k := range tns {
		defaultImports = append(defaultImports, filepath.Join(baseImport(a.Target), a.ServerPackage, a.APIPackage, k))
	}

	defaultConsumes := "application/json"
	rc := a.SpecDoc.RequiredConsumes()
	if !consumesJSON && len(rc) > 0 {
		defaultConsumes = rc[0]
	}

	defaultProduces := "application/json"
	rp := a.SpecDoc.RequiredProduces()
	if !producesJSON && len(rp) > 0 {
		defaultProduces = rp[0]
	}

	return genApp{
		Package:             a.Package,
		ReceiverName:        receiver,
		AppName:             util.ToGoName(a.Name),
		HumanAppName:        util.ToHumanNameLower(a.Name),
		Name:                util.ToJSONName(a.Name),
		ExternalDocs:        sw.ExternalDocs,
		Info:                sw.Info,
		Consumes:            consumes,
		Produces:            produces,
		DefaultConsumes:     defaultConsumes,
		DefaultProduces:     defaultProduces,
		DefaultImports:      defaultImports,
		SecurityDefinitions: security,
		Models:              genMods,
		Operations:          genOps,
		IncludeUI:           a.IncludeUI,
		Principal:           a.Principal,
		SwaggerJSON:         fmt.Sprintf("%#v", jsonb),
	}
}

type genApp struct {
	Package             string
	ReceiverName        string
	AppName             string
	HumanAppName        string
	Name                string
	Principal           string
	DefaultConsumes     string
	DefaultProduces     string
	Info                *spec.Info
	ExternalDocs        *spec.ExternalDocumentation
	Imports             map[string]string
	DefaultImports      []string
	Consumes            []genSerGroup
	Produces            []genSerGroup
	SecurityDefinitions []genSecurityScheme
	Models              []genModel
	Operations          []genOperation
	IncludeUI           bool
	SwaggerJSON         string
}

type genSerGroup struct {
	ReceiverName   string
	AppName        string
	ClassName      string
	HumanClassName string
	Name           string
	MediaType      string
	Implementation string
	AllSerializers []genSerializer
}

type genSerializer struct {
	ReceiverName   string
	AppName        string
	ClassName      string
	HumanClassName string
	Name           string
	MediaType      string
	Implementation string
}

type genSecurityScheme struct {
	AppName        string
	ClassName      string
	HumanClassName string
	Name           string
	ReceiverName   string
	IsBasicAuth    bool
	IsAPIKeyAuth   bool
	Source         string
	Principal      string
}

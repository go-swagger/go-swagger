package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-swagger/go-swagger/spec"
	"github.com/go-swagger/go-swagger/swag"
)

// GenerateSupport generates the supporting files for an API
func GenerateSupport(name string, modelNames, operationIDs []string, opts GenOpts) error {
	// Load the spec
	_, specDoc, err := loadSpec(opts.Spec)
	if err != nil {
		return err
	}

	models := gatherModels(specDoc, modelNames)
	operations := gatherOperations(specDoc, operationIDs)

	generator := appGenerator{
		Name:       appNameOrDefault(specDoc, name, "swagger"),
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
			pth, err = filepath.Rel(pp, p)
			if err != nil {
				log.Fatalln(err)
			}
			break
		}
	}

	if pth == "" {
		log.Fatalln("target must reside inside a location in the gopath")
	}
	return pth
}

func (a *appGenerator) Generate() error {
	app, err := a.makeCodegenApp()
	if err != nil {
		return err
	}

	if a.DumpData {
		bb, _ := json.MarshalIndent(swag.ToDynamicJSON(app), "", "  ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}

	if err := a.generateAPIBuilder(&app); err != nil {
		return err
	}
	importPath := filepath.ToSlash(filepath.Join(baseImport(a.Target), a.ServerPackage, a.APIPackage))
	app.DefaultImports = append(app.DefaultImports, importPath)

	if err := a.generateConfigureAPI(&app); err != nil {
		return err
	}

	if err := a.generateMain(&app); err != nil {
		return err
	}

	return nil
}

func (a *appGenerator) generateConfigureAPI(app *genApp) error {
	pth := filepath.Join(a.Target, "cmd", swag.ToCommandName(app.AppName+"Server"))
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
	return writeToFile(filepath.Join(a.Target, "cmd", swag.ToCommandName(app.AppName+"Server")), "main", buf.Bytes())
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

func (a *appGenerator) makeCodegenApp() (genApp, error) {
	sw := a.SpecDoc.Spec()
	receiver := "o" //strings.ToLower(a.Name[:1])
	appName := swag.ToGoName(a.Name)
	var defaultImports []string

	jsonb, _ := json.MarshalIndent(a.SpecDoc.Spec(), "", "  ")

	consumesJSON := false
	var consumes []genSerGroup
	for _, cons := range a.SpecDoc.RequiredConsumes() {
		cn, ok := mediaTypeNames[cons]
		if !ok {
			continue
		}
		nm := swag.ToJSONName(cn)
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
			ClassName:      swag.ToGoName(cn),
			HumanClassName: swag.ToHumanNameLower(cn),
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
		nm := swag.ToJSONName(pn)
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
			ClassName:      swag.ToGoName(pn),
			HumanClassName: swag.ToHumanNameLower(pn),
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
					ClassName:      swag.ToGoName(req.Name),
					HumanClassName: swag.ToHumanNameLower(req.Name),
					Name:           swag.ToJSONName(req.Name),
					IsBasicAuth:    strings.ToLower(req.Type) == "basic",
					IsAPIKeyAuth:   strings.ToLower(req.Type) == "apikey",
					Principal:      a.Principal,
					Source:         req.In,
				})
			}
		}
	}

	var genMods []genModel
	importPath := filepath.ToSlash(filepath.Join(baseImport(a.Target), a.ModelsPackage))
	defaultImports = append(defaultImports, importPath)
	for mn, m := range a.Models {
		mod, err := makeCodegenModel(
			mn,
			a.ModelsPackage,
			m,
			a.SpecDoc,
		)
		if err != nil {
			return genApp{}, err
		}
		mod.ReceiverName = receiver
		genMods = append(genMods, *mod)
	}

	var genOps []genOperation
	tns := make(map[string]struct{})
	var bldr codeGenOpBuilder
	bldr.ModelsPackage = a.ModelsPackage
	bldr.Principal = a.Principal
	bldr.Target = a.Target
	bldr.Doc = a.SpecDoc

	for on, o := range a.Operations {
		bldr.Name = on
		bldr.Operation = o
		bldr.Authed = len(a.SpecDoc.SecurityRequirementsFor(&o)) > 0
		ap := a.APIPackage
		if a.APIPackage == a.Package {
			ap = ""
		}
		if len(o.Tags) > 0 {
			for _, tag := range o.Tags {
				tns[tag] = struct{}{}
				bldr.APIPackage = tag
				op, err := makeCodegenOperation(bldr)
				if err != nil {
					return genApp{}, err
				}
				op.ReceiverName = receiver
				genOps = append(genOps, op)
			}
		} else {
			bldr.APIPackage = ap
			op, err := makeCodegenOperation(bldr)
			if err != nil {
				return genApp{}, err
			}
			op.ReceiverName = receiver
			genOps = append(genOps, op)
		}
	}
	for k := range tns {
		importPath := filepath.ToSlash(filepath.Join(baseImport(a.Target), a.ServerPackage, a.APIPackage, k))
		defaultImports = append(defaultImports, importPath)
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
		AppName:             swag.ToGoName(a.Name),
		HumanAppName:        swag.ToHumanNameLower(a.Name),
		Name:                swag.ToJSONName(a.Name),
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
		Principal:           a.Principal,
		SwaggerJSON:         fmt.Sprintf("%#v", jsonb),
	}, nil
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
	OperationGroups     []genOperationGroup
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

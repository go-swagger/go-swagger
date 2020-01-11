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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"sort"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
)

// GenerateServer generates a server application
func GenerateServer(name string, modelNames, operationIDs []string, opts *GenOpts) error {
	generator, err := newAppGenerator(name, modelNames, operationIDs, opts)
	if err != nil {
		return err
	}
	return generator.Generate()
}

// GenerateSupport generates the supporting files for an API
func GenerateSupport(name string, modelNames, operationIDs []string, opts *GenOpts) error {
	generator, err := newAppGenerator(name, modelNames, operationIDs, opts)
	if err != nil {
		return err
	}
	return generator.GenerateSupport(nil)
}

func newAppGenerator(name string, modelNames, operationIDs []string, opts *GenOpts) (*appGenerator, error) {
	if err := opts.CheckOpts(); err != nil {
		return nil, err
	}

	if err := opts.setTemplates(); err != nil {
		return nil, err
	}

	specDoc, analyzed, err := opts.analyzeSpec()
	if err != nil {
		return nil, err
	}

	models, err := gatherModels(specDoc, modelNames)
	if err != nil {
		return nil, err
	}

	operations := gatherOperations(analyzed, operationIDs)

	if len(operations) == 0 {
		return nil, errors.New("no operations were selected")
	}

	opts.Name = appNameOrDefault(specDoc, name, defaultServerName)
	apiPackage := opts.LanguageOpts.ManglePackagePath(opts.APIPackage, defaultOperationsTarget)
	return &appGenerator{
		Name:              opts.Name,
		Receiver:          "o",
		SpecDoc:           specDoc,
		Analyzed:          analyzed,
		Models:            models,
		Operations:        operations,
		Target:            opts.Target,
		DumpData:          opts.DumpData,
		Package:           opts.LanguageOpts.ManglePackageName(apiPackage, defaultOperationsTarget),
		APIPackage:        apiPackage,
		ModelsPackage:     opts.LanguageOpts.ManglePackagePath(opts.ModelPackage, defaultModelsTarget),
		ServerPackage:     opts.LanguageOpts.ManglePackagePath(opts.ServerPackage, defaultServerTarget),
		ClientPackage:     opts.LanguageOpts.ManglePackagePath(opts.ClientPackage, defaultClientTarget),
		OperationsPackage: filepath.Join(opts.LanguageOpts.ManglePackagePath(opts.ServerPackage, defaultServerTarget), apiPackage),
		Principal:         opts.Principal,
		DefaultScheme:     opts.DefaultScheme,
		DefaultProduces:   opts.DefaultProduces,
		DefaultConsumes:   opts.DefaultConsumes,
		GenOpts:           opts,
	}, nil
}

type appGenerator struct {
	Name              string
	Receiver          string
	SpecDoc           *loads.Document
	Analyzed          *analysis.Spec
	Package           string
	APIPackage        string
	ModelsPackage     string
	ServerPackage     string
	ClientPackage     string
	OperationsPackage string
	Principal         string
	Models            map[string]spec.Schema
	Operations        map[string]opRef
	Target            string
	DumpData          bool
	DefaultScheme     string
	DefaultProduces   string
	DefaultConsumes   string
	GenOpts           *GenOpts
}

func (a *appGenerator) Generate() error {
	app, err := a.makeCodegenApp()
	if err != nil {
		return err
	}

	if a.DumpData {
		return dumpData(app)
	}

	// NOTE: relative to previous implem with chan.
	// IPC removed concurrent execution because of the FuncMap that is being shared
	// templates are now lazy loaded so there is concurrent map access I can't guard
	if a.GenOpts.IncludeModel {
		log.Printf("rendering %d models", len(app.Models))
		for _, mod := range app.Models {
			modCopy := mod
			modCopy.IncludeValidator = true // a.GenOpts.IncludeValidator
			modCopy.IncludeModel = true
			if err := a.GenOpts.renderDefinition(&modCopy); err != nil {
				return err
			}
		}
	}

	if a.GenOpts.IncludeHandler {
		log.Printf("rendering %d operation groups (tags)", app.OperationGroups.Len())
		for _, opg := range app.OperationGroups {
			opgCopy := opg
			log.Printf("rendering %d operations for %s", opg.Operations.Len(), opg.Name)
			for _, op := range opgCopy.Operations {
				opCopy := op

				if err := a.GenOpts.renderOperation(&opCopy); err != nil {
					return err
				}
			}
			// optional OperationGroups templates generation
			opGroup := opg
			opGroup.DefaultImports = app.DefaultImports
			if err := a.GenOpts.renderOperationGroup(&opGroup); err != nil {
				return fmt.Errorf("error while rendering operation group: %v", err)
			}
		}
	}

	if a.GenOpts.IncludeSupport {
		log.Printf("rendering support")
		if err := a.GenerateSupport(&app); err != nil {
			return err
		}
	}
	return nil
}

func (a *appGenerator) GenerateSupport(ap *GenApp) error {
	app := ap
	if ap == nil {
		// allows for calling GenerateSupport standalone
		ca, err := a.makeCodegenApp()
		if err != nil {
			return err
		}
		app = &ca
	}
	baseImport := a.GenOpts.LanguageOpts.baseImport(a.Target)
	importPath := path.Join(baseImport, a.GenOpts.LanguageOpts.ManglePackagePath(a.OperationsPackage, ""))
	app.DefaultImports = append(
		app.DefaultImports,
		path.Join(baseImport, a.GenOpts.LanguageOpts.ManglePackagePath(a.ServerPackage, "")),
		importPath,
	)

	return a.GenOpts.renderApplication(app)
}

func (a *appGenerator) makeSecuritySchemes() GenSecuritySchemes {
	if a.Principal == "" {
		a.Principal = "interface{}"
	}
	requiredSecuritySchemes := make(map[string]spec.SecurityScheme, len(a.Analyzed.RequiredSecuritySchemes()))
	for _, scheme := range a.Analyzed.RequiredSecuritySchemes() {
		if req, ok := a.SpecDoc.Spec().SecurityDefinitions[scheme]; ok && req != nil {
			requiredSecuritySchemes[scheme] = *req
		}
	}
	return gatherSecuritySchemes(requiredSecuritySchemes, a.Name, a.Principal, a.Receiver)
}

func (a *appGenerator) makeCodegenApp() (GenApp, error) {
	log.Println("building a plan for generation")
	sw := a.SpecDoc.Spec()
	receiver := a.Receiver

	var defaultImports []string

	jsonb, _ := json.MarshalIndent(a.SpecDoc.OrigSpec(), "", "  ")
	flatjsonb, _ := json.MarshalIndent(a.SpecDoc.Spec(), "", "  ")

	consumes, _ := a.makeConsumes()
	produces, _ := a.makeProduces()
	security := a.makeSecuritySchemes()
	baseImport := a.GenOpts.LanguageOpts.baseImport(a.Target)
	var imports = make(map[string]string)

	var genMods GenDefinitions
	importPath := a.GenOpts.ExistingModels
	if a.GenOpts.ExistingModels == "" {
		imports[a.GenOpts.LanguageOpts.ManglePackageName(a.ModelsPackage, defaultModelsTarget)] = path.Join(
			baseImport,
			a.GenOpts.LanguageOpts.ManglePackagePath(a.GenOpts.ModelPackage, defaultModelsTarget))
	}
	if importPath != "" {
		defaultImports = append(defaultImports, importPath)
	}

	log.Println("planning definitions")
	for mn, m := range a.Models {
		mod, err := makeGenDefinition(
			mn,
			a.ModelsPackage,
			m,
			a.SpecDoc,
			a.GenOpts,
		)
		if err != nil {
			return GenApp{}, fmt.Errorf("error in model %s while planning definitions: %v", mn, err)
		}
		if mod != nil {
			if !mod.External {
				genMods = append(genMods, *mod)
			}

			// Copy model imports to operation imports
			for alias, pkg := range mod.Imports {
				target := a.GenOpts.LanguageOpts.ManglePackageName(alias, "")
				imports[target] = pkg
			}
		}
	}
	sort.Sort(genMods)

	log.Println("planning operations")
	tns := make(map[string]struct{})
	var genOps GenOperations
	for on, opp := range a.Operations {
		o := opp.Op
		o.Tags = pruneEmpty(o.Tags)
		o.ID = on

		bldr := codeGenOpBuilder{
			ModelsPackage:    a.ModelsPackage,
			Principal:        a.Principal,
			Target:           a.Target,
			DefaultImports:   defaultImports,
			Imports:          imports,
			DefaultScheme:    a.DefaultScheme,
			Doc:              a.SpecDoc,
			Analyzed:         a.Analyzed,
			BasePath:         a.SpecDoc.BasePath(),
			GenOpts:          a.GenOpts,
			Name:             on, // TODO: change operation name to something safe
			Operation:        *o,
			Method:           opp.Method,
			Path:             opp.Path,
			IncludeValidator: true,
			APIPackage:       a.APIPackage, // defaults to main operations package
		}

		bldr.Authed = len(a.Analyzed.SecurityRequirementsFor(o)) > 0
		bldr.Security = a.Analyzed.SecurityRequirementsFor(o)
		bldr.SecurityDefinitions = a.Analyzed.SecurityDefinitionsFor(o)
		bldr.RootAPIPackage = a.GenOpts.LanguageOpts.ManglePackageName(a.ServerPackage, defaultServerTarget)

		st := o.Tags
		if a.GenOpts != nil {
			st = a.GenOpts.Tags
		}
		intersected := intersectTags(o.Tags, st)
		if len(st) > 0 && len(intersected) == 0 {
			continue
		}

		if len(intersected) > 0 {
			tag := intersected[0]
			bldr.APIPackage = a.GenOpts.LanguageOpts.ManglePackagePath(tag, a.APIPackage)
			for _, t := range intersected {
				tns[t] = struct{}{}
			}
		}
		op, err := bldr.MakeOperation()
		if err != nil {
			return GenApp{}, err
		}
		op.ReceiverName = receiver
		op.Tags = intersected
		genOps = append(genOps, op)

	}
	for k := range tns {
		importPath := filepath.ToSlash(
			path.Join(
				baseImport,
				a.GenOpts.LanguageOpts.ManglePackagePath(a.OperationsPackage, ""),
				swag.ToFileName(k)))
		defaultImports = append(defaultImports, importPath)
	}
	sort.Sort(genOps)

	log.Println("grouping operations into packages")
	opsGroupedByPackage := make(map[string]GenOperations)
	for _, operation := range genOps {
		if operation.Package == "" {
			operation.Package = a.Package
		}
		opsGroupedByPackage[operation.Package] = append(opsGroupedByPackage[operation.Package], operation)
	}

	var opGroups GenOperationGroups
	for k, v := range opsGroupedByPackage {
		sort.Sort(v)
		// trim duplicate extra schemas within the same package
		vv := make(GenOperations, 0, len(v))
		seenExtraSchema := make(map[string]bool)
		for _, op := range v {
			uniqueExtraSchemas := make(GenSchemaList, 0, len(op.ExtraSchemas))
			for _, xs := range op.ExtraSchemas {
				if _, alreadyThere := seenExtraSchema[xs.Name]; !alreadyThere {
					seenExtraSchema[xs.Name] = true
					uniqueExtraSchemas = append(uniqueExtraSchemas, xs)
				}
			}
			op.ExtraSchemas = uniqueExtraSchemas
			vv = append(vv, op)
		}

		opGroup := GenOperationGroup{
			GenCommon: GenCommon{
				Copyright:        a.GenOpts.Copyright,
				TargetImportPath: baseImport,
			},
			Name:           k,
			Operations:     vv,
			DefaultImports: defaultImports,
			Imports:        imports,
			RootPackage:    a.APIPackage,
			GenOpts:        a.GenOpts,
		}
		opGroups = append(opGroups, opGroup)
		var importPath string
		if k == a.APIPackage {
			importPath = path.Join(baseImport, a.GenOpts.LanguageOpts.ManglePackagePath(a.OperationsPackage, ""))
		} else {
			importPath = path.Join(baseImport, a.GenOpts.LanguageOpts.ManglePackagePath(a.OperationsPackage, ""), k)
		}
		defaultImports = append(defaultImports, importPath)
	}
	sort.Sort(opGroups)

	log.Println("planning meta data and facades")

	var collectedSchemes []string
	var extraSchemes []string
	for _, op := range genOps {
		collectedSchemes = concatUnique(collectedSchemes, op.Schemes)
		extraSchemes = concatUnique(extraSchemes, op.ExtraSchemes)
	}
	sort.Strings(collectedSchemes)
	sort.Strings(extraSchemes)

	host := "localhost"
	if sw.Host != "" {
		host = sw.Host
	}

	basePath := "/"
	if sw.BasePath != "" {
		basePath = sw.BasePath
	}

	return GenApp{
		GenCommon: GenCommon{
			Copyright:        a.GenOpts.Copyright,
			TargetImportPath: baseImport,
		},
		APIPackage:          a.GenOpts.LanguageOpts.ManglePackageName(a.ServerPackage, defaultServerTarget),
		Package:             a.Package,
		ReceiverName:        receiver,
		Name:                a.Name,
		Host:                host,
		BasePath:            basePath,
		Schemes:             schemeOrDefault(collectedSchemes, a.DefaultScheme),
		ExtraSchemes:        extraSchemes,
		ExternalDocs:        sw.ExternalDocs,
		Info:                sw.Info,
		Consumes:            consumes,
		Produces:            produces,
		DefaultConsumes:     a.DefaultConsumes,
		DefaultProduces:     a.DefaultProduces,
		DefaultImports:      defaultImports,
		Imports:             imports,
		SecurityDefinitions: security,
		Models:              genMods,
		Operations:          genOps,
		OperationGroups:     opGroups,
		Principal:           a.Principal,
		SwaggerJSON:         generateReadableSpec(jsonb),
		FlatSwaggerJSON:     generateReadableSpec(flatjsonb),
		ExcludeSpec:         a.GenOpts != nil && a.GenOpts.ExcludeSpec,
		GenOpts:             a.GenOpts,
	}, nil
}

// generateReadableSpec makes swagger json spec as a string instead of bytes
// the only character that needs to be escaped is '`' symbol, since it cannot be escaped in the GO string
// that is quoted as `string data`. The function doesn't care about the beginning or the ending of the
// string it escapes since all data that needs to be escaped is always in the middle of the swagger spec.
func generateReadableSpec(spec []byte) string {
	buf := &bytes.Buffer{}
	for _, b := range string(spec) {
		if b == '`' {
			buf.WriteString("`+\"`\"+`")
		} else {
			buf.WriteRune(b)
		}
	}
	return buf.String()
}

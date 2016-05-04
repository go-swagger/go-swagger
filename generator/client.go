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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
)

// GenerateClient generates a client library for a swagger spec document.
func GenerateClient(name string, modelNames, operationIDs []string, opts GenOpts) error {

	if opts.TemplateDir != "" {
		if err := templates.LoadDir(opts.TemplateDir); err != nil {
			return err
		}
	}

	compileTemplates()

	// Load the spec
	_, specDoc, err := loadSpec(opts.Spec)
	if err != nil {
		return err
	}
	analyzed := analysis.New(specDoc.Spec())

	models, err := gatherModels(specDoc, modelNames)
	if err != nil {
		return err
	}
	operations := gatherOperations(analyzed, operationIDs)

	defaultScheme := opts.DefaultScheme
	if defaultScheme == "" {
		defaultScheme = sHTTP
	}

	defaultConsumes := opts.DefaultConsumes
	if defaultConsumes == "" {
		defaultConsumes = runtime.JSONMime
	}

	defaultProduces := opts.DefaultProduces
	if defaultProduces == "" {
		defaultProduces = runtime.JSONMime
	}

	generator := appGenerator{
		Name:            appNameOrDefault(specDoc, name, "rest"),
		SpecDoc:         specDoc,
		Analyzed:        analyzed,
		Models:          models,
		Operations:      operations,
		Target:          opts.Target,
		DumpData:        opts.DumpData,
		Package:         mangleName(swag.ToFileName(opts.ClientPackage), "client"),
		APIPackage:      mangleName(swag.ToFileName(opts.APIPackage), "api"),
		ModelsPackage:   mangleName(swag.ToFileName(opts.ModelPackage), "definitions"),
		ServerPackage:   mangleName(swag.ToFileName(opts.ServerPackage), "server"),
		ClientPackage:   mangleName(swag.ToFileName(opts.ClientPackage), "client"),
		Principal:       opts.Principal,
		DefaultScheme:   defaultScheme,
		DefaultProduces: defaultProduces,
		DefaultConsumes: defaultConsumes,
		GenOpts:         &opts,
	}
	generator.Receiver = "o"

	return (&clientGenerator{generator}).Generate()
}

type clientGenerator struct {
	appGenerator
}

func (c *clientGenerator) Generate() error {
	app, err := c.makeCodegenApp()
	if app.Name == "" {
		app.Name = "APIClient"
	}
	app.DefaultImports = []string{filepath.ToSlash(filepath.Join(baseImport(c.Target), c.ModelsPackage))}
	if err != nil {
		return err
	}

	if c.DumpData {
		bb, _ := json.MarshalIndent(swag.ToDynamicJSON(app), "", "  ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}

	if c.GenOpts.IncludeModel {
		for _, mod := range app.Models {
			mod.IncludeValidator = true // a.GenOpts.IncludeValidator
			gen := &definitionGenerator{
				Name:    mod.Name,
				SpecDoc: c.SpecDoc,
				Target:  filepath.Join(c.Target, c.ModelsPackage),
				Data:    &mod,
			}
			if err := gen.generateModel(); err != nil {
				return err
			}
		}
	}

	if c.GenOpts.IncludeHandler {

		for i := range app.OperationGroups {
			opGroup := app.OperationGroups[i]
			opGroup.DefaultImports = []string{filepath.ToSlash(filepath.Join(baseImport(c.Target), c.ModelsPackage))}
			opGroup.RootPackage = c.ClientPackage
			app.OperationGroups[i] = opGroup
			sort.Sort(opGroup.Operations)
			for _, op := range opGroup.Operations {
				if op.Package == "" {
					op.Package = c.Package
				}
				if err := c.generateParameters(&op); err != nil {
					return err
				}

				if err := c.generateResponses(&op); err != nil {
					return err
				}
			}
			app.DefaultImports = append(app.DefaultImports, filepath.ToSlash(filepath.Join(baseImport(c.Target), c.ClientPackage, opGroup.Name)))
			if err := c.generateGroupClient(opGroup); err != nil {
				return err
			}
		}

		sort.Sort(app.OperationGroups)
	}

	if c.GenOpts.IncludeSupport {
		if err := c.generateFacade(&app); err != nil {
			return err
		}
	}

	return nil
}

func (c *clientGenerator) generateParameters(op *GenOperation) error {
	buf := bytes.NewBuffer(nil)

	if err := clientParamTemplate.Execute(buf, op); err != nil {
		return err
	}
	log.Println("rendered client parameters template:", op.Package+"."+swag.ToGoName(op.Name)+"Parameters")

	fp := filepath.Join(c.Target, c.ClientPackage)
	if len(op.Package) > 0 {
		fp = filepath.Join(fp, op.Package)
	}
	return writeToFile(fp, swag.ToGoName(op.Name)+"Parameters", buf.Bytes())
}

func (c *clientGenerator) generateResponses(op *GenOperation) error {
	buf := bytes.NewBuffer(nil)

	if err := clientResponseTemplate.Execute(buf, op); err != nil {
		return err
	}
	log.Println("rendered client responses template:", op.Package+"."+swag.ToGoName(op.Name)+"Responses")

	fp := filepath.Join(c.Target, c.ClientPackage)
	if len(op.Package) > 0 {
		fp = filepath.Join(fp, op.Package)
	}
	return writeToFile(fp, swag.ToGoName(op.Name)+"Responses", buf.Bytes())
}

func (c *clientGenerator) generateGroupClient(opGroup GenOperationGroup) error {
	buf := bytes.NewBuffer(nil)

	if err := clientTemplate.Execute(buf, opGroup); err != nil {
		return err
	}
	log.Println("rendered operation group client template:", opGroup.Name+"."+swag.ToGoName(opGroup.Name)+"Client")

	fp := filepath.Join(c.Target, c.ClientPackage, opGroup.Name)
	return writeToFile(fp, swag.ToGoName(opGroup.Name)+"Client", buf.Bytes())
}

func (c *clientGenerator) generateFacade(app *GenApp) error {
	buf := bytes.NewBuffer(nil)

	if err := clientFacadeTemplate.Execute(buf, app); err != nil {
		return err
	}
	log.Println("rendered client facade template:", c.ClientPackage+"."+swag.ToGoName(app.Name)+"Client")

	fp := filepath.Join(c.Target, c.ClientPackage)
	return writeToFile(fp, swag.ToGoName(app.Name)+"Client", buf.Bytes())
}

func (c *clientGenerator) generateEmbeddedSwaggerJSON(app *GenApp) error {
	buf := bytes.NewBuffer(nil)

	if err := embeddedSpecTemplate.Execute(buf, app); err != nil {
		return err
	}
	log.Println("rendered client embedded swagger JSON template:", c.ClientPackage+"."+swag.ToGoName(app.Name)+"Client")

	fp := filepath.Join(c.Target, c.ClientPackage)
	return writeToFile(fp, swag.ToGoName(app.Name)+"EmbeddedSpec", buf.Bytes())
}

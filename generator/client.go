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

	"github.com/go-swagger/go-swagger/swag"
)

// GenerateClient generates a client library for a swagger spec document.
func GenerateClient(name string, modelNames, operationIDs []string, opts GenOpts) error {
	// Load the spec
	_, specDoc, err := loadSpec(opts.Spec)
	if err != nil {
		return err
	}

	models := gatherModels(specDoc, modelNames)
	operations := gatherOperations(specDoc, operationIDs)

	generator := appGenerator{
		Name:          appNameOrDefault(specDoc, name, "swagger"),
		SpecDoc:       specDoc,
		Models:        models,
		Operations:    operations,
		Target:        opts.Target,
		DumpData:      opts.DumpData,
		Package:       opts.APIPackage,
		APIPackage:    opts.APIPackage,
		ModelsPackage: opts.ModelPackage,
		ServerPackage: opts.ServerPackage,
		ClientPackage: opts.ClientPackage,
		Principal:     opts.Principal,
	}
	generator.Receiver = "o"

	return (&clientGenerator{generator}).Generate()
}

type clientGenerator struct {
	appGenerator
}

func (c *clientGenerator) Generate() error {
	app, err := c.makeCodegenApp()
	app.DefaultImports = []string{filepath.ToSlash(filepath.Join(baseImport(c.Target), c.ModelsPackage))}
	if err != nil {
		return err
	}

	if c.DumpData {
		bb, _ := json.MarshalIndent(swag.ToDynamicJSON(app), "", "  ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}

	opsGroupedByTag := make(map[string]GenOperations)
	for _, operation := range app.Operations {
		if operation.Package == "" {
			operation.Package = c.Package
		}
		if err := c.generateParameters(&operation); err != nil {
			return err
		}

		if err := c.generateResponses(&operation); err != nil {
			return err
		}
		opsGroupedByTag[operation.Package] = append(opsGroupedByTag[operation.Package], operation)
	}

	for k, v := range opsGroupedByTag {
		sort.Sort(v)
		opGroup := GenOperationGroup{
			Name:           k,
			Operations:     v,
			DefaultImports: []string{filepath.ToSlash(filepath.Join(baseImport(c.Target), c.ModelsPackage))},
		}
		app.OperationGroups = append(app.OperationGroups, opGroup)
		app.DefaultImports = append(app.DefaultImports, filepath.ToSlash(filepath.Join(baseImport(c.Target), c.ClientPackage, k)))
		if err := c.generateGroupClient(opGroup); err != nil {
			return err
		}
	}
	sort.Sort(app.OperationGroups)

	if err := c.generateEmbeddedSwaggerJSON(&app); err != nil {
		return err
	}

	if err := c.generateFacade(&app); err != nil {
		return err
	}

	return nil
}

func (c *clientGenerator) generateParameters(op *GenOperation) error {
	buf := bytes.NewBuffer(nil)

	if err := clientParamTemplate.Execute(buf, op); err != nil {
		return err
	}
	log.Println("rendered client parameters template:", op.Package+"."+swag.ToGoName(op.Name)+"Parameters")

	fp := filepath.Join(c.ClientPackage, c.Target)
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

	fp := filepath.Join(c.ClientPackage, c.Target)
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

	fp := filepath.Join(c.ClientPackage, c.Target, opGroup.Name)
	return writeToFile(fp, swag.ToGoName(opGroup.Name)+"Client", buf.Bytes())
}

func (c *clientGenerator) generateFacade(app *GenApp) error {
	buf := bytes.NewBuffer(nil)

	if err := clientFacadeTemplate.Execute(buf, app); err != nil {
		return err
	}
	log.Println("rendered client facade template:", c.ClientPackage+"."+swag.ToGoName(app.Name)+"Client")

	fp := filepath.Join(c.ClientPackage, c.Target)
	return writeToFile(fp, swag.ToGoName(app.Name)+"Client", buf.Bytes())
}

func (c *clientGenerator) generateEmbeddedSwaggerJSON(app *GenApp) error {
	buf := bytes.NewBuffer(nil)

	if err := embeddedSpecTemplate.Execute(buf, app); err != nil {
		return err
	}
	log.Println("rendered client embedded swagger JSON template:", c.ClientPackage+"."+swag.ToGoName(app.Name)+"Client")

	fp := filepath.Join(c.ClientPackage, c.Target)
	return writeToFile(fp, swag.ToGoName(app.Name)+"EmbeddedSpec", buf.Bytes())
}

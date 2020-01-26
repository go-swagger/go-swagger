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
	"errors"
	"path"
	"sort"

	"github.com/go-openapi/swag"
)

// GenerateClient generates a client library for a swagger spec document.
func GenerateClient(name string, modelNames, operationIDs []string, opts *GenOpts) error {
	if err := opts.CheckOpts(); err != nil {
		return err
	}

	if err := opts.setTemplates(); err != nil {
		return err
	}

	specDoc, analyzed, err := opts.analyzeSpec()
	if err != nil {
		return err
	}

	models, err := gatherModels(specDoc, modelNames)
	if err != nil {
		return err
	}

	operations := gatherOperations(analyzed, operationIDs)
	if len(operations) == 0 {
		return errors.New("no operations were selected")
	}

	generator := appGenerator{
		Name:              appNameOrDefault(specDoc, name, defaultClientName),
		SpecDoc:           specDoc,
		Analyzed:          analyzed,
		Models:            models,
		Operations:        operations,
		Target:            opts.Target,
		DumpData:          opts.DumpData,
		Package:           opts.LanguageOpts.ManglePackageName(opts.ClientPackage, defaultClientTarget),
		APIPackage:        opts.LanguageOpts.ManglePackagePath(opts.APIPackage, defaultOperationsTarget),
		ModelsPackage:     opts.LanguageOpts.ManglePackagePath(opts.ModelPackage, defaultModelsTarget),
		ServerPackage:     opts.LanguageOpts.ManglePackagePath(opts.ServerPackage, defaultServerTarget),
		ClientPackage:     opts.LanguageOpts.ManglePackagePath(opts.ClientPackage, defaultClientTarget),
		OperationsPackage: opts.LanguageOpts.ManglePackagePath(opts.ClientPackage, defaultClientTarget),
		Principal:         opts.Principal,
		DefaultScheme:     opts.DefaultScheme,
		DefaultProduces:   opts.DefaultProduces,
		DefaultConsumes:   opts.DefaultConsumes,
		GenOpts:           opts,
	}
	generator.Receiver = "o"
	return (&clientGenerator{generator}).Generate()
}

type clientGenerator struct {
	appGenerator
}

func (c *clientGenerator) Generate() error {
	app, err := c.makeCodegenApp()
	if err != nil {
		return err
	}

	if app.Name == "" {
		app.Name = "APIClient"
	}
	baseImport := c.GenOpts.LanguageOpts.baseImport(c.Target)

	if c.GenOpts.ExistingModels == "" {
		if app.Imports == nil {
			app.Imports = make(map[string]string)
		}
		pkgAlias := c.GenOpts.LanguageOpts.ManglePackageName(c.ModelsPackage, defaultModelsTarget)
		app.Imports[pkgAlias] = path.Join(
			baseImport,
			c.GenOpts.LanguageOpts.ManglePackagePath(c.GenOpts.ModelPackage, defaultModelsTarget))
	} else {
		app.DefaultImports = append(app.DefaultImports, c.GenOpts.LanguageOpts.ManglePackagePath(c.GenOpts.ExistingModels, ""))
	}

	if c.DumpData {
		return dumpData(swag.ToDynamicJSON(app))
	}

	if c.GenOpts.IncludeModel {
		for _, mod := range app.Models {
			modCopy := mod
			modCopy.IncludeValidator = true
			if mod.IsStream {
				continue
			}
			if err := c.GenOpts.renderDefinition(&modCopy); err != nil {
				return err
			}
		}
	}

	if c.GenOpts.IncludeHandler {
		sort.Sort(app.OperationGroups)
		for i := range app.OperationGroups {
			opGroup := app.OperationGroups[i]
			opGroup.DefaultImports = app.DefaultImports
			opGroup.RootPackage = c.ClientPackage
			opGroup.GenOpts = c.GenOpts
			app.OperationGroups[i] = opGroup
			sort.Sort(opGroup.Operations)
			for _, op := range opGroup.Operations {
				opCopy := op
				if opCopy.Package == "" {
					opCopy.Package = c.Package
				}
				if err := c.GenOpts.renderOperation(&opCopy); err != nil {
					return err
				}
			}
			app.DefaultImports = append(app.DefaultImports,
				path.Join(
					baseImport,
					c.GenOpts.LanguageOpts.ManglePackagePath(c.ClientPackage, defaultClientTarget),
					opGroup.Name))

			if err := c.GenOpts.renderOperationGroup(&opGroup); err != nil {
				return err
			}
		}
	}

	if c.GenOpts.IncludeSupport {
		if err := c.GenOpts.renderApplication(&app); err != nil {
			return err
		}
	}

	return nil
}

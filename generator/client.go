package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	return (&clientGenerator{generator}).Generate()
}

type clientGenerator struct {
	appGenerator
}

func (c *clientGenerator) Generate() error {
	app, err := c.makeCodegenApp()
	app.DefaultImports = []string{filepath.Join(baseImport(c.Target), c.ModelsPackage)}
	if err != nil {
		return err
	}

	if c.DumpData {
		bb, _ := json.MarshalIndent(swag.ToDynamicJSON(app), "", "  ")
		fmt.Fprintln(os.Stdout, string(bb))
		return nil
	}

	opsGroupedByTag := make(map[string][]genOperation)
	for _, operation := range app.Operations {
		if err := c.generateParameters(&operation); err != nil {
			return err
		}

		if err := c.generateResponses(&operation); err != nil {
			return err
		}
		opsGroupedByTag[operation.Package] = append(opsGroupedByTag[operation.Package], operation)
	}

	for k, v := range opsGroupedByTag {
		opGroup := genOperationGroup{
			Name:           k,
			Operations:     v,
			DefaultImports: []string{filepath.Join(baseImport(c.Target), c.ModelsPackage)},
		}
		app.OperationGroups = append(app.OperationGroups, opGroup)
		app.DefaultImports = append(app.DefaultImports, filepath.Join(baseImport(c.Target), c.ClientPackage, k))
		if err := c.generateGroupClient(opGroup); err != nil {
			return err
		}
	}

	if err := c.generateFacade(&app); err != nil {
		return err
	}

	return nil
}

func (c *clientGenerator) generateParameters(op *genOperation) error {
	buf := bytes.NewBuffer(nil)

	if err := clientParamTemplate.Execute(buf, op); err != nil {
		return err
	}
	log.Println("rendered client parameters template:", op.Package+"."+op.ClassName+"Parameters")

	fp := filepath.Join(c.ClientPackage, c.Target)
	if len(op.Package) > 0 {
		fp = filepath.Join(fp, op.Package)
	}
	return writeToFile(fp, op.ClassName+"Parameters", buf.Bytes())
}

func (c *clientGenerator) generateResponses(op *genOperation) error {
	buf := bytes.NewBuffer(nil)

	if err := clientResponseTemplate.Execute(buf, op); err != nil {
		return err
	}
	log.Println("rendered client responses template:", op.Package+"."+op.ClassName+"Responses")

	fp := filepath.Join(c.ClientPackage, c.Target)
	if len(op.Package) > 0 {
		fp = filepath.Join(fp, op.Package)
	}
	return writeToFile(fp, op.ClassName+"Responses", buf.Bytes())
}

func (c *clientGenerator) generateGroupClient(opGroup genOperationGroup) error {
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

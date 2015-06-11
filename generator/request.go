package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/swag"
)

// GenerateClientRequest generates a parameter model, parameter validator, http request builder implementations for a given operation
// It also generates an operation handler interface that uses the parameter model for handling a valid request.
// Allows for specifying a list of tags to include only certain tags for the generation
func GenerateClientRequest(operationNames, tags []string, includeHandler, includeParameters bool, opts GenOpts) error {
	// Load the spec
	specPath, specDoc, err := loadSpec(opts.Spec)
	if err != nil {
		return err
	}

	if len(operationNames) == 0 {
		operationNames = specDoc.OperationIDs()
	}

	for _, operationName := range operationNames {
		operation, ok := specDoc.OperationForName(operationName)
		if !ok {
			return fmt.Errorf("operation %q not found in %s", operationName, specPath)
		}

		var generator requestGenerator
		generator.Name = operationName
		generator.APIPackage = opts.APIPackage
		generator.ModelsPackage = opts.ModelPackage
		generator.ClientPackage = opts.ClientPackage
		generator.ServerPackage = opts.ServerPackage
		generator.Operation = *operation
		generator.SecurityRequirements = specDoc.SecurityRequirementsFor(operation)
		generator.Principal = opts.Principal
		generator.Target = filepath.Join(opts.Target, opts.APIPackage)
		generator.Tags = tags
		generator.IncludeHandler = includeHandler
		generator.IncludeParameters = includeParameters
		generator.DumpData = opts.DumpData
		generator.Doc = specDoc
		if err := generator.Generate(); err != nil {
			return err
		}
	}
	return nil
}

type requestGenerator struct {
	Name                 string
	Authorized           bool
	APIPackage           string
	ModelsPackage        string
	ServerPackage        string
	ClientPackage        string
	Operation            spec.Operation
	SecurityRequirements []spec.SecurityRequirement
	Principal            string
	Target               string
	Tags                 []string
	data                 interface{}
	pkg                  string
	cname                string
	IncludeHandler       bool
	IncludeParameters    bool
	DumpData             bool
	Doc                  *spec.Document
}

func (r *requestGenerator) Generate() error {
	// Build a list of codegen operations based on the tags,
	// the tag decides the actual package for an operation
	// the user specified package serves as root for generating the directory structure
	var operations []genOperation
	authed := len(r.SecurityRequirements) > 0

	var bldr codeGenOpBuilder
	bldr.Name = r.Name
	bldr.ModelsPackage = r.ModelsPackage
	bldr.Principal = r.Principal
	bldr.Target = r.Target
	bldr.Operation = r.Operation
	bldr.Authed = authed
	bldr.Doc = r.Doc

	for _, tag := range r.Operation.Tags {
		if len(r.Tags) == 0 {
			bldr.APIPackage = tag
			op, err := makeCodegenOperation(bldr)
			if err != nil {
				return err
			}
			operations = append(operations, op)
			continue
		}
		for _, ft := range r.Tags {
			if ft == tag {
				bldr.APIPackage = tag
				op, err := makeCodegenOperation(bldr)
				if err != nil {
					return err
				}
				operations = append(operations, op)
				break
			}
		}

	}
	if len(operations) == 0 {
		bldr.APIPackage = r.ClientPackage
		op, err := makeCodegenOperation(bldr)
		if err != nil {
			return err
		}
		operations = append(operations, op)
	}

	for _, op := range operations {
		if r.DumpData {
			bb, _ := json.MarshalIndent(swag.ToDynamicJSON(op), "", " ")
			fmt.Fprintln(os.Stdout, string(bb))
			continue
		}
		r.data = op
		r.pkg = op.Package
		r.cname = op.ClassName

		if r.IncludeHandler {
			if err := r.generateHandler(); err != nil {
				return fmt.Errorf("handler: %s", err)
			}
			log.Println("generated handler", op.Package+"."+op.ClassName)
		}

		if r.IncludeParameters && len(r.Operation.Parameters) > 0 {
			if err := r.generateParameterModel(); err != nil {
				return fmt.Errorf("parameters: %s", err)
			}
			log.Println("generated parameters", op.Package+"."+op.ClassName+"Parameters")
		}

		if len(r.Operation.Parameters) == 0 {
			log.Println("no parameters for operation", op.Package+"."+op.ClassName)
		}
	}

	return nil
}

func (r *requestGenerator) generateHandler() error {
	//buf := bytes.NewBuffer(nil)

	//if err := operationTemplate.Execute(buf, r.data); err != nil {
	//return err
	//}
	log.Println("rendered handler template:", r.pkg+"."+r.cname)

	//fp := filepath.Join(o.ServerPackage, o.Target)
	//if len(o.Operation.Tags) > 0 {
	//fp = filepath.Join(fp, o.pkg)
	//}
	//return writeToFile(fp, o.Name, buf.Bytes())
	return nil
}

func (r *requestGenerator) generateParameterModel() error {
	buf := bytes.NewBuffer(nil)

	if err := clientParamTemplate.Execute(buf, r.data); err != nil {
		return err
	}
	log.Println("rendered parameters template:", r.pkg+"."+r.cname+"Parameters")

	fp := filepath.Join(r.ClientPackage, r.Target)
	if len(r.Operation.Tags) > 0 {
		fp = filepath.Join(fp, r.pkg)
	}
	return writeToFile(fp, r.Name+"Parameters", buf.Bytes())
}

package spec

import (
	"encoding/json"
	"fmt"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/swagger/jsonschema"
)

// Document represents a swagger spec document
type Document struct {
	version    string
	specSchema *jsonschema.JsonSchemaDocument
	data       map[string]interface{}
	spec       *swagger.Spec
	analyzer   *specAnalyzer
}

// New creates a new shema document
func New(data json.RawMessage, version string) (*Document, error) {
	if version == "" {
		version = "2.0"
	}
	if version != "2.0" {
		return nil, fmt.Errorf("spec version %q is not supported", version)
	}

	specSchema := jsonschema.MustLoadSwagger20Schema()

	spec := new(swagger.Spec)
	if err := json.Unmarshal(data, spec); err != nil {
		return nil, err
	}

	var v map[string]interface{}
	json.Unmarshal(data, &v)

	return &Document{
		version:    version,
		specSchema: specSchema,
		data:       v,
		spec:       spec,
		analyzer:   newAnalyzer(spec),
	}, nil
}

// Version returns the version of this spec
func (d *Document) Version() string {
	return d.version
}

// Validate validates this spec document
func (d *Document) Validate() *jsonschema.ValidationResult {
	return d.specSchema.Validate(d.data)
}

// ParametersFor gets the parameters for the specified operation, collecting all the shared ones along the way
func (d *Document) ParametersFor(operation *swagger.Operation) map[string]swagger.Parameter {
	return d.analyzer.ParametersFor(operation)
}

// ProducesFor gets the mediatypes for the operation
func (d *Document) ProducesFor(operation *swagger.Operation) []string {
	return d.analyzer.ProducesFor(operation)
}

// ConsumesFor gets the mediatypes for the operation
func (d *Document) ConsumesFor(operation *swagger.Operation) []string {
	return d.analyzer.ConsumesFor(operation)
}

// AllPaths returns all the paths in the swagger spec
func (d *Document) AllPaths() map[string]swagger.PathItem {
	return d.spec.Paths.Paths
}

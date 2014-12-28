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
func (s *Document) Version() string {
	return s.version
}

// Validate validates this spec document
func (s *Document) Validate() *jsonschema.ValidationResult {
	return s.specSchema.Validate(s.data)
}

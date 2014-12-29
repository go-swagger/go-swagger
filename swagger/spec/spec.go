package spec

import (
	"encoding/json"
	"fmt"

	"github.com/casualjim/go-swagger"
	"github.com/casualjim/go-swagger/swagger/jsonschema"
)

// Document represents a swagger spec document
type Document struct {
	specAnalyzer
	specSchema *jsonschema.JsonSchemaDocument
	data       map[string]interface{}
	spec       *swagger.Spec
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

	d := &Document{
		specAnalyzer: specAnalyzer{
			spec:        spec,
			consumes:    make(map[string]struct{}),
			produces:    make(map[string]struct{}),
			authSchemes: make(map[string]struct{}),
			operations:  make(map[string]map[string]*swagger.Operation),
		},
		specSchema: specSchema,
		data:       v,
		spec:       spec,
	}
	d.initialize()
	return d, nil
}

func (d *Document) BasePath() string {
	return d.spec.BasePath
}

// Version returns the version of this spec
func (d *Document) Version() string {
	return d.spec.Swagger
}

// Validate validates this spec document
func (d *Document) Validate() *jsonschema.ValidationResult {
	return d.specSchema.Validate(d.data)
}

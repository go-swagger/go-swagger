package spec

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/casualjim/go-swagger/assets"
	"github.com/casualjim/go-swagger/jsonschema"
)

// MustLoadJSONSchemaDraft04 panics when Swagger20Schema returns an error
func MustLoadJSONSchemaDraft04() *jsonschema.Document {
	d, e := Swagger20Schema()
	if e != nil {
		panic(e)
	}
	return d
}

// JSONSchemaDraft04 loads the json schema document for json shema draft04
func JSONSchemaDraft04() (*jsonschema.Document, error) {
	b, err := assets.Asset("jsonschema-draft-04.json")
	if err != nil {
		return nil, err
	}
	loader := jsonschema.NewLoader(bytes.NewBuffer(b), "http://json-schema.org/draft-04/schema#")

	return jsonschema.Load(loader)
}

// MustLoadSwagger20Schema panics when Swagger20Schema returns an error
func MustLoadSwagger20Schema() *jsonschema.Document {
	d, e := Swagger20Schema()
	if e != nil {
		panic(e)
	}
	return d
}

// Swagger20Schema loads the swagger 2.0 schema from the embedded asses
func Swagger20Schema() (*jsonschema.Document, error) {

	b, err := assets.Asset("2.0/schema.json")
	if err != nil {
		return nil, err
	}
	loader := jsonschema.NewLoader(bytes.NewBuffer(b), "http://swagger.io/v2/schema.json#")

	return jsonschema.Load(loader)
}

// Document represents a swagger spec document
type Document struct {
	specAnalyzer
	specSchema *jsonschema.Document
	data       map[string]interface{}
	spec       *Swagger
}

var swaggerSchema *jsonschema.Document
var jsonSchema *jsonschema.Document

func init() {
	jsonSchema = MustLoadJSONSchemaDraft04()
	swaggerSchema = MustLoadSwagger20Schema()
}

func expandSpec(spec *Swagger) error {
	// TODO: add a spec expander,
	// this should be an external method that walks the tree of the swagger document
	// and loads all the references, it should make sure it doesn't get in an infinite loop
	// so it should track resolutions that are in progress and skip them
	// and use the same value when it is resolved later
	// it should keep a cache of resolutions already performed. the key for the cache of each item
	// is the json pointer string representation of the full path from root.
	//
	// things that can have a ref: schema, response, path item, parameter, items
	// note that items only have a ref to be resolved when used in a schema
	return nil
}

// New creates a new shema document
func New(data json.RawMessage, version string) (*Document, error) {
	if version == "" {
		version = "2.0"
	}
	if version != "2.0" {
		return nil, fmt.Errorf("spec version %q is not supported", version)
	}

	spec := new(Swagger)
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
			operations:  make(map[string]map[string]*Operation),
		},
		specSchema: swaggerSchema,
		data:       v,
		spec:       spec,
	}
	d.initialize()
	return d, nil
}

// BasePath the base path for this spec
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

package validate

import (
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/internal/validate"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
)

// Spec validates a spec document
func Spec(doc *spec.Document, formats strfmt.Registry) error {
	// TODO: add more validations beyond just jsonschema
	res := validate.NewSchemaValidator(doc.Schema(), nil, "", formats).Validate(doc.Spec())
	if len(res.Errors) > 0 {
		return errors.CompositeValidationError(res.Errors...)
	}
	return nil
}

// AgainstSchema validates the specified data with the provided schema, when no schema
// is provided it uses the json schema as default
func AgainstSchema(schema *spec.Schema, data interface{}, formats strfmt.Registry) error {
	res := validate.NewSchemaValidator(schema, nil, "", formats).Validate(data)
	if res.HasErrors() {
		return errors.CompositeValidationError(res.Errors...)
	}
	return nil
}

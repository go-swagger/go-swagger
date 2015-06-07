package validate

import (
	"github.com/casualjim/go-swagger/errors"
	"github.com/casualjim/go-swagger/internal/validate"
	"github.com/casualjim/go-swagger/spec"
	"github.com/casualjim/go-swagger/strfmt"
)

// Spec validates a spec document
// It validates the spec json against the json schema for swagger
// and then validates a number of extra rules that can't be expressed in json schema:
//
// 	- definition can't declare a property that's already defined by one of its ancestors
// 	- definition's ancestor can't be a descendant of the same model
// 	- each api path should be non-verbatim (account for path param names) unique per method
// 	- each security reference should contain only unique scopes
// 	- each security scope in a security definition should be unique
// 	- each path parameter should correspond to a parameter placeholder and vice versa
// 	- each referencable defintion must have references
// 	- each definition property listed in the required array must be defined in the properties of the model
// 	- each parameter should have a unique `name` and `type` combination
// 	- each operation should have only 1 parameter of type body
// 	- each reference must point to a valid object
// 	- every default value that is specified must validate against the schema for that property
// 	- items property is required for all schemas/definitions of type `array`
func Spec(doc *spec.Document, formats strfmt.Registry) error {
	errs, _ /*warns*/ := validate.NewSpecValidator(doc.Schema(), formats).Validate(doc)
	if errs.HasErrors() {
		return errors.CompositeValidationError(errs.Errors...)
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

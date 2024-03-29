// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	httpext "net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
)

// MapWithNoValidate A map of NoValidateExternal external types.
//
// If the "noValidation" hint is omitted in the definition above, this code won't build because `http.Request` has no `Validate` method.
//
// swagger:model MapWithNoValidate
type MapWithNoValidate map[string]httpext.Request

// Validate validates this map with no validate
func (m MapWithNoValidate) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this map with no validate based on context it is used
func (m MapWithNoValidate) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

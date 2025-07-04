// Code generated by go-swagger; DO NOT EDIT.

package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	stderrors "errors"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	"github.com/go-swagger/go-swagger/examples/contributed-templates/stratoscale/models"
)

// NewOrderCreateParams creates a new OrderCreateParams object
//
// There are no default values defined in the spec.
func NewOrderCreateParams() OrderCreateParams {

	return OrderCreateParams{}
}

// OrderCreateParams contains all the bound params for the order create operation
// typically these are obtained from a http.Request
//
// swagger:parameters OrderCreate
type OrderCreateParams struct {
	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*order placed for purchasing the pet
	  Required: true
	  In: body
	*/
	Body *models.Order
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewOrderCreateParams() beforehand.
func (o *OrderCreateParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer func() {
			_ = r.Body.Close()
		}()
		var body models.Order
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if stderrors.Is(err, io.EOF) {
				res = append(res, errors.Required("body", "body", ""))
			} else {
				res = append(res, errors.NewParseError("body", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(r.Context())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Body = &body
			}
		}
	} else {
		res = append(res, errors.Required("body", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// Code generated by go-swagger; DO NOT EDIT.

package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetTaskCommentsParams creates a new GetTaskCommentsParams object
// with the default values initialized.
func NewGetTaskCommentsParams() GetTaskCommentsParams {
	var (
		pageSizeDefault = int32(20)
	)
	return GetTaskCommentsParams{
		PageSize: &pageSizeDefault,
	}
}

// GetTaskCommentsParams contains all the bound params for the get task comments operation
// typically these are obtained from a http.Request
//
// swagger:parameters getTaskComments
type GetTaskCommentsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*The id of the item
	  Required: true
	  In: path
	*/
	ID int64
	/*Amount of items to return in a single page
	  In: query
	  Default: 20
	*/
	PageSize *int32
	/*The created time of the oldest seen comment
	  In: query
	*/
	Since *strfmt.DateTime
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetTaskCommentsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rID, rhkID, _ := route.Params.GetOK("id")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	qPageSize, qhkPageSize, _ := qs.GetOK("pageSize")
	if err := o.bindPageSize(qPageSize, qhkPageSize, route.Formats); err != nil {
		res = append(res, err)
	}

	qSince, qhkSince, _ := qs.GetOK("since")
	if err := o.bindSince(qSince, qhkSince, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetTaskCommentsParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("id", "path", "int64", raw)
	}
	o.ID = value

	return nil
}

func (o *GetTaskCommentsParams) bindPageSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		var pageSizeDefault int32 = int32(20)
		o.PageSize = &pageSizeDefault
		return nil
	}

	value, err := swag.ConvertInt32(raw)
	if err != nil {
		return errors.InvalidType("pageSize", "query", "int32", raw)
	}
	o.PageSize = &value

	return nil
}

func (o *GetTaskCommentsParams) bindSince(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}
	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := formats.Parse("date-time", raw)
	if err != nil {
		return errors.InvalidType("since", "query", "strfmt.DateTime", raw)
	}
	o.Since = (value.(*strfmt.DateTime))

	return nil
}

package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetTaskDetailsParams creates a new GetTaskDetailsParams object
// with the default values initialized.
func NewGetTaskDetailsParams() *GetTaskDetailsParams {
	var ()
	return &GetTaskDetailsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetTaskDetailsParamsWithTimeout creates a new GetTaskDetailsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetTaskDetailsParamsWithTimeout(timeout time.Duration) *GetTaskDetailsParams {
	var ()
	return &GetTaskDetailsParams{

		timeout: timeout,
	}
}

// NewGetTaskDetailsParamsWithContext creates a new GetTaskDetailsParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetTaskDetailsParamsWithContext(ctx context.Context) *GetTaskDetailsParams {
	var ()
	return &GetTaskDetailsParams{

		Context: ctx,
	}
}

// NewGetTaskDetailsParamsWithHTTPClient creates a new GetTaskDetailsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetTaskDetailsParamsWithHTTPClient(client *http.Client) *GetTaskDetailsParams {
	var ()
	return &GetTaskDetailsParams{
		HTTPClient: client,
	}
}

/*GetTaskDetailsParams contains all the parameters to send to the API endpoint
for the get task details operation typically these are written to a http.Request
*/
type GetTaskDetailsParams struct {

	/*ID
	  The id of the item
	  Required: true
	  In: path
	*/
	ID int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get task details params
func (o *GetTaskDetailsParams) WithTimeout(timeout time.Duration) *GetTaskDetailsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get task details params
func (o *GetTaskDetailsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get task details params
func (o *GetTaskDetailsParams) WithContext(ctx context.Context) *GetTaskDetailsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get task details params
func (o *GetTaskDetailsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get task details params
func (o *GetTaskDetailsParams) WithHTTPClient(client *http.Client) *GetTaskDetailsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get task details params
func (o *GetTaskDetailsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get task details params
func (o *GetTaskDetailsParams) WithID(id int64) *GetTaskDetailsParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get task details params
func (o *GetTaskDetailsParams) SetID(id int64) {
	o.ID = id
}

// Validate these params
func (o *GetTaskDetailsParams) Validate(formats strfmt.Registry) error {

	return nil
}

// WriteToRequest writes these params to a swagger request
func (o *GetTaskDetailsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	r.SetTimeout(o.timeout)
	var res []error

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

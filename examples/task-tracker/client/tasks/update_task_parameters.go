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

	"github.com/go-swagger/go-swagger/examples/task-tracker/models"
)

// NewUpdateTaskParams creates a new UpdateTaskParams object
// with the default values initialized.
func NewUpdateTaskParams() *UpdateTaskParams {
	var ()
	return &UpdateTaskParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateTaskParamsWithTimeout creates a new UpdateTaskParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUpdateTaskParamsWithTimeout(timeout time.Duration) *UpdateTaskParams {
	var ()
	return &UpdateTaskParams{

		timeout: timeout,
	}
}

// NewUpdateTaskParamsWithContext creates a new UpdateTaskParams object
// with the default values initialized, and the ability to set a context for a request
func NewUpdateTaskParamsWithContext(ctx context.Context) *UpdateTaskParams {
	var ()
	return &UpdateTaskParams{

		Context: ctx,
	}
}

// NewUpdateTaskParamsWithHTTPClient creates a new UpdateTaskParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUpdateTaskParamsWithHTTPClient(client *http.Client) *UpdateTaskParams {
	var ()
	return &UpdateTaskParams{
		HTTPClient: client,
	}
}

/*UpdateTaskParams contains all the parameters to send to the API endpoint
for the update task operation typically these are written to a http.Request
*/
type UpdateTaskParams struct {

	/*Body
	  The task to update
	  Required: true
	  In: body
	*/
	Body *models.Task
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

// WithTimeout adds the timeout to the update task params
func (o *UpdateTaskParams) WithTimeout(timeout time.Duration) *UpdateTaskParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update task params
func (o *UpdateTaskParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update task params
func (o *UpdateTaskParams) WithContext(ctx context.Context) *UpdateTaskParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update task params
func (o *UpdateTaskParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update task params
func (o *UpdateTaskParams) WithHTTPClient(client *http.Client) *UpdateTaskParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update task params
func (o *UpdateTaskParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the update task params
func (o *UpdateTaskParams) WithBody(body *models.Task) *UpdateTaskParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the update task params
func (o *UpdateTaskParams) SetBody(body *models.Task) {
	o.Body = body
}

// WithID adds the id to the update task params
func (o *UpdateTaskParams) WithID(id int64) *UpdateTaskParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the update task params
func (o *UpdateTaskParams) SetID(id int64) {
	o.ID = id
}

// Validate these params
func (o *UpdateTaskParams) Validate(formats strfmt.Registry) error {

	if o.Body != nil {
		if err := o.Body.Validate(formats); err != nil {
			return err
		}
	}

	return nil
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateTaskParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	r.SetTimeout(o.timeout)
	var res []error

	if o.Body == nil {
		o.Body = new(models.Task)
	}

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	// path param id
	if err := r.SetPathParam("id", swag.FormatInt64(o.ID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

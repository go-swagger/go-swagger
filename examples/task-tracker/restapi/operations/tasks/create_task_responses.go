// Code generated by go-swagger; DO NOT EDIT.

package tasks

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/go-swagger/go-swagger/examples/task-tracker/models"
)

// CreateTaskCreatedCode is the HTTP code returned for type CreateTaskCreated
const CreateTaskCreatedCode int = 201

/*
CreateTaskCreated Task created

swagger:response createTaskCreated
*/
type CreateTaskCreated struct {
	/*URL to the newly added Task

	 */
	Location strfmt.URI `json:"Location"`
}

// NewCreateTaskCreated creates CreateTaskCreated with default headers values
func NewCreateTaskCreated() *CreateTaskCreated {

	return &CreateTaskCreated{}
}

// WithLocation adds the location to the create task created response
func (o *CreateTaskCreated) WithLocation(location strfmt.URI) *CreateTaskCreated {
	o.Location = location
	return o
}

// SetLocation sets the location to the create task created response
func (o *CreateTaskCreated) SetLocation(location strfmt.URI) {
	o.Location = location
}

// WriteResponse to the client
func (o *CreateTaskCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header Location

	location := o.Location.String()
	if location != "" {
		rw.Header().Set("Location", location)
	}

	rw.Header().Del(runtime.HeaderContentType) // Remove Content-Type on empty responses

	rw.WriteHeader(201)
}

/*
CreateTaskDefault Error response

swagger:response createTaskDefault
*/
type CreateTaskDefault struct {
	_statusCode int
	/*

	 */
	XErrorCode string `json:"X-Error-Code"`

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateTaskDefault creates CreateTaskDefault with default headers values
func NewCreateTaskDefault(code int) *CreateTaskDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateTaskDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create task default response
func (o *CreateTaskDefault) WithStatusCode(code int) *CreateTaskDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create task default response
func (o *CreateTaskDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithXErrorCode adds the xErrorCode to the create task default response
func (o *CreateTaskDefault) WithXErrorCode(xErrorCode string) *CreateTaskDefault {
	o.XErrorCode = xErrorCode
	return o
}

// SetXErrorCode sets the xErrorCode to the create task default response
func (o *CreateTaskDefault) SetXErrorCode(xErrorCode string) {
	o.XErrorCode = xErrorCode
}

// WithPayload adds the payload to the create task default response
func (o *CreateTaskDefault) WithPayload(payload *models.Error) *CreateTaskDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create task default response
func (o *CreateTaskDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTaskDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	// response header X-Error-Code

	xErrorCode := o.XErrorCode
	if xErrorCode != "" {
		rw.Header().Set("X-Error-Code", xErrorCode)
	}

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

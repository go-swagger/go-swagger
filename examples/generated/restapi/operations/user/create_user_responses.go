// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

/*CreateUserDefault successful operation

swagger:response createUserDefault
*/
type CreateUserDefault struct {
	_statusCode int
}

// NewCreateUserDefault creates CreateUserDefault with default headers values
func NewCreateUserDefault(code int) *CreateUserDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create user default response
func (o *CreateUserDefault) WithStatusCode(code int) *CreateUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create user default response
func (o *CreateUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WriteResponse to the client
func (o *CreateUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
}

// Code generated by go-swagger; DO NOT EDIT.

package uploads

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	stderrors "errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// UploadFileMaxParseMemory sets the maximum size in bytes for
// the multipart form parser for this operation.
//
// The default value is 32 MB.
// The multipart parser stores up to this + 10MB.
var UploadFileMaxParseMemory int64 = 32 << 20

// NewUploadFileParams creates a new UploadFileParams object
//
// There are no default values defined in the spec.
func NewUploadFileParams() UploadFileParams {

	return UploadFileParams{}
}

// UploadFileParams contains all the bound params for the upload file operation
// typically these are obtained from a http.Request
//
// swagger:parameters uploadFile
type UploadFileParams struct {
	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: formData
	*/
	File io.ReadCloser
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUploadFileParams() beforehand.
func (o *UploadFileParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := r.ParseMultipartForm(UploadFileMaxParseMemory); err != nil {
		if !stderrors.Is(err, http.ErrNotMultipart) {
			return errors.New(400, "%v", err)
		} else if errParse := r.ParseForm(); errParse != nil {
			return errors.New(400, "%v", errParse)
		}
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		res = append(res, errors.New(400, "reading file %q failed: %v", "file", err))
	} else {
		if errBind := o.bindFile(file, fileHeader); errBind != nil {
			// Required: true
			res = append(res, errBind)
		} else {
			o.File = &runtime.File{Data: file, Header: fileHeader}
		}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindFile binds file parameter File.
//
// The only supported validations on files are MinLength and MaxLength
func (o *UploadFileParams) bindFile(file multipart.File, header *multipart.FileHeader) error {
	return nil
}

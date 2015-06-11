package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/casualjim/go-swagger/client"
	"github.com/casualjim/go-swagger/httpkit"
	"github.com/casualjim/go-swagger/strfmt"
)

// NewRequest creates a new swagger http client request
func NewRequest(method, pathPattern string, writer client.RequestWriter) (*Request, error) {
	return &Request{
		pathPattern: pathPattern,
		method:      method,
		writer:      writer,
		header:      http.Header(make(map[string][]string)),
		query:       url.Values(make(map[string][]string)),
	}, nil
}

// Request represents a swagger client request.
//
// This Request struct converts to a HTTP request.
// There might be others that convert to other transports.
// There is no error checking here, it is assumed to be used after a spec has been validated.
// so impossible combinations should not arise (hopefully).
//
// The main purpose of this struct is to hide the machinery of adding params to a transport request.
// The generated code only implements what is necessary to turn a param into a valid value for these methods.
type Request struct {
	pathPattern string
	method      string
	writer      client.RequestWriter
	runtime     *Runtime

	pathParams map[string]string
	header     http.Header
	query      url.Values
	formFields url.Values
	fileFields map[string]*os.File
	payload    interface{}
}

var (
	// ensure interface compliance
	_ client.Request = new(Request)
)

// BuildHTTP creates a new http request based on the data from the params
func (r *Request) BuildHTTP(producer httpkit.Producer, registry strfmt.Registry) (*http.Request, error) {
	// build the data
	if err := r.writer.WriteToRequest(r, registry); err != nil {
		return nil, err
	}

	// create http request
	path := r.pathPattern
	for k, v := range r.pathParams {
		path = strings.Replace(path, "{"+k+"}", v, -1)
	}

	// TODO: Support uploading huge bodies!
	// Not too excited about the buffer here, but it keeps me going for now
	body := bytes.NewBuffer(nil)
	req, err := http.NewRequest(r.method, path, body)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = r.query.Encode()
	req.Header = r.header

	// check if this is a form type request
	if r.formFields != nil {
		// check if this is multipart
		if r.fileFields != nil {

			mp := multipart.NewWriter(body)
			defer mp.Close()
			req.Header.Set(httpkit.HeaderContentType, mp.FormDataContentType())

			for fn, v := range r.formFields {
				if len(v) > 0 {
					if err := mp.WriteField(fn, v[0]); err != nil {
						return nil, err
					}
				}
			}

			for fn, f := range r.fileFields {
				wrtr, err := mp.CreateFormFile(fn, filepath.Base(f.Name()))
				if err != nil {
					return nil, err
				}
				defer func() {
					for _, ff := range r.fileFields {
						ff.Close()
					}
				}()
				// TODO: Support uploading huge files!
				// Stop this copy insanity!
				if _, err := io.Copy(wrtr, f); err != nil {
					return nil, err
				}
			}

			return req, nil
		}
		body.WriteString(r.formFields.Encode())
		return req, nil
	}

	// write the form values as body
	// if there is payload, use the producer to write the payload
	if err := producer.Produce(body, r.payload); err != nil {
		return nil, err
	}
	return req, nil
}

// SetHeaderParam adds a header param to the request
// when there is only 1 value provided for the varargs, it will set it.
// when there are several values provided for the varargs it will add it (no overriding)
func (r *Request) SetHeaderParam(name string, values ...string) error {
	if r.header == nil {
		r.header = make(map[string][]string)
	}
	r.header[name] = values
	return nil
}

// SetQueryParam adds a query param to the request
// when there is only 1 value provided for the varargs, it will set it.
// when there are several values provided for the varargs it will add it (no overriding)
func (r *Request) SetQueryParam(name string, values ...string) error {
	if r.header == nil {
		r.query = make(map[string][]string)
	}
	r.query[name] = values
	return nil
}

// SetFormParam adds a forn param to the request
// when there is only 1 value provided for the varargs, it will set it.
// when there are several values provided for the varargs it will add it (no overriding)
func (r *Request) SetFormParam(name string, values ...string) error {
	if r.formFields == nil {
		r.formFields = make(map[string][]string)
	}
	r.formFields[name] = values
	return nil
}

// SetPathParam adds a path param to the request
func (r *Request) SetPathParam(name string, value string) error {
	if r.pathParams == nil {
		r.pathParams = make(map[string]string)
	}

	r.pathParams[name] = value
	return nil
}

// SetFileParam adds a file param to the request
func (r *Request) SetFileParam(name string, toSend string) error {
	file, err := os.Open(toSend)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return fmt.Errorf("%q is a directory, only files are supported", toSend)
	}

	if r.fileFields == nil {
		r.fileFields = make(map[string]*os.File)
	}
	if r.formFields == nil {
		r.formFields = url.Values(make(map[string][]string))
	}

	r.fileFields[name] = file
	return nil
}

// SetBodyParam sets a body parameter on the request.
// This does not yet serialze the object, this happens as late as possible.
func (r *Request) SetBodyParam(payload interface{}) error {
	r.payload = payload
	return nil
}

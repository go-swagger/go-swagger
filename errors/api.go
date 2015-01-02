package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Error represents a error interface all swagger framework errors implement
type Error interface {
	error
	Code() int32
}

type apiError struct {
	code    int32
	message string
}

func (a *apiError) Error() string {
	return a.message
}

func (a *apiError) Code() int32 {
	return a.code
}

// New creates a new API error with a code and a message
func New(code int32, message string) Error {
	return &apiError{code, message}
}

// NotFound creates a new not found error
func NotFound(message string, args ...interface{}) Error {
	if message == "" {
		message = "Not found"
	}
	return New(http.StatusNotFound, fmt.Sprintf(message, args...))
}

// NotImplemented creates a new not implemented error
func NotImplemented(message string) Error {
	return New(http.StatusNotImplemented, message)
}

// MethodNotAllowedError represents an error for when the path matches but the method doesn't
type MethodNotAllowedError struct {
	code    int32
	Allowed []string
	message string
}

func (m *MethodNotAllowedError) Error() string {
	return m.message
}

// Code the error code
func (m *MethodNotAllowedError) Code() int32 {
	return m.code
}

func errorAsJSON(err Error) []byte {
	b, _ := json.Marshal(struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}{err.Code(), err.Error()})
	return b
}

// MethodNotAllowed creates a new method not allowed error
func MethodNotAllowed(requested string, allow []string) Error {
	msg := fmt.Sprintf("method %s is not allowed, but [%s] are", requested, strings.Join(allow, ","))
	return &MethodNotAllowedError{code: http.StatusMethodNotAllowed, Allowed: allow, message: msg}
}

// ServeError the error handler interface implemenation
func ServeError(rw http.ResponseWriter, r *http.Request, err error) {
	switch err.(type) {
	case *MethodNotAllowedError:
		e := err.(*MethodNotAllowedError)
		rw.Header().Set("content-type", "application/json")
		rw.Header().Add("Allow", strings.Join(err.(*MethodNotAllowedError).Allowed, ","))
		rw.WriteHeader(int(e.Code()))
		rw.Write(errorAsJSON(e))
	case Error:
		rw.Header().Set("content-type", "application/json")
		rw.WriteHeader(int(err.(Error).Code()))
		rw.Write(errorAsJSON(err.(Error)))
	default:
		rw.Header().Set("content-type", "application/json")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write(errorAsJSON(New(http.StatusInternalServerError, err.Error())))
	}

}

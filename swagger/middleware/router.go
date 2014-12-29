package middleware

import (
	"net/http"

	"github.com/casualjim/go-swagger/swagger/errors"
)

// NewRouter creates a new router middleware function
func NewRouter(context *Context) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// use context to lookup routes
		if _, ok := context.RouteInfo(r); ok {
			next(rw, r)
			return
		}

		// Not found, check if it exists in the other methods first
		if others := context.AllowedMethods(r); len(others) > 0 {
			context.Respond(rw, r, errors.MethodNotAllowed(r.Method, others))
			return
		}
		context.Respond(rw, r, errors.NotFound("path %s was not found", r.URL.Path))
	}
}

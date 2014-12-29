package middleware

import (
	"net/http"

	"github.com/casualjim/go-swagger/swagger/errors"
	"github.com/gorilla/context"
)

func newRouter(ctx *Context) func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		defer context.Clear(r)
		// use context to lookup routes
		if _, ok := ctx.RouteInfo(r); ok {
			next(rw, r)
			return
		}

		// Not found, check if it exists in the other methods first
		if others := ctx.AllowedMethods(r); len(others) > 0 {
			ctx.Respond(rw, r, errors.MethodNotAllowed(r.Method, others))
			return
		}
		ctx.Respond(rw, r, errors.NotFound("path %s was not found", r.URL.Path))
	}
}
